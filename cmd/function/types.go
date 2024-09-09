package function

type FunctionSpec struct {
	FnImage   string        `yaml:"fn_image"`
	Functions []FunctionDef `yaml:"functions"`
}

type FunctionDef struct {
	FnName                    string      `yaml:"fn_name"`
	ContainerArgs             string      `yaml:"containerArgs"`
	InferenceURL              string      `yaml:"inferenceURL"`
	InferencePort             int64       `yaml:"inferencePort"`
	APIBodyFormat             string      `yaml:"apiBodyFormat,omitempty"`
	FunctionType              string      `yaml:"functionType,omitempty"`
	Description               string      `yaml:"description"`
	Tags                      []string    `yaml:"tags"`
	Health                    HealthCheck `yaml:"health"`
	Env                       []EnvVar    `yaml:"env"`
	InstBackend               string      `yaml:"inst_backend"`
	InstGPUType               string      `yaml:"inst_gpu_type"`
	InstMax                   int64       `yaml:"inst_max"`
	InstMaxRequestConcurrency int64       `yaml:"inst_max_request_concurrency"`
	InstMin                   int64       `yaml:"inst_min"`
	InstType                  string      `yaml:"inst_type"`
	Models                    []ModelDef  `yaml:"models"`
}

type HealthCheck struct {
	URI        string `yaml:"uri"`
	Protocol   string `yaml:"protocol"`
	Port       int64  `yaml:"port"`
	Timeout    string `yaml:"timeout"`
	StatusCode int64  `yaml:"statusCode"`
}

type EnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type ModelDef struct {
	Name    string `yaml:"name"`
	URI     string `yaml:"uri"`
	Version string `yaml:"version"`
}
