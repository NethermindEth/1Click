package mevboostrelaylist

import (
	"io"
	"log"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

type RelayData struct {
	Mainnet []Relay `yaml:"mainnet"`
	Holesky []Relay `yaml:"holesky"`
}

func loadRelays(filename string) (map[string][]Relay, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var relayData RelayData
	err = yaml.Unmarshal(data, &relayData)
	if err != nil {
		return nil, err
	}

	return map[string][]Relay{
		"mainnet": relayData.Mainnet,
		"holesky": relayData.Holesky,
	}, nil
}

func TestGetRelays(t *testing.T) {
	// Silence logger
	log.SetOutput(io.Discard)

	// Load relays from the YAML file
	relaysMap, err := loadRelays("relays.yaml")
	if err != nil {
		t.Fatalf("Failed to load relays: %v", err)
	}

	tcs := []struct {
		name           string
		network        string
		expectedRelays []Relay
	}{
		{
			"GetRelays Mainnet", "mainnet", relaysMap["mainnet"],
		},
		{
			"GetRelays Holesky", "holesky", relaysMap["holesky"],
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			relays, err := GetRelays(tc.network)
			if err != nil {
				t.Fatalf("Failed to call GetRelays %v", err)
			}

			if !reflect.DeepEqual(relays, tc.expectedRelays) {
				t.Errorf("Relays do not match expected values\nGot: %+v\nExpected: %+v", relays, tc.expectedRelays)
			}
		})
	}
}