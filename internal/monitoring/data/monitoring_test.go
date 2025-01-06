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
package data

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"text/template"
	"time"

	mocks "github.com/NethermindEth/sedge/internal/monitoring/locker/mocks"
	"github.com/NethermindEth/sedge/internal/monitoring/services/templates"
	"github.com/NethermindEth/sedge/internal/monitoring/utils"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	// Create an in-memory filesystem
	afs := afero.NewMemMapFs()

	// Create a mock locker
	ctrl := gomock.NewController(t)
	locker := mocks.NewMockLocker(ctrl)

	// Expect the lock to initialized
	locker.EXPECT().New(utils.PathMatcher{Expected: "/.lock"}).Return(locker)

	// Create a new MonitoringStack with the in-memory filesystem
	stack := &MonitoringStack{
		path: "/",
		l:    locker,
		fs:   afs,
	}

	err := stack.Init()
	assert.NoError(t, err)

	// Check that the file was created
	exists, err := afero.Exists(afs, "/.lock")
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.Equal(t, stack.l, locker)
}

func TestSetup(t *testing.T) {
	t.Parallel()

	okLocker := func(t *testing.T) *mocks.MockLocker {
		ctrl := gomock.NewController(t)
		locker := mocks.NewMockLocker(ctrl)
		gomock.InOrder(
			locker.EXPECT().Lock().Return(nil),
			locker.EXPECT().Locked().Return(true),
			locker.EXPECT().Unlock().Return(nil),
		)
		return locker
	}

	mockTemplate := `{{/* docker-compose_base.tmpl */}}
{{ define "docker-compose" }}
services:
  service1:
    container_name: service1
    image: ${IMAGE}
    ports:
      - ${PORT}:9090
    networks:
      - sedge
networks:
  sedge:
    name: sedge-network
    external: true
{{ end }}
`

	tests := []struct {
		name    string
		env     map[string]string
		mocker  func(*testing.T) *mocks.MockLocker
		wantErr bool
	}{
		{
			name: "success",
			env: map[string]string{
				"IMAGE": "myimage:latest",
				"PORT":  "8088",
			},
			mocker:  okLocker,
			wantErr: false,
		},
		{
			name:    "empty .env",
			env:     map[string]string{},
			mocker:  okLocker,
			wantErr: false,
		},
		{
			name: "unlock error",
			env: map[string]string{
				"ERROR": "error",
			},
			mocker: func(t *testing.T) *mocks.MockLocker {
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)
				gomock.InOrder(
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(false),
				)
				return locker
			},
			wantErr: true,
		},
		{
			name: "lock error",
			env: map[string]string{
				"ERROR": "error",
			},
			mocker: func(t *testing.T) *mocks.MockLocker {
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)
				locker.EXPECT().Lock().Return(errors.New("lock error"))
				return locker
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			afs := afero.NewMemMapFs()
			mockFs := &afero.MemMapFs{}
			afero.WriteFile(mockFs, "services/docker-compose_base.tmpl", []byte(mockTemplate), 0o644)

			stack := &MonitoringStack{
				path: "/",
				l:    tt.mocker(t),
				fs:   afs,
			}

			err := stack.Setup(tt.env, afero.NewIOFS(mockFs))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check .env file
			exists, err := afero.Exists(afs, "/.env")
			assert.NoError(t, err)
			assert.True(t, exists)

			env, err := afero.ReadFile(afs, "/.env")
			require.NoError(t, err)
			gotEnv := make(map[string]string)
			for _, line := range strings.Split(string(env), "\n") {
				parts := strings.Split(line, "=")
				if len(parts) == 2 {
					gotEnv[parts[0]] = parts[1]
				}
			}
			assert.EqualValues(t, tt.env, gotEnv)

			// Check docker-compose.yml
			exists, err = afero.Exists(afs, "/docker-compose.yml")
			assert.NoError(t, err)
			assert.True(t, exists)

			gotCmp, err := afero.ReadFile(afs, "/docker-compose.yml")
			require.NoError(t, err)

			// Parse the template with the test data
			tmpl, err := template.New("test").Parse(mockTemplate)
			require.NoError(t, err)
			var expectedBuf bytes.Buffer
			err = tmpl.ExecuteTemplate(&expectedBuf, "docker-compose", nil)
			require.NoError(t, err)

			assert.Equal(t, expectedBuf.String(), string(gotCmp))
		})
	}
}

