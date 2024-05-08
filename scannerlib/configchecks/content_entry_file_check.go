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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"sync"

	"github.com/google/localtoast/scannerlib/fileset"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

var (
	regexCache = make(map[string]*regexp.Regexp)
	mutex = &sync.Mutex{}
)

// contentEntryFileChecker performs checks about whether files have specific entries in their content.
type contentEntryFileChecker struct {
	fc            *fileCheck
	delimiter     []byte
	matchCriteria []*matchCriterion
}

type matchCriterion struct {
	matchType     ipb.ContentEntryCheck_MatchType
	filterRegex   string
	expectedRegex string
	groupCriteria *groupCriteria
	matched       bool
}

func (m *matchCriterion) String() string {
	if gcStr := m.groupCriteria.String(); gcStr != "" {
		return fmt.Sprintf("%s with group criteria %s", m.expectedRegex, gcStr)
	}
	return m.expectedRegex
}

// compiledRegex returns the (potentially cached) compiled regex pattern.
// It assumes that the regex can be successfully compiled without errors.
func compiledRegex(pattern string) *regexp.Regexp {
	mutex.Lock()
	defer mutex.Unlock()
	if re, ok := regexCache[pattern]; ok {
		return re
	}
	re := regexp.MustCompile(pattern)
	regexCache[pattern] = re
	return re
}

func clearRegexCache() {
	mutex.Lock()
	defer mutex.Unlock()
	regexCache = make(map[string]*regexp.Regexp)
}

func newContentEntryFileChecker(fc *fileCheck) (*contentEntryFileChecker, error) {
	delimiter := fc.checkInstruction.GetContentEntry().GetDelimiter()
	if len(delimiter) == 0 {
		// Split by lines if nothing else is specified.
		delimiter = []byte{'\n'}
	}

	matchCriteriaProtos := fc.checkInstruction.GetContentEntry().GetMatchCriteria()
	matchCriteria := make([]*matchCriterion, 0, len(matchCriteriaProtos))
	for _, mc := range matchCriteriaProtos {
		mode := "(?s)" // '.' matches '\n' too
		filterRegex := mode + "^" + mc.GetFilterRegex() + "$"
		expectedRegex := mode + "^" + mc.GetExpectedRegex() + "$"
		// Check regexes for errors.
		if _, err := regexp.Compile(filterRegex); err != nil {
			return nil, err
		}
		compiledExpectedRegex, err := regexp.Compile(expectedRegex)
		if err != nil {
			return nil, err
		}
		groupCriteria, err := newGroupCriteria(
			expectedRegex,
			compiledExpectedRegex.NumSubexp(),
			mc.GetGroupCriteria(),
			fc.checkInstruction.GetContentEntry().GetMatchType(),
		)
		if err != nil {
			return nil, err
		}
		matchCriteria = append(matchCriteria, &matchCriterion{
			matchType:     fc.checkInstruction.GetContentEntry().GetMatchType(),
			filterRegex:   filterRegex,
			expectedRegex: expectedRegex,
			groupCriteria: groupCriteria,
			matched:       false,
		})
	}
	return &contentEntryFileChecker{
		fc:            fc,
		delimiter:     delimiter,
		matchCriteria: matchCriteria,
	}, nil
}

func (c *fileCheckers) execContentEntryChecksOnFile(path string, openError error, f io.ReadCloser) error {
	if len(c.contentEntryFileCheckers) == 0 {
		return nil
	}

	exists, err := fileExists(openError)
	if err != nil {
		return err
	}
	if !exists {
		for _, checker := range c.contentEntryFileCheckers {
			if checker.fc.checkInstruction.GetContentEntry().GetMatchType() != ipb.ContentEntryCheck_NONE_MATCH {
				checker.fc.addNonCompliantFile(path, "File doesn't exist")
			}
		}
		return nil
	}

	scanner := bufio.NewScanner(f)
	// Define a split function that splits the entries based on the delimiter char.
	delimiter := c.contentEntryFileCheckers[0].delimiter
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.Index(data, delimiter); i >= 0 {
			return i + len(delimiter), data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	})

	// Run all checks on each entry.
	for scanner.Scan() {
		if scanner.Err() != nil {
			return scanner.Err()
		}
		entry := scanner.Text()
		for _, checker := range c.contentEntryFileCheckers {
			if err := matchEntryAgainstCriteria(entry, path, checker.fc, checker.matchCriteria); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *fileCheckers) execContentEntryChecksAfterFileTraversal(filesToCheck *ipb.FileSet) {
	// Check for any remaining unmatched criteria.
	for _, checker := range c.contentEntryFileCheckers {
		for _, mc := range checker.matchCriteria {
			if !mc.matched && mc.matchType != ipb.ContentEntryCheck_NONE_MATCH {
				checker.fc.addNonCompliantFile(
					fileset.FileSetToString(filesToCheck),
					fmt.Sprintf("No entry matching %q found among files", mc.filterRegex))
			}
		}
	}
	// Clear the regex cache after execution to keep the memory usage low.
	// The next check will likely use different regexes anyway.
	defer clearRegexCache()

}

func matchEntryAgainstCriteria(entry string, filePath string, fc *fileCheck, matchCriteria []*matchCriterion) error {
	for i, mc := range matchCriteria {
		if !compiledRegex(mc.filterRegex).MatchString(entry) {
			continue
		}

		mc.matched = true
		satisfiesCriterion := compiledRegex(mc.expectedRegex).MatchString(entry) && mc.groupCriteria.check(entry)

		if !satisfiesCriterion && mc.matchType != ipb.ContentEntryCheck_NONE_MATCH {
			fc.addNonCompliantFile(filePath, fmt.Sprintf("File contains entry %q, expected %q", fc.redactContent(entry, filePath), mc))
		}
		if satisfiesCriterion {
			switch mc.matchType {
			case ipb.ContentEntryCheck_NONE_MATCH:
				fc.addNonCompliantFile(filePath, fmt.Sprintf("File contains entry %q, didn't expect any entries matching %q", fc.redactContent(entry, filePath), mc))
			case ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER:
				// Match was expected
			case ipb.ContentEntryCheck_ALL_MATCH_STRICT_ORDER:
				verifyCriterionMatchInStrictOrder(entry, filePath, fc, matchCriteria, i)
			default:
				return fmt.Errorf("unexpected match type %s", mc.matchType)
			}
		}
	}
	return nil
}

func verifyCriterionMatchInStrictOrder(entry string, filePath string, fc *fileCheck, matchCriteria []*matchCriterion, criterionPos int) {
	if len(fc.nonCompliantFiles) > 0 {
		// Avoid duplicate non-compliance messages about out-of-order matches.
		return
	}
	mc := matchCriteria[criterionPos]
	var prev, next *matchCriterion
	if criterionPos > 0 {
		prev = matchCriteria[criterionPos-1]
	}
	if criterionPos < len(matchCriteria)-1 {
		next = matchCriteria[criterionPos+1]
	}

	if prev != nil && !prev.matched {
		fc.addNonCompliantFile(filePath, fmt.Sprintf("Criteria expected to match in order but file entry %q, matched %q before %q was matched", fc.redactContent(entry, filePath), mc, prev))
	} else if next != nil && next.matched {
		fc.addNonCompliantFile(filePath, fmt.Sprintf("Criteria expected to match in order but file entry %q, matched %q after %q was matched", fc.redactContent(entry, filePath), mc, next))
	}
}
