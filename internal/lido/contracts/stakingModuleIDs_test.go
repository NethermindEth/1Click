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
	"math/big"
	"testing"

	"github.com/NethermindEth/sedge/configs"
)

func TestStakingModuleID(t *testing.T) {
	tcs := []struct {
		network string
		want    *big.Int
		wantErr bool
	}{
		{configs.NetworkHolesky, big.NewInt(4), false},
		{"unknown", nil, true},
		{configs.NetworkMainnet, big.NewInt(2), false},
	}

	for _, tc := range tcs {
		got, err := StakingModuleID(tc.network)
		if tc.wantErr && err == nil {
			t.Errorf("StakingModuleID(%s) returned no error, want error", tc.network)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("StakingModuleID(%s) returned error: %v", tc.network, err)
		}
		if got.Cmp(tc.want) != 0 {
			t.Errorf("StakingModuleID(%s) = %v, want %v", tc.network, got, tc.want)
		}
	}
}
