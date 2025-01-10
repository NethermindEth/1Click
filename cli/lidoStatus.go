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
	"errors"
	"flag"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/NethermindEth/sedge/internal/lido/contracts"
	bonds "github.com/NethermindEth/sedge/internal/lido/contracts/csaccounting"
	rewards "github.com/NethermindEth/sedge/internal/lido/contracts/csfeedistributor"
	"github.com/NethermindEth/sedge/internal/lido/contracts/csmodule"
	"github.com/NethermindEth/sedge/internal/ui"
	"github.com/gosuri/uiprogress"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type lidoData struct {
	nodeID   *big.Int
	nodeInfo csmodule.NodeOperator
	keys     csmodule.Keys
	bondInfo bonds.BondInfo
	rewards  *big.Int
}

var (
	rewardAddress    string
	networkName      string
	longDescriptions bool
	nodeIDInt        int64
)

const (
	nodeOpInfo  = `Node Operator Info`
	keysInfo    = `Keys`
	queueInfo   = `Queue`
	bondInfo    = `Bond`
	rewardsInfo = `Rewards`
)

func LidoStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lido-status [flags] [args]",
		Short: "Display status and information of Lido Node Operator",
		Long: `This command retrieves and displays the status and detailed information of Lido Node Operators. 

This information includes:
- Node Operator ID.
- Keys and queue information: available for deposit (in the queue), stuck, refunded, exited, deposited.
- Bond and rewards information: total amount, amounts lower and higher than required, non-claimed rewards.

Valid args: reward address of Node Operator (rewards recipient)`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if cobra.ExactArgs(1)(cmd, args) != nil {
					return errors.New("requires one argument")
				}
				rewardAddress = args[0]
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ui.EthAddressValidator(rewardAddress, false); err != nil && len(args) != 0 {
				return err
			}
			if len(args) == 0 && nodeIDInt < 0 {
				return errors.New("must provide reward address or node ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runListLidoStatusCmd(cmd, args); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&networkName, "network", "n", "holesky", "Target network. e.g. holesky, mainnet etc.")
	cmd.Flags().BoolVar(&longDescriptions, "l", false, "Show detailed descriptions for each value")
	cmd.Flags().Int64VarP(&nodeIDInt, "nodeID", "i", -1, "Your Node Operator ID (optional)")
	cmd.Flags().SortFlags = false
	return cmd
}

func runListLidoStatusCmd(cmd *cobra.Command, args []string) error {
	log.Infof("Retrieving Lido Node Operator Information\n")

	nodeData, err := nodeData()
	if err != nil {
		return err
	}
	dataMap := buildLidoData(nodeData)

	// Extract headers and sort by weight
	headers := make([]string, 0, len(dataMap))
	for header := range dataMap {
		headers = append(headers, header)
	}

	// Sort headers by their weights
	sort.SliceStable(headers, func(i, j int) bool {
		return dataMap[headers[i]].weight < dataMap[headers[j]].weight
	})

	log.Infof("Listing Node Operator Information")
	for _, header := range headers {
		ui.WriteLidoStatusTable(cmd.OutOrStdout(), dataMap[header].data, header)
	}
	return nil
}

// Get the data for the Node Operator
func nodeData() (*lidoData, error) {
	nodeData := &lidoData{}
	var nodeID *big.Int
	var err error

	steps := []string{
		"Fetching NO Info",
		"Fetching Keys & Queue",
		"Fetching Bond Info",
		"Fetching Rewards Data",
	}

	if !isTestEnv() {
		uiprogress.Start()
	}

	bar := uiprogress.AddBar(len(steps)).AppendCompleted()
	// Progress bar label setup
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		if b.Current() > 0 {
			return steps[b.Current()-1]
		}
		return "Retrieving Node Operator.."
	})

	if nodeIDInt < 0 {
		nodeID, err = csmodule.NodeID(networkName, rewardAddress)
		if err != nil {
			return nodeData, err
		}
	} else {
		nodeID = big.NewInt(nodeIDInt)
	}
	bar.Incr()

	nodeInfo, err := csmodule.NodeOperatorInfo(networkName, nodeID)
	if err != nil {
		return nodeData, err
	}
	time.Sleep(time.Second / 10)
	bar.Incr()

	keys, err := csmodule.KeysStatus(networkName, nodeID)
	if err != nil {
		return nodeData, err
	}
	time.Sleep(time.Second / 10)
	bar.Incr()

	bond, err := bonds.BondSummary(networkName, nodeID)
	if err != nil {
		return nodeData, err
	}
	time.Sleep(time.Second / 10)
	bar.Incr()

	reward, err := rewards.Rewards(networkName, nodeID)
	if err != nil {
		return nodeData, err
	}
	time.Sleep(time.Second / 10)
	bar.Incr()

	if !isTestEnv() {
		uiprogress.Stop()
	}

	nodeData.nodeID = nodeID
	nodeData.nodeInfo = nodeInfo
	nodeData.keys = keys
	nodeData.bondInfo = bond
	nodeData.rewards = reward

	return nodeData, nil
}

