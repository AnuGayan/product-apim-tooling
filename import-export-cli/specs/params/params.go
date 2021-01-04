package params

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
)

// Configuration represents endpoint config
type Configuration struct {
	// RetryTimeOut for endpoint
	RetryTimeOut *int `yaml:"retryTimeOut" json:"retryTimeOut,string"`
	// RetryDelay for endpoint
	RetryDelay *int `yaml:"retryDelay" json:"retryDelay,string"`
	// Factor used for config
	Factor *int `yaml:"factor" json:"factor,string"`
	// RetryErroCode used for config
	RetryErroCode *int `yaml:"retryErroCode" json:"retryErroCode,string"`
	// SuspendErrorCode used for config
	SuspendErrorCode *int `yaml:"suspendErrorCode" json:"suspendErrorCode,string"`
	// SuspendDuration used for config
	SuspendDuration *int `yaml:"suspendDuration" json:"suspendDuration,string"`
	// SuspendMaxDuration used for config
	SuspendMaxDuration *int `yaml:"suspendMaxDuration" json:"suspendMaxDuration,string"`
	// ActionSelect used for config
	ActionSelect *string `yaml:"actionSelect" json:"actionSelect"`
	// ActionDuration used for config
	ActionDuration *int `yaml:"actionDuration" json:"actionDuration,string"`
}

// Endpoint details
type Endpoint struct {
	// Url of the endpoint
	Url *string `yaml:"url" json:"url"`
	// Config of endpoint
	Config *Configuration `yaml:"config" json:"config"`
}

// EndpointData contains details about endpoints
type EndpointData struct {
	// Production endpoint
	Production *Endpoint `yaml:"production" json:"production_endpoints,omitempty"`
	// Sandbox endpoint
	Sandbox *Endpoint `yaml:"sandbox" json:"sandbox_endpoints,omitempty"`
}

// LoadBalanceEndpointsData contains details about endpoints mainly to be used in load balancing
type LoadBalanceEndpointsData struct {
	EndpointType string `yaml:"endpoint_type" json:"endpoint_type"`
	// Production endpoints list for load balancing
	Production []Endpoint `yaml:"production" json:"production_endpoints,omitempty"`
	// Sandbox endpoints list for load balancing
	Sandbox []Endpoint `yaml:"sandbox" json:"sandbox_endpoints,omitempty"`
	// Session management method from the load balancing group. Values can be "none", "transport" (by default), "soap", "simpleClientSession" (Client ID)
	SessionManagement string `yaml:"sessionManagement" json:"sessionManagement,omitempty"`
	// Session timeout means the number of milliseconds after which the session would time out
	SessionTimeout int `yaml:"sessionTimeOut" json:"sessionTimeOut,omitempty"`
	// Class name for algorithm to be used if load balancing should be done
	AlgorithmClassName string `yaml:"algoClassName" json:"algoClassName,omitempty"`
	// Combo for algorithm to be used if load balancing should be done
	AlgorithmCombName string `yaml:"algoCombo" json:"algoCombo,omitempty"`
}

// FailoverEndpointsData contains details about endpoints mainly to be used in load balancing
type FailoverEndpointsData struct {
	EndpointType string `yaml:"endpoint_type" json:"endpoint_type"`
	// Primary production endpoint for failover
	Production *Endpoint `yaml:"production" json:"production_endpoints,omitempty"`
	// Production failover endpoints list for failover
	ProductionFailovers []Endpoint `yaml:"productionFailovers" json:"production_failovers,omitempty"`
	// Primary sandbox endpoint for failover
	Sandbox *Endpoint `yaml:"sandbox" json:"sandbox_endpoints,omitempty"`
	// Production failover endpoints list for failover endpoint types
	SandboxFailovers []Endpoint `yaml:"sandboxFailovers" json:"sandbox_failovers,omitempty"`
	// To enable failover endpoints
	Failover bool `yaml:"failOver" json:"failOver,omitempty"`
}

// SecurityData contains the details about endpoint security from api_params.yaml
type SecurityData struct {
	// Decides whether the endpoint security is enabled
	Enabled string `yaml:"enabled" json:"enabled,omitempty"`
	// Type of the endpoint security (can be Basic or Digest)
	Type string `yaml:"type" json:"type,omitempty"`
	// Username for the endpoint
	Username string `yaml:"username" json:"username,omitempty"`
	// Password for the endpoint
	Password string `yaml:"password" json:"password,omitempty"`
}

