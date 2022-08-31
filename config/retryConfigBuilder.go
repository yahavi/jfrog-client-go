package config

func NewRetryConfigBuilder() *retryConfigBuilder {
	var httpRetryStop bool
	retryConfigBuilder := &retryConfigBuilder{
		serviceRetryConfig: serviceRetryConfig{
			httpRetries:            3,
			httpRetryWaitMilliSecs: 0,
			httpRetryStop:          &httpRetryStop,
		},
	}
	return retryConfigBuilder
}

type retryConfigBuilder struct {
	serviceRetryConfig
}

func (builder *retryConfigBuilder) SetHttpRetries(httpRetries int) *retryConfigBuilder {
	builder.httpRetries = httpRetries
	return builder
}

func (builder *retryConfigBuilder) SetHttpRetryWaitMilliSecs(httpRetryWaitMilliSecs int) *retryConfigBuilder {
	builder.httpRetryWaitMilliSecs = httpRetryWaitMilliSecs
	return builder
}

func (builder *retryConfigBuilder) SetHttpRetryStop(httpRetryStop *bool) *retryConfigBuilder {
	builder.httpRetryStop = httpRetryStop
	return builder
}

func (builder *retryConfigBuilder) BuildRetryConfig() *serviceRetryConfig {
	retryConfig := builder.serviceRetryConfig
	return &retryConfig
}
