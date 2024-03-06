// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repeatconfig_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/google/localtoast/scanapi"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
	"github.com/google/localtoast/scannerlib/repeatconfig"
)

type fakeFileReader struct {
	content          string
	loginDefsContent string
}

func (f *fakeFileReader) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	if path == "/etc/passwd" || path == "/proc/net/tcp" || path == "/proc/net/tcp6" {
		return io.NopCloser(bytes.NewReader([]byte(f.content))), nil
	}
	if path == "/etc/login.defs" {
		return io.NopCloser(bytes.NewReader([]byte(f.loginDefsContent))), nil
	}
	return nil, errors.New("file not found")
}

func (fakeFileReader) FilePermissions(ctx context.Context, path string) (*apb.PosixPermissions, error) {
	return nil, errors.New("Not implemented")
}

func (fakeFileReader) OpenDir(ctx context.Context, path string) (scanapi.DirReader, error) {
	return nil, errors.New("Not implemented")
}

func configHasError(config []*repeatconfig.RepeatConfig) bool {
	for _, c := range config {
		if c.Err != nil {
			return true
		}
	}
	return false
}

func TestCreateRepeatConfigsOnce(t *testing.T) {
	config := &ipb.RepeatConfig{Type: ipb.RepeatConfig_ONCE}
	got, err := repeatconfig.CreateRepeatConfigs(context.Background(), config, &fakeFileReader{})
	if err != nil {
		t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) returned an error: %v", config, err)
	}
	want := []*repeatconfig.RepeatConfig{&repeatconfig.RepeatConfig{}}
	if diff := cmp.Diff(want, got, cmp.AllowUnexported(repeatconfig.RepeatConfig{}, repeatconfig.TokenReplacement{})); diff != "" {
		t.Errorf("repeatconfig.CreateRepeatConfigs(ONCE) returned unexpected diff (-want +got):\n%s", diff)
	}
}

func TestCreateRepeatConfigsForEachUser(t *testing.T) {
	passwd := "user1:x:1337:1338::/home/user1:/bin/bash\n" +
		"nologinuser1:x:2337:2338::/dev/null:/bin/false\n" +
		"nologinuser2:x:3337:3338::/bin/nologin:/bin/false\n" +
		"user2:x:4337:4338::/home/user2:/bin/bash"

	testCases := []struct {
		desc       string
		configType ipb.RepeatConfig_RepeatType
		loginDefs  string
		want       []*repeatconfig.RepeatConfig
	}{
		{
			desc:       "Users without login",
			configType: ipb.RepeatConfig_FOR_EACH_USER,
			want: []*repeatconfig.RepeatConfig{
				{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "user1"},
						{TextToReplace: "$uid", ReplaceWith: "1337"},
						{TextToReplace: "$gid", ReplaceWith: "1338"},
						{TextToReplace: "$home", ReplaceWith: "/home/user1"},
					},
				},
				{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "nologinuser1"},
						{TextToReplace: "$uid", ReplaceWith: "2337"},
						{TextToReplace: "$gid", ReplaceWith: "2338"},
						{TextToReplace: "$home", ReplaceWith: "/dev/null"},
					},
				},
				{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "nologinuser2"},
						{TextToReplace: "$uid", ReplaceWith: "3337"},
						{TextToReplace: "$gid", ReplaceWith: "3338"},
						{TextToReplace: "$home", ReplaceWith: "/bin/nologin"},
					},
				},
				{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "user2"},
						{TextToReplace: "$uid", ReplaceWith: "4337"},
						{TextToReplace: "$gid", ReplaceWith: "4338"},
						{TextToReplace: "$home", ReplaceWith: "/home/user2"},
					},
				},
			},
		},
		{
			desc:       "Users with login",
			configType: ipb.RepeatConfig_FOR_EACH_USER_WITH_LOGIN,
			want: []*repeatconfig.RepeatConfig{
				{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "user1"},
						{TextToReplace: "$uid", ReplaceWith: "1337"},
						{TextToReplace: "$gid", ReplaceWith: "1338"},
						{TextToReplace: "$home", ReplaceWith: "/home/user1"},
					},
				},
				{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "user2"},
						{TextToReplace: "$uid", ReplaceWith: "4337"},
						{TextToReplace: "$gid", ReplaceWith: "4338"},
						{TextToReplace: "$home", ReplaceWith: "/home/user2"},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			config := &ipb.RepeatConfig{Type: tc.configType}
			got, err := repeatconfig.CreateRepeatConfigs(context.Background(), config, &fakeFileReader{content: passwd})
			if err != nil {
				t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) returned an error: %v", config, err)
			}
			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(repeatconfig.RepeatConfig{}, repeatconfig.TokenReplacement{})); diff != "" {
				t.Errorf("repeatconfig.CreateRepeatConfigs(%v) returned unexpected diff (-want +got):\n%s", config, diff)
			}
		})
	}
}

