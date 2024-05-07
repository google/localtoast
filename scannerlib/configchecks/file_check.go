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

// Package configchecks provides the implementation of the various checks the scanner can perform
// behind a single general interface.
package configchecks

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"google.golang.org/protobuf/proto"
	cpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	"github.com/google/localtoast/scanapi"
	"github.com/google/localtoast/scannerlib/fileset"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
	"github.com/google/localtoast/scannerlib/repeatconfig"
)

// MaxNonCompliantFiles is the maximum number of non-compliant files to be displayed for a single finding.
const MaxNonCompliantFiles = 10

// FileCheckBatch is an implementation of scanner.Check that performs various
// combined file checks on a set of files. By batching the checks together, we
// can avoid opening and reading the files multiple times.
type FileCheckBatch struct {
	fileChecks   []*fileCheck
	fileCheckers *fileCheckers
	ctx          context.Context
	filesToCheck *ipb.FileSet
	timeout      *timeoutOptions
	fs           scanapi.Filesystem
	benchmarkIDs []string
}

// Exec executes the file checks batched by the FileCheckBatch.
// The method takes as input the Previous check result output as string, if any
func (b *FileCheckBatch) Exec(prvRes string) (ComplianceMap, string, error) {

	fileset.ApplyPipelineTokenReplacement(b.filesToCheck, prvRes)

	err := fileset.WalkFiles(b.ctx, b.filesToCheck, b.fs, b.timeout.benchmarkCheckTimeoutNow(),
		func(path string, isDir bool, traversingDir bool) error {
			return b.fileCheckers.execChecksOnFile(b.ctx, path, isDir, traversingDir, b.fs)
		})
	if err != nil {
		return nil, "", err
	}
	b.fileCheckers.execChecksAfterFileTraversal(b.filesToCheck)
	return aggregateComplianceResults(b.fileChecks)
}

// BenchmarkIDs returns the IDs of the benchmarks associated with this check.
func (b *FileCheckBatch) BenchmarkIDs() []string {
	return b.benchmarkIDs
}

func (b *FileCheckBatch) String() string {
	return fmt.Sprintf("[file check on %s]", fileset.FileSetToString(b.filesToCheck))
}

// fileCheck is a single file check in the batch.
type fileCheck struct {
	benchmarkID           string
	alternativeID         int
	checkInstruction      *ipb.FileCheck
	filesToCheck          *ipb.FileSet
	contentOptoutRegexes  []*regexp.Regexp
	filenameOptoutRegexes []*regexp.Regexp
	nonCompliantFiles     []*cpb.NonCompliantFile
	err                   error
}

func (fc *fileCheck) addNonCompliantFile(path string, reason string) {
	fc.nonCompliantFiles = append(fc.nonCompliantFiles, &cpb.NonCompliantFile{Path: fc.redactPath(path), Reason: reason})
}

// redactContent redacts a given file content if the file path is in the list of opt-out regexes.
func (fc *fileCheck) redactContent(content, path string) string {
	for _, re := range fc.contentOptoutRegexes {
		if re.MatchString(path) {
			return "[redacted due to opt-out config]"
		}
	}
	return content
}

// redactPath redacts a given file path if it's in the list of opt-out regexes.
func (fc *fileCheck) redactPath(path string) string {
	for _, re := range fc.filenameOptoutRegexes {
		if re.MatchString(path) {
			return "[redacted due to opt-out config]"
		}
	}
	return path
}

// File checks are placed into the same batch if they have these properties in common.
type fileCheckBatchCommonProps struct {
	filesToCheck string // The set of files to check.
}
type fileCheckBatchMap map[fileCheckBatchCommonProps][]*fileCheck

