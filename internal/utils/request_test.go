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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Method != "GET" {
				t.Errorf("Unexpected HTTP method, expected GET, got %s", req.Method)
			}

			query := req.URL.Query().Get("test")
			if query == "ERROR" {
				rw.WriteHeader(http.StatusBadRequest)
				return
			} else if query == "OK" {
				rw.WriteHeader(http.StatusOK)
				rw.Write([]byte("OK"))
			}
		}))
	defer server.Close()

	tcs := []struct {
		name          string
		url           string
		want          string
		retryDuration time.Duration
		isError       bool
	}{
		{
			"Good request",
			server.URL + "/?test=OK",
			"OK",
			time.Second,
			false,
		},
		{
			"Bad request",
			server.URL + "/?test=ERROR",
			"",
			time.Second,
			false,
		},
		{
			"No response",
			"http://127.0.0.1:8088" + "/",
			"",
			time.Second,
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := GetRequest(tc.url, tc.retryDuration)
			descr := fmt.Sprintf("GetRequest(%s)", tc.url)
			if err = CheckErr(descr, tc.isError, err); err != nil {
				t.Error(err)
			}

			if resp != nil {
				defer resp.Body.Close()
				contents, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Reading response body failed. Error: %v", err)
				}

				if tc.want != string(contents) {
					t.Errorf("DoRequest(%s) response body is %s (got) != %s (want)", tc.url, string(contents), tc.want)
				}
			}
		})
	}
}
