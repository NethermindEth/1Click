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
package services_test

import (
	"errors"
	"testing"

	"github.com/NethermindEth/sedge/internal/pkg/services"
	mock_client "github.com/NethermindEth/sedge/test/mock_docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestContainerId(t *testing.T) {
	ctrl := gomock.NewController(t)
	dockerClient := mock_client.NewMockAPIClient(ctrl)
	defer ctrl.Finish()

	wantId := "container-id"
	containerName := "container-name"

	dockerClient.EXPECT().
		ContainerList(gomock.Any(), types.ContainerListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("name", containerName)),
		}).
		Return([]types.Container{
			{
				ID: wantId,
			},
		}, nil).
		Times(1)
	serviceManager := services.NewServiceManager(dockerClient)
	id, err := serviceManager.ContainerId(containerName)
	assert.Nil(t, err)
	assert.Equal(t, wantId, id)
}

func TestContainerIdError(t *testing.T) {
	ctrl := gomock.NewController(t)
	dockerClient := mock_client.NewMockAPIClient(ctrl)
	defer ctrl.Finish()

	wantErr := errors.New("error")
	containerName := "container-name"

	dockerClient.EXPECT().
		ContainerList(gomock.Any(), types.ContainerListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("name", containerName)),
		}).
		Return(nil, wantErr).
		Times(1)
	serviceManager := services.NewServiceManager(dockerClient)
	id, err := serviceManager.ContainerId(containerName)
	assert.ErrorIs(t, err, wantErr)
	assert.Equal(t, "", id)
}

func TestContainerIdNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	dockerClient := mock_client.NewMockAPIClient(ctrl)
	defer ctrl.Finish()

	containerName := "container-name"

	dockerClient.EXPECT().
		ContainerList(gomock.Any(), types.ContainerListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("name", containerName)),
		}).
		Return(make([]types.Container, 0), nil).
		Times(1)
	serviceManager := services.NewServiceManager(dockerClient)
	id, err := serviceManager.ContainerId(containerName)
	assert.ErrorIs(t, err, services.ErrContainerNotFound)
	assert.Equal(t, "", id)
}

func TestContainerIdMultiple(t *testing.T) {
	ctrl := gomock.NewController(t)
	dockerClient := mock_client.NewMockAPIClient(ctrl)
	defer ctrl.Finish()

	containerName := "container-name"

	dockerClient.EXPECT().
		ContainerList(gomock.Any(), types.ContainerListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("name", containerName)),
		}).
		Return(make([]types.Container, 2), nil).
		Times(1)
	serviceManager := services.NewServiceManager(dockerClient)
	id, err := serviceManager.ContainerId(containerName)
	assert.ErrorIs(t, err, services.ErrMultipleContainers)
	assert.Equal(t, "", id)
}
