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

// Package repeatconfig provides a utility for applying the RepeatConfig proto to config checks.
package repeatconfig

import (
	"bufio"
	"context"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"github.com/google/localtoast/scanapi"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

const (
	usernameWildcard = "$user"
	uidWildcard      = "$uid"
	gidWildcard      = "$gid"
	homeDirWildcard  = "$home"
	portWildcard     = "$port"
	shellWildcard    = "$shell"
	defaultUIDMin    = 1000
)

var (
	// Regexp for capturing the local ip:port and TCP state from the
	// /proc/net/tcp(6) file. See https://regex101.com/r/MFk3qd/1 for how to
	// read this regex.
	procTCPRe = regexp.MustCompile("^\\s*[0-9]+:\\s+([0-9A-F]+):([0-9A-F]+)\\s+[0-9A-F]+:[0-9A-F]+\\s+([0-9A-F]+)\\s+.*$")
	// Regexp for capturing the UID_MIN value from the /etc/login.defs file.
	uidMinRe = regexp.MustCompile("^UID_MIN\\s+(\\d+)$")
	// Regexp for capturing the SYS_UID_MIN value from the /etc/login.defs file.
	sysUIDMinRe = regexp.MustCompile("^SYS_UID_MIN\\s+(\\d+)$")
	// Regexp for capturing the SYS_UID_MAX value from the /etc/login.defs file.
	sysUIDMaxRe = regexp.MustCompile("^SYS_UID_MAX\\s+(\\d+)$")
)

// RepeatConfig is a single repeat config that specifies what tokens to replace
// in the files/instructions in this iteration of the check.
type RepeatConfig struct {
	TokenReplacements []*TokenReplacement
	Err               error
}

// TokenReplacement describes a single token to replace.
type TokenReplacement struct {
	TextToReplace string
	ReplaceWith   string
}

// CreateRepeatConfigs creates a list of configs with the appropriate token
// substitutions based on the supplied repeat config enum.
func CreateRepeatConfigs(ctx context.Context, repeatOptions *ipb.RepeatConfig, fs scanapi.Filesystem) ([]*RepeatConfig, error) {
	var rc []*RepeatConfig
	var err error
	switch repeatOptions.GetType() {
	case ipb.RepeatConfig_ONCE:
		rc, err = []*RepeatConfig{&RepeatConfig{}}, nil
	case ipb.RepeatConfig_FOR_EACH_USER:
		rc, err = createRepeatConfigForEachUser(ctx, userRepeatConfigOptions{
			fileReader: fs, loginOnly: false, systemOnly: false,
		})
	case ipb.RepeatConfig_FOR_EACH_USER_WITH_LOGIN:
		rc, err = createRepeatConfigForEachUser(ctx, userRepeatConfigOptions{
			fileReader: fs, loginOnly: true, systemOnly: false,
		})
	case ipb.RepeatConfig_FOR_EACH_SYSTEM_USER_WITH_LOGIN:
		rc, err = createRepeatConfigForEachUser(ctx, userRepeatConfigOptions{
			fileReader: fs, loginOnly: true, systemOnly: true,
		})
	case ipb.RepeatConfig_FOR_EACH_OPEN_IPV4_PORT:
		rc, err = createRepeatConfigForEachOpenTCPPort(ctx, fs, false)
	case ipb.RepeatConfig_FOR_EACH_OPEN_IPV6_PORT:
		rc, err = createRepeatConfigForEachOpenTCPPort(ctx, fs, true)
	default:
		return nil, fmt.Errorf("unknown repeat option type %s", repeatOptions)
	}
	if err != nil {
		return repeatConfigWithError(err), nil
	}
	return applyOptOutConfig(rc, repeatOptions.GetOptOut()), nil
}

func repeatConfigWithError(err error) []*RepeatConfig {
	return []*RepeatConfig{{Err: fmt.Errorf("error creating RepeatConfig: %v", err)}}
}

type userRepeatConfigOptions struct {
	fileReader scanapi.Filesystem
	loginOnly  bool
	systemOnly bool
}

// createRepeatConfigForEachUser finds all users in /etc/passwd and creates
// repeat configs that have the usernames as the substitution. if systemOnly is
// true, only the system users are included in the config.
func createRepeatConfigForEachUser(ctx context.Context, opt userRepeatConfigOptions) ([]*RepeatConfig, error) {
	uidMin, sysUIDMin, sysUIDMax := -1, -1, -1
	if opt.systemOnly {
		var err error
		uidMin, sysUIDMin, sysUIDMax, err = readUID(ctx, opt.fileReader)
		if err != nil {
			return nil, err
		}
	}

	r, err := opt.fileReader.OpenFile(ctx, "/etc/passwd")
	if err != nil {
		return nil, err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	result := []*RepeatConfig{}

	i := 0
	for scanner.Scan() {
		i++
		if scanner.Err() != nil {
			return nil, scanner.Err()
		}
		line := scanner.Text()
		tokens := strings.Split(line, ":")
		if len(tokens) < 7 {
			return nil, fmt.Errorf("can't parse line %d in /etc/passwd: expected at least 7 tokens, got %d", i, len(tokens))
		}
		user := tokens[0]
		uid := tokens[2]
		gid := tokens[3]
		homeDir := tokens[5]
		shell := tokens[6]

		if opt.systemOnly {
			uidInt, err := strconv.Atoi(uid)
			if err != nil {
				return nil, err
			}
			if sysUIDMin == -1 || sysUIDMax == -1 {
				if uidInt >= uidMin { // Non-system users' uid starts from UID_MIN.
					continue
				}
			} else if uidInt < sysUIDMin || uidInt > sysUIDMax {
				continue
			}
		}
		// Ignore users with no shell.
		if opt.loginOnly && (shell == "/bin/false" || path.Base(shell) == "nologin") {
			continue
		}

		result = append(result, &RepeatConfig{
			TokenReplacements: []*TokenReplacement{
				{
					TextToReplace: usernameWildcard,
					ReplaceWith:   user,
				},
				{
					TextToReplace: uidWildcard,
					ReplaceWith:   uid,
				},
				{
					TextToReplace: gidWildcard,
					ReplaceWith:   gid,
				},
				{
					TextToReplace: homeDirWildcard,
					ReplaceWith:   homeDir,
				},
				{
					TextToReplace: shellWildcard,
					ReplaceWith:   shell,
				},
			},
		})
	}
	return result, nil
}

func readUID(ctx context.Context, f scanapi.Filesystem) (uidMin, sysUIDMin, sysUIDMax int, err error) {
	r, err := f.OpenFile(ctx, "/etc/login.defs")
	if err != nil {
		return 0, 0, 0, err
	}
	defer r.Close()

	uidMin = defaultUIDMin
	sysUIDMin = -1
	sysUIDMax = -1

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return 0, 0, 0, scanner.Err()
		}

		line := scanner.Text()

		// Look for UID_MIN value
		groups := uidMinRe.FindStringSubmatch(line)
		if groups != nil {
			uidMin, err = strconv.Atoi(groups[1])
			if err != nil {
				return 0, 0, 0, err
			}
		}

		// Look for SYS_UID_MIN value
		groups = sysUIDMinRe.FindStringSubmatch(line)
		if groups != nil {
			sysUIDMin, err = strconv.Atoi(groups[1])
			if err != nil {
				return 0, 0, 0, err
			}
		}

		// Look for SYS_UID_MAX value
		groups = sysUIDMaxRe.FindStringSubmatch(line)
		if groups != nil {
			sysUIDMax, err = strconv.Atoi(groups[1])
			if err != nil {
				return 0, 0, 0, err
			}
		}
	}
	return uidMin, sysUIDMin, sysUIDMax, nil
}

