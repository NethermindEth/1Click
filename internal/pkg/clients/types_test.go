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
package clients

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetImageOrDefault_Execution(t *testing.T) {
	tests := []struct {
		client        Client
		expectedImage regexp.Regexp
	}{
		{
			client: Client{
				Name: "geth",
				Type: "execution",
			},
			expectedImage: *regexp.MustCompile(`^ethereum/client-go:v\d+\.\d+\.\d+$`),
		},
		{
			client: Client{
				Name: "besu",
				Type: "execution",
			},
			expectedImage: *regexp.MustCompile(`^hyperledger/besu:\d+\.\d+\.\d+$`),
		},
		{
			client: Client{
				Name: "nethermind",
				Type: "execution",
			},
			expectedImage: *regexp.MustCompile(`^nethermind/nethermind:\d+\.\d+\.\d+$`),
		},
		{
			client: Client{
				Name: "erigon",
				Type: "execution",
			},
			expectedImage: *regexp.MustCompile(`^erigontech/erigon:v\d+\.\d+\.\d+$`),
		},
	}
	for _, test := range tests {
		t.Run(test.client.Name, func(t *testing.T) {
			test.client.SetImageOrDefault("")
			assert.True(t, test.expectedImage.Match([]byte(test.client.Image)))
		})
	}
}

func TestSetImageOrDefault_Consensus(t *testing.T) {
	tests := []struct {
		client        Client
		expectedImage regexp.Regexp
	}{
		{
			client: Client{
				Name: "lighthouse",
				Type: "consensus",
			},
			expectedImage: *regexp.MustCompile(`^sigp/lighthouse:v\d+\.\d+\.\d+$`),
		},
		{
			client: Client{
				Name: "prysm",
				Type: "consensus",
			},
			expectedImage: *regexp.MustCompile(`^gcr.io/prysmaticlabs/prysm/beacon-chain:v\d+\.\d+\.\d+$`),
		},
		{
			client: Client{
				Name: "teku",
				Type: "consensus",
			},
			expectedImage: *regexp.MustCompile(`^consensys/teku:\d+\.\d+\.\d+$`),
		},
		{
			client: Client{
				Name: "lodestar",
				Type: "consensus",
			},
			expectedImage: *regexp.MustCompile(`^chainsafe/lodestar:v\d+\.\d+\.\d+$`),
		},
	}
	for _, test := range tests {
		t.Run(test.client.Name, func(t *testing.T) {
			test.client.SetImageOrDefault("")
			assert.True(t, test.expectedImage.Match([]byte(test.client.Image)))
		})
	}
}

func TestSetImageOrDefault_Validator(t *testing.T) {
	tests := []struct {
		client        Client
		expectedImage regexp.Regexp
	}{
		{
			client: Client{
				Name: "lighthouse",
				Type: "validator",
			},
			expectedImage: *regexp.MustCompile(`^sigp/lighthouse:v\d+\.\d+\.\d+$`),
		},
		{
			client: Client{
				Name: "prysm",
				Type: "validator",
			},
			expectedImage: *regexp.MustCompile(`^gcr.io/prysmaticlabs/prysm/validator:v\d+\.\d+\.\d+$`),
		},
		{
			client: Client{
				Name: "teku",
				Type: "validator",
			},
			expectedImage: *regexp.MustCompile(`^consensys/teku:\d+\.\d+\.\d+$`),
		},
		{
			client: Client{
				Name: "lodestar",
				Type: "validator",
			},
			expectedImage: *regexp.MustCompile(`^chainsafe/lodestar:v\d+\.\d+\.\d+$`),
		},
	}
	for _, test := range tests {
		t.Run(test.client.Name, func(t *testing.T) {
			test.client.SetImageOrDefault("")
			assert.True(t, test.expectedImage.Match([]byte(test.client.Image)))
		})
	}
}

func TestSetImageOrDefault_DistributedValidator(t *testing.T) {
	tests := []struct {
		client        Client
		expectedImage regexp.Regexp
	}{
		{
			client: Client{
				Name: "charon",
				Type: "distributedValidator",
			},
			expectedImage: *regexp.MustCompile(`^ghcr.io/obolnetwork/charon:v\d+\.\d+\.\d+$`),
		},
	}
	for _, test := range tests {
		t.Run(test.client.Name, func(t *testing.T) {
			test.client.SetImageOrDefault("")
			t.Logf("print %s", test.client.Image)
			assert.True(t, test.expectedImage.Match([]byte(test.client.Image)))
		})
	}
}

func TestSetImageOrDefault_CustomImage(t *testing.T) {
	tests := []struct {
		client      Client
		customImage string
	}{
		{
			client: Client{
				Name: "geth",
				Type: "execution",
			},
			customImage: "my/geth-image:v1.0.0",
		},
		{
			client: Client{
				Name: "lighthouse",
				Type: "consensus",
			},
			customImage: "my/lighthouse-image:v1.0.0",
		},
		{
			client: Client{
				Name: "prysm",
				Type: "validator",
			},
			customImage: "my/prysm-image:v1.0.0",
		},
		{
			client: Client{
				Name: "charon",
				Type: "distributedValidator",
			},
			customImage: "my/charon-image:v1.0.0",
		},
	}
	for _, test := range tests {
		t.Run(test.client.Name, func(t *testing.T) {
			test.client.SetImageOrDefault(test.customImage)
			assert.Equal(t, test.customImage, test.client.Image)
		})
	}
}
