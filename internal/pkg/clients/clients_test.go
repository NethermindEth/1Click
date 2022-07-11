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

	"github.com/spf13/viper"
)

// TODO: Add testcases for other networks

func validateSupportedClients(t *testing.T, clientType string, supportedClients []string) {
	//TODO: validate supported clients
}

func TestSupportedClients(t *testing.T) {
	inputs := [...]struct {
		clientType string
		network    string
		isErr      bool
	}{
		{"execution", "mainnet", false},
		{"consensus", "mainnet", false},
		{"validator", "mainnet", false},
		{"random", "mainnet", true},
		{"", "mainnet", false}, // This must fail at validation
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
				validateSupportedClients(t, input.clientType, res)
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

func prepareGetConfigClientsTestCase(testCase clientsTestCase) {
	for key, value := range testCase.configClientsTypes {
		viper.Set(key+"Clients", value)
	}
}

func cleanGetConfigClientsTestCase(_ clientsTestCase) {
	viper.Reset()
}

func cleanAll() {
	viper.Reset()
}

func validateClients(resultClients OrderedClients, tc clientsTestCase) bool {
	//Check if all query clients types are in the result types
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
		{},
		{
			map[string][]string{
				"consensus": {"lighthouse"},
				"execution": {"nethermind"},
				"validator": {"lighthouse"},
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
				"consensus": {"lighthouse"},
				"execution": {"nethermind"},
				"validator": {"lighthouse"},
			},
			[]string{"execution", "validator"},
			"mainnet",
			false,
		},
		{
			map[string][]string{
				"consensus": {"lighthouse"},
				"execution": {"nethermind"},
				"validator": {"lighthouse"},
			},
			[]string{"consensus", "other"},
			"mainnet",
			true,
		},
	}

	t.Cleanup(cleanAll)

	for _, input := range inputs {
		descr := fmt.Sprintf("Clients(%s)", input.query)

		prepareGetConfigClientsTestCase(input)

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

		cleanGetConfigClientsTestCase(input)
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
				"nethermind",
				"execution",
				true,
			},
			clientType: "execution",
			isErr:      false,
		},
		{
			client: Client{
				"nethermind",
				"execution",
				false,
			},
			clientType: "execution",
			isErr:      true,
		},
	}

	for _, input := range inputs {
		descr := fmt.Sprintf("ValidateClient(%v, %s)", input.client, input.clientType)

		if err := ValidateClient(input.client, input.clientType); input.isErr && err == nil {
			t.Errorf("%s expected to fail", descr)
		} else if !input.isErr && err != nil {
			t.Errorf("%s failed: %v", descr, err)
		}
	}
}
