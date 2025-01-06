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
package main

import (
	"runtime"

	"github.com/NethermindEth/sedge/cli"
	"github.com/NethermindEth/sedge/cli/actions"
	"github.com/NethermindEth/sedge/internal/compose"
	"github.com/NethermindEth/sedge/internal/locker"
	"github.com/NethermindEth/sedge/internal/monitoring"
	"github.com/NethermindEth/sedge/internal/monitoring/services/grafana"
	"github.com/NethermindEth/sedge/internal/monitoring/services/node_exporter"
	"github.com/NethermindEth/sedge/internal/monitoring/services/prometheus"
	"github.com/NethermindEth/sedge/internal/pkg/commands"
	"github.com/NethermindEth/sedge/internal/pkg/dependencies"
	"github.com/NethermindEth/sedge/internal/pkg/services"
	"github.com/NethermindEth/sedge/internal/ui"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func main() {
	// Commands Runner
	cmdRunner := commands.NewCMDRunner(commands.CMDRunnerOptions{
		RunAsAdmin: runtime.GOOS == "linux",
	})

	// Prompt used to interact with the user input
	prompt := ui.NewPrompter()

	// Docker client
	dockerClient, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer dockerClient.Close()

	// Docker service
	dockerServiceManager := services.NewDockerServiceManager(dockerClient)

	// Init dependencies manager
	depsMgr := dependencies.NewDependenciesManager(cmdRunner)
	// Init compose manager
	composeManager := compose.NewComposeManager(cmdRunner)

	// Set filesystem
	// fs := afero.NewMemMapFs() // Uncomment this line if you want to use the in-memory filesystem
	// fs := afero.NewBasePathFs(afero.NewOsFs(), "/tmp") // Uncomment this line if you want to use the real filesystem with a base path
	fs := afero.NewOsFs() // Uncomment this line if you want to use the real filesystem

	// Set locker
	locker := locker.NewFLock()

	// Get the monitoring manager
	monitoringServices := []monitoring.ServiceAPI{
		grafana.NewGrafana(),
		prometheus.NewPrometheus(),
		node_exporter.NewNodeExporter(),
	}
	monitoringMgr := monitoring.NewMonitoringManager(
		monitoringServices,
		composeManager,
		dockerServiceManager,
		fs,
		locker,
	)
	// Init Sedge Actions
	sdgOpts := actions.SedgeActionsOptions{
		DockerClient:         dockerClient,
		DockerServiceManager: dockerServiceManager,
		CommandRunner:        cmdRunner,
		ComposeManager:       *composeManager,
	}
	sedgeActions := actions.NewSedgeActions(sdgOpts)

	sedgeCmd := cli.RootCmd()
	sedgeCmd.AddCommand(
		cli.CliCmd(prompt, sedgeActions, depsMgr, monitoringMgr),
		cli.KeysCmd(cmdRunner, prompt),
		cli.DownCmd(*composeManager, cmdRunner, sedgeActions, depsMgr),
		cli.ClientsCmd(),
		cli.NetworksCmd(),
		cli.LogsCmd(*composeManager, cmdRunner, sedgeActions, depsMgr),
		cli.VersionCmd(),
		cli.SlashingExportCmd(sedgeActions, depsMgr),
		cli.SlashingImportCmd(sedgeActions, depsMgr),
		cli.RunCmd(sedgeActions, depsMgr),
		cli.ImportKeysCmd(sedgeActions, depsMgr),
		cli.GenerateCmd(sedgeActions),
		cli.DependenciesCommand(depsMgr),
		cli.ShowCmd(cmdRunner, sedgeActions, depsMgr),
		cli.LidoStatusCmd(),
		cli.MonitoringCmd(monitoringMgr),
	)
	sedgeCmd.SilenceErrors = true
	sedgeCmd.SilenceUsage = true
	if err := sedgeCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