// createFileCheckBatchesFromConfig parses the benchmark config and creates the
// file check batches defined by it.
func createFileCheckBatchesFromConfig(
	ctx context.Context, benchmarks []*benchmark, optOut *apb.OptOutConfig, replacement *apb.ReplacementConfig, timeout *timeoutOptions, fs scanapi.Filesystem) ([]*FileCheckBatch, error) {
	batchMap := make(fileCheckBatchMap)

	for _, b := range benchmarks {
		for _, alt := range b.alts {
			for _, fileCheckInstruction := range alt.proto.GetFileChecks() {
				if err := validateFileCheckInstruction(fileCheckInstruction); err != nil {
					return nil, err
				}
				options := addFileCheckToBatchMapOptions{
					fileCheckInstruction,
					batchMap,
					fs,
					optOut,
					replacement,
					b.id,
					alt.id,
				}
				if err := addFileCheckToBatchMap(ctx, options); err != nil {
					return nil, err
				}
			}
		}
	}

	fileCheckBatches := make([]*FileCheckBatch, 0, len(batchMap))
	for _, fileChecks := range batchMap {
		batch, err := newFileCheckBatch(ctx, fileChecks, fileChecks[0].filesToCheck, timeout, fs)
		if err != nil {
			return nil, err
		}
		fileCheckBatches = append(fileCheckBatches, batch)
	}
	return fileCheckBatches, nil
}

func validateFileCheckInstruction(instruction *ipb.FileCheck) error {
	if len(instruction.GetFileDisplayCommand()) > 0 && len(instruction.GetNonComplianceMsg()) == 0 {
		return fmt.Errorf("check instruction %v has a file display command set but no non-compliance message", instruction)
	}
	return nil
}

type addFileCheckToBatchMapOptions struct {
	fc            *ipb.FileCheck
	batchMap      fileCheckBatchMap
	fs            scanapi.Filesystem
	optOut        *apb.OptOutConfig
	replacement   *apb.ReplacementConfig
	benchmarkID   string
	alternativeID int
}

func addFileCheckToBatchMap(ctx context.Context, options addFileCheckToBatchMapOptions) error {
	repeatConfigs, err := repeatconfig.CreateRepeatConfigs(ctx, options.fc.GetRepeatConfig(), options.fs)
	if err != nil {
		return err
	}
	for _, repeatConfig := range repeatConfigs {
		fc := repeatconfig.ApplyRepeatConfigToInstruction(options.fc, repeatConfig)
		for _, filesToCheck := range fc.GetFilesToCheck() {
			filesToCheck := repeatconfig.ApplyRepeatConfigToFile(filesToCheck, repeatConfig)
			fileset.ApplyOptOutConfig(filesToCheck, options.optOut)
			fileset.ApplyReplacementConfig(filesToCheck, options.replacement)

			fileSetAsBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(filesToCheck)
			if err != nil {
				return err
			}
			key := fileCheckBatchCommonProps{
				filesToCheck: string(fileSetAsBytes),
			}
			contentOptoutRegexes, err := strToRegex(options.optOut.GetContentOptoutRegexes())
			if err != nil {
				return err
			}
			filenameOptoutRegexes, err := strToRegex(options.optOut.GetFilenameOptoutRegexes())
			if err != nil {
				return err
			}
			options.batchMap[key] = append(options.batchMap[key],
				&fileCheck{
					benchmarkID:           options.benchmarkID,
					alternativeID:         options.alternativeID,
					checkInstruction:      fc,
					filesToCheck:          filesToCheck,
					contentOptoutRegexes:  contentOptoutRegexes,
					filenameOptoutRegexes: filenameOptoutRegexes,
					err:                   repeatConfig.Err,
				})
		}
	}
	return nil
}

type fileCheckers struct {
	existenceFileCheckers    []*existenceFileChecker
	permissionFileCheckers   []*permissionFileChecker
	contentFileCheckers      []*contentFileChecker
	contentEntryFileCheckers []*contentEntryFileChecker
}

