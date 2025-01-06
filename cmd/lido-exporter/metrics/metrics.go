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
package metrics

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

func InitMetrics(nodeOperatorID, network string) {
	// Node operator metrics
	nodeOperatorIDGauge.WithLabelValues(nodeOperatorID, network)
	nodeOperatorManagerAddressGauge.WithLabelValues(nodeOperatorID, network)
	nodeOperatorRewardAddressGauge.WithLabelValues(nodeOperatorID, network)

	// Keys metrics
	keysStuckValidatorsCountGauge.WithLabelValues(nodeOperatorID, network)
	keysRefundedValidatorsCountGauge.WithLabelValues(nodeOperatorID, network)
	keysExitedValidatorsCountGauge.WithLabelValues(nodeOperatorID, network)
	keysDepositedValidatorsCountGauge.WithLabelValues(nodeOperatorID, network)
	keysDepositableValidatorsCountGauge.WithLabelValues(nodeOperatorID, network)
	addedKeysCountGauge.WithLabelValues(nodeOperatorID, network)
	withdrawnKeysCountGauge.WithLabelValues(nodeOperatorID, network)
	vettedKeysCountGauge.WithLabelValues(nodeOperatorID, network)
	enqueuedKeysCountGauge.WithLabelValues(nodeOperatorID, network)

	// Penalties metrics
	penaltiesTotal.WithLabelValues(nodeOperatorID, "el_rewards_stealing", "")
	penaltiesTotal.WithLabelValues(nodeOperatorID, "withdrawal", "")
	penaltiesTotal.WithLabelValues(nodeOperatorID, "initial_slashing", "")

	// Exit requests metrics
	exitRequestsTotal.WithLabelValues(nodeOperatorID, network, "")

	// Bond metrics
	bondCurrentGauge.WithLabelValues(nodeOperatorID, network)
	bondRequiredGauge.WithLabelValues(nodeOperatorID, network)
	bondExcessGauge.WithLabelValues(nodeOperatorID, network)
	bondMissedGauge.WithLabelValues(nodeOperatorID, network)

	// Rewards metrics
	nonClaimedRewardsGauge.WithLabelValues(nodeOperatorID, network)
}

func CollectMetrics(ctx context.Context, client *ethclient.Client, wsClient *ethclient.Client, nodeOperatorID *big.Int, network string, scrapeTime time.Duration) {
	var filterOpts *bind.WatchOpts

	// Get start block
	// Try to use current block - 50000, if not then use 0
	currentBlock, err := client.BlockNumber(ctx)
	if err != nil {
		log.Errorf("Failed to get current block: %v", err)
		zero := uint64(0)
		filterOpts = &bind.WatchOpts{
			Start: &zero,
		}
	} else {
		startBlock := currentBlock - 50000
		if startBlock < 0 {
			startBlock = 0
		}
		filterOpts = &bind.WatchOpts{
			Start: &startBlock,
		}
	}

	log.Infof("Collecting metrics for network %s with start block %d", network, filterOpts.Start)
	go collectPenalties(ctx, network, nodeOperatorID, wsClient, filterOpts)
	go collectExitRequests(ctx, network, nodeOperatorID, wsClient, filterOpts)
	go collectNodeOperatorInfo(ctx, network, nodeOperatorID, scrapeTime)
	go collectKeysInfo(ctx, network, nodeOperatorID, scrapeTime)
	go collectBondInfo(ctx, network, nodeOperatorID, scrapeTime)
	go collectRewardsInfo(ctx, network, nodeOperatorID, scrapeTime)
}
