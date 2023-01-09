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
package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type WindowsCMDRunner struct {
	terminal string
}

func NewCMDRunner(options CMDRunnerOptions) CommandRunner {
	cr := new(WindowsCMDRunner)
	cr.terminal = "powershell"
	return cr
}

func (cr *WindowsCMDRunner) BuildDockerComposeUpCMD(options DockerComposeUpOptions) Command {
	command := fmt.Sprintf("docker compose -f %s up -d", options.Path)
	if len(options.Services) > 0 {
		command += " " + strings.Join(options.Services, " ")
	}
	return Command{Cmd: command}
}

func (cr *WindowsCMDRunner) BuildDockerComposePullCMD(options DockerComposePullOptions) Command {
	command := fmt.Sprintf("docker compose -f %s pull", options.Path)
	if len(options.Services) > 0 {
		command += " " + strings.Join(options.Services, " ")
	}
	return Command{Cmd: command}
}

func (cr *WindowsCMDRunner) BuildDockerComposeCreateCMD(options DockerComposeCreateOptions) Command {
	command := fmt.Sprintf("docker compose -f %s create", options.Path)
	if len(options.Services) > 0 {
		command += " " + strings.Join(options.Services, " ")
	}
	return Command{Cmd: command}
}

func (cr *WindowsCMDRunner) BuildDockerComposeBuildCMD(options DockerComposeBuildOptions) Command {
	command := fmt.Sprintf("docker compose -f %s build", options.Path)
	if len(options.Services) > 0 {
		command += " " + strings.Join(options.Services, " ")
	}
	return Command{Cmd: command}
}

func (cr *WindowsCMDRunner) BuildDockerPSCMD(options DockerPSOptions) Command {
	command := "docker ps"
	if options.All {
		command += " -a"
	} else {
		log.Debug(`Command "docker ps" built without the "--all" flag.`)
	}
	return Command{Cmd: command}
}

func (cr *WindowsCMDRunner) BuildDockerComposePSCMD(options DockerComposePsOptions) Command {
	flags := ""
	name := ""

	if options.Services {
		log.Debug(`Command "docker compose ps" built with "--service" flag.`)
		flags += " --services"
	} else if options.Quiet {
		log.Debug(`Command "docker compose ps" built with "--quiet" flag.`)
		flags += " --quiet"
	}

	if options.FilterRunning {
		flags += " --filter status=running"
	}

	if options.ServiceName != "" {
		name += " " + options.ServiceName
	}

	var command string
	if options.Path != "" {
		command = fmt.Sprintf("docker compose -f %s ps%s%s", options.Path, flags, name)
	} else {
		command = fmt.Sprintf("docker compose ps%s%s", flags, name)
	}
	return Command{Cmd: command}
}

func (cr *WindowsCMDRunner) BuildDockerComposeLogsCMD(options DockerComposeLogsOptions) Command {
	command := fmt.Sprintf("docker compose -f %s logs", options.Path)
	servs := strings.Join(options.Services, " ")
	if options.Follow {
		log.Debug(`Command "docker compose logs" built with "--follow" flag.`)
		command += fmt.Sprintf(" --follow %s", servs)
	} else if options.Tail > 0 {
		log.Debugf(`Command "docker compose logs" built with "--tail=%d" flag.`, options.Tail)
		command += fmt.Sprintf(" --tail=%d %s", options.Tail, servs)
	} else {
		log.Warn(`Command "docker compose logs" built without "--follow" or "--tail" flags. Add follow argument or make tail argument greater than 0.`)
		command += " " + servs
	}
	return Command{Cmd: command}
}

func (cr *WindowsCMDRunner) BuildDockerBuildCMD(options DockerBuildOptions) Command {
	command := fmt.Sprintf("docker build %s", options.Path)
	if len(options.Tag) > 0 {
		log.Debug(`Command "docker build" built with "-t" flag.`)
		command += " -t " + options.Tag
	} else {
		log.Debug(`Command "docker build" built withot "-t" flag.`)
	}
	return Command{Cmd: command}
}

func (cr *WindowsCMDRunner) BuildDockerPullCMD(options DockerBuildOptions) Command {
	command := fmt.Sprintf("docker pull %s", options.Tag)
	return Command{Cmd: command}
}

