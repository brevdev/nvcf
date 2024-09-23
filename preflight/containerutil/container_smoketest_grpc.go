package containerutil

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func (cst *ContainerSmokeTest) CheckGRPCHealthEndpoint(secondsToWaitForHealthy int) error {
	healthURL := fmt.Sprintf("localhost:%s", cst.HostPort)
	fmt.Printf("Looking for gRPC health signal at %s\n", healthURL)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(secondsToWaitForHealthy)*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, healthURL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	healthClient := grpc_health_v1.NewHealthClient(conn)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("gRPC health check did not complete successfully in time")
		case <-ticker.C:
			resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
			if err == nil && resp.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING {
				fmt.Printf("%s is SERVING\n", healthURL)
				return nil
			}
			fmt.Printf("Health check status: %v, error: %v\n", resp.GetStatus(), err)
		}
	}
}

func (cst *ContainerSmokeTest) TestGRPCInference(serviceName, methodName string, inputData map[string]interface{}, useAlphaReflection bool) error {
	inferenceURL := fmt.Sprintf("localhost:%s", cst.HostPort)
	fmt.Printf("Connecting to gRPC server at %s...\n", inferenceURL)

	conn, err := grpc.Dial(inferenceURL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connected to gRPC server successfully")

	reflectClient := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := reflectClient.ServerReflectionInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to create reflection stream: %v", err)
	}

	// List services
	err = stream.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{},
	})
	if err != nil {
		return fmt.Errorf("failed to send list services request: %v", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive list services response: %v", err)
	}

	fmt.Printf("Available services: %+v\n", resp.GetListServicesResponse().Service)

	// Find the correct service name
	var fullServiceName string
	for _, service := range resp.GetListServicesResponse().Service {
		if service.Name == serviceName || service.Name == fmt.Sprintf("%s.%s", serviceName, serviceName) {
			fullServiceName = service.Name
			break
		}
	}

	if fullServiceName == "" {
		return fmt.Errorf("service %s not found", serviceName)
	}

	// Now request file descriptor for the service
	err = stream.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol{
			FileContainingSymbol: fullServiceName,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send file descriptor request: %v", err)
	}

	resp, err = stream.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive file descriptor response: %v", err)
	}

	// Parse the file descriptor proto
	var fileDescProto descriptorpb.FileDescriptorProto
	if err := proto.Unmarshal(resp.GetFileDescriptorResponse().FileDescriptorProto[0], &fileDescProto); err != nil {
		return fmt.Errorf("failed to unmarshal file descriptor proto: %v", err)
	}

	fd, err := protodesc.NewFile(&fileDescProto, nil)
	if err != nil {
		return fmt.Errorf("failed to create file descriptor: %v", err)
	}
	// Prepare the gRPC method to be called
	method := fmt.Sprintf("/%s/%s", fullServiceName, methodName)
	fmt.Printf("Calling gRPC method: %s\n", method)

	// Get method descriptors
	serviceDesc := fd.Services().ByName(protoreflect.Name(serviceName))
	if serviceDesc == nil {
		return fmt.Errorf("service %s not found in file descriptor", serviceName)
	}
	methodDesc := serviceDesc.Methods().ByName(protoreflect.Name(methodName))
	if methodDesc == nil {
		return fmt.Errorf("method %s not found in service %s", methodName, serviceName)
	}

	// Create dynamic messages for input and output
	inputMsg := dynamicpb.NewMessage(methodDesc.Input())
	outputMsg := dynamicpb.NewMessage(methodDesc.Output())

	// Set fields in the input message
	for fieldName, fieldValue := range inputData {
		field := inputMsg.Descriptor().Fields().ByName(protoreflect.Name(fieldName))
		if field == nil {
			return fmt.Errorf("field %s not found in input message", fieldName)
		}
		if err := setDynamicMessageField(inputMsg, field, fieldValue); err != nil {
			return fmt.Errorf("failed to set field %s: %v", fieldName, err)
		}
	}

	fmt.Printf("Input message: %+v\n", inputMsg)

	// Invoke the gRPC method
	err = conn.Invoke(ctx, method, inputMsg, outputMsg)
	if err != nil {
		return fmt.Errorf("failed to call %s: %v", method, err)
	}

	fmt.Printf("Raw response: %+v\n", outputMsg)

	// Print all fields of the output message
	outputMsg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		fmt.Printf("Field: %s, Value: %v\n", fd.Name(), v.Interface())
		return true
	})

	return nil
}

func setDynamicMessageField(msg *dynamicpb.Message, field protoreflect.FieldDescriptor, value interface{}) error {
	switch field.Kind() {
	case protoreflect.StringKind:
		if strVal, ok := value.(string); ok {
			msg.Set(field, protoreflect.ValueOfString(strVal))
		} else {
			return fmt.Errorf("expected string value for field %s", field.Name())
		}
	case protoreflect.Int32Kind, protoreflect.Int64Kind:
		if intVal, ok := value.(int64); ok {
			msg.Set(field, protoreflect.ValueOfInt64(intVal))
		} else if intVal, ok := value.(int); ok {
			msg.Set(field, protoreflect.ValueOfInt64(int64(intVal)))
		} else {
			return fmt.Errorf("expected int value for field %s", field.Name())
		}
	case protoreflect.BoolKind:
		if boolVal, ok := value.(bool); ok {
			msg.Set(field, protoreflect.ValueOfBool(boolVal))
		} else {
			return fmt.Errorf("expected bool value for field %s", field.Name())
		}
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		if floatVal, ok := value.(float64); ok {
			msg.Set(field, protoreflect.ValueOfFloat64(floatVal))
		} else {
			return fmt.Errorf("expected float value for field %s", field.Name())
		}
	// Add more cases for other protobuf types as needed
	default:
		return fmt.Errorf("unsupported field type %s for field %s", field.Kind(), field.Name())
	}
	return nil
}