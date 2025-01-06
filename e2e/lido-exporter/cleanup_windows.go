//go:build windows
// +build windows

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

	"golang.org/x/sys/windows"
)

func (e *e2eLidoExporterTestCase) Cleanup() {
	if e.pid != 0 {
		// Check if the process is still running
		handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, uint32(e.pid))
		if err == nil {
			defer func() {
				if closeErr := windows.CloseHandle(handle); closeErr != nil {
					e.T.Fatalf("error closing handle: %v", closeErr)
				}
			}()

			// Terminate Process
			err = windows.TerminateProcess(handle, 0) // Exit code 0
			if err != nil {
				e.T.Fatalf("error terminating process %d: %v", e.pid, err)
			}
		}
	}

	// Restore environment variables
	for _, env := range e.envVars {
		os.Unsetenv(env)
	}
}