func TestCreateRepeatConfigsForEachSystemUser(t *testing.T) {
	testCases := []struct {
		desc      string
		passwd    string
		loginDefs string
		want      []*repeatconfig.RepeatConfig
	}{
		{
			desc: "Non-default UID_MIN",
			passwd: "systemuser1:x:1000:1000::/home/systemuser1:/bin/bash\n" +
				"systemuser2:x:1999:1999::/home/systemuser2:/bin/bash\n" +
				"normaluser:x:2000:2000::/home/normaluser:/bin/bash\n" +
				"nologin1:x:1:1::/home/nologin1:/bin/false\n" +
				"nologin2:x:2:2::/home/nologin1:/bin/nologin",
			loginDefs: "UID_MIN    2000\n" +
				"SYS_UID_MIN    5000",
			want: []*repeatconfig.RepeatConfig{
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "systemuser1"},
						{TextToReplace: "$uid", ReplaceWith: "1000"},
						{TextToReplace: "$gid", ReplaceWith: "1000"},
						{TextToReplace: "$home", ReplaceWith: "/home/systemuser1"},
					},
				},
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "systemuser2"},
						{TextToReplace: "$uid", ReplaceWith: "1999"},
						{TextToReplace: "$gid", ReplaceWith: "1999"},
						{TextToReplace: "$home", ReplaceWith: "/home/systemuser2"},
					},
				},
			},
		},
		{
			desc: "Default UID_MIN",
			passwd: "systemuser1:x:500:500::/home/systemuser1:/bin/bash\n" +
				"systemuser2:x:999:999::/home/systemuser2:/bin/bash\n" +
				"normaluser:x:1000:1000::/home/normaluser:/bin/bash",
			loginDefs: "",
			want: []*repeatconfig.RepeatConfig{
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "systemuser1"},
						{TextToReplace: "$uid", ReplaceWith: "500"},
						{TextToReplace: "$gid", ReplaceWith: "500"},
						{TextToReplace: "$home", ReplaceWith: "/home/systemuser1"},
					},
				},
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "systemuser2"},
						{TextToReplace: "$uid", ReplaceWith: "999"},
						{TextToReplace: "$gid", ReplaceWith: "999"},
						{TextToReplace: "$home", ReplaceWith: "/home/systemuser2"},
					},
				},
			},
		},
		{
			desc: "SYS_UID_MIN and SYS_UID_MAX present but UID_MIN is not",
			passwd: "systemuser1:x:120:120::/home/systemuser1:/bin/bash\n" +
				"systemuser2:x:456:456::/home/systemuser2:/bin/bash\n" +
				"normaluser:x:5000:5000::/home/normaluser:/bin/bash",
			loginDefs: "SYS_UID_MIN    101\n" +
				"SYS_UID_MAX    999",
			want: []*repeatconfig.RepeatConfig{
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "systemuser1"},
						{TextToReplace: "$uid", ReplaceWith: "120"},
						{TextToReplace: "$gid", ReplaceWith: "120"},
						{TextToReplace: "$home", ReplaceWith: "/home/systemuser1"},
					},
				},
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "systemuser2"},
						{TextToReplace: "$uid", ReplaceWith: "456"},
						{TextToReplace: "$gid", ReplaceWith: "456"},
						{TextToReplace: "$home", ReplaceWith: "/home/systemuser2"},
					},
				},
			},
		},
		{
			desc: "All SYS_UID_MIN, SYS_UID_MAX, and UID_MIN present",
			passwd: "systemuser1:x:120:120::/home/systemuser1:/bin/bash\n" +
				"systemuser2:x:456:456::/home/systemuser2:/bin/bash\n" +
				"normaluser:x:2000:2000::/home/normaluser:/bin/bash",
			loginDefs: "UID_MIN    5000\n" +
				"SYS_UID_MIN    101\n" +
				"SYS_UID_MAX    999",
			want: []*repeatconfig.RepeatConfig{
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "systemuser1"},
						{TextToReplace: "$uid", ReplaceWith: "120"},
						{TextToReplace: "$gid", ReplaceWith: "120"},
						{TextToReplace: "$home", ReplaceWith: "/home/systemuser1"},
					},
				},
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$user", ReplaceWith: "systemuser2"},
						{TextToReplace: "$uid", ReplaceWith: "456"},
						{TextToReplace: "$gid", ReplaceWith: "456"},
						{TextToReplace: "$home", ReplaceWith: "/home/systemuser2"},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			config := &ipb.RepeatConfig{Type: ipb.RepeatConfig_FOR_EACH_SYSTEM_USER_WITH_LOGIN}
			got, err := repeatconfig.CreateRepeatConfigs(context.Background(), config, &fakeFileReader{content: tc.passwd, loginDefsContent: tc.loginDefs})
			if err != nil {
				t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) returned an error: %v", config, err)
			}
			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(repeatconfig.RepeatConfig{}, repeatconfig.TokenReplacement{})); diff != "" {
				t.Errorf("repeatconfig.CreateRepeatConfigs(%v) returned unexpected diff (-want +got):\n%s", config, diff)
			}
		})
	}
}

