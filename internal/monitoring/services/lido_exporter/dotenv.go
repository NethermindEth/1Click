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

var dotEnv map[string]string = map[string]string{
	"LIDO_EXPORTER_IMAGE":            "nethermindeth/lido-exporter:v1.0.2",
	"LIDO_EXPORTER_PORT":             "",
	"LIDO_EXPORTER_NODE_OPERATOR_ID": "",
	"LIDO_EXPORTER_REWARD_ADDRESS":   "",
	"LIDO_EXPORTER_NETWORK":          "",
	"LIDO_EXPORTER_RPC_ENDPOINTS":    "",
	"LIDO_EXPORTER_WS_ENDPOINTS":     "",
	"LIDO_EXPORTER_SCRAPE_TIME":      "",
	"LIDO_EXPORTER_LOG_LEVEL":        "",
}
