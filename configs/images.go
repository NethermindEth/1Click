package configs

import (
	_ "embed"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Image struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

func (i *Image) String() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Version)
}

var ClientImages struct {
	Execution struct {
		Geth       Image `yaml:"geth"`
		Besu       Image `yaml:"besu"`
		Nethermind Image `yaml:"nethermind"`
		Erigon     Image `yaml:"erigon"`
	}
	Consensus struct {
		Lighthouse Image `yaml:"lighthouse"`
		Lodestar   Image `yaml:"lodestar"`
		Teku       Image `yaml:"teku"`
		Prysm      Image `yaml:"prysm"`
		Nimbus     Image `yaml:"nimbus"`
	}
	Validator struct {
		Lighthouse Image `yaml:"lighthouse"`
		Lodestar   Image `yaml:"lodestar"`
		Teku       Image `yaml:"teku"`
		Prysm      Image `yaml:"prysm"`
		Nimbus     Image `yaml:"nimbus"`
		Charon     Image `yaml:"charon"`
	}
	Distributed struct {
		Charon Image `yaml:"charon"`
	}
	Optimism struct {
		OpNode Image `yaml:"opnode"`
	}
	OpExecution struct {
		OpGeth Image `yaml:"opgeth"`
		OpNeth Image `yaml:"opnethermind"`
		OpReth Image `yaml:"opreth"`
	}
}

//go:embed client_images.yaml
var clientImages string

func init() {
	err := yaml.Unmarshal([]byte(clientImages), &ClientImages)
	if err != nil {
		log.Fatal(err)
	}
}