func TestCreateRepeatConfigsInvalidRepeatOption(t *testing.T) {
	config := &ipb.RepeatConfig{Type: 10}
	if _, err := repeatconfig.CreateRepeatConfigs(context.Background(), config, &fakeFileReader{}); err == nil {
		t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) didn't return an error: %v", config, err)
	}
}

func TestCreateRepeatConfigsInvalidPasswd(t *testing.T) {
	passwd := "invalid"
	config := &ipb.RepeatConfig{Type: ipb.RepeatConfig_FOR_EACH_USER_WITH_LOGIN}
	got, err := repeatconfig.CreateRepeatConfigs(context.Background(), config, &fakeFileReader{content: passwd})
	if err != nil {
		t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) returned an error: %v", config, err)
	}
	if !configHasError(got) {
		t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) didn't return a config with errors", config)
	}
}

func TestCreateRepeatConfigForEachOpenIpv4Port(t *testing.T) {
	testCases := []struct {
		desc       string
		tcpContent string
		want       []*repeatconfig.RepeatConfig
	}{
		{
			desc: "Non-local listening ports included",
			tcpContent: "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n" +
				"  0: 00000000:006F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 15458 1 0000000000000000 100 0 0 10 0\n" +
				"  1: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12934200 1 0000000000000000 100 0 0 10 0",
			want: []*repeatconfig.RepeatConfig{
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$port", ReplaceWith: "111"},
					},
				},
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$port", ReplaceWith: "22"},
					},
				},
			},
		},
		{
			desc: "Local ports skipped",
			tcpContent: "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n" +
				"  0: 0100007F:098C 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12977297 1 0000000000000000 100 0 0 10 0\n" +
				"  1: 0100007F:2555 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 25317325 1 0000000000000000 100 0 0 10 0",
			want: []*repeatconfig.RepeatConfig{},
		},
		{
			desc: "Non-listening ports skipped",
			tcpContent: "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n" +
				"  0: 00000000:006F 00000000:0000 01 00000000:00000000 00:00000000 00000000     0        0 15458 1 0000000000000000 100 0 0 10 0\n" +
				"  1: 00000000:0016 00000000:0000 01 00000000:00000000 00:00000000 00000000     0        0 12934200 1 0000000000000000 100 0 0 10 0",
			want: []*repeatconfig.RepeatConfig{},
		},
		{
			desc: "Ports deduplicated",
			tcpContent: "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n" +
				"  0: 00000000:006F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 15458 1 0000000000000000 100 0 0 10 0\n" +
				"  1: 00000000:006F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12934200 1 0000000000000000 100 0 0 10 0",
			want: []*repeatconfig.RepeatConfig{
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$port", ReplaceWith: "111"},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			config := &ipb.RepeatConfig{Type: ipb.RepeatConfig_FOR_EACH_OPEN_IPV4_PORT}
			got, err := repeatconfig.CreateRepeatConfigs(
				context.Background(),
				config,
				&fakeFileReader{content: tc.tcpContent},
			)
			if err != nil {
				t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) returned an error: %v", config, err)
			}
			sort := func(r1, r2 *repeatconfig.RepeatConfig) bool {
				return r1.TokenReplacements[0].ReplaceWith < r2.TokenReplacements[0].ReplaceWith
			}
			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(repeatconfig.RepeatConfig{}, repeatconfig.TokenReplacement{}), cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("repeatconfig.CreateRepeatConfigs(%v) returned unexpected diff (-want +got):\n%s", config, diff)
			}
		})
	}
}