func (cr *WindowsCMDRunner) BuildDockerInspectCMD(options DockerInspectOptions) Command {
	flags := ""
	if options.Format != "" {
		flags += " --format " + options.Format
	}
	command := fmt.Sprintf("docker inspect%s %s", flags, options.Name)
	return Command{Cmd: command}
}

func (cr *WindowsCMDRunner) BuildDockerComposeDownCMD(options DockerComposeDownOptions) Command {
	command := fmt.Sprintf("docker compose -f %s down", options.Path)
	return Command{Cmd: command}
}

func (cr *WindowsCMDRunner) BuildCreateFileCMD(options CreateFileOptions) Command {
	return Command{Cmd: fmt.Sprintf("echo $null >> %s", options.FileName)}
}

func (cr *WindowsCMDRunner) BuildEchoToFileCMD(options EchoToFileOptions) Command {
	return Command{Cmd: fmt.Sprintf("echo %s > %s", options.Content, options.FileName)}
}

func (cr *WindowsCMDRunner) BuildOpenTextEditor(options OpenTextEditorOptions) Command {
	return Command{Cmd: fmt.Sprintf("notepad %s", options.FilePath), IgnoreTerminal: true}
}

func (cr *WindowsCMDRunner) RunCMD(cmd Command) (string, error) {
	var out string
	r := strings.ReplaceAll(cmd.Cmd, "\n", "")

	var exc *exec.Cmd
	if !cmd.IgnoreTerminal {
		exc = exec.Command(cr.terminal, r)
	} else {
		s := strings.Split(r, " ")
		exc = exec.Command(s[0], s[1:]...)
	}

	var combinedOut bytes.Buffer
	exc.Stdin = os.Stdin
	if cmd.GetOutput {
		exc.Stdout = &combinedOut
		exc.Stderr = &combinedOut
	} else {
		exc.Stdout = os.Stdout
		exc.Stderr = os.Stderr
	}

	if err := exc.Start(); err != nil {
		return out, err
	}

	err := exc.Wait()
	if cmd.GetOutput {
		out = combinedOut.String()
		out = strings.ReplaceAll(out, "\r", "") // Remove windows \r character
	}

	return out, err
}

func (cr *WindowsCMDRunner) RunScript(script ScriptFile) (string, error) {
	var out string
	var scriptBuffer, combinedOut bytes.Buffer
	scriptBuffer.WriteString("$ErrorActionPreference = \"Stop\"\n") // Force powershell to stop on errors
	if err := script.Tmp.Execute(&scriptBuffer, script.Data); err != nil {
		return out, err
	}

	tempFileDir := os.TempDir()
	rawScript, err := io.ReadAll(&scriptBuffer)
	if err != nil {
		return out, err
	}

	tempScript := filepath.Join(tempFileDir, "temp.ps1")
	if err := os.WriteFile(tempScript, rawScript, os.ModePerm); err != nil {
		return out, err
	}

	cmd := exec.Command(cr.terminal, tempScript)

	// Prepare pipes for stdin, stdout and stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return out, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return out, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return out, err
	}

	wait := sync.WaitGroup{}

	// Prepare channel to receive errors from goroutines
	errChans := make([]<-chan error, 0)
	// cmd executes any instructions coming from stdin
	errChans = append(errChans, goCopy(&wait, stdin, &scriptBuffer, true))

	if script.GetOutput {
		// If the script is to get the output, then use an unified buffer to combine stdout and stderr
		cmd.Stdout = &combinedOut
		cmd.Stderr = &combinedOut
	} else {
		// If the script is not to get the output, then pipe the output to stdout and stderr
		errChans = append(errChans, goCopy(&wait, os.Stdout, stdout, false))
		errChans = append(errChans, goCopy(&wait, os.Stderr, stderr, false))
	}

	if err := cmd.Start(); err != nil {
		return out, err
	}

	// Check for errors from goroutines
	for _, errChan := range errChans {
		err = <-errChan
		if err != nil {
			return out, err
		}
	}

	wait.Wait()

	if err := cmd.Wait(); err != nil {
		return out, err
	}

	if script.GetOutput {
		out = combinedOut.String()
	}

	out = strings.ReplaceAll(out, "\r", "")
	return out, nil
}
