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
)

func TestKeysStatus(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)
	tcs := []struct {
		name      string
		network   string
		nodeID    *big.Int
		invalidID bool
	}{
		{
			"Valid NodeID, Holesky #1", "holesky", big.NewInt(13), false,
		},
		{
			"Valid NodeID, Holesky #2", "holesky", big.NewInt(4), false,
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
			keys, err := KeysStatus(tc.network, tc.nodeID)
			if err != nil && !tc.invalidID {
				t.Fatalf("failed to call KeysStatus: %v", err)
			}
			nodeInfo, err := NodeOperatorInfo(tc.network, tc.nodeID)
			if err != nil {
				t.Fatalf("failed to call NodeOperatorInfo: %v", err)
			}
			expectedDeposited := big.NewInt(int64(nodeInfo.TotalDepositedKeys))
			deposited := keys.DepositedValidators
			if deposited == nil && tc.invalidID {
				t.Skipf("Expected nil value for deposited keys")
			} else if deposited == nil && !tc.invalidID {
				t.Fatalf("invalid deposited value: expected a value, got nil")
			}
			if deposited.Cmp(expectedDeposited) != 0 {
				t.Errorf("Not same nodeID, expected %v, got: %v", expectedDeposited, deposited)
			}
		})
	}
}

func FuzzTestSigningKeys(f *testing.F) {
	testcases := []struct {
		network string
		nodeID  *big.Int
	}{
		{"holesky", big.NewInt(13)},
		{"holesky", big.NewInt(-1)},
		{"holesky", big.NewInt(40000)},
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
		signingKeys, err := SigningKeys(network, nodeID)
		if err != nil {
			t.Logf("Expected failure in SigningKeys: %v", err)
			return
		}
		keysCount, err := nodeOpNonWithdrawnKeys(network, nodeID)
		if err != nil {
			t.Logf("Expected failure in nodeOpNonWithdrawnKeys: %v", err)
			return
		}
		if len(signingKeys) != int(keysCount.Int64()) {
			t.Errorf("mismatch: keys size (%d) != keysCount (%d)", len(signingKeys), keysCount.Int64())
		}
	})
}
