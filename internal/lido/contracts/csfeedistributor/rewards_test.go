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
package csfeedistributor

import (
	"io"
	"log"
	"math/big"
	"testing"

	bond "github.com/NethermindEth/sedge/internal/lido/contracts/csaccounting"
)

func TestRewards(t *testing.T) {
	log.SetOutput(io.Discard)
	tcs := []struct {
		name    string
		network string
		nodeID  *big.Int
		wantErr bool
	}{
		{
			name:    "Rewards for nodeID 10, Holesky",
			network: "holesky",
			nodeID:  big.NewInt(10),
			wantErr: false,
		},
		{
			name:    "Rewards for nodeID 250, Holesky",
			network: "holesky",
			nodeID:  big.NewInt(250),
			wantErr: false,
		},
		{
			name:    "Rewards for nodeID 113, Holesky",
			network: "holesky",
			nodeID:  big.NewInt(113),
			wantErr: false,
		},
		{
			name:    "Invalid nodeID, Holesky",
			network: "holesky",
			nodeID:  big.NewInt(-3),
			wantErr: true,
		},
		{
			name:    "Invalid nodeID, Mainnet",
			network: "mainnet",
			nodeID:  big.NewInt(-15),
			wantErr: true,
		},
		{
			name:    "Rewards for nodeID 1, Mainnet",
			network: "mainnet",
			nodeID:  big.NewInt(1),
			wantErr: false,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			treeCID, err := treeCID(tc.network)
			if err != nil {
				t.Fatalf("failed to call treeCID: %v", err)
			}

			rewards, err := Rewards(tc.network, tc.nodeID)
			if err != nil && !tc.wantErr {
				t.Fatalf("failed to call Rewards: %v", err)
			}

			fees, err := cumulativeFeeShares(treeCID, tc.nodeID)
			if err != nil && !tc.wantErr {
				t.Fatalf("failed to call cumulativeFeeShares: %v", err)
			}

			bond, err := bond.BondSummary(tc.network, tc.nodeID)
			if err != nil && !tc.wantErr {
				t.Fatalf("failed to call BondSummary: %v", err)
			}

			if rewards == nil && tc.wantErr {
				t.Skipf("Expected nil value for rewards")
			} else if rewards == nil && !tc.wantErr {
				t.Fatalf("invalid rewards value: expected a value, got nil")
			}
			expectedRewards := new(big.Int).Add(bond.Excess, fees)
			if rewards.Cmp(expectedRewards) != 0 {
				t.Errorf("invalid rewards amount, expected %v, got: %v", expectedRewards, rewards)
			}
		})
	}
}

func TestConvertTreeValuesToBigInt(t *testing.T) {
	tcs := []struct {
		input    interface{}
		expected *big.Int
	}{
		{12345.0, big.NewInt(12345)},
		{67890.0, big.NewInt(67890)},
		{0.0, big.NewInt(0)},
	}

	for _, tc := range tcs {
		result, err := convertTreeValuesToBigInt(tc.input)
		if err != nil {
			t.Errorf("convertTreeValuesToBigInt(%v) returned error: %v", tc.input, err)
		}

		if result.Cmp(tc.expected) != 0 {
			t.Errorf("convertTreeValuesToBigInt(%v) = %v; expected %v", tc.input, result, tc.expected)
		}
	}
}
