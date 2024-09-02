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

const (
	PrometheusServiceName     = "prometheus"
	PrometheusContainerName   = "sedge_prometheus"
	GrafanaServiceName        = "grafana"
	GrafanaContainerName      = "sedge_grafana"
	NodeExporterServiceName   = "node_exporter"
	NodeExporterContainerName = "sedge_node_exporter"
	monitoringPath            = "monitoring"
	InstanceIDLabel           = "instance_id"
	CommitHashLabel           = "instance_commit_hash"
	AVSNameLabel              = "avs_name"
	AVSVersionLabel           = "avs_version"
	SpecVersionLabel          = "spec_version"
)