func TestCreateRepeatConfigForEachOpenIpv6Port(t *testing.T) {
	testCases := []struct {
		desc       string
		tcpContent string
		want       []*repeatconfig.RepeatConfig
	}{
		{
			desc: "Non-local listening ports included",
			tcpContent: "  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n" +
				"   0: 00000000000000000000000000000000:312D 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000   114        0 27677486 1 0000000000000000 100 0 0 10 0\n" +
				"   1: 00000000000000000000000000000000:AB6F 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000 454117        0 28573562 1 0000000000000000 100 0 0 10 0",
			want: []*repeatconfig.RepeatConfig{
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$port", ReplaceWith: "12589"},
					},
				},
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$port", ReplaceWith: "43887"},
					},
				},
			},
		},
		{
			desc: "Local ports skipped",
			tcpContent: "  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n" +
				"   0: 00000000000000000000000001000000:82D5 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000 454117        0 28610339 1 0000000000000000 100 0 0 10 0\n" +
				"   1: 00000000000000000000000001000000:2555 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 25318393 1 0000000000000000 100 0 0 10 0",
			want: []*repeatconfig.RepeatConfig{},
		},
		{
			desc: "Non-listening ports skipped",
			tcpContent: "  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n" +
				"  0: E079002A02021B00CAAC1C7E77ADFD91:A972 5014002A030C13400000000081000000:13A7 01 00000000:00000000 02:0000003C 00000000     0        0 28580605 2 0000000000000000 21 4 30 10 -1\n" +
				"  1: 0000000000000000FFFF00000100007F:842C 0000000000000000FFFF00000100007F:2555 01 00000000:00000000 00:00000000 00000000 454117        0 28620360 1 0000000000000000 20 4 24 10 -1",
			want: []*repeatconfig.RepeatConfig{},
		},
		{
			desc: "Ports deduplicated",
			tcpContent: "  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n" +
				"   0: 00000000000000000000000000000000:312D 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000   114        0 27677486 1 0000000000000000 100 0 0 10 0\n" +
				"   1: 00000000000000000000000000000000:312D 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000 454117        0 28573562 1 0000000000000000 100 0 0 10 0",
			want: []*repeatconfig.RepeatConfig{
				&repeatconfig.RepeatConfig{
					TokenReplacements: []*repeatconfig.TokenReplacement{
						{TextToReplace: "$port", ReplaceWith: "12589"},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			config := &ipb.RepeatConfig{Type: ipb.RepeatConfig_FOR_EACH_OPEN_IPV6_PORT}
			got, err := repeatconfig.CreateRepeatConfigs(
				context.Background(),
				config,
				&fakeFileReader{content: tc.tcpContent},
			)
			if err != nil {
				t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) returned an error: %v", config, err)
			}
			sort := func(r1, r2 *repeatconfig.RepeatConfig) bool {
				return r1.TokenReplacements[0].ReplaceWith < r2.TokenReplacements[0].ReplaceWith
			}
			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(repeatconfig.RepeatConfig{}, repeatconfig.TokenReplacement{}), cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("repeatconfig.CreateRepeatConfigs(%v) returned unexpected diff (-want +got):\n%s", config, diff)
			}
		})
	}
}

