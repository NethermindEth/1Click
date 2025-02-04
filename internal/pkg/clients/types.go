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
package clients

import (
	"github.com/NethermindEth/sedge/configs"
)

// Client : Struct Represent a client like geth, prysm, etc
type Client struct {
	Name      string
	Type      string
	Image     string
	Modified  bool
	Endpoint  string
	Supported bool
}

func (c *Client) SetImageOrDefault(image string) {
	switch c.Type {
	case "validator":
		c.setValidatorImage(image)
	case "consensus":
		c.setConsensusImage(image)
	case "execution":
		c.setExecutionImage(image)
	case "distributedValidator":
		c.setDistributedValidatorImage(image)
	case "optimism":
		c.setOptimismImage(image)
	case "opexecution":
		c.SetOpExecutionImage(image)
	}
}

func (c *Client) setExecutionImage(image string) {
	switch c.Name {
	case "geth":
		c.Image = valueOrDefault(image, configs.ClientImages.Execution.Geth.String())
	case "besu":
		c.Image = valueOrDefault(image, configs.ClientImages.Execution.Besu.String())
	case "nethermind":
		c.Image = valueOrDefault(image, configs.ClientImages.Execution.Nethermind.String())
	case "erigon":
		c.Image = valueOrDefault(image, configs.ClientImages.Execution.Erigon.String())
	}
}

func (c *Client) setConsensusImage(image string) {
	switch c.Name {
	case "lighthouse":
		c.Image = valueOrDefault(image, configs.ClientImages.Consensus.Lighthouse.String())
	case "prysm":
		c.Image = valueOrDefault(image, configs.ClientImages.Consensus.Prysm.String())
	case "teku":
		c.Image = valueOrDefault(image, configs.ClientImages.Consensus.Teku.String())
	case "lodestar":
		c.Image = valueOrDefault(image, configs.ClientImages.Consensus.Lodestar.String())
	case "grandine":
		c.Image = valueOrDefault(image, configs.ClientImages.Consensus.Grandine.String())
	case "nimbus":
		c.Image = valueOrDefault(image, configs.ClientImages.Consensus.Nimbus.String())
	}
}

func (c *Client) setValidatorImage(image string) {
	switch c.Name {
	case "lighthouse":
		c.Image = valueOrDefault(image, configs.ClientImages.Validator.Lighthouse.String())
	case "prysm":
		c.Image = valueOrDefault(image, configs.ClientImages.Validator.Prysm.String())
	case "teku":
		c.Image = valueOrDefault(image, configs.ClientImages.Validator.Teku.String())
	case "lodestar":
		c.Image = valueOrDefault(image, configs.ClientImages.Validator.Lodestar.String())
	case "nimbus":
		c.Image = valueOrDefault(image, configs.ClientImages.Validator.Nimbus.String())
	}
}

func (c *Client) setDistributedValidatorImage(image string) {
	switch c.Name {
	case "charon":
		c.Image = valueOrDefault(image, configs.ClientImages.Distributed.Charon.String())
	default:
		c.Image = valueOrDefault(image, configs.ClientImages.Distributed.Charon.String())
	}
}

func (c *Client) setOptimismImage(image string) {
	switch c.Name {
	case "opnode":
		c.Image = valueOrDefault(image, configs.ClientImages.Optimism.OpNode.String())
	}
}

func (c *Client) SetOpExecutionImage(image string) {
	switch c.Name {
	case "opgeth":
		c.Image = valueOrDefault(image, configs.ClientImages.OpExecution.OpGeth.String())
	case "opnethermind":
		c.Image = valueOrDefault(image, configs.ClientImages.OpExecution.OpNeth.String())
	case "opreth":
		c.Image = valueOrDefault(image, configs.ClientImages.OpExecution.OpReth.String())
	}
}

func valueOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// Clients : Struct Represent a combination of execution, consensus, validator and distributed validator clients
type Clients struct {
	Execution            *Client
	Consensus            *Client
	Validator            *Client
	Optimism             *Client
	ExecutionOP          *Client
	DistributedValidator *Client
}

type ClientMap map[string]*Client

type OrderedClients map[string]ClientMap
