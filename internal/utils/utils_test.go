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
package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/NethermindEth/sedge/configs"
)

func TestSkipLines(t *testing.T) {
	inputs := []struct {
		content string
		symbol  string
		result  string
	}{
		{},
		{"bbbb", "a", "bbbb"},
		{"aaaa", "a", ""},
		{"bbbb\naaaa", "a", "bbbb"},
		{"aaaa\nbbbb", "a", "bbbb"},
	}

	for _, input := range inputs {
		descr := fmt.Sprintf("SkipLines(%s, %s)", input.content, input.symbol)
		got := SkipLines(input.content, input.symbol)

		if got != input.result {
			t.Errorf("%s expected %s but got %s", descr, input.result, got)
		}
	}
}

func TestContains(t *testing.T) {
	inputs := [...]struct {
		list []string
		str  string
		want bool
	}{
		{},
		{[]string{"a", "A", "b", "B"}, "b", true},
		{[]string{"a", "A", "b", "B"}, "c", false},
		{[]string{}, "a", false},
	}

	for _, input := range inputs {
		descr := fmt.Sprintf("Contains(%s)", input.list)
		if got := Contains(input.list, input.str); got != input.want {
			t.Errorf("%s expected %t but got %t", descr, input.want, got)
		}
	}
}

func TestContainsOnly(t *testing.T) {
	inputs := [...]struct {
		list   []string
		target []string
		want   bool
	}{
		{[]string{}, []string{}, true},
		{[]string{}, []string{"a"}, true},
		{[]string{"b"}, []string{"a", "A", "b", "B"}, true},
		{[]string{"c"}, []string{"a", "A", "b", "B"}, false},
		{[]string{"a"}, []string{}, false},
		{[]string{"execution", "validator"}, []string{"execution", "consensus"}, false},
		{[]string{"execution", "consensus", "validator"}, []string{"execution", "consensus"}, false},
		{[]string{"execution", "validator"}, []string{"execution", "consensus", "validator"}, true},
		{[]string{"execution", "consensus", "validator"}, []string{"execution", "consensus", "validator"}, true},
	}

	for _, input := range inputs {
		descr := fmt.Sprintf("Contains(%s)", input.list)
		if got := ContainsOnly(input.list, input.target); got != input.want {
			t.Errorf("%s expected %t but got %t", descr, input.want, got)
		}
	}
}

func TestIsAddress(t *testing.T) {
	tcs := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"2131", false},
		{"dasd31gsd1231", false},
		{"0x2312313aaef2312312", false},
		{"0x5c00ABEf07604C59Ac72E859E5F93D5abZXCVF83", false},
		{"5c00ABEf07604C59Ac72E859E5F93D5ab8546F83", false},
		{"0x5c00ABEf07604C59Ac72E859E5F93D5ab8546F83", true},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("IsAddress(%s)", tc.input), func(t *testing.T) {
			if got := IsAddress(tc.input); got != tc.want {
				t.Errorf("got != want. Expected %v, got %v", tc.want, tc.input)
			}
		})
	}
}

func TestPortAvailable(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
		}))
	defer server.Close()
	split := strings.Split(server.URL, ":")
	host, strPort := split[1][2:], split[2]
	port64, err := strconv.ParseUint(strPort, 10, 16)
	if err != nil {
		t.Fatalf("cannot convert http server port: %v", err)
	}
	port := uint16(port64)

	tcs := []struct {
		name string
		host string
		port uint16
		want bool
	}{
		{
			"Test case 1, good host and unavailable port",
			host, port,
			false,
		},
		{
			"Test case 2, bad host and port",
			"b@dh0$t", port,
			true,
		},
		{
			"Test case 3, good host and bad port",
			host, 9999,
			true,
		},
		{
			"Test case 4, good host and available port",
			"localhost", 9999,
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if got := portAvailable(tc.host, tc.port); tc.want != got {
				t.Errorf("portAvailable(%s, %d) failed; expected: %v, got: %v", tc.host, tc.port, tc.want, got)
			}
		})
	}
}