// createRepeatConfigForEachPort creates repeat configs that have the currently
// open TCP ports as the substitution.
func createRepeatConfigForEachOpenTCPPort(ctx context.Context, f scanapi.Filesystem, isIpv6 bool) ([]*RepeatConfig, error) {
	tcpFile := "/proc/net/tcp"
	if isIpv6 {
		tcpFile = "/proc/net/tcp6"
	}
	r, err := f.OpenFile(ctx, tcpFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	scanner.Scan() // Skip the header line.
	ports := make(map[int]bool)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return nil, scanner.Err()
		}
		line := scanner.Text()
		if !procTCPRe.MatchString(line) {
			return nil, fmt.Errorf("unable to parse %s", tcpFile)
		}
		groups := procTCPRe.FindStringSubmatch(line)
		address := groups[1]
		port, err := strconv.ParseInt(groups[2], 16, 32)
		if err != nil {
			return nil, err
		}
		state, err := strconv.ParseInt(groups[3], 16, 32)
		if err != nil {
			return nil, err
		}
		// We only care about ports in state LISTENING (0xa).
		if state != 0xa {
			continue
		}
		// Skip loopback addresses.
		// 127.0.0.0/8 for IPv4. The address consists of one little-endian word.
		if !isIpv6 && strings.HasSuffix(address, "7F") {
			continue
		}
		// ::1 for IPv6. The address consists of 4 little-endian words.
		if isIpv6 && address == "00000000000000000000000001000000" {
			continue
		}

		ports[int(port)] = true
	}

	result := []*RepeatConfig{}
	for port, _ := range ports {
		result = append(result, &RepeatConfig{
			TokenReplacements: []*TokenReplacement{
				&TokenReplacement{
					TextToReplace: portWildcard,
					ReplaceWith:   strconv.Itoa(port),
				},
			},
		})
	}
	return result, nil
}

