package function

type FunctionSpec struct {
	FnImage   string        `yaml:"fn_image"`
	Functions []FunctionDef `yaml:"functions"`
}

type FunctionDef struct {
	FnName                    string      `yaml:"fn_name"`
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
	Protocol           string `yaml:"protocol,omitempty"`
	Port               int64  `yaml:"port,omitempty"`
	Timeout            int64  `yaml:"timeout,omitempty"`
	ExpectedStatusCode int64  `yaml:"expectedStatusCode,omitempty"`
	Uri                string `yaml:"uri,omitempty"`
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
