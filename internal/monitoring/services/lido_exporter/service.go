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
package lido_exporter

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/NethermindEth/sedge/internal/monitoring"
	"github.com/NethermindEth/sedge/internal/monitoring/services/types"
)

var _ monitoring.ServiceAPI = &LidoExporterService{}

type LidoExporterParams struct {
	NodeOperatorID string
	RewardAddress  string
	Network        string
	RPCEndpoints   []string
	WSEndpoints    []string
	Port           uint16
	ScrapeTime     time.Duration
	LogLevel       string
}

type LidoExporterService struct {
	containerIP net.IP
	params      LidoExporterParams
}

func NewLidoExporter(params LidoExporterParams) *LidoExporterService {
	// Set other Lido Exporter parameters
	dotEnv["LIDO_EXPORTER_NODE_OPERATOR_ID"] = params.NodeOperatorID
	dotEnv["LIDO_EXPORTER_REWARD_ADDRESS"] = params.RewardAddress
	dotEnv["LIDO_EXPORTER_NETWORK"] = params.Network
	dotEnv["LIDO_EXPORTER_RPC_ENDPOINTS"] = strings.Join(params.RPCEndpoints, ",")
	dotEnv["LIDO_EXPORTER_WS_ENDPOINTS"] = strings.Join(params.WSEndpoints, ",")
	dotEnv["LIDO_EXPORTER_SCRAPE_TIME"] = params.ScrapeTime.String()
	dotEnv["LIDO_EXPORTER_LOG_LEVEL"] = params.LogLevel
	dotEnv["LIDO_EXPORTER_PORT"] = strconv.Itoa(int(params.Port))

	return &LidoExporterService{
		params: params,
	}
}

func (n *LidoExporterService) Init(opts types.ServiceOptions) error {
	// Validate dotEnv
	lidoExporterPort, ok := opts.Dotenv["LIDO_EXPORTER_PORT"]
	if !ok {
		return fmt.Errorf("%w: %s missing in options", ErrInvalidOptions, "LIDO_EXPORTER_PORT")
	} else if lidoExporterPort == "" {
		return fmt.Errorf("%w: %s can't be empty", ErrInvalidOptions, "LIDO_EXPORTER_PORT")
	}

	port, err := strconv.ParseUint(opts.Dotenv["LIDO_EXPORTER_PORT"], 10, 16)
	if err != nil {
		return fmt.Errorf("%w: %s is not a valid port", ErrInvalidOptions, "LIDO_EXPORTER_PORT")
	}
	n.params.Port = uint16(port)
	return nil
}

func (l *LidoExporterService) AddTarget(target types.MonitoringTarget, labels map[string]string, jobName string) error {
	return nil
}

func (l *LidoExporterService) RemoveTarget(instanceID string) (string, error) {
	return "", nil
}

func (l *LidoExporterService) DotEnv() map[string]string {
	return dotEnv
}

func (l *LidoExporterService) Setup(options map[string]string) error {
	return nil
}

func (l *LidoExporterService) SetContainerIP(ip net.IP) {
	l.containerIP = ip
}

func (l *LidoExporterService) ContainerName() string {
	return monitoring.LidoExporterContainerName
}

func (l *LidoExporterService) Endpoint() string {
	return fmt.Sprintf("http://%s:%d", l.containerIP, l.params.Port)
}

func (l *LidoExporterService) Name() string {
	return monitoring.LidoExporterServiceName
}
