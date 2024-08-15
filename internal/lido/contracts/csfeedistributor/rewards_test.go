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
			name:    "Rewards for nodeID 1, Holesky",
			network: "holesky",
			nodeID:  big.NewInt(1),
			wantErr: false,
		},
		{
			name:    "Rewards for nodeID 182, Holesky",
			network: "holesky",
			nodeID:  big.NewInt(182),
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

func TestBinarySearchNodeID(t *testing.T) {
	tcs := []struct {
		name      string
		treeData  Tree
		nodeID    *big.Int
		expected  int
		expectErr bool
	}{
		{
			name: "Node ID found at index 1",
			treeData: Tree{
				Values: []struct {
					Value     []interface{} `json:"value"`
					TreeIndex int           `json:"treeIndex"`
				}{
					{Value: []interface{}{1.0, 5000.0}, TreeIndex: 0},
					{Value: []interface{}{2.0, 6000.0}, TreeIndex: 1},
					{Value: []interface{}{30.0, 7000.0}, TreeIndex: 2},
				},
			},
			nodeID:   big.NewInt(2),
			expected: 1,
		},
		{
			name: "Node ID found at index 0",
			treeData: Tree{
				Values: []struct {
					Value     []interface{} `json:"value"`
					TreeIndex int           `json:"treeIndex"`
				}{
					{Value: []interface{}{100.0, 5000.0}, TreeIndex: 0},
					{Value: []interface{}{200.0, 6000.0}, TreeIndex: 1},
					{Value: []interface{}{330.0, 7000.0}, TreeIndex: 2},
				},
			},
			nodeID:   big.NewInt(100),
			expected: 0,
		},
		{
			name: "Node ID not found",
			treeData: Tree{
				Values: []struct {
					Value     []interface{} `json:"value"`
					TreeIndex int           `json:"treeIndex"`
				}{
					{Value: []interface{}{10.0, 5000}, TreeIndex: 0},
					{Value: []interface{}{20.0, 6000}, TreeIndex: 1},
					{Value: []interface{}{30.0, 7000}, TreeIndex: 2},
				},
			},
			nodeID:    big.NewInt(400),
			expected:  0,
			expectErr: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			index, err := binarySearchNodeID(tc.nodeID, tc.treeData)
			if (err != nil) != tc.expectErr {
				t.Errorf("expected error: %v, got: %v", tc.expectErr, err)
			}
			if err == nil && index != tc.expected {
				t.Errorf("expected index: %v, got: %v", tc.expected, index)
			}
		})
	}
}
