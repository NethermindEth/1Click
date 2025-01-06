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
package monitoring

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/NethermindEth/sedge/internal/common"
	"github.com/NethermindEth/sedge/internal/monitoring/data"
	"github.com/NethermindEth/sedge/internal/monitoring/locker"
	"github.com/NethermindEth/sedge/internal/monitoring/services/templates"
	"github.com/NethermindEth/sedge/internal/monitoring/services/types"
	"github.com/NethermindEth/sedge/internal/monitoring/utils"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	funk "github.com/thoas/go-funk"
)

// MonitoringManager manages the monitoring services. It provides methods for initializing the monitoring stack,
// adding and removing targets, running and stopping the monitoring stack, and checking the status of the monitoring stack.
type MonitoringManager struct {
	services             []ServiceAPI
	composeManager       ComposeManager
	dockerServiceManager DockerServiceManager
	stack                *data.MonitoringStack
}

// NewMonitoringManager creates a new MonitoringManager with the given services, compose manager, docker manager, file system, and locker.
func NewMonitoringManager(
	services []ServiceAPI,
	cmpMgr ComposeManager,
	dockerMgr DockerServiceManager,
	fs afero.Fs,
	locker locker.Locker,
) *MonitoringManager {
	// Create stack
	datadir, err := data.NewDataDirDefault(fs, locker)
	if err != nil {
		log.Fatal(err)
	}
	stack, err := datadir.MonitoringStack()
	if err != nil {
		log.Fatal(err)
	}

	return &MonitoringManager{
		services:             services,
		composeManager:       cmpMgr,
		dockerServiceManager: dockerMgr,
		stack:                stack,
	}
}

// Init initializes the monitoring stack. Assumes that the stack is already installed.
func (m *MonitoringManager) Init() error {
	// Read installed .env
	rawDotEnv, err := m.stack.ReadFile(filepath.Join(".env"))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitializingMonitoringMngr, err)
	}

	dotEnv := make(map[string]string)
	for _, line := range bytes.Split(rawDotEnv, []byte("\n")) {
		split := bytes.Split(line, []byte("="))
		if len(split) != 2 {
			continue
		}
		dotEnv[string(split[0])] = string(split[1])
	}

	// Initialize stack
	for _, service := range m.services {
		if err := service.Init(types.ServiceOptions{
			Stack:  m.stack,
			Dotenv: dotEnv,
		}); err != nil {
			return fmt.Errorf("%w: %w", ErrInitializingMonitoringMngr, err)
		}
	}

	// Save container IPs of monitoring services
	if err := m.saveServiceIP(); err != nil {
		return fmt.Errorf("%w: %w", ErrInitializingMonitoringMngr, err)
	}

	return nil
}

// InitStack initializes the monitoring stack by merging all environment variables, checking ports, setting up the stack and services, and creating containers.
func (m *MonitoringManager) InstallStack() error {
	// Merge all dotEnv
	dotEnv := make(map[string]string)
	defaultPorts := make(map[string]uint16)

	for _, service := range m.services {
		for k, v := range service.DotEnv() {
			dotEnv[k] = v
			// Grab default ports
			if strings.HasSuffix(k, "_PORT") {
				// Cast string to uint16
				p, err := strconv.ParseUint(v, 10, 16)
				if err != nil {
					return fmt.Errorf("%w: %w", ErrInstallingMonitoringMngr, err)
				}
				defaultPorts[k] = uint16(p)
			}
		}
	}

	// Check ports
	ports, err := utils.AssignPorts("localhost", defaultPorts)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInstallingMonitoringMngr, err)
	}
	for k, v := range ports {
		dotEnv[k] = strconv.Itoa(int(v))
	}

	// Intialize stack
	for _, service := range m.services {
		if err := service.Init(types.ServiceOptions{
			Stack:  m.stack,
			Dotenv: dotEnv,
		}); err != nil {
			return fmt.Errorf("%w: %w", ErrInstallingMonitoringMngr, err)
		}
	}

	if err = m.stack.Setup(dotEnv, templates.Services); err != nil {
		return fmt.Errorf("%w: %w", ErrInstallingMonitoringMngr, err)
	}

	// Setup services
	log.Debug("Setting up monitoring stack...")
	for _, service := range m.services {
		if err = service.Setup(dotEnv); err != nil {
			return fmt.Errorf("%w: %w", ErrInstallingMonitoringMngr, err)
		}
	}

	// Create containers
	if err = m.composeManager.Create(commands.DockerComposeCreateOptions{Path: filepath.Join(m.stack.Path(), "docker-compose.yml")}); err != nil {
		return fmt.Errorf("%w: %w", ErrInstallingMonitoringMngr, err)
	}

	log.Debug("Starting monitoring stack...")
	if err := m.composeManager.Up(commands.DockerComposeUpOptions{Path: filepath.Join(m.stack.Path(), "docker-compose.yml")}); err != nil {
		return fmt.Errorf("%w: %w", ErrRunningMonitoringStack, err)
	}

	// Save container IPs of monitoring services
	if err := m.saveServiceIP(); err != nil {
		return fmt.Errorf("%w: %w", ErrInitializingMonitoringMngr, err)
	}

	return nil
}

