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

package configchecks

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

var (
	umaskRe = regexp.MustCompile("^0?[0-7][0-7][0-7]$")
)

var (
	swVersRe = regexp.MustCompile("^\\d+.\\d+.\\d+-\\d+\\+.*$")
)

type groupCriteria struct {
	regex    string
	criteria []groupCriterion
}

func newGroupCriteria(regex string, numSubexp int, gcs []*ipb.GroupCriterion, matchType ipb.ContentEntryCheck_MatchType) (*groupCriteria, error) {
	criteria := make([]groupCriterion, 0, len(gcs))
	for _, gc := range gcs {
		i := int(gc.GetGroupIndex())
		if i <= 0 || i > numSubexp {
			return nil, fmt.Errorf("group criteria index %d out of bounds", i)
		}

		var m groupCriterionMatcher
		switch t := gc.GetType(); t {
		case ipb.GroupCriterion_LESS_THAN:
			m = &lessThanMatcher{cmp: getCmp(gc), cmpStr: getCmpStr(gc)}
		case ipb.GroupCriterion_GREATER_THAN:
			m = &greaterThanMatcher{cmp: getCmp(gc), cmpStr: getCmpStr(gc)}
		case ipb.GroupCriterion_NO_LESS_RESTRICTIVE_UMASK:
			if gc.GetToday() != nil { // today set instead of const
				return nil, errors.New("GroupCriterion_NO_LESS_RESTRICTIVE_UMASK requires a constant to compare to")
			}
			m = &umaskMatcher{wantMask: gc.GetConst(), cmpStr: getCmpStr(gc)}
		case ipb.GroupCriterion_UNIQUE:
			if matchType == ipb.ContentEntryCheck_NONE_MATCH {
				return nil, errors.New("GroupCriterion_UNIQUE and ContentEntryCheck_NONE_MATCH are incompatible")
			}
			m = &uniqueMatcher{seen: make(map[string]bool)}
		case ipb.GroupCriterion_VERSION_LESS_THAN:
			m = &lessThanVersionMatcher{cmpStr: gc.GetVersion()}
		case ipb.GroupCriterion_VERSION_GREATER_THAN:
			m = &greaterThanVersionMatcher{cmpStr: gc.GetVersion()}
		default:
			return nil, fmt.Errorf("unrecognized group criterion type %v", t)
		}

		criteria = append(criteria, groupCriterion{
			index:   i,
			matcher: m,
		})
	}

	return &groupCriteria{
		regex:    regex,
		criteria: criteria,
	}, nil
}

// check returns whether all of the group criteria are met by the entry.
// The check assumes that entry matches gc.re so groups can safely be extracted.
// Also it assumes that all indices of the contained groupCriterion objects have
// been bounds checked before.
func (gc *groupCriteria) check(entry string) bool {
	if len(gc.criteria) == 0 {
		return true
	}
	groups := compiledRegex(gc.regex).FindStringSubmatch(entry)
	for _, c := range gc.criteria {
		g := groups[c.index]
		if !c.matcher.match(g) {
			return false
		}
	}
	return true
}

func (gc *groupCriteria) String() string {
	if len(gc.criteria) == 0 {
		return ""
	}
	strs := make([]string, 0, len(gc.criteria))
	for _, crit := range gc.criteria {
		strs = append(strs, crit.String())
	}
	return "{" + strings.Join(strs, ", ") + "}"
}

type groupCriterion struct {
	index   int
	matcher groupCriterionMatcher
}

func (gc *groupCriterion) String() string {
	return fmt.Sprintf("[group#%d %s]", gc.index, gc.matcher)
}

type groupCriterionMatcher interface {
	match(group string) bool
	String() string
}

type lessThanMatcher struct {
	cmp    int32
	cmpStr string
}

func (m *lessThanMatcher) match(group string) bool {
	parsed, err := strconv.ParseInt(group, 10, 32)
	if err != nil {
		log.Printf("unable to parse %q as int32: %v", group, err)
		return false
	}
	val := int32(parsed)
	return val < m.cmp
}

func (m *lessThanMatcher) String() string {
	return "< " + m.cmpStr
}

type greaterThanMatcher struct {
	cmp    int32
	cmpStr string
}

