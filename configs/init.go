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
package configs

import (
	"fmt"
	"strings"
)

// TODO: Remove public level to this variable (NetworksConfigs), use getters to access instead
var NetworksConfigs map[string]NetworkConfig

func init() {
	// TODO: This initialization can be made in the variable declaration
	NetworksConfigs = map[string]NetworkConfig{
		"mainnet": {
			Name:               "mainnet",
			RequireJWT:         true,
			NetworkService:     "merge",
			GenesisForkVersion: "0x00000000",
		},
		"goerli": {
			Name:               "goerli",
			RequireJWT:         true,
			NetworkService:     "merge",
			GenesisForkVersion: "0x00001020",
		},
		"sepolia": {
			Name:               "sepolia",
			RequireJWT:         true,
			NetworkService:     "merge",
			GenesisForkVersion: "0x90000069",
		},
		"chiado": {
			Name:               "chiado",
			RequireJWT:         true,
			NetworkService:     "merge",
			GenesisForkVersion: "0x0000006f",
		},
		"gnosis": {
			Name:               "gnosis",
			RequireJWT:         true,
			NetworkService:     "merge",
			GenesisForkVersion: "0x00000064",
		},
	}
}

// TODO: add doc
// TODO: test
func CheckNetwork(networkName string) error {
	if _, ok := NetworksConfigs[strings.ToLower(networkName)]; !ok {
		return fmt.Errorf(UnknownNetworkError, networkName)
	}
	return nil
}
