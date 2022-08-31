package config

import (
	"context"
	"net/http"
	"time"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type Config interface {
	RetryConfig
	GetCertificatesPath() string
	GetThreads() int
	IsDryRun() bool
	GetServiceDetails() auth.ServiceDetails
	GetLogger() log.Log
	IsInsecureTls() bool
	GetContext() context.Context
	GetHttpTimeout() time.Duration
	GetHttpClient() *http.Client
}

type servicesConfig struct {
	auth.ServiceDetails
	serviceRetryConfig
	certificatesPath string
	dryRun           bool
	threads          int
	logger           log.Log
	insecureTls      bool
	ctx              context.Context
	httpTimeout      time.Duration
	httpClient       *http.Client
}

func (config *servicesConfig) IsDryRun() bool {
	return config.dryRun
}

func (config *servicesConfig) GetCertificatesPath() string {
	return config.certificatesPath
}

func (config *servicesConfig) GetThreads() int {
	return config.threads
}

func (config *servicesConfig) GetServiceDetails() auth.ServiceDetails {
	return config.ServiceDetails
}

func (config *servicesConfig) GetLogger() log.Log {
	return config.logger
}

func (config *servicesConfig) IsInsecureTls() bool {
	return config.insecureTls
}

func (config *servicesConfig) GetContext() context.Context {
	return config.ctx
}

func (config *servicesConfig) GetHttpTimeout() time.Duration {
	return config.httpTimeout
}

func (config *servicesConfig) GetHttpClient() *http.Client {
	return config.httpClient
}
