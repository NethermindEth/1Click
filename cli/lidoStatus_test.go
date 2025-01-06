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
package cli

import (
	"bytes"
	"io"
	"math/big"
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type flags struct {
	rewardAddress    string
	networkName      string
	longDescriptions bool
	nodeIDInt        int64
}

type lidoStatusCmdTestCase struct {
	name  string
	flags flags
	isErr bool
}

func (f flags) argsList() []string {
	s := make([]string, 0)
	s = append(s, f.rewardAddress)
	if f.networkName != "" {
		s = append(s, "--network", f.networkName)
	}
	if f.longDescriptions {
		s = append(s, "--l")
	}
	if f.nodeIDInt >= 0 {
		s = append(s, "--nodeID", strconv.FormatInt(f.nodeIDInt, 10))
	}

	return s
}

func TestLidoStatusCmd(t *testing.T) {
	tcs := []lidoStatusCmdTestCase{
		{
			name: "valid lido-status flags with long descriptions, Holesky",
			flags: flags{
				rewardAddress:    "0xC870Fd7316956C1582A2c8Fd2c42552cCEC70C88",
				networkName:      "holesky",
				longDescriptions: true,
			},
			isErr: false,
		},
		{
			name: "valid lido-status flags, Holesky",
			flags: flags{
				rewardAddress: "0xe6b5A31d8bb53D2C769864aC137fe25F4989f1fd",
				networkName:   "holesky",
			},
			isErr: false,
		},
		{
			name: "Invalid: missing address, Holesky",
			flags: flags{
				networkName: "holesky",
			},
			isErr: true,
		},
		{
			name: "Invalid: incorrect address, Holesky",
			flags: flags{
				rewardAddress: "0xC870Fd",
				networkName:   "holesky",
			},
			isErr: true,
		},
		{
			name: "Invalid: address missing 0x prefix, Holesky",
			flags: flags{
				rewardAddress: "22bA5CaFB5E26E6Fe51f330294209034013A5A4c",
				networkName:   "holesky",
			},
			isErr: true,
		},
		{
			name: "Valid node ID, Holesky",
			flags: flags{
				rewardAddress: "0xe6b5A31d8bb53D2C769864aC137fe25F4989f1fd", // rewardAddress should be ignored
				networkName:   "holesky",
				nodeIDInt:     1,
			},
			isErr: false,
		},
		{
			name: "Invalid: negative node ID, Holesky",
			flags: flags{
				networkName: "holesky",
				nodeIDInt:   -2,
			},
			isErr: true,
		},
		{
			name: "Invalid: missing address, Mainnet",
			flags: flags{
				networkName: "mainnet",
			},
			isErr: true,
		},
		{
			name: "Invalid: incorrect address, Mainnet",
			flags: flags{
				rewardAddress: "0xC870Fd10dd",
				networkName:   "mainnet",
			},
			isErr: true,
		},
		{
			name: "Invalid: address missing 0x prefix, Mainnet",
			flags: flags{
				rewardAddress: "bA99F374C20A3475De737B466ee68Ad9C38c26AF",
				networkName:   "mainnet",
			},
			isErr: true,
		},
		// {
		// 	name: "Valid node ID, Mainnet",
		// 	flags: flags{
		// 		rewardAddress: "0xe6b5A31d8bb53D2C769864aC137fe25F4989f1fd", // rewardAddress should be ignored
		// 		networkName: "mainnet",
		// 		nodeIDInt:   1,
		// 	},
		// 	isErr: false,
		// },
		// {
		// 	name: "Valid node ID with long description, Mainnet",
		// 	flags: flags{
		// 		rewardAddress: "0xe6b5A31d8bb53D2C769864aC137fe25F4989f1fd", // rewardAddress should be ignored
		// 		networkName: "mainnet",
		// 		nodeIDInt:   1,
		// 		longDescriptions: true,
		// 	},
		// 	isErr: false,
		// },
		{
			name: "Invalid: negative node ID, Mainnet",
			flags: flags{
				networkName: "mainnet",
				nodeIDInt:   -14,
			},
			isErr: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			log.SetOutput(io.Discard)
			logsOut := new(bytes.Buffer)
			tableOut := new(bytes.Buffer)
			rootCmd := RootCmd()
			rootCmd.SetOut(tableOut)
			rootCmd.AddCommand(LidoStatusCmd())
			argsL := append([]string{"lido-status"}, tc.flags.argsList()...)
			rootCmd.SetArgs(argsL)
			initLogging()
			log.SetOutput(logsOut)
			rootCmd.SetOutput(io.Discard)
			err := rootCmd.Execute()
			if tc.isErr && err == nil {
				t.Error("sedge lido-status expected to fail")
			} else if !tc.isErr && err != nil {
				t.Errorf("sedge lido-status failed: %v", err)
			}
		})
	}
}

func TestWeiToEth(t *testing.T) {
	tests := []struct {
		name     string
		wei      *big.Int
		expected decimal.Decimal
	}{
		{
			name:     "100 Wei",
			wei:      big.NewInt(100),
			expected: decimal.NewFromFloat(0.0000000000000001),
		},
		{
			name:     "1,000 Wei",
			wei:      big.NewInt(1000),
			expected: decimal.NewFromFloat(0.000000000000001),
		},
		{
			name:     "50,000 Wei",
			wei:      big.NewInt(50000),
			expected: decimal.NewFromFloat(0.00000000000005),
		},
		{
			name:     "0 Wei",
			wei:      big.NewInt(0),
			expected: decimal.NewFromFloat(0),
		},
		{
			name:     "expected 1 ETH",
			wei:      new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
			expected: decimal.NewFromFloat(1),
		},
		{
			name:     "expected 1.5 ETH",
			wei:      new(big.Int).Set(new(big.Int).Mul(big.NewInt(15), new(big.Int).Exp(big.NewInt(10), big.NewInt(17), nil))),
			expected: decimal.NewFromFloat(1.5),
		},
		{
			name:     "expected 200 ETH",
			wei:      new(big.Int).Set(new(big.Int).Mul(big.NewInt(200), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))),
			expected: decimal.NewFromFloat(200),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := weiToEth(tt.wei)
			if !result.Equal(tt.expected) {
				t.Errorf("weiToEth(%v) = %v; expected %v", tt.wei, result, tt.expected)
			}
		})
	}
}