func TestAssignPorts(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
		}))
	defer server.Close()
	split := strings.Split(server.URL, ":")
	host, strPort := split[1][2:], split[2]
	port64, err := strconv.ParseUint(strPort, 10, 16)
	if err != nil {
		t.Fatalf("cannot convert http server port: %v", err)
	}
	port := uint16(port64)

	tcs := []struct {
		name     string
		host     string
		defaults map[string]uint16
		want     map[string]uint16
		isErr    bool
	}{
		{
			"Test case 1, good host and defaults",
			host,
			map[string]uint16{"EL": 8545, "CL": port},
			map[string]uint16{"EL": 8545, "CL": port + 1},
			false,
		},
		{
			"Test case 2, good host and bad defaults",
			host,
			map[string]uint16{"EL": 8545, "CL": 0},
			map[string]uint16{},
			true,
		},
		{
			"Test case 3, good host and bad defaults",
			host,
			map[string]uint16{"CL": 0, "EL": 8545},
			map[string]uint16{},
			true,
		},
		{
			"Test case 4, bad host and good defaults",
			"b@dh0$t",
			map[string]uint16{"CL": 9000, "EL": 8545},
			map[string]uint16{"CL": 9000, "EL": 8545},
			false,
		},
		{
			"Test case 5, good host and successive increments",
			host,
			map[string]uint16{"CL": port, "EL": port + 1},
			map[string]uint16{"CL": port + 1, "EL": port + 2},
			false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := AssignPorts(tc.host, tc.defaults)

			descr := fmt.Sprintf("AssingPorts(%s, %+v)", tc.host, tc.defaults)
			if cerr := CheckErr(descr, tc.isErr, err); cerr != nil {
				t.Error(cerr)
			}

			if err == nil {
				for k := range tc.want {
					if tc.want[k] != got[k] {
						t.Errorf("A mismatch in the result has been found. Expected (key: %s, value: %d); got (key: %s, value %d). Call: %s. Expected object: %+v, Got: %+v", k, tc.want[k], k, got[k], descr, tc.want, got)
					}
				}
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tcs := []struct {
		name   string
		in     []string
		want   []string
		filter func(string) bool
	}{
		{
			"Test case 1, no filter",
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
			func(s string) bool {
				return true
			},
		},
		{
			"Test case 2, filter",
			[]string{"a", "b", "c"},
			[]string{"a", "c"},
			func(s string) bool {
				return s != "b"
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := Filter(tc.in, tc.filter)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Filter(%+v) failed; expected: %+v, got: %+v", tc.in, tc.want, got)
			}
		})
	}
}

func TestUriValidator(t *testing.T) {
	tcs := []struct {
		name string
		in   []string
		want bool
	}{
		{
			"good uris",
			[]string{
				"http://localhost:8545",
				"https://localhost:8545",
				"https://0xac6e77dfe25ecd6110b8e780608cce0dab71fdd5ebea22a16c0205200f2f8e2e3ad3b71d3499c54ad14d6c21b41a37ae@boost-relay.flashbots.net",
				"https://boost-relay.flashbots.net",
				"https://localhost:8545/api/v1/eth1",
				"https://localhost:8545/api/v1/eth1,.{}",
				"https://nethermind/api/v1/eth1",
				"http://sedge",
				"https://192.168.0.1",
			},
			true,
		},
		{
			"good uri",
			[]string{"http://banana/api/monkey/[spliat]"},
			true,
		},
		{
			"bad uri",
			[]string{"https://192.168.0.1:8545", "localhost:8545", "https:/boost-relay.flashbots.net"},
			false,
		},
		{
			"bad uri",
			[]string{"localhost:8545"},
			false,
		},
		{
			"bad uri",
			[]string{"localhost/545"},
			false,
		},
		{
			"bad uri",
			[]string{"https:/boost-relay.flashbots.net"},
			false,
		},
		{
			"bad uri",
			[]string{"htp://localhost:8545"},
			false,
		},
		{
			"bad uri",
			[]string{"./8088"},
			false,
		},
		{
			"bad uri",
			[]string{"44.33.55.66:8088"},
			false,
		},
		{
			"bad uri",
			[]string{""},
			false,
		},
	}

	for i, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			t.Logf("Test case %d: %s", i+1, tc.name)

			uri, got := UriValidator(tc.in)

			if got != tc.want {
				t.Errorf("UriValidator(%s) failed; expected: %v, got: %v. Bad uri: %s", tc.in, tc.want, got, uri)
			}
			if !got && !Contains(tc.in, uri) {
				t.Errorf("UriValidator(%s) returned a different uri that provided; expected: %s, got: %s", tc.in, tc.in, uri)
			}
		})
	}
}

