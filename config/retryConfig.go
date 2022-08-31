package config

type RetryConfig interface {
	GetHttpRetries() int
	GetHttpRetryWaitMilliSecs() int
	GetHttpRetryStop() *bool
}

type serviceRetryConfig struct {
	httpRetries            int
	httpRetryWaitMilliSecs int
	httpRetryStop          *bool
}

func (config *serviceRetryConfig) GetHttpRetries() int {
	return config.httpRetries
}

func (config *serviceRetryConfig) GetHttpRetryWaitMilliSecs() int {
	return config.httpRetryWaitMilliSecs
}

func (config *serviceRetryConfig) GetHttpRetryStop() *bool {
	return config.httpRetryStop
}
