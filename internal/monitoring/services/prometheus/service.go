/*
Copyright 2022 Nethermind

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package prometheus

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/NethermindEth/sedge/internal/monitoring"
	"github.com/NethermindEth/sedge/internal/monitoring/data"
	"github.com/NethermindEth/sedge/internal/monitoring/services/types"
	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

//go:embed config
var config embed.FS

//go:embed rules
var rules embed.FS

//go:embed alertmanager
var alertmanager embed.FS

// Config represents the Prometheus configuration.
type Config struct {
	Global        GlobalConfig   `yaml:"global"`
	RuleFiles     []string       `yaml:"rule_files"`
	Alerting      AlertingConfig `yaml:"alerting"`
	ScrapeConfigs []ScrapeConfig `yaml:"scrape_configs"`
}

// AlertingConfig represents the alerting configuration for Prometheus.
type AlertingConfig struct {
	Alertmanagers []AlertmanagerConfig `yaml:"alertmanagers"`
}

// AlertmanagerConfig represents the configuration for an Alertmanager.
type AlertmanagerConfig struct {
	StaticConfigs []StaticConfig `yaml:"static_configs"`
}

// GlobalConfig represents the global configuration for Prometheus.
type GlobalConfig struct {
	ScrapeInterval     string `yaml:"scrape_interval"`
	EvaluationInterval string `yaml:"evaluation_interval"`
}

// ScrapeConfig represents the configuration for a Prometheus scrape job.
type ScrapeConfig struct {
	JobName       string         `yaml:"job_name"`
	StaticConfigs []StaticConfig `yaml:"static_configs"`
	MetricsPath   string         `yaml:"metrics_path,omitempty"`
}

// StaticConfig represents the static configuration for a Prometheus scrape job.
type StaticConfig struct {
	Targets []string          `yaml:"targets"`
	Labels  map[string]string `yaml:"labels,omitempty"`
}

// Verify that PrometheusService implements the ServiceAPI interface.
var _ monitoring.ServiceAPI = &PrometheusService{}

// PrometheusService implements the ServiceAPI interface for a Prometheus service.
type PrometheusService struct {
	stack       *data.MonitoringStack
	containerIP net.IP
	port        uint16
}

// NewPrometheus creates a new PrometheusService.
func NewPrometheus() *PrometheusService {
	return &PrometheusService{}
}

// Init initializes the Prometheus service with the given options.
func (p *PrometheusService) Init(opts types.ServiceOptions) error {
	// Validate dotEnv
	promPort, ok := opts.Dotenv["PROM_PORT"]
	if !ok {
		return fmt.Errorf("%w: %s missing in options", ErrInvalidOptions, "PROM_PORT")
	} else if promPort == "" {
		return fmt.Errorf("%w: %s can't be empty", ErrInvalidOptions, "PROM_PORT")
	}

	port, err := strconv.ParseUint(opts.Dotenv["PROM_PORT"], 10, 16)
	if err != nil {
		return fmt.Errorf("%w: %s is not a valid port", ErrInvalidOptions, "PROM_PORT")
	}
	p.port = uint16(port)
	p.stack = opts.Stack
	return nil
}

// AddTarget adds a new target to the Prometheus config and reloads the Prometheus configuration.
// Assumes endpoint is in the form http://<ip/domain>:<port>
func (p *PrometheusService) AddTarget(target types.MonitoringTarget, labels map[string]string, jobName string) error {
	path := filepath.Join("prometheus", "prometheus.yml")
	// Read the existing config
	rawConfig, err := p.stack.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal the YAML data into the Config struct
	var config Config
	if err = yaml.Unmarshal(rawConfig, &config); err != nil {
		return err
	}

	// Add a new job for the new endpoint
	// Check if the job already exists
	for _, job := range config.ScrapeConfigs {
		if job.JobName == jobName {
			// There is no need to add the job if it already exists
			return nil
		}
	}

	// Default to /metrics if no path is provided
	metricsPath := "/metrics"
	if target.Path != "" {
		metricsPath = target.Path
	}
	job := ScrapeConfig{
		JobName: jobName,
		StaticConfigs: []StaticConfig{
			{
				Targets: []string{target.Endpoint()},
				Labels:  labels,
			},
		},
		MetricsPath: metricsPath,
	}
	config.ScrapeConfigs = append(config.ScrapeConfigs, job)

	// Marshal the updated config back to YAML
	newConfig, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	// Write the updated YAML data back to the file
	if err = p.stack.WriteFile(path, newConfig); err != nil {
		return err
	}

	// Reload the config
	if err = p.reloadConfig(); err != nil {
		return err
	}

	return nil
}

// RemoveTarget removes a target from the Prometheus config and reloads the Prometheus configuration.
func (p *PrometheusService) RemoveTarget(instanceID string) (string, error) {
	path := filepath.Join("prometheus", "prometheus.yml")
	// Read the existing config
	rawConfig, err := p.stack.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Unmarshal the YAML data into the Config struct
	var config Config
	if err = yaml.Unmarshal(rawConfig, &config); err != nil {
		return "", err
	}

	// Remove the target from the jobs
	var network string
	config.ScrapeConfigs = funk.Filter(config.ScrapeConfigs, func(job ScrapeConfig) bool {
		if strings.Contains(job.JobName, instanceID) {
			network = strings.Split(strings.TrimPrefix(job.JobName, instanceID), "++")[1]
			return false
		}
		return true
	}).([]ScrapeConfig)

	// Check if the target was removed
	if network == "" {
		// The target was not removed because it was not in the targets
		return "", fmt.Errorf("%w: %s", monitoring.ErrNonexistingTarget, instanceID)
	}

	// Marshal the updated config back to YAML
	newConfig, err := yaml.Marshal(&config)
	if err != nil {
		return network, err
	}

	// Write the updated YAML data back to the file
	if err = p.stack.WriteFile(path, newConfig); err != nil {
		return network, err
	}

	// Reload the config
	if err = p.reloadConfig(); err != nil {
		return network, err
	}

	return network, nil
}

// DotEnv returns the dotenv variables and default values for the Prometheus service.
func (p *PrometheusService) DotEnv() map[string]string {
	return dotEnv
}

// Setup sets up the Prometheus service configuration files with the given dotenv values.
func (p *PrometheusService) Setup(options map[string]string) error {
	// Validate options
	nodeExporterPort, ok := options["NODE_EXPORTER_PORT"]
	if !ok {
		return fmt.Errorf("%w: %s missing in options", ErrInvalidOptions, "NODE_EXPORTER_PORT")
	} else if nodeExporterPort == "" {
		return fmt.Errorf("%w: %s can't be empty", ErrInvalidOptions, "NODE_EXPORTER_PORT")
	}

	// Read config from the embedded FS
	rawConfig, err := config.ReadFile("config/prometheus.yml")
	if err != nil {
		return err
	}

	// Unmarshal the YAML data into the Config struct
	var config Config
	if err = yaml.Unmarshal(rawConfig, &config); err != nil {
		return err
	}

	// Add node exporter target
	endpoint := fmt.Sprintf("%s:%s", monitoring.NodeExporterContainerName, options["NODE_EXPORTER_PORT"])
	config.ScrapeConfigs = []ScrapeConfig{
		{
			JobName: endpoint,
			StaticConfigs: []StaticConfig{
				{
					Targets: []string{endpoint},
				},
			},
		},
	}

	// Marshal the updated config back to YAML
	newConfig, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	// Create config directory
	if err = p.stack.CreateDir("prometheus"); err != nil {
		return err
	}

	// Write the updated YAML data to datadir
	if err = p.stack.WriteFile(filepath.Join("prometheus", "prometheus.yml"), newConfig); err != nil {
		return err
	}

	// Copy rules
	if err = p.copyRules(filepath.Join("prometheus")); err != nil {
		return err
	}

	// Copy alertmanager config
	if err = p.copyAlertManagerConfig(filepath.Join("prometheus")); err != nil {
		return err
	}

	return nil
}

// SetContainerIP sets the container IP for the Prometheus service.
func (p *PrometheusService) SetContainerIP(ip net.IP) {
	p.containerIP = ip
}

func (p *PrometheusService) ContainerName() string {
	return monitoring.PrometheusContainerName
}

func (p *PrometheusService) Endpoint() string {
	return fmt.Sprintf("http://%s:%d", p.containerIP, p.port)
}

// reloadConfig reloads the Prometheus config by making a POST request to the /-/reload endpoint
func (p *PrometheusService) reloadConfig() error {
	// Adding exponential retry
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Minute

	err := backoff.Retry(func() (err error) {
		resp, err := http.Post(fmt.Sprintf("http://%s:%d/-/reload", p.containerIP, p.port), "", nil)
		if err != nil {
			// TODO: Use fields to log the error
			log.Debug("Retrying request...")
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// TODO: Use fields to log the error
			log.Debug("Retrying request...")
			return fmt.Errorf("%w: %s", ErrReloadFailed, resp.Status)
		}
		return nil
	}, b)

	return err
}

func (p *PrometheusService) Name() string {
	return monitoring.PrometheusServiceName
}

// copyRules copy rules to $DATA_DIR/rules
func (p *PrometheusService) copyRules(dst string) (err error) {
	return fs.WalkDir(rules, "rules", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rule, err := rules.Open(path)
			if err != nil {
				return err
			}
			defer func() {
				cerr := rule.Close()
				if err == nil {
					err = cerr
				}
			}()
			data, err := io.ReadAll(rule)
			if err != nil {
				return err
			}
			if err = p.stack.WriteFile(filepath.Join(dst, path), data); err != nil {
				return err
			}
		} else {
			if err = p.stack.CreateDir(filepath.Join(dst, path)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *PrometheusService) copyAlertManagerConfig(dst string) (err error) {
	return fs.WalkDir(alertmanager, "alertmanager", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			alertmanagerConfig, err := alertmanager.Open(path)
			if err != nil {
				return err
			}
			defer func() {
				cerr := alertmanagerConfig.Close()
				if err == nil {
					err = cerr
				}
			}()
			data, err := io.ReadAll(alertmanagerConfig)
			if err != nil {
				return err
			}
			if err = p.stack.WriteFile(filepath.Join(dst, path), data); err != nil {
				return err
			}
		} else {
			if err = p.stack.CreateDir(filepath.Join(dst, path)); err != nil {
				return err
			}
		}
		return nil
	})
}
