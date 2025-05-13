package config

import (
	"os"

	"k8s.io/apimachinery/pkg/api/resource"
)

type Config struct {
	VirtioFSImage    string
	ResourceRequests struct {
		CPU    resource.Quantity
		Memory resource.Quantity
	}
	ResourceLimits struct {
		CPU    resource.Quantity
		Memory resource.Quantity
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	cfg.VirtioFSImage = getEnvOrDefault("VIRTIOFS_IMAGE", "quay.io/kubevirt/virt-launcher:v1.5.1")

	cfg.ResourceRequests.CPU = resource.MustParse(getEnvOrDefault("RESOURCE_REQUESTS_CPU", "100m"))
	cfg.ResourceRequests.Memory = resource.MustParse(getEnvOrDefault("RESOURCE_REQUESTS_MEMORY", "128Mi"))

	cfg.ResourceLimits.CPU = resource.MustParse(getEnvOrDefault("RESOURCE_LIMITS_CPU", "200m"))
	cfg.ResourceLimits.Memory = resource.MustParse(getEnvOrDefault("RESOURCE_LIMITS_MEMORY", "256Mi"))

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
