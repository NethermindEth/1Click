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
	"math/big"

	"github.com/NethermindEth/sedge/configs"
)

var stakingModuleIDs = map[string]*big.Int{
	configs.NetworkHolesky: big.NewInt(4),
	configs.NetworkMainnet: big.NewInt(2),
}

func StakingModuleID(network string) (*big.Int, error) {
	stakingModuleID, found := stakingModuleIDs[network]
	if !found {
		return nil, fmt.Errorf("no staking module ID found for network %s", network)
	}
	return stakingModuleID, nil
}