func newFileCheckers(fileChecks []*fileCheck) (*fileCheckers, error) {
	result := &fileCheckers{}
	for _, fc := range fileChecks {
		if fc.err != nil { // The check couldn't properly be created because of an error.
			continue
		}
		if fc.checkInstruction.GetExistence() != nil {
			result.existenceFileCheckers = append(result.existenceFileCheckers, newExistenceFileChecker(fc))
		} else if fc.checkInstruction.GetPermission() != nil {
			result.permissionFileCheckers = append(result.permissionFileCheckers, &permissionFileChecker{fc: fc})
		} else if fc.checkInstruction.GetContent() != nil {
			result.contentFileCheckers = append(result.contentFileCheckers, &contentFileChecker{fc: fc})
		} else if fc.checkInstruction.GetContentEntry() != nil {
			checker, err := newContentEntryFileChecker(fc)
			if err != nil {
				return nil, err
			}
			result.contentEntryFileCheckers = append(result.contentEntryFileCheckers, checker)
		} else {
			return nil, fmt.Errorf("Received FileCheck with unexpected type: %v", fc.checkInstruction)
		}
	}

	// Check for invalid configurations.
	if len(result.contentEntryFileCheckers) > 0 {
		if len(result.contentFileCheckers) > 0 {
			return nil, fmt.Errorf("file %v has both content and content entry file checks", fileChecks[0].filesToCheck)
		}
		delimiter := result.contentEntryFileCheckers[0].delimiter
		for _, checker := range result.contentEntryFileCheckers {
			if string(delimiter) != string(checker.delimiter) {
				return nil, fmt.Errorf("file %v has content entry checks with different delimiters: %s %v",
					fileChecks[0].filesToCheck, delimiter, checker.fc.checkInstruction)
			}
		}
	}
	return result, nil
}

func (c *fileCheckers) execChecksOnFile(ctx context.Context, path string, isDir bool, traversingDir bool, fs scanapi.Filesystem) error {
	f, openError := c.openFileForCheckExec(ctx, path, fs)
	if f != nil {
		defer f.Close()
	}

	for _, checker := range c.existenceFileCheckers {
		if err := checker.exec(path, openError, isDir, traversingDir); err != nil {
			return err
		}
	}
	for _, checker := range c.permissionFileCheckers {
		if err := checker.exec(ctx, path, fs, openError); err != nil {
			return err
		}
	}
	if err := c.execContentChecksOnFile(path, openError, f); err != nil {
		return err
	}
	if err := c.execContentEntryChecksOnFile(path, openError, f); err != nil {
		return err
	}
	return nil
}

func (c *fileCheckers) execChecksAfterFileTraversal(filesToCheck *ipb.FileSet) {
	for _, checker := range c.existenceFileCheckers {
		checker.execAfterFileTraversal()
	}
	c.execContentEntryChecksAfterFileTraversal(filesToCheck)
}

func (c *fileCheckers) openFileForCheckExec(ctx context.Context, path string, fs scanapi.Filesystem) (io.ReadCloser, error) {
	if len(c.contentFileCheckers) == 0 && len(c.contentEntryFileCheckers) == 0 {
		// We won't read the file, we only care about whether it could successfully be opened.
		f, openError := fs.OpenFile(ctx, path)
		if f != nil {
			f.Close()
		}
		return nil, openError
	}
	return openFileForReading(ctx, path, fs)
}

func strToRegex(strs []string) ([]*regexp.Regexp, error) {
	result := make([]*regexp.Regexp, 0, len(strs))
	for _, s := range strs {
		re, err := regexp.Compile("^" + s + "$")
		if err != nil {
			return nil, err
		}
		result = append(result, re)
	}
	return result, nil
}

// newFileCheckBatch creates a FileCheckBatch from fileChecks that perform checks on the same files.
func newFileCheckBatch(
	ctx context.Context, fileChecks []*fileCheck, filesToCheck *ipb.FileSet, timeout *timeoutOptions, fs scanapi.Filesystem) (*FileCheckBatch, error) {
	// De-duplicate the benchmark IDs.
	benchmarkIDMap := make(map[string]bool)
	for _, fc := range fileChecks {
		benchmarkIDMap[fc.benchmarkID] = true
	}
	benchmarkIDs := make([]string, 0, len(benchmarkIDMap))
	for id := range benchmarkIDMap {
		benchmarkIDs = append(benchmarkIDs, id)
	}

	fileCheckers, err := newFileCheckers(fileChecks)
	if err != nil {
		return nil, err
	}

	return &FileCheckBatch{
		fileChecks:   fileChecks,
		fileCheckers: fileCheckers,
		ctx:          ctx,
		filesToCheck: filesToCheck,
		timeout:      timeout,
		fs:           fs,
		benchmarkIDs: benchmarkIDs,
	}, nil
}

// existenceFileChecker performs checks about whether files exist or not.
type existenceFileChecker struct {
	fc        *fileCheck
	foundFile string
}

func newExistenceFileChecker(fc *fileCheck) *existenceFileChecker {
	return &existenceFileChecker{fc: fc, foundFile: ""}
}

