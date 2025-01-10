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
package csmodule

import (
	"io"
	"log"
	"math/big"
	"testing"
	"unicode/utf8"
)

func TestNodeOpIDs(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)
	tcs := []struct {
		name    string
		network string
	}{
		{
			"NodeOpIDs, Holesky", "holesky",
		},
		{
			"NodeOpIDs, Mainnet", "mainnet",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			nodeOperatorIDs, err := nodeOpIDs(tc.network)
			if err != nil {
				t.Errorf("failed to call  NodeOpIDs: %v", err)
			}

			nodeOperatorsCount, err := nodeOpsCount(tc.network)
			if err != nil {
				t.Errorf("failed to call  nodeOpsCount: %v", err)
			}

			if len(nodeOperatorIDs) != int(nodeOperatorsCount.Int64()) {
				t.Errorf("mismatch: nodeOperatorIDs size (%d) != nodeOperatorsCount (%d)", len(nodeOperatorIDs), nodeOperatorsCount.Int64())
			}
		})
	}
}

func TestNodeOperatorInfo(t *testing.T) {
	log.SetOutput(io.Discard)
	tcs := []struct {
		name            string
		network         string
		nodeID          *big.Int
		expectedAddress string
		wantErr         bool
	}{
		{
			"Valid NodeID, Holesky", "holesky", big.NewInt(13), "0xC870Fd7316956C1582A2c8Fd2c42552cCEC70C88", false,
		},
		{
			"Valid Address, Holesky", "holesky", big.NewInt(4), "0xbA99F374C20A3475De737B466ee68Ad9C38c26AF", false,
		},
		{
			"Invalid Address, Holesky", "holesky", big.NewInt(4), "0xC870Fd7316956C1582A2c8Fd2c46752cCEC70C99", true,
		},
		{
			"Invalid Address, Mainnet", "mainnet", big.NewInt(4), "0xC870Fd7316956C1582A2c8Fd2c46752", true,
		},
		{
			"Valid Address, Mainnet", "mainnet", big.NewInt(1), "0x556fedf2213A31c7Ab9F8bc8Db5B2254261A5B0b", false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			nodeOp, err := NodeOperatorInfo(tc.network, tc.nodeID)
			if err != nil {
				t.Fatalf("failed to call NodeOperatorInfo: %v", err)
			}
			if nodeOp.ManagerAddress.String() != tc.expectedAddress && !tc.wantErr {
				t.Errorf("Not same Manager Address, expected %v, got: %v", tc.expectedAddress, nodeOp.ManagerAddress.Hex())
			}
		})
	}
}

func TestNodeID(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)
	tcs := []struct {
		name           string
		network        string
		expectedNodeID *big.Int
		wantErr        bool
	}{
		{
			"Invalid NodeID, Mainnet", "mainnet", big.NewInt(-2), true,
		},
		{
			"Valid NodeID, Mainnet", "mainnet", big.NewInt(1), false,
		},
		{
			"Valid NodeID, Mainnet #2", "mainnet", big.NewInt(12), false,
		},
		{
			"Valid NodeID, Holesky", "holesky", big.NewInt(4), false,
		},
		{
			"Invalid NodeID, Holesky #1", "holesky", big.NewInt(-4), true,
		},
		{
			"Invalid NodeID, Holesky #2", "holesky", big.NewInt(20000), true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			nodeOp, err := NodeOperatorInfo(tc.network, tc.expectedNodeID)
			if err != nil {
				t.Fatalf("failed to call NodeOperatorInfo: %v", err)
			}
			nodeID, err := NodeID(tc.network, nodeOp.RewardAddress.Hex())
			if err != nil && !tc.wantErr {
				t.Fatalf("failed to call NodeID: %v", err)
			}
			if nodeID != nil && nodeID.Cmp(tc.expectedNodeID) != 0 {
				t.Errorf("not same nodeID, expected %v, got: %v", tc.expectedNodeID, nodeID)
			}
		})
	}
}

func FuzzTestNodeID(f *testing.F) {
	testcases := []struct {
		network string
		nodeID  *big.Int
	}{
		{"holesky", big.NewInt(13)},
		{"holesky", big.NewInt(-1)},
		{"holesky", big.NewInt(40000)},
		{"mainnet", big.NewInt(12)},
		{"mainnet", big.NewInt(-5)},
		{"mainnet", big.NewInt(500000)},
	}

	for _, tc := range testcases {
		f.Add(tc.network, tc.nodeID.String())
	}

	f.Fuzz(func(t *testing.T, network string, nodeIDStr string) {
		// Convert nodeIDStr back to *big.Int
		nodeID, ok := new(big.Int).SetString(nodeIDStr, 10)
		if !ok {
			t.Skip("Skipping invalid big.Int string")
		}

		// Silence logger
		log.SetOutput(io.Discard)

		nodeOp, err := NodeOperatorInfo(network, nodeID)
		if err != nil {
			t.Logf("Expected failure in NodeOperatorInfo: %v", err)
			return
		}

		nodeIDReturned, err := NodeID(network, nodeOp.RewardAddress.Hex())
		if err != nil {
			t.Logf("Expected failure in NodeID: %v", err)
			return
		}

		if nodeIDReturned != nil && utf8.ValidString(network) && nodeIDReturned.Cmp(nodeID) != 0 {
			t.Errorf("not same nodeID, expected %v, got: %v", nodeID, nodeIDReturned)
		}
	})
}