func TestCreateDir(t *testing.T) {
	t.Parallel()

	// Create an in-memory filesystem
	afs := afero.NewMemMapFs()

	// Create a mock locker
	ctrl := gomock.NewController(t)
	locker := mocks.NewMockLocker(ctrl)

	// Expect the lock to be acquired
	gomock.InOrder(
		locker.EXPECT().Lock().Return(nil),
		locker.EXPECT().Locked().Return(true),
		locker.EXPECT().Unlock().Return(nil),
	)
	gomock.InOrder(
		locker.EXPECT().Lock().Return(nil),
		locker.EXPECT().Locked().Return(true),
		locker.EXPECT().Unlock().Return(nil),
	)

	// Create a new MonitoringStack with the in-memory filesystem
	stack := &MonitoringStack{
		path: "/",
		l:    locker,
		fs:   afs,
	}

	// Create a wait group to synchronize the goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Create a channel to communicate the results of the directory creation attempts
	results := make(chan string, 2)

	// Goroutine 1
	go func() {
		t.Helper()
		defer wg.Done()
		err := stack.CreateDir("testdir1")
		assert.NoError(t, err)
		results <- "goroutine 1 created directory"
	}()

	// Give goroutine 1 a head start to ensure it creates the directory first
	time.Sleep(1 * time.Second)

	// Goroutine 2
	go func() {
		t.Helper()
		defer wg.Done()
		err := stack.CreateDir("testdir2")
		assert.NoError(t, err)
		results <- "goroutine 2 created directory"
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	// Check the results
	require.Equal(t, "goroutine 1 created directory", <-results)
	require.Equal(t, "goroutine 2 created directory", <-results)

	// Check that the directories were created
	exists, err := afero.DirExists(afs, "/testdir1")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = afero.DirExists(afs, "/testdir2")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCreateDirUnlockError(t *testing.T) {
	// Trying to unlock but the lock is not locked

	// Create an in-memory filesystem
	afs := afero.NewMemMapFs()

	// Create a mock locker
	ctrl := gomock.NewController(t)
	locker := mocks.NewMockLocker(ctrl)

	// Expect the lock to be acquired
	gomock.InOrder(
		locker.EXPECT().Lock().Return(nil),
		locker.EXPECT().Locked().Return(false),
	)

	// Create a new MonitoringStack with the in-memory filesystem
	stack := &MonitoringStack{
		path: "/",
		l:    locker,
		fs:   afs,
	}

	err := stack.CreateDir("testdir1")
	assert.Error(t, err)
	assert.EqualError(t, err, "monitoring stack is not locked")

	// Check that the directory was created
	exists, err := afero.DirExists(afs, "/testdir1")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCreateDirLockError(t *testing.T) {
	// Trying to unlock but the lock is not locked

	// Create an in-memory filesystem
	afs := afero.NewMemMapFs()

	// Create a mock locker
	ctrl := gomock.NewController(t)
	locker := mocks.NewMockLocker(ctrl)

	// Expect the lock to be acquired
	locker.EXPECT().Lock().Return(errors.New("lock error"))

	// Create a new MonitoringStack with the in-memory filesystem
	stack := &MonitoringStack{
		path: "/",
		l:    locker,
		fs:   afs,
	}

	err := stack.CreateDir("testdir1")
	assert.Error(t, err)
	assert.EqualError(t, err, "lock error")
}

func TestCreate(t *testing.T) {
	t.Parallel()

	okLocker := func(t *testing.T) *mocks.MockLocker {
		// Create a mock locker
		ctrl := gomock.NewController(t)
		locker := mocks.NewMockLocker(ctrl)

		// Expect the lock to be acquired
		gomock.InOrder(
			locker.EXPECT().Lock().Return(nil),
			locker.EXPECT().Locked().Return(true),
			locker.EXPECT().Unlock().Return(nil),
		)
		return locker
	}

	tests := []struct {
		name    string
		path    string
		mocker  func(*testing.T) *mocks.MockLocker
		wantErr bool
	}{
		{
			name:    "ok",
			path:    "/testfile",
			mocker:  okLocker,
			wantErr: false,
		},
		{
			name:    "already exists",
			path:    "/existingfile",
			mocker:  okLocker,
			wantErr: false,
		},
		{
			name: "unlock error",
			path: "/testfile",
			mocker: func(t *testing.T) *mocks.MockLocker {
				// Create a mock locker
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)

				// Expect the lock to be acquired
				gomock.InOrder(
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(false),
				)
				return locker
			},
			wantErr: true,
		},
		{
			name: "lock error",
			path: "/testfile",
			mocker: func(t *testing.T) *mocks.MockLocker {
				// Create a mock locker
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)

				// Expect the lock to be acquired
				locker.EXPECT().Lock().Return(errors.New("lock error"))
				return locker
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an in-memory filesystem
			afs := afero.NewMemMapFs()

			// Create an existing file
			_, err := afs.Create("/existingfile")
			require.NoError(t, err)

			// Write some content to the file
			err = afero.WriteFile(afs, "/existingfile", []byte("existing content"), 0o644)
			require.NoError(t, err)

			// Create a new MonitoringStack with the in-memory filesystem
			stack := &MonitoringStack{
				path: "/",
				l:    tt.mocker(t),
				fs:   afs,
			}

			_, err = stack.Create(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Check that the file was created
				exists, err := afero.Exists(afs, tt.path)
				assert.NoError(t, err)
				assert.True(t, exists)

				// Check that the file has the expected content
				content, err := afero.ReadFile(afs, "/existingfile")
				assert.NoError(t, err)
				if tt.path == "/existingfile" {
					assert.Len(t, content, 0)
				} else {
					assert.Equal(t, "existing content", string(content))
				}
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	t.Parallel()

	okLocker := func(t *testing.T) *mocks.MockLocker {
		// Create a mock locker
		ctrl := gomock.NewController(t)
		locker := mocks.NewMockLocker(ctrl)

		// Expect the lock to be acquired
		gomock.InOrder(
			locker.EXPECT().Lock().Return(nil),
			locker.EXPECT().Locked().Return(true),
			locker.EXPECT().Unlock().Return(nil),
		)
		return locker
	}

	tests := []struct {
		name    string
		path    string
		mocker  func(*testing.T) *mocks.MockLocker
		want    []byte
		wantErr bool
	}{
		{
			name:    "ok",
			path:    "/testfile",
			mocker:  okLocker,
			want:    []byte("test content"),
			wantErr: false,
		},
		{
			name:    "empty file",
			path:    "/emptyfile",
			mocker:  okLocker,
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "not found",
			path:    "/notfound",
			mocker:  okLocker,
			wantErr: true,
		},
		{
			name: "unlock error",
			path: "/testfile",
			mocker: func(t *testing.T) *mocks.MockLocker {
				// Create a mock locker
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)

				// Expect the lock to be acquired
				gomock.InOrder(
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(false),
				)
				return locker
			},
			wantErr: true,
		},
		{
			name: "lock error",
			path: "/testfile",
			mocker: func(t *testing.T) *mocks.MockLocker {
				// Create a mock locker
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)

				// Expect the lock to be acquired
				locker.EXPECT().Lock().Return(errors.New("lock error"))
				return locker
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an in-memory filesystem
			afs := afero.NewMemMapFs()

			// Create files
			_, err := afs.Create("/testfile")
			require.NoError(t, err)
			err = afero.WriteFile(afs, "/testfile", []byte("test content"), 0o644)
			require.NoError(t, err)

			_, err = afs.Create("/emptyfile")
			require.NoError(t, err)

			// Create a new MonitoringStack with the in-memory filesystem
			stack := &MonitoringStack{
				path: "/",
				l:    tt.mocker(t),
				fs:   afs,
			}

			got, err := stack.ReadFile(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	t.Parallel()

	okLocker := func(t *testing.T) *mocks.MockLocker {
		// Create a mock locker
		ctrl := gomock.NewController(t)
		locker := mocks.NewMockLocker(ctrl)

		// Expect the lock to be acquired
		gomock.InOrder(
			locker.EXPECT().Lock().Return(nil),
			locker.EXPECT().Locked().Return(true),
			locker.EXPECT().Unlock().Return(nil),
		)
		return locker
	}

	tests := []struct {
		name    string
		path    string
		content []byte
		mocker  func(*testing.T) *mocks.MockLocker
		wantErr bool
	}{
		{
			name:    "empty content",
			path:    "/existingfile",
			content: []byte{},
			mocker:  okLocker,
			wantErr: false,
		},
		{
			name:    "already exists",
			path:    "/existingfile",
			content: []byte("overwritten?"),
			mocker:  okLocker,
			wantErr: false,
		},
		{
			name:    "not found",
			path:    "/notfound",
			content: []byte("test content"),
			mocker:  okLocker,
			wantErr: false,
		},
		{
			name:    "unlock error",
			path:    "/existingfile",
			content: []byte("test content"),
			mocker: func(t *testing.T) *mocks.MockLocker {
				// Create a mock locker
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)

				// Expect the lock to be acquired
				locker.EXPECT().Locked().Return(false).After(locker.EXPECT().Lock().Return(nil))
				return locker
			},
			wantErr: true,
		},
		{
			name:    "lock error",
			path:    "/existingfile",
			content: []byte("test content"),
			mocker: func(t *testing.T) *mocks.MockLocker {
				// Create a mock locker
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)

				// Expect the lock to be acquired
				locker.EXPECT().Lock().Return(errors.New("lock error"))
				return locker
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an in-memory filesystem
			afs := afero.NewMemMapFs()

			// Create an existing file
			_, err := afs.Create("/existingfile")
			require.NoError(t, err)
			err = afero.WriteFile(afs, "/existingfile", []byte("test content"), 0o644)
			require.NoError(t, err)

			// Create a new MonitoringStack with the in-memory filesystem
			stack := &MonitoringStack{
				path: "/",
				l:    tt.mocker(t),
				fs:   afs,
			}

			err = stack.WriteFile(tt.path, tt.content)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				content, err := afero.ReadFile(afs, tt.path)
				require.NoError(t, err)
				assert.Equal(t, tt.content, content)
			}
		})
	}
}

func TestInstalled(t *testing.T) {
	t.Parallel()

	okLocker := func(t *testing.T) *mocks.MockLocker {
		// Create a mock locker
		ctrl := gomock.NewController(t)
		locker := mocks.NewMockLocker(ctrl)

		// Expect the lock to be acquired
		gomock.InOrder(
			locker.EXPECT().Lock().Return(nil),
			locker.EXPECT().Locked().Return(true),
			locker.EXPECT().Unlock().Return(nil),
		)
		return locker
	}

	type testcase struct {
		name    string
		fs      afero.Fs
		mocker  func(*testing.T) *mocks.MockLocker
		want    bool
		wantErr bool
	}

	tests := []testcase{
		func() testcase {
			afs := afero.NewMemMapFs()
			_, err := afs.Create("/.env")
			require.NoError(t, err)
			_, err = afs.Create("/docker-compose.yml")
			require.NoError(t, err)

			return testcase{
				name:   "ok",
				fs:     afs,
				mocker: okLocker,
				want:   true,
			}
		}(),
		func() testcase {
			afs := afero.NewMemMapFs()
			_, err := afs.Create("/docker-compose.yml")
			require.NoError(t, err)

			return testcase{
				name:   "no .env",
				fs:     afs,
				mocker: okLocker,
				want:   false,
			}
		}(),
		func() testcase {
			afs := afero.NewMemMapFs()
			_, err := afs.Create("/.env")
			require.NoError(t, err)

			return testcase{
				name:   "no docker-compose.yml",
				fs:     afs,
				mocker: okLocker,
				want:   false,
			}
		}(),
		func() testcase {
			afs := afero.NewMemMapFs()
			_, err := afs.Create("/.env")
			require.NoError(t, err)
			_, err = afs.Create("/docker-compose.yml")
			require.NoError(t, err)

			return testcase{
				name: "unlock error",
				fs:   afs,
				mocker: func(t *testing.T) *mocks.MockLocker {
					// Create a mock locker
					ctrl := gomock.NewController(t)
					locker := mocks.NewMockLocker(ctrl)

					// Expect the lock to be acquired
					locker.EXPECT().Locked().Return(false).After(locker.EXPECT().Lock().Return(nil))
					return locker
				},
				wantErr: true,
			}
		}(),
		func() testcase {
			afs := afero.NewMemMapFs()

			return testcase{
				name: "lock error, stack initialized",
				fs:   afs,
				mocker: func(t *testing.T) *mocks.MockLocker {
					// Create a mock locker
					ctrl := gomock.NewController(t)
					locker := mocks.NewMockLocker(ctrl)

					// Expect the lock to be acquired
					locker.EXPECT().Lock().Return(errors.New("lock error"))
					return locker
				},
				wantErr: true,
			}
		}(),
		func() testcase {
			afs := afero.NewMemMapFs()

			return testcase{
				name: "lock error, stack not initialized",
				fs:   afs,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new MonitoringStack with the in-memory filesystem
			var stack *MonitoringStack

			if tt.mocker != nil {
				stack = &MonitoringStack{
					path: "/",
					l:    tt.mocker(t),
					fs:   tt.fs,
				}
			} else {
				stack = &MonitoringStack{
					path: "/",
					fs:   tt.fs,
				}
			}

			got, err := stack.Installed()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPath(t *testing.T) {
	// dumbest test of all times lol
	t.Parallel()

	stack := &MonitoringStack{
		path: "/",
	}

	assert.Equal(t, "/", stack.Path())
}

func TestCleanup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		mocker     func(*testing.T) (*mocks.MockLocker, afero.Fs)
		force      bool
		notInstall bool
		wantErr    bool
	}{
		{
			name: "ok, force false",
			mocker: func(t *testing.T) (*mocks.MockLocker, afero.Fs) {
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join("monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(nil),
				)
				locker.EXPECT().Lock().Return(nil)
				return locker, afero.NewMemMapFs()
			},
			force: false,
		},
		{
			name: "ok, force true",
			mocker: func(t *testing.T) (*mocks.MockLocker, afero.Fs) {
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join("monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(nil),
				)
				return locker, afero.NewMemMapFs()
			},
			force: true,
		},
		{
			name: "not installed",
			mocker: func(t *testing.T) (*mocks.MockLocker, afero.Fs) {
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)
				locker.EXPECT().Lock().Return(nil)
				return locker, afero.NewMemMapFs()
			},
			notInstall: true,
		},
		{
			name: "lock error",
			mocker: func(t *testing.T) (*mocks.MockLocker, afero.Fs) {
				ctrl := gomock.NewController(t)
				locker := mocks.NewMockLocker(ctrl)
				gomock.InOrder(
					locker.EXPECT().New(utils.PathMatcher{Expected: filepath.Join("monitoring", ".lock")}).Return(locker),
					locker.EXPECT().Lock().Return(nil),
					locker.EXPECT().Locked().Return(true),
					locker.EXPECT().Unlock().Return(nil),
				)
				locker.EXPECT().Lock().Return(errors.New("lock error"))
				return locker, afero.NewMemMapFs()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locker, fs := tt.mocker(t)
			stack := &MonitoringStack{
				path: "/monitoring",
				l:    locker,
				fs:   fs,
			}

			// Install the stack
			var err error
			if !tt.notInstall {
				err = stack.Init()
				require.NoError(t, err)

				// Write the template file to the in-memory filesystem
				afero.WriteFile(fs, "services/docker-compose_base.tmpl", []byte(`
services:
  {{.ServiceName}}:
    image: {{.Image}}
    ports:
      - "{{.Port}}:{{.Port}}"
`), 0o644)

				err = stack.Setup(map[string]string{"ServiceName": "myservice", "Image": "myimage:latest", "Port": "8080"}, templates.Services)
				require.NoError(t, err)
			}

			err = stack.Cleanup(tt.force)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Check that monitoring stack has been removed
				exists, err := afero.DirExists(fs, "/monitoring")
				assert.NoError(t, err)
				assert.False(t, exists)
			}
		})
	}
}
