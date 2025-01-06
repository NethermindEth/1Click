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
package cli

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/NethermindEth/sedge/internal/common"
	"github.com/NethermindEth/sedge/internal/monitoring"
	lidoExporter "github.com/NethermindEth/sedge/internal/monitoring/services/lido_exporter"
	"github.com/NethermindEth/sedge/internal/utils"
)

func MonitoringCmd(mgr MonitoringManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitoring [init|clean]",
		Short: "Manage the monitoring stack",
		Long:  "Manage the monitoring stack. Use 'init' to install and run, or 'clean' to stop and uninstall.",
	}
	cmd.AddCommand(InitSubCmd(mgr))
	cmd.AddCommand(CleanSubCmd(mgr))

	return cmd
}

func InitSubCmd(mgr MonitoringManager) *cobra.Command {
	var additionalServices []monitoring.ServiceAPI
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the monitoring stack",
		Long: `This command initializes the monitoring stack (Grafana, Prometheus, etc.) for Lido CSM or general node monitoring.

The monitoring stack includes:
- Grafana dashboards for real-time monitoring of Lido CSM node metrics.
- Prometheus for collecting and displaying key metrics about your node operations.`,
	}
	cmd.AddCommand(DefaultSubCmd(mgr, additionalServices))
	cmd.AddCommand(LidoSubCmd(mgr, additionalServices))

	return cmd
}

func CleanSubCmd(mgr MonitoringManager) *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Clean and uninstall the monitoring stack",
		RunE: func(cmd *cobra.Command, args []string) error {
			return CleanMonitoring(mgr)
		},
	}
}

func LidoSubCmd(mgr MonitoringManager, additionalServices []monitoring.ServiceAPI) *cobra.Command {
	lido := &lidoExporter.LidoExporterParams{}
	cmd := &cobra.Command{
		Use:   "lido",
		Short: "Configure Lido CSM Node monitoring",
		Long:  "Configure Lido CSM node monitoring (Prometheus, Grafana, Node Exporter,Lido Exporter)",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if lido.NodeOperatorID == "" && lido.RewardAddress == "" {
				return errors.New("Node Operator ID or Reward Address is required")
			}
			if lido.NodeOperatorID != "" {
				var nodeOperatorIDBigInt *big.Int
				var ok bool
				nodeOperatorIDBigInt, ok = new(big.Int).SetString(lido.NodeOperatorID, 10)
				if !ok {
					return errors.New("Failed to convert Node Operator ID to big.Int")
				}
				if nodeOperatorIDBigInt.Sign() < 0 {
					return errors.New("Node Operator ID cannot be negative")
				}
			} else {
				if !utils.IsAddress(rewardAddress) {
					return errors.New("Invalid reward address")
				}
			}
			additionalServices = append(additionalServices, lidoExporter.NewLidoExporter(*lido))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return InitMonitoring(true, true, mgr, additionalServices)
		},
	}
	cmd.Flags().StringVar(&lido.NodeOperatorID, "node-operator-id", "", "Node Operator ID")
	cmd.Flags().StringVar(&lido.RewardAddress, "reward-address", "", "Reward address of Node Operator. It is used to calculate Node Operator ID if not set")
	cmd.Flags().StringVar(&lido.Network, "network", "holesky", "Network name")
	cmd.Flags().StringSliceVar(&lido.RPCEndpoints, "rpc-endpoints", nil, "List of Ethereum HTTP RPC endpoints")
	cmd.Flags().StringSliceVar(&lido.WSEndpoints, "ws-endpoints", nil, "List of Ethereum WebSocket RPC endpoints")
	cmd.Flags().Uint16Var(&lido.Port, "port", 8080, "Port where the metrics will be exported.")
	cmd.Flags().DurationVar(&lido.ScrapeTime, "scrape-time", 30*time.Second, "Time interval for scraping metrics. Values should be in the format of 10s, 1m, 1h, etc.")
	cmd.Flags().StringVar(&logLevel, "log-level", "info", "Set Log Level, e.g panic, fatal, error, warn, warning, info, debug, trace")

	return cmd
}

func DefaultSubCmd(mgr MonitoringManager, additionalServices []monitoring.ServiceAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default",
		Short: "Default monitoring configuration",
		Long:  "Default monitoring configuration (Prometheus, Grafana, Node Exporter)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return InitMonitoring(true, true, mgr, nil)
		},
	}
	return cmd
}

// Init initializes the Monitoring Stack. If install is true, it will install the Monitoring Stack if it is not installed.
// If run is true, it will run the Monitoring Stack if it is not running.
func InitMonitoring(install, run bool, monitoringMgr MonitoringManager, additionalServices []monitoring.ServiceAPI) error {
	// Check if the monitoring stack is installed.
	installStatus, err := monitoringMgr.InstallationStatus()
	if err != nil {
		return err
	}
	log.Debugf("Monitoring stack installation status: %v", installStatus == common.Installed)
	// If the monitoring stack is not installed, install it.
	if installStatus == common.NotInstalled && install {
		err = monitoringMgr.InstallStack()
		if errors.Is(err, monitoring.ErrInstallingMonitoringMngr) {
			// If the monitoring stack installation fails, remove the monitoring stack directory.
			if cerr := monitoringMgr.Cleanup(); cerr != nil {
				return fmt.Errorf("install failed: %w. Failed to cleanup monitoring stack after installation failure: %w", err, cerr)
			}
			return err
		} else if err != nil {
			return err
		}
	}

	// Check if the monitoring stack is running.
	status, err := monitoringMgr.Status()
	if err != nil {
		log.Errorf("Monitoring stack status: unknown. Got error: %v", err)
	}
	// If the monitoring stack is not running, start it.
	if status != common.Running && status != common.Restarting && run {
		if err := monitoringMgr.Run(); err != nil {
			return err
		}
	} else if status != common.Running && status != common.Restarting && !run {
		// If the monitoring stack is not supposed to be running then exit.
		return nil
	}

	// Initialize monitoring stack if it is running.
	if err := monitoringMgr.Init(); err != nil {
		return err
	}

	// Add additional services to the monitoring manager
	for _, service := range additionalServices {
		if err := monitoringMgr.AddService(service); err != nil {
			return fmt.Errorf("failed to add service %s: %w", service.Name(), err)
		}
	}

	return nil
}

// CleanMonitoring stops and uninstalls the Monitoring Stack
func CleanMonitoring(monitoringMgr MonitoringManager) error {
	// Check if the monitoring stack is installed.
	installStatus, err := monitoringMgr.InstallationStatus()
	if err != nil {
		return err
	}
	log.Debugf("Monitoring stack installation status: %v", installStatus == common.Installed)
	// If the monitoring stack is installed, uninstall it.
	if installStatus == common.Installed {
		if err := monitoringMgr.Cleanup(); err != nil {
			return err
		}
	}
	return nil
}
