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
package actions

import (
	"github.com/NethermindEth/sedge/configs"
	"github.com/NethermindEth/sedge/internal/pkg/generate"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type GenerateComposeOptions struct {
	GenerationData *generate.GenData
	GenerationPath string
}

func (s *sedgeActions) GenerateCompose(options GenerateComposeOptions) error {
	// Create scripts directory if not exists
	if _, err := os.Stat(options.GenerationPath); os.IsNotExist(err) {
		err = os.MkdirAll(options.GenerationPath, 0o755)
		if err != nil {
			return err
		}
	}

	log.Info(configs.GeneratingDockerComposeScript)
	// open output file
	out, err := os.Create(filepath.Join(options.GenerationPath, configs.DefaultDockerComposeScriptName))
	if err != nil {
		return err
	}
	err = generate.ComposeFile(options.GenerationData, out)
	if err != nil {
		return err
	}

	log.Info(configs.GeneratingEnvFile)
	// open output file
	out, err = os.Create(filepath.Join(options.GenerationPath, configs.DefaultEnvFileName))
	err = generate.EnvFile(options.GenerationData, out)
	if err != nil {
		return err
	}

	err = generate.CleanGenerated(options.GenerationPath)
	if err != nil {
		return err
	}

	return nil
}