// Structure the data to be displayed
func buildLidoData(node *lidoData) map[string]struct {
	data   []string
	weight int
} {
	var nodeOpDetailed, keysDetailed, queueDetailed, bondDetailed, rewardsDetailed string
	var currentBond, requiredBond, excessBond, missedBond, rewards decimal.Decimal
	rewardAddressLink := fmt.Sprintf(`https://etherscan.io/address/%s`, node.nodeInfo.RewardAddress)
	var prefix string
	if networkName == "mainnet" {
		prefix = ""
	} else {
		prefix = networkName + "."
	}
	claimRewardsLink := fmt.Sprintf(`https://%setherscan.io/address/%s#writeProxyContract#F10`, prefix, contracts.DeployedAddresses(contracts.CSModule)[networkName])

	detailedDescriptions := map[string]string{
		nodeOpInfo: `
## Description 
- Node Operator ID: Unique identifier for the node operator.
- Reward Address: Address that is the ultimate recipient of the rewards
- Manager Address: Address used to perform routine management operations regarding the CSM Node Operator.`,

		keysInfo: `
## Description 
- Stuck Keys Count: Number of keys stuck in the system. A validator is considered to be "stuck" if it has not been exited timely following an exit signal from the protocol.
- Refunded Keys Count: Number of keys that were refunded.
- Exited Keys Count: Number of keys that have exited.
- Deposited Keys Count: Number of keys currently deposited.
- Depositable Keys Count: Number of keys eligible for deposits.`,

		queueInfo: `
## Description 
- Keys in the deposit queue: Number of the depositable keys that are in the deposit queue.`,

		bondInfo: `
## Description 
- Bond : a security collateral that Node Operators must submit before uploading validator keys into CSM. It covers possible losses caused by inappropriate actions on the Node Operator's side.
- Current Bond: The current amount of bonded ETH.
- Required Bond: The required amount of ETH to maintain.
- Excess Bond: The amount of excess bond over the required amount.
- Missed Bond: The amount of bond that is missing.`,

		rewardsInfo: `
## Description 
- Non-claimed Rewards: The amount of rewards available for claiming.`,
	}

	if longDescriptions {
		nodeOpDetailed = detailedDescriptions[nodeOpInfo]
		keysDetailed = detailedDescriptions[keysInfo]
		queueDetailed = detailedDescriptions[queueInfo]
		bondDetailed = detailedDescriptions[bondInfo]
		rewardsDetailed = detailedDescriptions[rewardsInfo]
		currentBond = weiToEth(node.bondInfo.Current)
		excessBond = weiToEth(node.bondInfo.Excess)
		missedBond = weiToEth(node.bondInfo.Missed)
		requiredBond = weiToEth(node.bondInfo.Required)
		rewards = weiToEth(node.rewards)
	} else {
		currentBond = weiToEth(node.bondInfo.Current).Round(1)
		excessBond = weiToEth(node.bondInfo.Excess).Round(1)
		missedBond = weiToEth(node.bondInfo.Missed).Round(1)
		requiredBond = weiToEth(node.bondInfo.Required).Round(1)
		rewards = weiToEth(node.rewards).Round(1)
	}

	data := map[string]struct {
		data   []string
		weight int
	}{
		nodeOpInfo: {
			data: []string{
				fmt.Sprintf(`- **Node Operator ID:** %s`, node.nodeID.String()),
				fmt.Sprintf(`- **Reward Address:** %s`, node.nodeInfo.RewardAddress.String()),
				fmt.Sprintf(`- **Manager Address:** %s`, node.nodeInfo.ManagerAddress.String()),
				fmt.Sprintf(`- [Reward Address Link on etherscan](%s)`, rewardAddressLink),
				nodeOpDetailed,
			},
			weight: 1,
		},
		keysInfo: {
			data: []string{
				fmt.Sprintf(`- **Stuck Keys Count:** %s`, node.keys.StuckValidatorsCount.String()),
				fmt.Sprintf(`- **Refunded Keys Count:** %s`, node.keys.RefundedValidatorsCount.String()),
				fmt.Sprintf(`- **Exited Keys Count:** %s`, node.keys.ExitedValidators.String()),
				fmt.Sprintf(`- **Deposited Keys Count:** %s`, node.keys.DepositedValidators.String()),
				fmt.Sprintf(`- **Depositable Keys Count:** %s`, node.keys.DepositableValidatorsCount.String()),
				keysDetailed,
			},
			weight: 2,
		},
		queueInfo: {
			data: []string{
				fmt.Sprintf(`- **Keys in the deposit queue:** %d`, node.nodeInfo.EnqueuedCount),
				queueDetailed,
			},
			weight: 3,
		},
		bondInfo: {
			data: []string{
				fmt.Sprintf(`- **Current Bond:** %s ETH`, currentBond.String()),
				fmt.Sprintf(`- **Required Bond:** %s ETH`, requiredBond.String()),
				fmt.Sprintf(`- **Excess Bond:** %s ETH`, excessBond.String()),
				fmt.Sprintf(`- **Missed Bond:** %s ETH`, missedBond.String()),
				bondDetailed,
			},
			weight: 4,
		},
		rewardsInfo: {
			data: []string{
				fmt.Sprintf(`- **Non-claimed Rewards:** %s ETH`, rewards.String()),
				fmt.Sprintf(`- [Claim your rewards here!](%s)`, claimRewardsLink),
				rewardsDetailed,
			},
			weight: 5,
		},
	}

	return data
}

// Convert Wei to Ether
func weiToEth(wei *big.Int) decimal.Decimal {
	weiToEther := decimal.NewFromBigInt(big.NewInt(1e18), 0)
	weiDecimal := decimal.NewFromBigInt(wei, 0)
	return weiDecimal.Div(weiToEther)
}

// Used to disable progress bar when running tests
func isTestEnv() bool {
	return flag.Lookup("test.v") != nil
}
