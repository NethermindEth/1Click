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

var dotEnv map[string]string = map[string]string{
	"PROM_IMAGE":         "prom/prometheus:v2.54.1",
	"PROM_PORT":          "9090",
	"PROM_CONF":          "./prometheus/prometheus.yml",
	"PROM_RULES":         "./prometheus/rules",
	"ALERTMANAGER_IMAGE": "prom/alertmanager:v0.27.0",
	"ALERTMANAGER_CONF":  "./prometheus/alertmanager/alertmanager.yml",
	"ALERTMANAGER_PORT":  "9093",
}