func TestCreateRepeatConfigForEachOpenPortInvalidTCPContent(t *testing.T) {
	tcpContent := "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n" +
		"invalid"
	config := &ipb.RepeatConfig{Type: ipb.RepeatConfig_FOR_EACH_OPEN_IPV4_PORT}
	got, err := repeatconfig.CreateRepeatConfigs(
		context.Background(),
		config,
		&fakeFileReader{content: tcpContent},
	)
	if err != nil {
		t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) returned an error: %v", config, err)
	}
	if !configHasError(got) {
		t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) didn't return a config with errors", config)
	}

	config = &ipb.RepeatConfig{Type: ipb.RepeatConfig_FOR_EACH_OPEN_IPV6_PORT}
	got, err = repeatconfig.CreateRepeatConfigs(
		context.Background(),
		config,
		&fakeFileReader{content: tcpContent},
	)
	if err != nil {
		t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) returned an error: %v", config, err)
	}
	if !configHasError(got) {
		t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) didn't return a config with errors", config)
	}
}

func TestCreateRepeatConfigsWithOptOut(t *testing.T) {
	passwd := "user1:x:1337:1338::/home/user1:/bin/bash\n" +
		"user2:x:2337:2338::/home/user2:/bin/bash\n" +
		"user3:x:3337:3338::/home/user3:/bin/bash\n" +
		"user4:x:4337:4338::/home/user4:/bin/bash"
	repeatOptions := &ipb.RepeatConfig{
		Type: ipb.RepeatConfig_FOR_EACH_USER_WITH_LOGIN,
		OptOut: []*ipb.RepeatConfig_OptOutSubstitution{
			{
				Wildcard: "$user",
				Value:    "user1",
			},
			{
				Wildcard: "$home",
				Value:    "/home/user3",
			},
		},
	}
	got, err := repeatconfig.CreateRepeatConfigs(context.Background(), repeatOptions, &fakeFileReader{content: passwd})
	if err != nil {
		t.Fatalf("repeatconfig.CreateRepeatConfigs(%v) returned an error: %v", repeatOptions, err)
	}
	want := []*repeatconfig.RepeatConfig{
		{
			TokenReplacements: []*repeatconfig.TokenReplacement{
				{TextToReplace: "$user", ReplaceWith: "user2"},
				{TextToReplace: "$uid", ReplaceWith: "2337"},
				{TextToReplace: "$gid", ReplaceWith: "2338"},
				{TextToReplace: "$home", ReplaceWith: "/home/user2"},
			},
		},
		{
			TokenReplacements: []*repeatconfig.TokenReplacement{
				{TextToReplace: "$user", ReplaceWith: "user4"},
				{TextToReplace: "$uid", ReplaceWith: "4337"},
				{TextToReplace: "$gid", ReplaceWith: "4338"},
				{TextToReplace: "$home", ReplaceWith: "/home/user4"},
			},
		},
	}
	if diff := cmp.Diff(want, got, cmp.AllowUnexported(repeatconfig.RepeatConfig{}, repeatconfig.TokenReplacement{})); diff != "" {
		t.Errorf("repeatconfig.CreateRepeatConfigs(%v) returned unexpected diff (-want +got):\n%s", repeatOptions, diff)
	}
}

