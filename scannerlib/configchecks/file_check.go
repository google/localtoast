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
	"time"

	"google.golang.org/protobuf/proto"
	cpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	"github.com/google/localtoast/scannerlib/fileset"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
	"github.com/google/localtoast/scannerlib/repeatconfig"
)

// MaxNonCompliantFiles is the maximum number of non-compliant files to be displayed for a single finding.
const MaxNonCompliantFiles = 10

// FileSystemReader is an interface that gives the checkers read access to the
// filesystem of the scanned machine.
type FileSystemReader interface {
	OpenFile(ctx context.Context, path string) (io.ReadCloser, error)
	FilesInDir(ctx context.Context, path string) ([]*apb.DirContent, error)
	FilePermissions(ctx context.Context, path string) (*apb.PosixPermissions, error)
}

// FileCheckBatch is an implementation of scanner.Check that performs various
// combined file checks on a set of files. By batching the checks together, we
// can avoid opening and reading the files multiple times.
type FileCheckBatch struct {
	checker      fileCheckBatchChecker
	benchmarkIDs []string
}

// Exec executes the file checks defined by the FileCheckBatch.
func (b *FileCheckBatch) Exec() (ComplianceMap, error) {
	return b.checker.exec()
}

// BenchmarkIDs returns the IDs of the benchmarks associated with this check.
func (b *FileCheckBatch) BenchmarkIDs() []string {
	return b.benchmarkIDs
}

func (b *FileCheckBatch) String() string {
	return b.checker.String()
}

// fileCheckBatchChecker is an interface used by FileCheckBatch to perform a
// specific type of batched file checks on a set of files.
type fileCheckBatchChecker interface {
	exec() (ComplianceMap, error)
	String() string
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
	checkType    string // The type of file check to perform.
	delimiter    string // For ContentEntryChecks, the delimiter of the entries.
	err          string // Any error that occurred during check creation.
}
type fileCheckBatchMap map[fileCheckBatchCommonProps][]*fileCheck

