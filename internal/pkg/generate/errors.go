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
package generate

import "errors"

// ErrTemplateNotFound is returned when the template is not found
var ErrTemplateNotFound = errors.New("template not found")

// ErrEmptyData is returned when the data is nil
var ErrEmptyData = errors.New("data is nil")

// ErrUnableToGetClientsInfo is returned when the client information cannot be retrieved
var ErrUnableToGetClientsInfo = errors.New("unable to get clients information")

// ErrConsensusClientNotValid is returned when the consensus client is not valid
var ErrConsensusClientNotValid = errors.New("invalid consensus client")

// ErrExecutionClientNotValid is returned when the execution client is not valid
var ErrExecutionClientNotValid = errors.New("invalid execution client")

// ErrValidatorClientNotValid is returned when the validator client is not valid
var ErrValidatorClientNotValid = errors.New("invalid validator client")

// ErrDistributedValidatorClientNotValid is returned when the distributed validator client is not valid
var ErrDistributedValidatorClientNotValid = errors.New("invalid distributed validator client")