// Cert stores certificate details
type Cert struct {
	// Host of the certificate
	Host string `yaml:"host" json:"hostName"`
	// Alias for certificate
	Alias string `yaml:"alias" json:"alias"`
	// Path for certificate file
	Path string `yaml:"path" json:"-"`
	// Certificate is used for internal purposes, it contains secret in base64
	Certificate string `json:"certificate"`
}

// MutualSslCert stores mutualssl certificate details
type MutualSslCert struct {
	// TierName of the certificate (eg:- Unlimited, Gold, Silver, Bronze)
	TierName string `yaml:"tierName" json:"tierName"`
	// Alias for certificate
	Alias string `yaml:"alias" json:"alias"`
	// Path for certificate file
	Path string `yaml:"path" json:"-"`
	// Certificate is used for internal purposes, it contains secret in base64
	Certificate string `json:"certificate"`
	// ApiIdentifier is used for internal purposes, it contains details of the API to be stored in client_certificates file
	APIIdentifier APIIdentifier `json:"apiIdentifier"`
}

// ApiIdentifier stores API Identifier details
type APIIdentifier struct {
	// Name of the provider of the API
	ProviderName string `json:"providerName"`
	// Name of the API
	APIName string `json:"apiName"`
	// Version of the API
	Version string `json:"version"`
}

// Environment represents an api environment
type Environment struct {
	// Name of the environment
	Name string `yaml:"name"`
	// Endpoints contain details about endpoints in a configuration
	Endpoints *EndpointData `yaml:"endpoints"`
	// LoadBalanceEndpoints contain details about endpoints in a configuration for load balancing scenarios
	LoadBalanceEndpoints *LoadBalanceEndpointsData `yaml:"loadBalanceEndpoints"`
	// FailoverEndpoints contain details about endpoints in a configuration for failover scenarios
	FailoverEndpoints *FailoverEndpointsData `yaml:"failoverEndpoints"`
	// Security contains the details about endpoint security
	Security *SecurityData `yaml:"security"`
	// GatewayEnvironments contains environments that used to deploy API
	GatewayEnvironments []string `yaml:"gatewayEnvironments"`
	// Certs for environment
	Certs          []Cert          `yaml:"certs"`
	MutualSslCerts []MutualSslCert `yaml:"mutualSslCerts"`
	// Policies contains the available subscription policies in an environment that can be enforced to an API
	Policies []string `yaml:"policies"`
}

// ApiParams represents environments defined in configuration file
type ApiParams struct {
	// Environments contains all environments in a configuration
	Environments []Environment `yaml:"environments"`
}

// APIEndpointConfig contains details about endpoints in an API
type APIEndpointConfig struct {
	// EPConfig is representing endpoint configuration
	EPConfig string `json:"endpointConfig"`
}

// LoadApiParams loads an configuration from a reader. It returns an error or a valid ApiParams
func LoadApiParams(r io.Reader) (*ApiParams, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	str, err := utils.EnvSubstitute(string(data))
	if err != nil {
		return nil, err
	}

	apiParams := &ApiParams{}
	err = yaml.Unmarshal([]byte(str), &apiParams)
	if err != nil {
		return nil, err
	}

	return apiParams, nil
}

// LoadApiParamsFromFile loads a configuration YAML file located in path. It returns an error or a valid ApiParams
func LoadApiParamsFromFile(path string) (*ApiParams, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	apiConfig, err := LoadApiParams(r)
	_ = r.Close()

	return apiConfig, err
}

// ExtractAPIEndpointConfig extracts API endpoint information from a slice of byte b
func ExtractAPIEndpointConfig(b []byte) (string, error) {
	apiConfig := &APIEndpointConfig{}
	err := json.Unmarshal(b, &apiConfig)
	if err != nil {
		return "", err
	}

	return apiConfig.EPConfig, err
}

// GetEnv returns the EndpointData associated for key in the ApiParams, if not found returns nil
func (config ApiParams) GetEnv(key string) *Environment {
	for index, env := range config.Environments {
		if env.Name == key {
			return &config.Environments[index]
		}
	}
	return nil
}