func TestApplyRepeatConfigToInstruction(t *testing.T) {
	config := &repeatconfig.RepeatConfig{
		TokenReplacements: []*repeatconfig.TokenReplacement{
			{TextToReplace: "$user", ReplaceWith: "root"},
			{TextToReplace: "$home", ReplaceWith: "/root"},
		},
	}
	testCases := []struct {
		desc        string
		instruction *ipb.FileCheck
		want        *ipb.FileCheck
	}{
		{
			desc: "permission with user",
			instruction: &ipb.FileCheck{
				CheckType: &ipb.FileCheck_Permission{Permission: &ipb.PermissionCheck{
					User: &ipb.PermissionCheck_OwnerCheck{Name: "$user"},
				},
				}},
			want: &ipb.FileCheck{
				CheckType: &ipb.FileCheck_Permission{Permission: &ipb.PermissionCheck{
					User: &ipb.PermissionCheck_OwnerCheck{Name: "root"},
				}},
			},
		},
		{
			desc: "permission with group",
			instruction: &ipb.FileCheck{
				CheckType: &ipb.FileCheck_Permission{Permission: &ipb.PermissionCheck{
					Group: &ipb.PermissionCheck_OwnerCheck{Name: "$user"},
				}},
			},
			want: &ipb.FileCheck{
				CheckType: &ipb.FileCheck_Permission{Permission: &ipb.PermissionCheck{
					Group: &ipb.PermissionCheck_OwnerCheck{Name: "root"},
				}},
			},
		},
		{
			desc: "content",
			instruction: &ipb.FileCheck{
				CheckType: &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Home dir of $user is $home"}},
			},
			want: &ipb.FileCheck{
				CheckType: &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Home dir of root is /root"}},
			},
		},
		{
			desc: "content entry",
			instruction: &ipb.FileCheck{
				CheckType: &ipb.FileCheck_ContentEntry{ContentEntry: &ipb.ContentEntryCheck{
					MatchCriteria: []*ipb.MatchCriterion{
						&ipb.MatchCriterion{
							FilterRegex:   "$user=.*",
							ExpectedRegex: "$user=2",
						},
					},
				}},
			},
			want: &ipb.FileCheck{
				CheckType: &ipb.FileCheck_ContentEntry{ContentEntry: &ipb.ContentEntryCheck{
					MatchCriteria: []*ipb.MatchCriterion{
						&ipb.MatchCriterion{
							FilterRegex:   "root=.*",
							ExpectedRegex: "root=2",
						},
					},
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := repeatconfig.ApplyRepeatConfigToInstruction(tc.instruction, config)
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("repeatconfig.ApplyRepeatConfigToInstruction(%v %v) returned unexpected diff (-want +got):\n%s",
					tc.instruction, config, diff)
			}
		})
	}
}

func TestApplyRepeatConfigToFile(t *testing.T) {
	config := &repeatconfig.RepeatConfig{
		TokenReplacements: []*repeatconfig.TokenReplacement{
			{TextToReplace: "$user", ReplaceWith: "sundar"},
		},
	}
	testCases := []struct {
		desc string
		file *ipb.FileSet
		want *ipb.FileSet
	}{
		{
			desc: "single file",
			file: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/home/$user"}},
			},
			want: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/home/sundar"}},
			},
		},
		{
			desc: "files in dir",
			file: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{DirPath: "/home/$user"}},
			},
			want: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{DirPath: "/home/sundar"}},
			},
		},
		{
			desc: "files in dir with filename regex",
			file: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:       "/path",
					FilenameRegex: "$user-.*",
				}},
			},
			want: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:       "/path",
					FilenameRegex: "sundar-.*",
				}},
			},
		},
		{
			desc: "files in dir with path regex",
			file: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:           "/path",
					OptOutPathRegexes: []string{"$user-.*"},
				}},
			},
			want: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:           "/path",
					OptOutPathRegexes: []string{"sundar-.*"},
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := repeatconfig.ApplyRepeatConfigToFile(tc.file, config)
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("repeatconfig.ApplyRepeatConfigToFile(%v %v) returned unexpected diff (-want +got):\n%s", tc.file, config, diff)
			}
		})
	}
}
