package function

import "time"

type FunctionSpec struct {
	FnImage   string        `yaml:"fn_image"`
	Functions []FunctionDef `yaml:"functions"`
}

type FunctionDef struct {
	FnName                    string      `yaml:"name"`
	ExistingFunctionID        string      `yaml:"existingFunctionID,omitempty"`
	InferenceURL              string      `yaml:"inferenceUrl"`
	InferencePort             int64       `yaml:"inferencePort,omitempty"`
	HealthUri                 string      `yaml:"healthUri,omitempty"`
	ContainerImage            string      `yaml:"containerImage,omitempty"`
	ContainerArgs             string      `yaml:"containerArgs,omitempty"`
	Custom                    bool        `yaml:"custom,omitempty"`
	Description               string      `yaml:"description,omitempty"`
	Streaming                 bool        `yaml:"streaming,omitempty"`
	Tags                      []string    `yaml:"tags,omitempty"`
	Health                    HealthCheck `yaml:"health"`
	InstBackend               string      `yaml:"inst_backend"`
	InstGPUType               string      `yaml:"inst_gpu_type"`
	InstType                  string      `yaml:"inst_type"`
	InstMin                   int64       `yaml:"inst_min,omitempty"`
	InstMax                   int64       `yaml:"inst_max,omitempty"`
	InstMaxRequestConcurrency int64       `yaml:"inst_max_request_concurrency,omitempty"`
	ContainerEnvironment      []EnvVar    `yaml:"containerEnvironment,omitempty"`
	Models                    []ModelDef  `yaml:"models,omitempty"`
}

type HealthCheck struct {
	Protocol           string        `yaml:"protocol,omitempty"`
	Port               int64         `yaml:"port,omitempty"`
	Timeout            time.Duration `yaml:"timeout,omitempty"`
	ExpectedStatusCode int64         `yaml:"expectedStatusCode,omitempty"`
	Uri                string        `yaml:"uri,omitempty"`
}

type EnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type ModelDef struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Uri     string `yaml:"uri"`
}

func (h *HealthCheck) UnmarshalYAML(unmarshal func(interface{}) error) error {
	aux := &struct {
		Protocol           string `yaml:"protocol"`
		Port               int64  `yaml:"port"`
		Timeout            int    `yaml:"timeout"`
		ExpectedStatusCode int64  `yaml:"expectedStatusCode"`
		Uri                string `yaml:"uri"`
	}{}

	if err := unmarshal(aux); err != nil {
		return err
	}

	h.Protocol = aux.Protocol
	h.Port = aux.Port
	h.Timeout = time.Duration(aux.Timeout) * time.Second
	h.ExpectedStatusCode = aux.ExpectedStatusCode
	h.Uri = aux.Uri

	return nil
}

type LogParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LogsPayload struct {
	Parameters []LogParameter `json:"parameters"`
}

type NVCFLogsResponse struct {
	JobID    string    `json:"jobId"`
	Data     []NVCFLog `json:"data"`
	Metadata struct {
		TotalPages int `json:"totalPages"`
		Page       int `json:"page"`
	} `json:"metadata"`
}

type NVCFLog struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type DeploymentLog struct {
	Level     string
	Timestamp string
	Message   string
}
