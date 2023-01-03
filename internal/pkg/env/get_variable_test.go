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
package env

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/NethermindEth/sedge/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetCCBootnodes(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name         string
		testdataCase string
		want         []string
		isErr        bool
	}{
		{
			name:         "Test case 1, no bootnodes",
			testdataCase: "case_1",
			want:         nil,
			isErr:        false,
		},
		{
			name:         "Test case 2, bootnodes",
			testdataCase: "case_2",
			want:         []string{"enr:-LK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo", "enr:-KG4QCIzJZTY_fs_2vqWEatJL9RrtnPwDCv-jRBuO5FQ2qBrfJubWOWazri6s9HsyZdu-fRUfEzkebhf1nvO42_FVzwDhGV0aDKQed8EKAAAECD__________4JpZIJ2NIJpcISHtbYziXNlY3AyNTZrMaED4m9AqVs6F32rSCGsjtYcsyfQE2K8nDiGmocUY_iq-TSDdGNwgiMog3VkcIIjKA", "enr:-Ku4QFmUkNp0g9bsLX2PfVeIyT-9WO-PZlrqZBNtEyofOOfLMScDjaTzGxIb1Ns9Wo5Pm_8nlq-SZwcQfTH2cgO-s88Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpDkvpOTAAAQIP__________gmlkgnY0gmlwhBLf22SJc2VjcDI1NmsxoQLV_jMOIxKbjHFKgrkFvwDvpexo6Nd58TK5k7ss4Vt0IoN1ZHCCG1g", "enr:-LK4QLINdtobGquK7jukLDAKmsrH2ZuHM4k0TklY5jDTD4ZgfxR9weZmo5Jwu81hlKu3qPAvk24xHGBDjYs4o8f1gZ0Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhDRN_P6Jc2VjcDI1NmsxoQJuNujTgsJUHUgVZML3pzrtgNtYg7rQ4K1tkWERgl0DdoN0Y3CCIyiDdWRwgiMo", "enr:-LK4QMzPq4Q7w5R-rnGQDcI8BYky6oPVBGQTbS1JJLVtNi_8PzBLV7Bdzsoame9nJK5bcJYpGHn4SkaDN2CM6tR5G_4Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhAN4yvyJc2VjcDI1NmsxoQKa8Qnp_P2clLIP6VqLKOp_INvEjLszalEnW0LoBZo4YYN0Y3CCI4yDdWRwgiOM", "enr:-LK4QLM_pPHa78R8xlcU_s40Y3XhFjlb3kPddW9lRlY67N5qeFE2Wo7RgzDgRs2KLCXODnacVHMFw1SfpsW3R474RZEBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhANBY-yJc2VjcDI1NmsxoQNsZkFXgKbTzuxF7uwxlGauTGJelE6HD269CcFlZ_R7A4N0Y3CCI4yDdWRwgiOM", "enr:-KK4QH0RsNJmIG0EX9LSnVxMvg-CAOr3ZFF92hunU63uE7wcYBjG1cFbUTvEa5G_4nDJkRhUq9q2ck9xY-VX1RtBsruBtIRldGgykIL0pysBABAg__________-CaWSCdjSCaXCEEnXQ0YlzZWNwMjU2azGhA1grTzOdMgBvjNrk-vqWtTZsYQIi0QawrhoZrsn5Hd56g3RjcIIjKIN1ZHCCIyg"},
			isErr:        false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			testCaseEnvFilePath := filepath.Join("testdata", tc.testdataCase, ".env")
			got, err := GetCCBootnodes(testCaseEnvFilePath)

			descr := fmt.Sprintf("GetClBootnodes(%s)", testCaseEnvFilePath)
			if err = utils.CheckErr(descr, tc.isErr, err); err != nil {
				t.Error(err)
			}

			assert.ElementsMatch(t, got, tc.want)
		})
	}
}

func TestGetECBootnodes(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name         string
		testdataCase string
		want         []string
		isErr        bool
	}{
		{
			name:         "Test case 1, no bootnodes",
			testdataCase: "case_1",
			want:         nil,
			isErr:        false,
		},
		{
			name:         "Test case 2, bootnodes",
			testdataCase: "case_2",
			want:         []string{"enode://011f758e6552d105183b1761c5e2dea0111bc20fd5f6422bc7f91e0fabbec9a6595caf6239b37feb773dddd3f87240d99d859431891e4a642cf2a0a9e6cbb98a@51.141.78.53:30303", "enode://176b9417f511d05b6b2cf3e34b756cf0a7096b3094572a8f6ef4cdcb9d1f9d00683bf0f83347eebdf3b81c3521c2332086d9592802230bf528eaf606a1d9677b@13.93.54.137:30303", "enode://46add44b9f13965f7b9875ac6b85f016f341012d84f975377573800a863526f4da19ae2c620ec73d11591fa9510e992ecc03ad0751f53cc02f7c7ed6d55c7291@94.237.54.114:30313", "enode://b5948a2d3e9d486c4d75bf32713221c2bd6cf86463302339299bd227dc2e276cd5a1c7ca4f43a0e9122fe9af884efed563bd2a1fd28661f3b5f5ad7bf1de5949@18.218.250.66:30303"},
			isErr:        false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			testCaseEnvFilePath := filepath.Join("testdata", tc.testdataCase, ".env")
			got, err := GetECBootnodes(testCaseEnvFilePath)

			descr := fmt.Sprintf("GetElBootnodes(%s)", testCaseEnvFilePath)
			if err = utils.CheckErr(descr, tc.isErr, err); err != nil {
				t.Error(err)
			}

			assert.ElementsMatch(t, got, tc.want)
		})
	}
}

func TestGetTTD(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name         string
		testdataCase string
		want         string
		isErr        bool
	}{
		{
			name:         "Test case 1, no TTD",
			testdataCase: "case_1",
			want:         "",
			isErr:        false,
		},
		{
			name:         "Test case 2, TTD",
			testdataCase: "case_2",
			want:         "8000000000000000000000001000000000000",
			isErr:        false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			testCaseEnvFilePath := filepath.Join("testdata", tc.testdataCase, ".env")
			got, err := GetTTD(testCaseEnvFilePath)

			descr := fmt.Sprintf("GetTTD(%s)", testCaseEnvFilePath)
			if err = utils.CheckErr(descr, tc.isErr, err); err != nil {
				t.Error(err)
			}

			if got != tc.want {
				t.Errorf("Expected %v, got %v. Function call: %s", tc.want, got, descr)
			}
		})
	}
}

func TestCheckVariableBase(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		regex   *regexp.Regexp
		network string
		want    bool
		isErr   bool
	}{
		{
			"Test case 1, gnosis, EL_NETWORK",
			ReSPLITTED,
			"gnosis",
			true,
			false,
		},
		{
			"Test case 2, invalid network, error",
			ReSPLITTED,
			"testnet",
			false,
			true,
		},
		{
			"Test case 3, mainnet, no EL_NETWORK",
			ReSPLITTED,
			"mainnet",
			false,
			false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CheckVariableBase(tc.regex, tc.network)

			descr := fmt.Sprintf("CheckVariableBase(re, %s) with regex %v", tc.network, tc.regex)
			if err = utils.CheckErr(descr, tc.isErr, err); err != nil {
				t.Error(err)
			}

			if tc.want != got {
				t.Errorf("Expected %v, got %v. Function call: %s", tc.want, got, descr)
			}
		})
	}
}