// AddTarget adds a new target to all services in the monitoring stack.
// It also connects the target to the docker network of the monitoring stack if it isn't already connected.
// The labels are added to the service's metrics.
func (m *MonitoringManager) AddTarget(target types.MonitoringTarget, labels map[string]string, dockerNetwork string) error {
	for _, service := range m.services {
		// Check if network was already added to service
		containerName := service.ContainerName()
		if containerName == PrometheusContainerName {
			networks, err := m.dockerServiceManager.ContainerNetworks(containerName)
			if err != nil {
				return err
			}
			if !funk.Contains(networks, dockerNetwork) {
				if err := m.dockerServiceManager.NetworkConnect(containerName, dockerNetwork); err != nil {
					return err
				}
			}
		}
		if err := service.AddTarget(target, labels, labels[InstanceIDLabel]+"--"+containerName+"++"+dockerNetwork); err != nil {
			return err
		}
	}
	return nil
}

// RemoveTarget removes a target from all services in the monitoring stack.
// It also disconnects the target from the docker network of the monitoring stack if it isn't already disconnected.
func (m *MonitoringManager) RemoveTarget(instanceID string) error {
	for _, service := range m.services {
		network, err := service.RemoveTarget(instanceID)
		if err != nil {
			return err
		}
		// Disconnect may fail if the network was already disconnected or if the container was already removed
		// so we ignore the error
		serviceName := service.ContainerName()
		if err := m.dockerServiceManager.NetworkDisconnect(serviceName, network); err != nil {
			log.Debugf("Error disconnecting %s from %s: %s", serviceName, network, err)
		}
	}
	return nil
}

// Run starts the monitoring stack by shutting down any existing stack and starting a new one.
func (m *MonitoringManager) Run() error {
	log.Info("Starting monitoring stack...")
	if err := m.composeManager.Up(commands.DockerComposeUpOptions{Path: filepath.Join(m.stack.Path(), "docker-compose.yml")}); err != nil {
		return fmt.Errorf("%w: %w", ErrRunningMonitoringStack, err)
	}

	// Save container IPs of monitoring services
	if err := m.saveServiceIP(); err != nil {
		return fmt.Errorf("%w: %w", ErrRunningMonitoringStack, err)
	}

	return nil
}

// Stop shuts down the monitoring stack.
func (m *MonitoringManager) Stop() error {
	log.Info("Shutting down monitoring stack...")
	if err := m.composeManager.Down(commands.DockerComposeDownOptions{Path: filepath.Join(m.stack.Path(), "docker-compose.yml")}); err != nil {
		return fmt.Errorf("%w: %w", ErrRunningMonitoringStack, err)
	}

	return nil
}

// Status checks the status of the containers in the monitoring stack and returns the status.
func (m *MonitoringManager) Status() (status common.Status, err error) {
	var containers []string
	for _, service := range m.services {
		containers = append(containers, service.ContainerName())
	}

	for _, container := range containers {
		status, err = m.dockerServiceManager.ContainerStatus(container)
		if err != nil {
			return common.Unknown, fmt.Errorf("%w: %w", ErrCheckingMonitoringStack, err)
		}
		// running or restarting means the stack is running
		if status != common.Running && status != common.Restarting {
			return common.Broken,
				fmt.Errorf("%w: %s container is either paused, exited, or dead", ErrCheckingMonitoringStack, container)
		}
	}

	return status, nil
}

