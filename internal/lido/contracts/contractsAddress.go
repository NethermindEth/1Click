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
package contracts

import (
	"fmt"

	"github.com/NethermindEth/sedge/configs"
)

type address = map[string]string

const (
	// Contract names
	CSModule                 = "csmodule"
	CSAccounting             = "csaccounting"
	CSFeeDistributor         = "csfeedistributor"
	MEVBoostRelayAllowedList = "mevboostrelayallowedlist"
	Vebo                     = "vebo"
)

var deployedAddresses = map[string]address{
	CSModule: {
		configs.NetworkHolesky: "0x4562c3e63c2e586cD1651B958C22F88135aCAd4f",
		configs.NetworkMainnet: "0xdA7dE2ECdDfccC6c3AF10108Db212ACBBf9EA83F",
	},
	CSAccounting: {
		configs.NetworkHolesky: "0xc093e53e8F4b55A223c18A2Da6fA00e60DD5EFE1",
		configs.NetworkMainnet: "0x4d72BFF1BeaC69925F8Bd12526a39BAAb069e5Da",
	},
	CSFeeDistributor: {
		configs.NetworkHolesky: "0xD7ba648C8F72669C6aE649648B516ec03D07c8ED",
		configs.NetworkMainnet: "0xD99CC66fEC647E68294C6477B40fC7E0F6F618D0",
	},
	MEVBoostRelayAllowedList: {
		configs.NetworkHolesky: "0x2d86C5855581194a386941806E38cA119E50aEA3",
		configs.NetworkMainnet: "0xF95f069F9AD107938F6ba802a3da87892298610E",
	},
	Vebo: {
		configs.NetworkHolesky: "0xffDDF7025410412deaa05E3E1cE68FE53208afcb",
		configs.NetworkMainnet: "0x0De4Ea0184c2ad0BacA7183356Aea5B8d5Bf5c6e",
	},
}

func DeployedAddresses(contractName string) address {
	return deployedAddresses[contractName]
}

func ContractAddressByNetwork(contractName, network string) (string, error) {
	address, found := deployedAddresses[contractName][network]
	if !found {
		return "", fmt.Errorf("no contract code at network %s, please double check a smart contract is deployed on given network", network)
	}
	return address, nil
}