func (c *existenceFileChecker) exec(path string, openError error, isDir bool, traversingDir bool) error {
	exists := false
	// If this file was listed while traversing a directory
	// we know it exists without needing to open it.
	if traversingDir && !isDir {
		exists = true
	} else {
		var err error
		exists, err = fileExists(openError)
		if err != nil {
			return err
		}
	}
	if exists {
		c.foundFile = path
	}
	return nil
}

func (c *existenceFileChecker) execAfterFileTraversal() {
	se := c.fc.checkInstruction.GetExistence().GetShouldExist()
	switch {
	case c.foundFile == "" && se:
		c.fc.addNonCompliantFile(fileset.FileSetToString(c.fc.filesToCheck), "File doesn't exist but it should")
	case c.foundFile != "" && !se:
		c.fc.addNonCompliantFile(c.foundFile, "File exists but it shouldn't")
	}
}

// permissionFileChecker performs a checks on the permissions of files.
type permissionFileChecker struct {
	fc *fileCheck
}

func (c *permissionFileChecker) exec(ctx context.Context, path string, fs scanapi.Filesystem, openError error) error {
	perms, err := fs.FilePermissions(ctx, path)
	if err != nil {
		// Return a non-compliance instead of an error if if the file doesn't exist.
		if exists, err := fileExists(openError); err == nil && !exists {
			c.fc.addNonCompliantFile(path, "File doesn't exist")
			return nil
		}
		return err
	}
	pc := c.fc.checkInstruction.GetPermission()
	wantSetBits := pc.GetSetBits() > 0
	wantClearBits := pc.GetClearBits() > 0
	if wantSetBits || wantClearBits {
		correctBitsSet := perms.GetPermissionNum()&pc.GetSetBits() == pc.GetSetBits()
		correctBitsClear := (^perms.GetPermissionNum())&pc.GetClearBits() == pc.GetClearBits()
		var matches bool
		if wantSetBits && wantClearBits {
			switch pc.GetBitsShouldMatch() {
			case ipb.PermissionCheck_BOTH_SET_AND_CLEAR:
				matches = correctBitsSet && correctBitsClear
			case ipb.PermissionCheck_EITHER_SET_OR_CLEAR:
				matches = correctBitsSet || correctBitsClear
			default:
				return fmt.Errorf("invalid BitMatchCriterion enum value: %v", pc.GetBitsShouldMatch())
			}
		} else if wantSetBits {
			matches = correctBitsSet
		} else { // wantClearBits
			matches = correctBitsClear
		}
		if !matches {
			reason := fmt.Sprintf("File permission is %04o, expected ", perms.GetPermissionNum())
			if wantSetBits {
				reason += fmt.Sprintf("the following bits to be set: %04o", pc.GetSetBits())
			}
			if wantSetBits && wantClearBits {
				if pc.GetBitsShouldMatch() == ipb.PermissionCheck_BOTH_SET_AND_CLEAR {
					reason += " and "
				} else { // EITHER_SET_OR_CLEAR
					reason += " or "
				}
			}
			if wantClearBits {
				reason += fmt.Sprintf("the following bits to be clear: %04o", pc.GetClearBits())
			}
			c.fc.addNonCompliantFile(path, reason)
		}
	}

	if pc.GetUser() != nil {
		if pc.GetUser().GetShouldOwn() && perms.GetUser() != pc.GetUser().GetName() {
			c.fc.addNonCompliantFile(
				path, fmt.Sprintf("Owner is %s, expected it to be %s",
					perms.GetUser(), pc.GetUser().GetName()))
		} else if !pc.GetUser().GetShouldOwn() && perms.GetUser() == pc.GetUser().GetName() {
			c.fc.addNonCompliantFile(
				path, fmt.Sprintf("Owner is %s, expected it to be a different user", perms.GetUser()))
		}
	}

	if pc.GetGroup() != nil {
		if pc.GetGroup().GetShouldOwn() && perms.GetGroup() != pc.GetGroup().GetName() {
			c.fc.addNonCompliantFile(
				path, fmt.Sprintf("Group is %s, expected it to be %s",
					perms.GetGroup(), pc.GetGroup().GetName()))
		} else if !pc.GetGroup().GetShouldOwn() && perms.GetGroup() == pc.GetGroup().GetName() {
			c.fc.addNonCompliantFile(
				path, fmt.Sprintf("Group is %s, expected it to be a different group", perms.GetGroup()))
		}
	}

	return nil
}

