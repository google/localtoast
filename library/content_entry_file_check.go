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
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/localtoast/library/fileset"
	ipb "github.com/google/localtoast/library/proto/scan_instructions_go_proto"
)

// contentEntryFileCheckBatch performs a series of checks about whether files have specific
// entries in their content.
type contentEntryFileCheckBatch struct {
	ctx                    context.Context
	fileChecks             []*fileCheck
	filesToCheck           *ipb.FileSet
	fs                     FileSystemReader
	contentEntryFileChecks []*contentEntryFileCheck
	delimiter              byte
}

type contentEntryFileCheck struct {
	fileCheck     *fileCheck
	matchCriteria []*matchCriterion
}

type matchCriterion struct {
	matchType     ipb.ContentEntryCheck_MatchType
	filterRegex   *regexp.Regexp
	expectedRegex *regexp.Regexp
	groupCriteria *groupCriteria
	matched       bool
}

func (m *matchCriterion) String() string {
	if gcStr := m.groupCriteria.String(); gcStr != "" {
		return fmt.Sprintf("%s with group criteria %s", m.expectedRegex, gcStr)
	}
	return m.expectedRegex.String()
}

func newContentEntryFileCheckBatch(ctx context.Context, fileChecks []*fileCheck, filesToCheck *ipb.FileSet, fs FileSystemReader) (*contentEntryFileCheckBatch, error) {
	if len(fileChecks) == 0 {
		return nil, errors.New("attempted to create content entry check batch without any file checks")
	}
	delimiter := fileChecks[0].checkInstruction.GetContentEntry().GetDelimiter()
	if len(delimiter) == 0 {
		// Split by lines if nothing else is specified.
		delimiter = []byte{'\n'}
	}
	if len(delimiter) > 1 {
		// TODO(b/181930060): Add support for multi-char delimiters.
		return nil, fmt.Errorf("invalid delimiter for content entry check: %v", delimiter)
	}

	contentEntryFileChecks := make([]*contentEntryFileCheck, 0, len(fileChecks))
	for _, fc := range fileChecks {
		matchCriteriaProtos := fc.checkInstruction.GetContentEntry().GetMatchCriteria()
		matchCriteria := make([]*matchCriterion, 0, len(matchCriteriaProtos))
		for _, mc := range matchCriteriaProtos {
			mode := "(?s)" // '.' matches '\n' too
			filterRegex, err := regexp.Compile(mode + "^" + mc.GetFilterRegex() + "$")
			if err != nil {
				return nil, err
			}
			expectedRegex, err := regexp.Compile(mode + "^" + mc.GetExpectedRegex() + "$")
			if err != nil {
				return nil, err
			}
			groupCriteria, err := newGroupCriteria(
				expectedRegex,
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
		contentEntryFileChecks = append(contentEntryFileChecks, &contentEntryFileCheck{
			fileCheck:     fc,
			matchCriteria: matchCriteria,
		})
	}

	return &contentEntryFileCheckBatch{
		ctx:                    ctx,
		fileChecks:             fileChecks,
		filesToCheck:           filesToCheck,
		fs:                     fs,
		contentEntryFileChecks: contentEntryFileChecks,
		// Delimiters are expected to be single-char for now.
		// TODO(b/181930060): Add support for multi-char delimiters.
		delimiter: delimiter[0],
	}, nil
}

func (c *contentEntryFileCheckBatch) exec() (ComplianceMap, error) {
	err := fileset.WalkFiles(c.ctx, c.filesToCheck, c.fs,
		func(path string, isDir bool) error {
			if isDir {
				return nil
			}

			exists, err := fileExists(c.ctx, path, c.fs)
			if err != nil {
				return err
			}
			if !exists {
				for _, fc := range c.fileChecks {
					if fc.checkInstruction.GetContentEntry().GetMatchType() != ipb.ContentEntryCheck_NONE_MATCH {
						fc.addNonCompliantFile(path, "File doesn't exist")
					}
				}
				return nil
			}

			f, err := openFileForReading(c.ctx, path, c.fs)
			if err != nil {
				return err
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			// Define a split function that splits the entries based on the delimiter char.
			scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
				if atEOF && len(data) == 0 {
					return 0, nil, nil
				}
				if i := bytes.IndexByte(data, c.delimiter); i >= 0 {
					return i + 1, data[0:i], nil
				}
				if atEOF {
					return len(data), data, nil
				}
				return 0, nil, nil
			})

			for scanner.Scan() {
				if scanner.Err() != nil {
					return scanner.Err()
				}
				entry := scanner.Text()
				for _, fc := range c.contentEntryFileChecks {
					if err := matchEntryAgainstCriteria(entry, path, fc); err != nil {
						return err
					}
				}
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	// Check for any remaining unmatched criteria
	for _, fc := range c.contentEntryFileChecks {
		for _, mc := range fc.matchCriteria {
			if !mc.matched && mc.matchType != ipb.ContentEntryCheck_NONE_MATCH {
				fc.fileCheck.addNonCompliantFile(
					fileset.FileSetToString(c.filesToCheck),
					fmt.Sprintf("No entry matching %q found among files", mc.filterRegex))
			}
		}
	}

	return aggregateComplianceResults(c.fileChecks)
}

func matchEntryAgainstCriteria(entry string, filePath string, check *contentEntryFileCheck) error {
	for i, mc := range check.matchCriteria {
		if !mc.filterRegex.MatchString(entry) {
			continue
		}

		mc.matched = true
		satisfiesCriterion := mc.expectedRegex.MatchString(entry) && mc.groupCriteria.check(entry)

		if !satisfiesCriterion && mc.matchType != ipb.ContentEntryCheck_NONE_MATCH {
			check.fileCheck.addNonCompliantFile(filePath,
				fmt.Sprintf("File contains entry %q, expected %q", entry, mc))
		}
		if satisfiesCriterion {
			switch mc.matchType {
			case ipb.ContentEntryCheck_NONE_MATCH:
				check.fileCheck.addNonCompliantFile(filePath,
					fmt.Sprintf("File contains entry %q, didn't expect any entries matching %q", entry, mc))
			case ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER:
				// Match was expected
			case ipb.ContentEntryCheck_ALL_MATCH_STRICT_ORDER:
				verifyCriterionMatchInStrictOrder(entry, filePath, check, i)
			default:
				return fmt.Errorf("unexpected match type %s", mc.matchType)
			}
		}
	}
	return nil
}

func verifyCriterionMatchInStrictOrder(entry string, filePath string, check *contentEntryFileCheck, criterionPos int) {
	if len(check.fileCheck.nonCompliantFiles) > 0 {
		// Avoid duplicate non-compliance messages about out-of-order matches.
		return
	}
	mc := check.matchCriteria[criterionPos]
	var prev, next *matchCriterion
	if criterionPos > 0 {
		prev = check.matchCriteria[criterionPos-1]
	}
	if criterionPos < len(check.matchCriteria)-1 {
		next = check.matchCriteria[criterionPos+1]
	}

	if prev != nil && !prev.matched {
		check.fileCheck.addNonCompliantFile(filePath,
			fmt.Sprintf("Criteria expected to match in order but file entry %q, matched %q before %q was matched", entry, mc, prev))
	} else if next != nil && next.matched {
		check.fileCheck.addNonCompliantFile(filePath,
			fmt.Sprintf("Criteria expected to match in order but file entry %q, matched %q after %q was matched", entry, mc, next))
	}
}

func (c *contentEntryFileCheckBatch) String() string {
	return fmt.Sprintf("[content entry check on %s]", fileset.FileSetToString(c.filesToCheck))
}