func TestENodesValidator(t *testing.T) {
	tcs := []struct {
		name string
		in   []string
		want error
	}{
		{
			"good enodes",
			[]string{
				"enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303",
				"enode://22a8232c3abc76a16ae9d6c3b164f98775fe226f0917b0ca871128a74a8e9630b458460865bab457221f1d448dd9791d24c4e5d88786180ac185df813a68d4de@3.209.45.79:30303",
				"enode://715171f50508aba88aecd1250af392a45a330af91d7b90701c436b618c86aaa1589c9184561907bebbb56439b8f8787bc01f49a7c77276c58c1b09822d75e8e8@@52.231.165.108:30303",
				"enode://5d6d7cd20d6da4bb83a1d28cadb5d409b64edf314c0335df658c1a54e32c7c4a7ab7823d57c39b6a757556e68ff1df17c748b698544a55cb488b52479a92b60f@the-second-most-cool-enode:666",
			},
			nil,
		},
		{
			"bad enode",
			[]string{
				"enode://0x0f6b",
			},
			fmt.Errorf(configs.InvalidEnodeError, "enode://0x0f6b"),
		},
		{
			"bad enode",
			[]string{
				"enode://2b252ab6a1d0f971d9722cb839a42cb81db019ba44c08754628ab4a823487071b5695317c8ccd085219c3a03af063495b2f1da8d18218da2d6a82981b45e6ffc@the-most-cool-enode",
			},
			fmt.Errorf(configs.InvalidEnodeError, "enode://2b252ab6a1d0f971d9722cb839a42cb81db019ba44c08754628ab4a823487071b5695317c8ccd085219c3a03af063495b2f1da8d18218da2d6a82981b45e6ffc@the-most-cool-enode"),
		},
		{
			"bad enode",
			[]string{
				"enode://4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166",
			},
			fmt.Errorf(configs.InvalidEnodeError, "enode://4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166"),
		},
		{
			"bad enode",
			[]string{
				"enode:4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166",
				"enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303",
			},
			fmt.Errorf(configs.InvalidEnodeError, "enode:4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166"),
		},
	}

	for i, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			t.Logf("Test case %d: %s", i+1, tc.name)

			got := ENodesValidator(tc.in)

			if err := CheckErr("ENodesValidator", tc.want != nil, got); err != nil {
				t.Error(err)
			}

			if got != nil && tc.want != nil && got.Error() != tc.want.Error() {
				t.Errorf("ENodesValidator(%s) returned a different error; expected: %s, got: %s", tc.in, tc.want, got)
			}
		})
	}
}

func TestENRValidator(t *testing.T) {
	tcs := []struct {
		name string
		in   []string
		want error
	}{
		{
			"good enrs",
			[]string{
				"enr:-LK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo", "enr:-KG4QCIzJZTY_fs_2vqWEatJL9RrtnPwDCv-jRBuO5FQ2qBrfJubWOWazri6s9HsyZdu-fRUfEzkebhf1nvO42_FVzwDhGV0aDKQed8EKAAAECD__________4JpZIJ2NIJpcISHtbYziXNlY3AyNTZrMaED4m9AqVs6F32rSCGsjtYcsyfQE2K8nDiGmocUY_iq-TSDdGNwgiMog3VkcIIjKA", "enr:-Ku4QFmUkNp0g9bsLX2PfVeIyT-9WO-PZlrqZBNtEyofOOfLMScDjaTzGxIb1Ns9Wo5Pm_8nlq-SZwcQfTH2cgO-s88Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpDkvpOTAAAQIP__________gmlkgnY0gmlwhBLf22SJc2VjcDI1NmsxoQLV_jMOIxKbjHFKgrkFvwDvpexo6Nd58TK5k7ss4Vt0IoN1ZHCCG1g",
				"enr:-LK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNH",
			},
			nil,
		},
		{
			"bad enr",
			[]string{
				"enr:LK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo",
			},
			fmt.Errorf(configs.InvalidEnrError, "enr:LK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo"),
		},
		{
			"bad enr",
			[]string{
				"enr-LK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo",
			},
			fmt.Errorf(configs.InvalidEnrError, "enr-LK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo"),
		},
		{
			"bad enr",
			[]string{
				"en:-LK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo",
			},
			fmt.Errorf(configs.InvalidEnrError, "en:-LK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo"),
		},
		{
			"bad enr",
			[]string{
				"enrLK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo",
				"enr:-LK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo",
			},
			fmt.Errorf(configs.InvalidEnrError, "enrLK4QH1xnjotgXwg25IDPjrqRGFnH1ScgNHA3dv1Z8xHCp4uP3N3Jjl_aYv_WIxQRdwZvSukzbwspXZ7JjpldyeVDzMCh2F0dG5ldHOIAAAAAAAAAACEZXRoMpB53wQoAAAQIP__________gmlkgnY0gmlwhIe1te-Jc2VjcDI1NmsxoQOkcGXqbCJYbcClZ3z5f6NWhX_1YPFRYRRWQpJjwSHpVIN0Y3CCIyiDdWRwgiMo"),
		},
	}

	for i, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			t.Logf("Test case %d: %s", i+1, tc.name)

			got := ENRValidator(tc.in)

			if err := CheckErr("ENRValidator", tc.want != nil, got); err != nil {
				t.Error(err)
			}

			if got != nil && tc.want != nil && got.Error() != tc.want.Error() {
				t.Errorf("ENRValidator(%s) returned a different error; expected: %s, got: %s", tc.in, tc.want, got)
			}
		})
	}
}