// contentFileChecker performs checks on the full content of files.
type contentFileChecker struct {
	fc *fileCheck
}

func (c *fileCheckers) execContentChecksOnFile(path string, openError error, f io.ReadCloser) error {
	if len(c.contentFileCheckers) == 0 {
		return nil
	}
	exists, err := fileExists(openError)
	if err != nil {
		return err
	}
	if !exists {
		for _, checker := range c.contentFileCheckers {
			checker.fc.addNonCompliantFile(path, "File doesn't exist")
		}
		return nil
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	for _, checker := range c.contentFileCheckers {
		if err := checker.exec(path, content); err != nil {
			return err
		}
	}
	return nil
}

func (c *contentFileChecker) exec(path string, content []byte) error {
	expectedContent := c.fc.checkInstruction.GetContent().GetContent()
	if string(content) != expectedContent {
		c.fc.addNonCompliantFile(path, fmt.Sprintf("Got content %q, expected %q",
			c.fc.redactContent(string(content), path), expectedContent))
	}
	return nil
}

// aggregateComplianceResults merges the non-compliant files of the
// specified fileChecks for each check alternative.
func aggregateComplianceResults(fileChecks []*fileCheck) (ComplianceMap, string, error) {
	result := make(map[int]*apb.ComplianceResult) // Key: The CheckAlternative ID.
	for _, fc := range fileChecks {
		nonCompliantFiles := fc.nonCompliantFiles

		// Report only the first N non-compliant files.
		if len(nonCompliantFiles) > MaxNonCompliantFiles {
			nonCompliantFiles = nonCompliantFiles[:MaxNonCompliantFiles]
		}
		// Replace the non-compliance reason and files with custom values if available.
		if len(nonCompliantFiles) > 0 {
			if len(fc.checkInstruction.GetFileDisplayCommand()) > 0 {
				nonCompliantFiles = []*cpb.NonCompliantFile{&cpb.NonCompliantFile{
					DisplayCommand: fc.checkInstruction.GetFileDisplayCommand(),
					Reason:         fc.checkInstruction.GetNonComplianceMsg(),
				}}
			} else if len(fc.checkInstruction.GetNonComplianceMsg()) > 0 {
				for _, f := range nonCompliantFiles {
					f.Reason = fc.checkInstruction.GetNonComplianceMsg()
				}
			}
		}

		if prev, ok := result[fc.alternativeID]; ok {
			prev.GetComplianceOccurrence().NonCompliantFiles =
				append(prev.GetComplianceOccurrence().NonCompliantFiles, nonCompliantFiles...)
			if fc.err != nil {
				prev.GetComplianceOccurrence().NonComplianceReason = fc.err.Error()
			}
		} else {
			result[fc.alternativeID] = &apb.ComplianceResult{
				Id: fc.benchmarkID,
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: nonCompliantFiles,
				},
			}
			if fc.err != nil {
				result[fc.alternativeID].GetComplianceOccurrence().NonComplianceReason = fc.err.Error()
			}
		}
	}

	return result, "", nil
}

// openFileForReading opens the specified path and returns a ReadCloser.
// If the file is a gzip file, the returned ReadCloser reads the unzipped
// contents of the file.
func openFileForReading(ctx context.Context, filePath string, fs scanapi.Filesystem) (io.ReadCloser, error) {
	reader, err := fs.OpenFile(ctx, filePath)
	if err != nil {
		return nil, err
	}
	if path.Ext(filePath) != ".gz" {
		return reader, nil
	}
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		reader.Close()
		return nil, err
	}
	return &gzipReadCloser{
		regularReader: reader,
		gzipReader:    gzipReader,
	}, nil
}

type gzipReadCloser struct {
	regularReader io.ReadCloser
	gzipReader    *gzip.Reader
}

func (r *gzipReadCloser) Close() error {
	err1 := r.gzipReader.Close()
	err2 := r.regularReader.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

func (r *gzipReadCloser) Read(p []byte) (n int, err error) {
	return r.gzipReader.Read(p)
}

func fileExists(err error) (bool, error) {
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	default:
		return false, err
	}
}
