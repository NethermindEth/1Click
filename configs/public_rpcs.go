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

type RPC struct {
	NetworkName string
	PublicRPCs  []string
	PublicWSs   []string
}

var networkRPCs = map[string]RPC{
	NetworkMainnet: {
		NetworkName: NetworkMainnet,
		PublicRPCs: []string{
			"https://eth.llamarpc.com",
			"https://eth-pokt.nodies.app",
			"https://rpc.mevblocker.io",
			"https://ethereum-rpc.publicnode.com",
			"https://rpc.flashbots.net",
			"https://eth.drpc.org",
		},
		PublicWSs: []string{
			"wss://ethereum-rpc.publicnode.com",
		},
	},
	NetworkHolesky: {
		NetworkName: NetworkHolesky,
		PublicRPCs: []string{
			"https://ethereum-holesky-rpc.publicnode.com",
			"https://endpoints.omniatech.io/v1/eth/holesky/public",
			"https://ethereum-holesky.blockpi.network/v1/rpc/public",
			"https://1rpc.io/holesky",
			"https://holesky.drpc.org",
		},
		PublicWSs: []string{
			"wss://ethereum-holesky-rpc.publicnode.com",
		},
	},
}
