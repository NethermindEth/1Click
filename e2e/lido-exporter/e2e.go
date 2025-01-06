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
package e2e

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	base "github.com/NethermindEth/sedge/e2e"
)

type (
	e2eArranger func(t *testing.T, appPath string) (map[string]string, error)
	e2eAct      func(t *testing.T, appPath string) *exec.Cmd
	e2eAssert   func(t *testing.T)
)

type e2eLidoExporterTestCase struct {
	base.E2ETestCase
	arranger e2eArranger
	act      e2eAct
	assert   e2eAssert
	pid      int
	ctx      context.Context
	envVars  []string
}

func newE2ELidoExporterTestCase(t *testing.T, arranger e2eArranger, act e2eAct, assert e2eAssert) *e2eLidoExporterTestCase {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tc := &e2eLidoExporterTestCase{
		E2ETestCase: base.E2ETestCase{
			T:          t,
			TestDir:    t.TempDir(),
			RepoPath:   filepath.Dir(filepath.Dir(wd)),
			BinaryName: "lido-exporter",
		},
		arranger: arranger,
		act:      act,
		assert:   assert,
	}
	t.Logf("Creating new E2E test case (%p). TestDir: %s", tc, tc.TestDir)
	base.CheckGoInstalled(t)
	tc.E2ETestCase.InstallGoModules()
	tc.E2ETestCase.Build()
	return tc
}

func (e *e2eLidoExporterTestCase) run() {
	// Cleanup environment before and after test
	e.Cleanup()
	defer e.Cleanup()
	if e.arranger != nil {
		env, err := e.arranger(e.T, e.BinaryPath())
		if err != nil {
			e.T.Fatalf("error in Arrange step: %v", err)
		}

		for key, value := range env {
			e.setEnvVariable(key, value)
		}
	}
	if e.act != nil {
		cmd := e.act(e.T, e.BinaryPath())
		e.pid = cmd.Process.Pid
	}
	if e.assert != nil {
		e.assert(e.T)
	}
}

func (e *e2eLidoExporterTestCase) setEnvVariable(key, value string) {
	e.T.Helper()
	e.envVars = append(e.envVars, key)
	os.Setenv(key, value)
}

func (e *e2eLidoExporterTestCase) BinaryPath() string {
	binaryName := "lido-exporter"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	return filepath.Join(e.TestDir, binaryName)
}