// ApplyRepeatConfigToInstruction applies the substitutions in the given repeat config to the
// given instruction. The returned instruction proto is a copy of the original.
func ApplyRepeatConfigToInstruction(instruction *ipb.FileCheck, config *RepeatConfig) *ipb.FileCheck {
	if len(config.TokenReplacements) == 0 {
		return instruction
	}
	result := proto.Clone(instruction).(*ipb.FileCheck)
	for _, r := range config.TokenReplacements {
		switch {
		case instruction.GetPermission() != nil:
			applyRepeatConfigToPermissionCheck(result.GetPermission(), r)
		case instruction.GetContent() != nil:
			applyRepeatConfigToContentCheck(result.GetContent(), r)
		case instruction.GetContentEntry() != nil:
			applyRepeatConfigToContentEntryCheck(result.GetContentEntry(), r)
		}
	}
	return result
}

// ApplyRepeatConfigToFile applies the substitutions in the given repeat config
// to the given FileSet. The returned FileSet proto is a copy of the original.
func ApplyRepeatConfigToFile(fileSet *ipb.FileSet, config *RepeatConfig) *ipb.FileSet {
	if len(config.TokenReplacements) == 0 {
		return fileSet
	}
	result := proto.Clone(fileSet).(*ipb.FileSet)
	for _, r := range config.TokenReplacements {
		switch {
		case fileSet.GetSingleFile() != nil:
			result.GetSingleFile().Path = applyReplacement(result.GetSingleFile().GetPath(), r)
		case fileSet.GetFilesInDir() != nil:
			f := result.GetFilesInDir()
			f.DirPath = applyReplacement(f.GetDirPath(), r)
			f.FilenameRegex = applyReplacement(f.GetFilenameRegex(), r)
			optOut := f.GetOptOutPathRegexes()
			for i, o := range optOut {
				optOut[i] = applyReplacement(o, r)
			}
			f.OptOutPathRegexes = optOut
		}
	}
	return result
}

func applyRepeatConfigToPermissionCheck(check *ipb.PermissionCheck, replacement *TokenReplacement) {
	if check.GetUser() != nil {
		check.GetUser().Name = applyReplacement(check.GetUser().GetName(), replacement)
	}
	if check.GetGroup() != nil {
		check.GetGroup().Name = applyReplacement(check.GetGroup().GetName(), replacement)
	}
}

func applyRepeatConfigToContentCheck(check *ipb.ContentCheck, replacement *TokenReplacement) {
	check.Content = applyReplacement(check.GetContent(), replacement)
}

func applyRepeatConfigToContentEntryCheck(check *ipb.ContentEntryCheck, replacement *TokenReplacement) {
	for _, mc := range check.GetMatchCriteria() {
		mc.FilterRegex = applyReplacement(mc.GetFilterRegex(), replacement)
		mc.ExpectedRegex = applyReplacement(mc.GetExpectedRegex(), replacement)
	}
}

func applyReplacement(str string, replacement *TokenReplacement) string {
	return strings.ReplaceAll(str, replacement.TextToReplace, replacement.ReplaceWith)
}

func applyOptOutConfig(rcs []*RepeatConfig, optOut []*ipb.RepeatConfig_OptOutSubstitution) []*RepeatConfig {
	filteredRcs := make([]*RepeatConfig, 0, len(rcs))
	for _, rc := range rcs {
		filter := false
		for _, r := range rc.TokenReplacements {
			for _, o := range optOut {
				if o.GetWildcard() == r.TextToReplace && o.GetValue() == r.ReplaceWith {
					filter = true
					break
				}
			}
		}
		if !filter {
			filteredRcs = append(filteredRcs, rc)
		}
	}
	return filteredRcs
}