// InstallationStatus checks whether the monitoring stack is installed and returns the installation status.
func (m *MonitoringManager) InstallationStatus() (common.Status, error) {
	installed, err := m.stack.Installed()
	if err != nil {
		return common.Unknown, fmt.Errorf("%w: %w", ErrCheckingMonitoringStack, err)
	}
	if installed {
		return common.Installed, nil
	}

	return common.NotInstalled, nil
}

// Cleanup removes the monitoring stack.
func (m *MonitoringManager) Cleanup() error {
	log.Info("Shutting down monitoring stack...")
	if err := m.composeManager.Down(commands.DockerComposeDownOptions{Path: filepath.Join(m.stack.Path(), "docker-compose.yml")}); err != nil {
		return fmt.Errorf("%w: %w", ErrRunningMonitoringStack, err)
	}

	log.Info("Cleaning up monitoring stack...")
	if err := m.stack.Cleanup(true); err != nil {
		return fmt.Errorf("%w: %w", ErrRunningMonitoringStack, err)
	}

	return nil
}

// ServiceEndpoints returns a map of the service's container names and their endpoints.
func (m *MonitoringManager) ServiceEndpoints() map[string]string {
	endpoints := make(map[string]string)
	for _, service := range m.services {
		endpoints[service.ContainerName()] = service.Endpoint()
	}
	return endpoints
}

func (m *MonitoringManager) idToIP(id string) (string, error) {
	ip, err := m.dockerServiceManager.ContainerIP(id)
	if err != nil {
		return "", err
	}
	return ip, nil
}

func (m *MonitoringManager) saveServiceIP() error {
	for _, service := range m.services {
		name := service.ContainerName()
		ip, err := m.idToIP(name)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInitializingMonitoringMngr, err)
		}
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			return fmt.Errorf("%w: failed to save the IP address of the monitoring service %s: %s is not a valid IP address", ErrInitializingMonitoringMngr, name, ip)
		}
		service.SetContainerIP(parsedIP)
	}
	return nil
}

// AddService adds a new service to the monitoring stack dynamically.
func (m *MonitoringManager) AddService(service ServiceAPI) error {
	err := m.validateNewService(service)
	if err != nil {
		return fmt.Errorf("error validating service %s: %w", service.Name(), err)
	}

	// Add the new service to the list
	m.services = append(m.services, service)

	// Get the new service's environment variables
	dotEnv := service.DotEnv()

	err = m.setupNewService(service, dotEnv)
	if err != nil {
		return fmt.Errorf("failed to update .env file: %w", err)
	}

	// Update the .env file in the stack
	if err := m.updateEnvFile(dotEnv); err != nil {
		return fmt.Errorf("failed to update .env file: %w", err)
	}

	// Update the docker-compose.yml file
	if err := m.updateDockerComposeFile(service, templates.Services); err != nil {
		return fmt.Errorf("failed to update docker-compose.yml: %w", err)
	}

	// Create and start the new service's container
	if err := m.composeManager.Create(commands.DockerComposeCreateOptions{Path: filepath.Join(m.stack.Path(), "docker-compose.yml")}); err != nil {
		return fmt.Errorf("failed to create service container: %w", err)
	}
	if err := m.composeManager.Up(commands.DockerComposeUpOptions{Path: filepath.Join(m.stack.Path(), "docker-compose.yml")}); err != nil {
		return fmt.Errorf("failed to start service container: %w", err)
	}

	// Save the new service's IP
	if err := m.saveServiceIP(); err != nil {
		return fmt.Errorf("failed to save service IP: %w", err)
	}

	monitoringTarget, labels, err := m.makeTarget(service)
	if err != nil {
		return fmt.Errorf("error making target of service %s: %w", service.Name(), err)
	}
	// Add this new service as a target to the monitoring manager
	if err := m.AddTarget(monitoringTarget, labels, SedgeNetworkName); err != nil {
		return fmt.Errorf("failed to add target for service %s: %w", service.Name(), err)
	}

	return nil
}

