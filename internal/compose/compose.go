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
package compose

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	log "github.com/sirupsen/logrus"
)

// DockerComposeCmdError represents an error that occurs when running a Docker Compose command.
type DockerComposeCmdError struct {
	Cmd string
}

// Error returns a string representation of the DockerComposeCmdError.
func (e DockerComposeCmdError) Error() string {
	return fmt.Sprintf("Docker Compose Manager running 'docker compose %s'", e.Cmd)
}

// ComposeManager manages Docker Compose operations.
type ComposeManager struct {
	cmdRunner commands.CommandRunner
}

// NewComposeManager creates a new instance of ComposeManager.
func NewComposeManager(runner commands.CommandRunner) *ComposeManager {
	return &ComposeManager{
		cmdRunner: runner,
	}
}

// Up runs the Docker Compose 'up' command for the specified options.
func (cm *ComposeManager) Up(opts commands.DockerComposeUpOptions) error {
	upCmd := cm.cmdRunner.BuildDockerComposeUpCMD(opts)

	log.Infof(configs.RunningCommand, upCmd.Cmd)
	if out, exitCode, err := cm.cmdRunner.RunCMD(upCmd); err != nil || exitCode != 0 {
		return fmt.Errorf("%w: %s. Output: %s", DockerComposeCmdError{Cmd: "up"}, err, out)
	}
	return nil
}

// Pull runs the Docker Compose 'pull' command for the specified options.
func (cm *ComposeManager) Pull(opts commands.DockerComposePullOptions) error {
	pullCmd := cm.cmdRunner.BuildDockerComposePullCMD(opts)

	log.Infof(configs.RunningCommand, pullCmd.Cmd)
	if out, exitCode, err := cm.cmdRunner.RunCMD(pullCmd); err != nil || exitCode != 0 {
		return fmt.Errorf("%w: %s. Output: %s", DockerComposeCmdError{Cmd: "pull"}, err, out)
	}
	return nil
}

// Create runs the Docker Compose 'create' command for the specified options.
func (cm *ComposeManager) Create(opts commands.DockerComposeCreateOptions) error {
	createCmd := cm.cmdRunner.BuildDockerComposeCreateCMD(opts)

	log.Infof(configs.RunningCommand, createCmd.Cmd)
	if out, exitCode, err := cm.cmdRunner.RunCMD(createCmd); err != nil || exitCode != 0 {
		return fmt.Errorf("%w: %s. Output: %s", DockerComposeCmdError{Cmd: "create"}, err, out)
	}
	return nil
}

// Build runs the Docker Compose 'build' command for the specified options.
func (cm *ComposeManager) Build(opts commands.DockerComposeBuildOptions) error {
	buildCmd := cm.cmdRunner.BuildDockerComposeBuildCMD(opts)

	log.Infof(configs.RunningCommand, buildCmd.Cmd)
	if out, exitCode, err := cm.cmdRunner.RunCMD(buildCmd); err != nil || exitCode != 0 {
		return fmt.Errorf("%w: %s. Output: %s", DockerComposeCmdError{Cmd: "build"}, err, out)
	}
	return nil
}

// PS runs the Docker Compose 'ps' command for the specified options and returns
// the list of services.
func (c *ComposeManager) PS(opts commands.DockerComposePsOptions) ([]ComposeService, error) {
	psCmd := c.cmdRunner.BuildDockerComposePSCMD(opts)

	log.Infof(configs.RunningCommand, psCmd.Cmd)
	out, exitCode, err := c.cmdRunner.RunCMD(psCmd)
	if err != nil || exitCode != 0 {
		return nil, fmt.Errorf("%w: %s. Output: %s", DockerComposeCmdError{Cmd: "ps"}, err, out)
	}
	outList := make([]ComposeService, 0)
	if len(out) == 0 {
		return outList, nil
	}
	// Following `if` cases are necessary to handle the different output formats
	// of the `docker compose ps` command. Some times it returns a list of
	// services, other times it returns a single service for the edge case of
	// only one service being present. Depending on the docker compose version.
	if out[0] == '[' {
		// Multiple services
		err = json.Unmarshal([]byte(out), &outList)
		if err != nil {
			return outList, fmt.Errorf("%w: %s. Output: %s", DockerComposeCmdError{Cmd: "ps"}, err, out)
		}
	} else if out[0] == '{' {
		// Single service
		var s ComposeService
		err = json.Unmarshal([]byte(out), &s)
		if err != nil {
			return outList, fmt.Errorf("%w: %s. Output: %s", DockerComposeCmdError{Cmd: "ps"}, err, out)
		}
		outList = append(outList, s)
	} else if strings.HasPrefix(out, "null") {
		// No services
		return outList, nil
	} else {
		// Unexpected output
		return outList, fmt.Errorf("%w: %s. Output: %s", DockerComposeCmdError{Cmd: "ps"}, "unknown output format", out)
	}
	return outList, nil
}

// Logs runs the Docker Compose 'logs' command for the specified options.
func (cm *ComposeManager) Logs(opts commands.DockerComposeLogsOptions) error {
	logsCmd := cm.cmdRunner.BuildDockerComposeLogsCMD(opts)

	log.Infof(configs.RunningCommand, logsCmd.Cmd)
	if out, exitCode, err := cm.cmdRunner.RunCMD(logsCmd); err != nil || exitCode != 0 {
		if exitCode == 130 {
			return nil
		}
		return fmt.Errorf("%w: %s. Output: %s", DockerComposeCmdError{Cmd: "logs"}, err, out)
	}
	return nil
}

// Stop runs the Docker Compose 'stop' command for the specified options.
func (cm *ComposeManager) Stop(opts DockerComposeStopOptions) error {
	stopCmd := fmt.Sprintf("docker compose -f %s stop", opts.Path)

	log.Infof(configs.RunningCommand, stopCmd)
	if out, exitCode, err := cm.cmdRunner.RunCMD(commands.Command{Cmd: stopCmd, GetOutput: true}); err != nil || exitCode != 0 {
		return fmt.Errorf("%w: %s. Output: %s", DockerComposeCmdError{Cmd: "stop"}, err, out)
	}
	return nil
}

// Down runs the Docker Compose 'down' command for the specified options.
func (cm *ComposeManager) Down(opts commands.DockerComposeDownOptions) error {
	downCmd := cm.cmdRunner.BuildDockerComposeDownCMD(opts)

	log.Infof(configs.RunningCommand, downCmd.Cmd)
	if out, exitCode, err := cm.cmdRunner.RunCMD(downCmd); err != nil || exitCode != 0 {
		return fmt.Errorf("%w: %s. Output: %s", DockerComposeCmdError{Cmd: "down"}, err, out)
	}
	return nil
}