func (m *greaterThanMatcher) match(group string) bool {
	parsed, err := strconv.ParseInt(group, 10, 32)
	if err != nil {
		log.Printf("unable to parse %q as int32: %v", group, err)
		return false
	}
	val := int32(parsed)
	return val > m.cmp
}

func (m *greaterThanMatcher) String() string {
	return "> " + m.cmpStr
}

type umaskMatcher struct {
	wantMask int32
	cmpStr   string
}

func (m *umaskMatcher) match(group string) bool {
	if !umaskRe.MatchString(group) {
		log.Printf("unable to parse %q as a umask", group)
		return false
	}
	parsed, err := strconv.ParseInt(group, 8, 32)
	if err != nil {
		log.Printf("unable to parse %q as a umask: %v", group, err)
		return false
	}
	mask := int32(parsed)
	// At least all bits that are set in the expected mask have to be set in the actual mask.
	// More set bits would mean more restrictive permissions.
	return (m.wantMask & mask) == m.wantMask
}

func (m *umaskMatcher) String() string {
	return "not less restrictive than " + m.cmpStr
}

type uniqueMatcher struct {
	seen map[string]bool
}

func (m *uniqueMatcher) match(group string) bool {
	if m.seen[group] {
		return false
	}
	m.seen[group] = true
	return true
}

func (m *uniqueMatcher) String() string {
	return "is unique"
}

// getCmp returns the value to compare the group against for a given LESS_THAN or GREATER_THAN group criterion.
// In case the today field is set on the proto, the number of days since the epoch for the current time is returned.
func getCmp(gc *ipb.GroupCriterion) int32 {
	if gc.GetToday() != nil {
		now := time.Now()
		epoch := time.Unix(0, 0)
		days := now.Sub(epoch).Hours() / 24
		return int32(days)
	}
	return gc.GetConst()
}

func getCmpStr(gc *ipb.GroupCriterion) string {
	if gc.GetToday() != nil {
		return "today"
	}
	return strconv.Itoa(int(gc.GetConst()))
}
type lessThanVersionMatcher struct {
	cmpStr string
}
func (m *lessThanVersionMatcher) match(group string) bool {
	if !swVersRe.MatchString(group) {
		log.Printf("unable to parse %q as a software version", group)
		return false
	}
	version := strings.ReplaceAll(m.cmpStr, "-", ".")
	re := regexp.MustCompile("\\W+")
	chunks_group := re.Split(group, -1)
	chunks_version := re.Split(version,-1)
	min_len := 0
	if len(chunks_group) < len(chunks_version) {
		min_len = len(chunks_group)
	} else {
		min_len = len(chunks_version)
	}
	for i:=0; i< min_len; i++ {
		chunks_group[i] = fmt.Sprintf("%06s", chunks_group[i])
		chunks_version[i] = fmt.Sprintf("%06s", chunks_version[i])
		if chunks_group[i] == chunks_version[i] {
			continue
		}
		return chunks_group[i] < chunks_version[i]
	}
	return false
}
func (m *lessThanVersionMatcher) String() string {
	return "< " + m.cmpStr
}
type greaterThanVersionMatcher struct {
	cmpStr string
}
func (m *greaterThanVersionMatcher) match(group string) bool {
	if !swVersRe.MatchString(group) {
		log.Printf("unable to parse %q as a software version", group)
		return false
	}
	version := strings.ReplaceAll(m.cmpStr, "-", ".")
	re := regexp.MustCompile("\\W+")
	chunks_group := re.Split(group, -1)
	chunks_version := re.Split(version,-1)
	min_len := 0
	if len(chunks_group) < len(chunks_version) {
		min_len = len(chunks_group)
	} else {
		min_len = len(chunks_version)
	}
	for i:=0; i< min_len; i++ {
		chunks_group[i] = fmt.Sprintf("%06s", chunks_group[i])
		chunks_version[i] = fmt.Sprintf("%06s", chunks_version[i])
		if chunks_group[i] == chunks_version[i] {
			continue
		}
		return chunks_group[i] > chunks_version[i]
	}
	return false
}
func (m *greaterThanVersionMatcher) String() string {
	return "< " + m.cmpStr
}