// createFileCheckBatchesFromConfig parses the benchmark config and creates the
// file check batches defined by it.
func createFileCheckBatchesFromConfig(
	ctx context.Context, benchmarks []*benchmark, optOut *apb.OptOutConfig, timeout time.Time, fs FileSystemReader) ([]*FileCheckBatch, error) {
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
	fs            FileSystemReader
	optOut        *apb.OptOutConfig
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

			fileSetAsBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(filesToCheck)
			if err != nil {
				return err
			}
			delimiter := ""
			if fc.GetContentEntry() != nil {
				delimiter = string(fc.GetContentEntry().GetDelimiter())
			}
			checkType, err := checkTypeStr(fc)
			if err != nil {
				return err
			}
			errStr := ""
			if repeatConfig.Err != nil {
				errStr = repeatConfig.Err.Error()
			}
			key := fileCheckBatchCommonProps{
				filesToCheck: string(fileSetAsBytes),
				checkType:    checkType,
				delimiter:    delimiter,
				err:          errStr,
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

func checkTypeStr(fc *ipb.FileCheck) (string, error) {
	switch {
	case fc.GetExistence() != nil:
		return "EXISTENCE", nil
	case fc.GetPermission() != nil:
		return "PERMISSION", nil
	case fc.GetContent() != nil:
		return "CONTENT", nil
	case fc.GetContentEntry() != nil:
		return "CONTENT_ENTRY", nil
	default:
		return "", fmt.Errorf("unknown file check type for %v", fc)
	}
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

// newFileCheckBatch creates a FileCheckBatch from several fileChecks that
// perform the same type of checks on the same files.
func newFileCheckBatch(
	ctx context.Context, fileChecks []*fileCheck, filesToCheck *ipb.FileSet, timeout time.Time, fs FileSystemReader) (*FileCheckBatch, error) {
	// De-duplicate the benchmark IDs.
	benchmarkIDMap := make(map[string]bool)
	for _, fc := range fileChecks {
		benchmarkIDMap[fc.benchmarkID] = true
	}
	benchmarkIDs := make([]string, 0, len(benchmarkIDMap))
	for id, _ := range benchmarkIDMap {
		benchmarkIDs = append(benchmarkIDs, id)
	}

	// Create the file check implementation corresponding to the message type.
	var checker fileCheckBatchChecker
	var err error
	if len(fileChecks) == 0 {
		return nil, errors.New("Attempted to create a batch without any file checks")
	}

	if fileChecks[0].err != nil { // The checks couldn't properly be created because of an error.
		checker, err = newErroredFileCheckBatch(fileChecks, fileChecks[0].err)
	} else if fileChecks[0].checkInstruction.GetExistence() != nil {
		checker, err = newExistenceFileCheckBatch(ctx, fileChecks, filesToCheck, timeout, fs)
	} else if fileChecks[0].checkInstruction.GetPermission() != nil {
		checker, err = newPermissionFileCheckBatch(ctx, fileChecks, filesToCheck, timeout, fs)
	} else if fileChecks[0].checkInstruction.GetContent() != nil {
		checker, err = newContentFileCheckBatch(ctx, fileChecks, filesToCheck, timeout, fs)
	} else if fileChecks[0].checkInstruction.GetContentEntry() != nil {
		checker, err = newContentEntryFileCheckBatch(ctx, fileChecks, filesToCheck, timeout, fs)
	} else {
		return nil, fmt.Errorf("Received FileCheck with unexpected type: %v",
			fileChecks[0].checkInstruction)
	}
	if err != nil {
		return nil, err
	}

	return &FileCheckBatch{checker: checker, benchmarkIDs: benchmarkIDs}, nil
}

// existenceFileCheckBatch performs a series of checks about whether files exist or not.
type existenceFileCheckBatch struct {
	ctx          context.Context
	fileChecks   []*fileCheck
	filesToCheck *ipb.FileSet
	timeout      time.Time
	fs           FileSystemReader
	foundFile    string
}

func newExistenceFileCheckBatch(
	ctx context.Context,
	fileChecks []*fileCheck,
	filesToCheck *ipb.FileSet,
	timeout time.Time,
	fs FileSystemReader) (*existenceFileCheckBatch, error) {
	return &existenceFileCheckBatch{
		ctx:          ctx,
		fileChecks:   fileChecks,
		filesToCheck: filesToCheck,
		timeout:      timeout,
		fs:           fs,
		foundFile:    "",
	}, nil
}

func (c *existenceFileCheckBatch) exec() (ComplianceMap, error) {
	err := fileset.WalkFiles(c.ctx, c.filesToCheck, c.fs, c.timeout,
		func(path string, isDir bool) error {
			exists, err := fileExists(c.ctx, path, c.fs)
			if err != nil {
				return err
			}
			if exists {
				c.foundFile = path
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	for _, fc := range c.fileChecks {
		se := fc.checkInstruction.GetExistence().GetShouldExist()
		switch {
		case c.foundFile == "" && se:
			fc.addNonCompliantFile(fileset.FileSetToString(c.filesToCheck), "File doesn't exist but it should")
		case c.foundFile != "" && !se:
			fc.addNonCompliantFile(c.foundFile, "File exists but it shouldn't")
		}
	}

	return aggregateComplianceResults(c.fileChecks)
}

func (c *existenceFileCheckBatch) String() string {
	return fmt.Sprintf("[existence check on %s]", fileset.FileSetToString(c.filesToCheck))
}

// permissionFileCheckBatch performs a series of checks about the permissions of files.
type permissionFileCheckBatch struct {
	ctx          context.Context
	fileChecks   []*fileCheck
	filesToCheck *ipb.FileSet
	timeout      time.Time
	fs           FileSystemReader
}

func newPermissionFileCheckBatch(
	ctx context.Context,
	fileChecks []*fileCheck,
	filesToCheck *ipb.FileSet,
	timeout time.Time,
	fs FileSystemReader) (*permissionFileCheckBatch, error) {
	return &permissionFileCheckBatch{
		ctx:          ctx,
		fileChecks:   fileChecks,
		filesToCheck: filesToCheck,
		timeout:      timeout,
		fs:           fs,
	}, nil
}

func (c *permissionFileCheckBatch) exec() (ComplianceMap, error) {
	err := fileset.WalkFiles(c.ctx, c.filesToCheck, c.fs, c.timeout,
		func(path string, isDir bool) error {
			perms, err := c.fs.FilePermissions(c.ctx, path)
			if err != nil {
				// Return a non-compliance instead of an error if if the file doesn't exist.
				if exists, err := fileExists(c.ctx, path, c.fs); err == nil && !exists {
					for _, fc := range c.fileChecks {
						fc.addNonCompliantFile(path, "File doesn't exist")
					}
					return nil
				}
				return err
			}
			for _, fc := range c.fileChecks {
				pc := fc.checkInstruction.GetPermission()
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
						fc.addNonCompliantFile(path, reason)
					}
				}

				if pc.GetUser() != nil {
					if pc.GetUser().GetShouldOwn() && perms.GetUser() != pc.GetUser().GetName() {
						fc.addNonCompliantFile(
							path, fmt.Sprintf("Owner is %s, expected it to be %s",
								perms.GetUser(), pc.GetUser().GetName()))
					} else if !pc.GetUser().GetShouldOwn() && perms.GetUser() == pc.GetUser().GetName() {
						fc.addNonCompliantFile(
							path, fmt.Sprintf("Owner is %s, expected it to be a different user", perms.GetUser()))
					}
				}

				if pc.GetGroup() != nil {
					if pc.GetGroup().GetShouldOwn() && perms.GetGroup() != pc.GetGroup().GetName() {
						fc.addNonCompliantFile(
							path, fmt.Sprintf("Group is %s, expected it to be %s",
								perms.GetGroup(), pc.GetGroup().GetName()))
					} else if !pc.GetGroup().GetShouldOwn() && perms.GetGroup() == pc.GetGroup().GetName() {
						fc.addNonCompliantFile(
							path, fmt.Sprintf("Group is %s, expected it to be a different group", perms.GetGroup()))
					}
				}
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	return aggregateComplianceResults(c.fileChecks)
}

func (c *permissionFileCheckBatch) String() string {
	return fmt.Sprintf("[permission check on %s]", fileset.FileSetToString(c.filesToCheck))
}

// contentFileCheckBatch performs a series of checks about the full content of files.
type contentFileCheckBatch struct {
	ctx          context.Context
	fileChecks   []*fileCheck
	filesToCheck *ipb.FileSet
	timeout      time.Time
	fs           FileSystemReader
}

func newContentFileCheckBatch(
	ctx context.Context,
	fileChecks []*fileCheck,
	filesToCheck *ipb.FileSet,
	timeout time.Time,
	fs FileSystemReader) (*contentFileCheckBatch, error) {
	return &contentFileCheckBatch{
		ctx:          ctx,
		fileChecks:   fileChecks,
		filesToCheck: filesToCheck,
		timeout:      timeout,
		fs:           fs,
	}, nil
}

func (c *contentFileCheckBatch) exec() (ComplianceMap, error) {
	err := fileset.WalkFiles(c.ctx, c.filesToCheck, c.fs, c.timeout,
		func(path string, isDir bool) error {
			exists, err := fileExists(c.ctx, path, c.fs)
			if err != nil {
				return err
			}
			if !exists {
				for _, fc := range c.fileChecks {
					fc.addNonCompliantFile(path, "File doesn't exist")
				}
				return nil
			}

			f, err := openFileForReading(c.ctx, path, c.fs)
			if err != nil {
				return err
			}
			var content []byte
			content, err = ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			defer f.Close()

			for _, fc := range c.fileChecks {
				expectedContent := fc.checkInstruction.GetContent().GetContent()
				if string(content) != expectedContent {
					fc.addNonCompliantFile(path, fmt.Sprintf("Got content %q, expected %q",
						fc.redactContent(string(content), path), expectedContent))
				}
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	return aggregateComplianceResults(c.fileChecks)
}

func (c *contentFileCheckBatch) String() string {
	return fmt.Sprintf("[content check on %s]", fileset.FileSetToString(c.filesToCheck))
}

// erroredFileCheckBatch always returns non-compliant results with the specified error message.
type erroredFileCheckBatch struct {
	fileChecks   []*fileCheck
	filesToCheck *ipb.FileSet
	err          error
}

func newErroredFileCheckBatch(fileChecks []*fileCheck, err error) (*erroredFileCheckBatch, error) {
	return &erroredFileCheckBatch{
		fileChecks: fileChecks,
		err:        err,
	}, nil
}

func (e *erroredFileCheckBatch) exec() (ComplianceMap, error) {
	result := make(ComplianceMap)
	for _, fc := range e.fileChecks {
		result[fc.alternativeID] = &apb.ComplianceResult{
			Id: fc.benchmarkID,
			ComplianceOccurrence: &cpb.ComplianceOccurrence{
				NonComplianceReason: e.err.Error(),
			},
		}
	}
	return result, nil
}

func (erroredFileCheckBatch) String() string {
	return fmt.Sprintf("[errored check]")
}

// aggregateComplianceResults merges the non-compliant files of the
// specified fileChecks for each check alternative.
func aggregateComplianceResults(fileChecks []*fileCheck) (ComplianceMap, error) {
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
		} else {
			result[fc.alternativeID] = &apb.ComplianceResult{
				Id: fc.benchmarkID,
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: nonCompliantFiles,
				},
			}
		}
	}

	return result, nil
}

// openFileForReading opens the specified path and returns a ReadCloser.
// If the file is a gzip file, the returned ReadCloser reads the unzipped
// contents of the file.
func openFileForReading(ctx context.Context, filePath string, fs FileSystemReader) (io.ReadCloser, error) {
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

func fileExists(ctx context.Context, path string, fs FileSystemReader) (bool, error) {
	f, err := fs.OpenFile(ctx, path)
	switch {
	case err == nil:
		f.Close()
		return true, nil
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	default:
		return false, err
	}
}
