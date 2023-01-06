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
	"testing"
)

func TestRootCmdExecute(t *testing.T) {
	loggingLevel := t.TempDir()

	rootCmd := RootCmd()
	rootCmd.SetArgs([]string{"--logLevel", loggingLevel})

	if err := rootCmd.Execute(); err != nil {
		t.Errorf("rootCmd.Execute() failed: %v", err)
	}
}