func (m *MonitoringManager) validateNewService(service ServiceAPI) error {
	// Check if the service already exists
	for _, existingService := range m.services {
		if existingService.ContainerName() == service.ContainerName() {
			return fmt.Errorf("service %s already exists", service.ContainerName())
		}
	}
	return nil
}

func (m *MonitoringManager) setupNewService(service ServiceAPI, dotEnv map[string]string) error {
	// Initialize the new service
	if err := service.Init(types.ServiceOptions{
		Stack:  m.stack,
		Dotenv: dotEnv,
	}); err != nil {
		return fmt.Errorf("failed to initialize service %s: %w", service.Name(), err)
	}

	// Setup the new service
	if err := service.Setup(dotEnv); err != nil {
		return fmt.Errorf("failed to setup service %s: %w", service.Name(), err)
	}
	return nil
}

func (m *MonitoringManager) makeTarget(service ServiceAPI) (types.MonitoringTarget, map[string]string, error) {
	var monitoringTarget types.MonitoringTarget
	// Split the service's Endpoint into host and port
	endpoint := service.Endpoint()
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	_, portStr, err := net.SplitHostPort(endpoint)
	if err != nil {
		return monitoringTarget, nil, fmt.Errorf("invalid service endpoint %s: %w", endpoint, err)
	}

	// Convert port to uint16
	port64, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return monitoringTarget, nil, fmt.Errorf("invalid port in service endpoint %s: %w", portStr, err)
	}
	// Set service as target
	monitoringTarget = types.MonitoringTarget{
		Host: service.ContainerName(),
		Port: uint16(port64),
		Path: "/metrics",
	}

	labels := map[string]string{
		InstanceIDLabel: service.ContainerName(),
	}
	return monitoringTarget, labels, nil
}

// Helper method to update the .env file
func (m *MonitoringManager) updateEnvFile(newEnv map[string]string) error {
	currentEnv, err := m.stack.ReadFile(".env")
	if err != nil {
		return err
	}

	// Parse current env
	env := make(map[string]string)
	for _, line := range bytes.Split(currentEnv, []byte("\n")) {
		parts := bytes.SplitN(line, []byte("="), 2)
		if len(parts) == 2 {
			env[string(parts[0])] = string(parts[1])
		}
	}

	// Merge new env
	for k, v := range newEnv {
		env[k] = v
	}

	// Write updated env
	var buf bytes.Buffer
	for k, v := range env {
		buf.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return m.stack.WriteFile(".env", buf.Bytes())
}

// Helper method to update the docker-compose.yml file
func (m *MonitoringManager) updateDockerComposeFile(service ServiceAPI, monitoringFs fs.FS) error {
	// Read the main Docker Compose template
	rawBaseTmp, err := monitoringFs.Open("services/docker-compose_base.tmpl")
	if err != nil {
		return fmt.Errorf("error opening docker-compose template: %w", err)
	}
	defer rawBaseTmp.Close()

	// Read the content of the base template file
	rawBaseTmpContent, err := io.ReadAll(rawBaseTmp)
	if err != nil {
		return fmt.Errorf("error reading docker-compose template: %w", err)
	}

	baseTmp, err := template.New("docker-compose").Parse(string(rawBaseTmpContent))
	if err != nil {
		return err
	}

	// Fetch the service template file
	serviceTmp, err := monitoringFs.Open("services/" + service.Name() + ".tmpl")
	if err != nil {
		return fmt.Errorf("error opening service %s template: %w", service.Name(), err)
	}
	defer serviceTmp.Close()

	serviceTmpContent, err := io.ReadAll(serviceTmp)
	if err != nil {
		return fmt.Errorf("error reading service %s template: %w", service.Name(), err)
	}

	// Merge both templates
	baseTmp, err = baseTmp.Parse(string(serviceTmpContent))
	if err != nil {
		return fmt.Errorf("error parsing service %s template: %w", service.Name(), err)
	}

	// Create a buffer to hold the merged content
	var buf bytes.Buffer

	data := types.ServiceTemplateData{
		LidoExporter: service.Name() == LidoExporterServiceName,
	}

	// Execute the main template with the service template as data
	if err := baseTmp.Execute(&buf, data); err != nil {
		return err
	}

	// Write the merged content to the final Docker Compose file
	return m.stack.WriteFile("docker-compose.yml", buf.Bytes())
}
