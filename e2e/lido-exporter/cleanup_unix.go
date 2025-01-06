//go:build !windows
// +build !windows

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
	"os"
	"syscall"
)

func (e *e2eLidoExporterTestCase) Cleanup() {
	if e.pid != 0 {
		process, err := os.FindProcess(e.pid)
		if err == nil {
			process.Signal(syscall.SIGTERM)
		}
	}

	// Restore environment variables
	for _, env := range e.envVars {
		os.Unsetenv(env)
	}
}
