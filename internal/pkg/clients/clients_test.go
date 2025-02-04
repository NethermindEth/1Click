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
	"fmt"
	"testing"

	"github.com/NethermindEth/sedge/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestSupportedClients(t *testing.T) {
	inputs := [...]struct {
		clientType string
		network    string
		want       []string
		isErr      bool
	}{
		{"execution", "gnosis", []string{"nethermind", "erigon"}, false},
		{"consensus", "gnosis", utils.Filter(AllClients["consensus"], func(c string) bool { return c != "prysm" }), false},
		{"execution", "mainnet", AllClients["execution"], false},
		{"consensus", "mainnet", AllClients["consensus"], false},
		{"validator", "mainnet", AllClients["validator"], false},
		{"distributedValidator", "holesky", AllClients["distributedValidator"], false},
		{"random", "mainnet", []string{}, true},
	}

	for _, input := range inputs {
		descr := fmt.Sprintf("SupportedClients(%s)", input.clientType)

		c := ClientInfo{Network: input.network}
		if res, err := c.SupportedClients(input.clientType); input.isErr && err == nil {
			t.Errorf("%s expected to fail", descr)
		} else if !input.isErr {
			if err != nil {
				t.Errorf("%s failed: %v", descr, err)
			} else {
				assert.EqualValues(t, input.want, res, descr)
			}
		}
	}
}

type clientsTestCase struct {
	configClientsTypes map[string][]string
	query              []string
	network            string
	isErr              bool
}

func validateClients(resultClients OrderedClients, tc clientsTestCase) bool {
	// Check if all query clients types are in the result types
Loop1:
	for _, queryType := range tc.query {
		for resultType := range resultClients {
			if queryType == resultType {
				continue Loop1
			}
		}
		return false
	}

	for resultType, resultTypeClients := range resultClients {
	Loop2:
		for _, resultClient := range resultTypeClients {
			for _, configClientName := range tc.configClientsTypes[resultType] {
				if resultClient.Name == configClientName {
					continue Loop2
				}
			}
			return false
		}
	}
	return true
}

func TestClients(t *testing.T) {
	inputs := [...]clientsTestCase{
		{
			map[string][]string{
				"consensus": {"lighthouse", "prysm", "teku", "lodestar", "nimbus", "grandine"},
				"validator": {"lighthouse", "prysm", "teku", "lodestar", "nimbus"},
				"execution": {"nethermind", "geth", "besu", "erigon"},
			},
			[]string{"consensus"},
			"mainnet",
			false,
		},
		{
			map[string][]string{
				"consensus": {"lighthouse"},
				"execution": {"nethermind"},
				"validator": {"lighthouse"},
			},
			[]string{"other"},
			"mainnet",
			true,
		},
		{
			map[string][]string{
				"validator": {"lighthouse", "prysm", "teku", "lodestar", "nimbus"},
				"execution": {"nethermind", "geth", "besu", "erigon"},
			},
			[]string{"execution", "validator"},
			"mainnet",
			false,
		},
		{
			map[string][]string{
				"validator": {"lighthouse", "prysm", "teku", "lodestar", "nimbus"},
				"consensus": {"lighthouse", "prysm", "teku", "lodestar", "nimbus"},
				"execution": {"nethermind", "geth", "besu", "erigon"},
			},
			[]string{"consensus", "other"},
			"mainnet",
			true,
		},
		{
			map[string][]string{
				"validator": {"lighthouse", "teku", "lodestar", "nimbus"},
				"consensus": {"lighthouse", "teku", "lodestar", "nimbus", "grandine"},
				"execution": {"nethermind", "erigon"},
			},
			[]string{"consensus", "execution", "validator"},
			"gnosis",
			false,
		},
		{
			map[string][]string{
				"validator":            {"lighthouse", "prysm", "teku", "lodestar", "nimbus"},
				"consensus":            {"lighthouse", "prysm", "teku", "lodestar", "nimbus", "grandine"},
				"execution":            {"nethermind", "geth", "besu", "erigon"},
				"distributedValidator": {"charon"},
			},
			[]string{"consensus", "execution", "validator", "distributedValidator"},
			"holesky",
			false,
		},
	}

	for i, input := range inputs {
		t.Run(fmt.Sprintf("Network %s, testcase: %d", input.network, i), func(t *testing.T) {
			descr := fmt.Sprintf("Clients(%s)", input.query)

			c := ClientInfo{Network: input.network}
			if res, err := c.Clients(input.query); input.isErr && err == nil {
				t.Errorf("%s expected to fail", descr)
			} else if !input.isErr {
				if err != nil {
					t.Errorf("%s failed: %v", descr, err)
				} else if !validateClients(res, input) {
					t.Errorf("%s got invalid result: %v", descr, res)
				}
			}
		})
	}
}

func TestValidateClient(t *testing.T) {
	inputs := [...]struct {
		client     Client
		clientType string
		isErr      bool
	}{
		{
			client:     Client{},
			clientType: "",
			isErr:      true,
		},
		{
			client: Client{
				Name:      "nethermind",
				Type:      "execution",
				Supported: true,
			},
			clientType: "execution",
			isErr:      false,
		},
		{
			client: Client{
				Name:      "nethermind",
				Type:      "execution",
				Supported: false,
			},
			clientType: "execution",
			isErr:      true,
		},
	}

	for _, input := range inputs {
		descr := fmt.Sprintf("ValidateClient(%v, %s)", input.client, input.clientType)

		if err := ValidateClient(&input.client, input.clientType); input.isErr && err == nil {
			t.Errorf("%s expected to fail", descr)
		} else if !input.isErr && err != nil {
			t.Errorf("%s failed: %v", descr, err)
		}
	}
}
