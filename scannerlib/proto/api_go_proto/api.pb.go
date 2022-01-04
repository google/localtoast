// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License"); you
// may not use this file except in compliance with the License.  You
// may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.  See the License for the specific language governing
// permissions and limitations under the License.

// Generated from scannerlib/proto/api using "bazel build"
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.14.0

package api_go_proto

import (
	reflect "reflect"
	sync "sync"

	compliance_go_proto "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ScanStatus_ScanStatusEnum int32

const (
	ScanStatus_UNSPECIFIED ScanStatus_ScanStatusEnum = 0
	ScanStatus_FAILED      ScanStatus_ScanStatusEnum = 1
	ScanStatus_SUCCEEDED   ScanStatus_ScanStatusEnum = 2
)

// Enum value maps for ScanStatus_ScanStatusEnum.
var (
	ScanStatus_ScanStatusEnum_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "FAILED",
		2: "SUCCEEDED",
	}
	ScanStatus_ScanStatusEnum_value = map[string]int32{
		"UNSPECIFIED": 0,
		"FAILED":      1,
		"SUCCEEDED":   2,
	}
)

func (x ScanStatus_ScanStatusEnum) Enum() *ScanStatus_ScanStatusEnum {
	p := new(ScanStatus_ScanStatusEnum)
	*p = x
	return p
}

func (x ScanStatus_ScanStatusEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ScanStatus_ScanStatusEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_scannerlib_proto_api_proto_enumTypes[0].Descriptor()
}

func (ScanStatus_ScanStatusEnum) Type() protoreflect.EnumType {
	return &file_scannerlib_proto_api_proto_enumTypes[0]
}

func (x ScanStatus_ScanStatusEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ScanStatus_ScanStatusEnum.Descriptor instead.
func (ScanStatus_ScanStatusEnum) EnumDescriptor() ([]byte, []int) {
	return file_scannerlib_proto_api_proto_rawDescGZIP(), []int{4, 0}
}

type ScanConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ScanTimeout           *durationpb.Duration `protobuf:"bytes,1,opt,name=scan_timeout,json=scanTimeout,proto3" json:"scan_timeout,omitempty"`
	BenchmarkCheckTimeout *durationpb.Duration `protobuf:"bytes,2,opt,name=benchmark_check_timeout,json=benchmarkCheckTimeout,proto3" json:"benchmark_check_timeout,omitempty"`
	OptOutConfig          *OptOutConfig        `protobuf:"bytes,3,opt,name=opt_out_config,json=optOutConfig,proto3" json:"opt_out_config,omitempty"`
	BenchmarkConfigs      []*BenchmarkConfig   `protobuf:"bytes,4,rep,name=benchmark_configs,json=benchmarkConfigs,proto3" json:"benchmark_configs,omitempty"`
}

func (x *ScanConfig) Reset() {
	*x = ScanConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanConfig) ProtoMessage() {}

func (x *ScanConfig) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanConfig.ProtoReflect.Descriptor instead.
func (*ScanConfig) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_api_proto_rawDescGZIP(), []int{0}
}

func (x *ScanConfig) GetScanTimeout() *durationpb.Duration {
	if x != nil {
		return x.ScanTimeout
	}
	return nil
}

func (x *ScanConfig) GetBenchmarkCheckTimeout() *durationpb.Duration {
	if x != nil {
		return x.BenchmarkCheckTimeout
	}
	return nil
}

func (x *ScanConfig) GetOptOutConfig() *OptOutConfig {
	if x != nil {
		return x.OptOutConfig
	}
	return nil
}

func (x *ScanConfig) GetBenchmarkConfigs() []*BenchmarkConfig {
	if x != nil {
		return x.BenchmarkConfigs
	}
	return nil
}

type OptOutConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContentOptoutRegexes   []string `protobuf:"bytes,1,rep,name=content_optout_regexes,json=contentOptoutRegexes,proto3" json:"content_optout_regexes,omitempty"`
	FilenameOptoutRegexes  []string `protobuf:"bytes,2,rep,name=filename_optout_regexes,json=filenameOptoutRegexes,proto3" json:"filename_optout_regexes,omitempty"`
	TraversalOptoutRegexes []string `protobuf:"bytes,3,rep,name=traversal_optout_regexes,json=traversalOptoutRegexes,proto3" json:"traversal_optout_regexes,omitempty"`
}

func (x *OptOutConfig) Reset() {
	*x = OptOutConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OptOutConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OptOutConfig) ProtoMessage() {}

func (x *OptOutConfig) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OptOutConfig.ProtoReflect.Descriptor instead.
func (*OptOutConfig) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_api_proto_rawDescGZIP(), []int{1}
}

func (x *OptOutConfig) GetContentOptoutRegexes() []string {
	if x != nil {
		return x.ContentOptoutRegexes
	}
	return nil
}

func (x *OptOutConfig) GetFilenameOptoutRegexes() []string {
	if x != nil {
		return x.FilenameOptoutRegexes
	}
	return nil
}

func (x *OptOutConfig) GetTraversalOptoutRegexes() []string {
	if x != nil {
		return x.TraversalOptoutRegexes
	}
	return nil
}

type BenchmarkConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                              `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ComplianceNote *compliance_go_proto.ComplianceNote `protobuf:"bytes,2,opt,name=compliance_note,json=complianceNote,proto3" json:"compliance_note,omitempty"`
}

func (x *BenchmarkConfig) Reset() {
	*x = BenchmarkConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BenchmarkConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BenchmarkConfig) ProtoMessage() {}

func (x *BenchmarkConfig) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BenchmarkConfig.ProtoReflect.Descriptor instead.
func (*BenchmarkConfig) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_api_proto_rawDescGZIP(), []int{2}
}

func (x *BenchmarkConfig) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BenchmarkConfig) GetComplianceNote() *compliance_go_proto.ComplianceNote {
	if x != nil {
		return x.ComplianceNote
	}
	return nil
}

type ScanResults struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTime              *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime                *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	ScannerVersion         string                 `protobuf:"bytes,3,opt,name=scanner_version,json=scannerVersion,proto3" json:"scanner_version,omitempty"`
	BenchmarkVersion       string                 `protobuf:"bytes,4,opt,name=benchmark_version,json=benchmarkVersion,proto3" json:"benchmark_version,omitempty"`
	Status                 *ScanStatus            `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	CompliantBenchmarks    []*ComplianceResult    `protobuf:"bytes,6,rep,name=compliant_benchmarks,json=compliantBenchmarks,proto3" json:"compliant_benchmarks,omitempty"`
	NonCompliantBenchmarks []*ComplianceResult    `protobuf:"bytes,7,rep,name=non_compliant_benchmarks,json=nonCompliantBenchmarks,proto3" json:"non_compliant_benchmarks,omitempty"`
}

func (x *ScanResults) Reset() {
	*x = ScanResults{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanResults) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanResults) ProtoMessage() {}

func (x *ScanResults) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanResults.ProtoReflect.Descriptor instead.
func (*ScanResults) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_api_proto_rawDescGZIP(), []int{3}
}

func (x *ScanResults) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *ScanResults) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

func (x *ScanResults) GetScannerVersion() string {
	if x != nil {
		return x.ScannerVersion
	}
	return ""
}

func (x *ScanResults) GetBenchmarkVersion() string {
	if x != nil {
		return x.BenchmarkVersion
	}
	return ""
}

func (x *ScanResults) GetStatus() *ScanStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *ScanResults) GetCompliantBenchmarks() []*ComplianceResult {
	if x != nil {
		return x.CompliantBenchmarks
	}
	return nil
}

func (x *ScanResults) GetNonCompliantBenchmarks() []*ComplianceResult {
	if x != nil {
		return x.NonCompliantBenchmarks
	}
	return nil
}

type ScanStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status        ScanStatus_ScanStatusEnum `protobuf:"varint,1,opt,name=status,proto3,enum=localtoast.ScanStatus_ScanStatusEnum" json:"status,omitempty"`
	FailureReason string                    `protobuf:"bytes,2,opt,name=failure_reason,json=failureReason,proto3" json:"failure_reason,omitempty"`
}

func (x *ScanStatus) Reset() {
	*x = ScanStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanStatus) ProtoMessage() {}

func (x *ScanStatus) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanStatus.ProtoReflect.Descriptor instead.
func (*ScanStatus) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_api_proto_rawDescGZIP(), []int{4}
}

func (x *ScanStatus) GetStatus() ScanStatus_ScanStatusEnum {
	if x != nil {
		return x.Status
	}
	return ScanStatus_UNSPECIFIED
}

func (x *ScanStatus) GetFailureReason() string {
	if x != nil {
		return x.FailureReason
	}
	return ""
}

type ComplianceResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                   string                                    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ComplianceOccurrence *compliance_go_proto.ComplianceOccurrence `protobuf:"bytes,3,opt,name=compliance_occurrence,json=complianceOccurrence,proto3" json:"compliance_occurrence,omitempty"`
}

func (x *ComplianceResult) Reset() {
	*x = ComplianceResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ComplianceResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComplianceResult) ProtoMessage() {}

func (x *ComplianceResult) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComplianceResult.ProtoReflect.Descriptor instead.
func (*ComplianceResult) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_api_proto_rawDescGZIP(), []int{5}
}

func (x *ComplianceResult) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ComplianceResult) GetComplianceOccurrence() *compliance_go_proto.ComplianceOccurrence {
	if x != nil {
		return x.ComplianceOccurrence
	}
	return nil
}

type DirContent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	IsDir     bool   `protobuf:"varint,2,opt,name=is_dir,json=isDir,proto3" json:"is_dir,omitempty"`
	IsSymlink bool   `protobuf:"varint,3,opt,name=is_symlink,json=isSymlink,proto3" json:"is_symlink,omitempty"`
}

func (x *DirContent) Reset() {
	*x = DirContent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DirContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DirContent) ProtoMessage() {}

func (x *DirContent) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DirContent.ProtoReflect.Descriptor instead.
func (*DirContent) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_api_proto_rawDescGZIP(), []int{6}
}

func (x *DirContent) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DirContent) GetIsDir() bool {
	if x != nil {
		return x.IsDir
	}
	return false
}

func (x *DirContent) GetIsSymlink() bool {
	if x != nil {
		return x.IsSymlink
	}
	return false
}

type PosixPermissions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PermissionNum int32  `protobuf:"varint,1,opt,name=permission_num,json=permissionNum,proto3" json:"permission_num,omitempty"`
	Uid           int32  `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	User          string `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	Gid           int32  `protobuf:"varint,4,opt,name=gid,proto3" json:"gid,omitempty"`
	Group         string `protobuf:"bytes,5,opt,name=group,proto3" json:"group,omitempty"`
}

func (x *PosixPermissions) Reset() {
	*x = PosixPermissions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PosixPermissions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PosixPermissions) ProtoMessage() {}

func (x *PosixPermissions) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PosixPermissions.ProtoReflect.Descriptor instead.
func (*PosixPermissions) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_api_proto_rawDescGZIP(), []int{7}
}

func (x *PosixPermissions) GetPermissionNum() int32 {
	if x != nil {
		return x.PermissionNum
	}
	return 0
}

func (x *PosixPermissions) GetUid() int32 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *PosixPermissions) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *PosixPermissions) GetGid() int32 {
	if x != nil {
		return x.Gid
	}
	return 0
}

func (x *PosixPermissions) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

type PerOsBenchmarkConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version     *compliance_go_proto.ComplianceVersion `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	BenchmarkId []string                               `protobuf:"bytes,2,rep,name=benchmark_id,json=benchmarkId,proto3" json:"benchmark_id,omitempty"`
}

func (x *PerOsBenchmarkConfig) Reset() {
	*x = PerOsBenchmarkConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_api_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PerOsBenchmarkConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PerOsBenchmarkConfig) ProtoMessage() {}

func (x *PerOsBenchmarkConfig) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_api_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PerOsBenchmarkConfig.ProtoReflect.Descriptor instead.
func (*PerOsBenchmarkConfig) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_api_proto_rawDescGZIP(), []int{8}
}

func (x *PerOsBenchmarkConfig) GetVersion() *compliance_go_proto.ComplianceVersion {
	if x != nil {
		return x.Version
	}
	return nil
}

func (x *PerOsBenchmarkConfig) GetBenchmarkId() []string {
	if x != nil {
		return x.BenchmarkId
	}
	return nil
}

var File_scannerlib_proto_api_proto protoreflect.FileDescriptor

var file_scannerlib_proto_api_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x6c, 0x69, 0x62, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x73, 0x63, 0x61, 0x6e, 0x6e,
	0x65, 0x72, 0x6c, 0x69, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x70,
	0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa7, 0x02, 0x0a,
	0x0a, 0x53, 0x63, 0x61, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3c, 0x0a, 0x0c, 0x73,
	0x63, 0x61, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x73, 0x63,
	0x61, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x51, 0x0a, 0x17, 0x62, 0x65, 0x6e,
	0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x15, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x3e, 0x0a, 0x0e,
	0x6f, 0x70, 0x74, 0x5f, 0x6f, 0x75, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73,
	0x74, 0x2e, 0x4f, 0x70, 0x74, 0x4f, 0x75, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0c,
	0x6f, 0x70, 0x74, 0x4f, 0x75, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x48, 0x0a, 0x11,
	0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74,
	0x6f, 0x61, 0x73, 0x74, 0x2e, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x52, 0x10, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x22, 0xb6, 0x01, 0x0a, 0x0c, 0x4f, 0x70, 0x74, 0x4f, 0x75,
	0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x34, 0x0a, 0x16, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x5f, 0x6f, 0x70, 0x74, 0x6f, 0x75, 0x74, 0x5f, 0x72, 0x65, 0x67, 0x65, 0x78, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x14, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x4f, 0x70, 0x74, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x67, 0x65, 0x78, 0x65, 0x73, 0x12, 0x36, 0x0a,
	0x17, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x6f, 0x70, 0x74, 0x6f, 0x75, 0x74,
	0x5f, 0x72, 0x65, 0x67, 0x65, 0x78, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x15,
	0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x4f, 0x70, 0x74, 0x6f, 0x75, 0x74, 0x52, 0x65,
	0x67, 0x65, 0x78, 0x65, 0x73, 0x12, 0x38, 0x0a, 0x18, 0x74, 0x72, 0x61, 0x76, 0x65, 0x72, 0x73,
	0x61, 0x6c, 0x5f, 0x6f, 0x70, 0x74, 0x6f, 0x75, 0x74, 0x5f, 0x72, 0x65, 0x67, 0x65, 0x78, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x16, 0x74, 0x72, 0x61, 0x76, 0x65, 0x72, 0x73,
	0x61, 0x6c, 0x4f, 0x70, 0x74, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x67, 0x65, 0x78, 0x65, 0x73, 0x22,
	0x66, 0x0a, 0x0f, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x43, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65,
	0x5f, 0x6e, 0x6f, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x72,
	0x61, 0x66, 0x65, 0x61, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61,
	0x6e, 0x63, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x52, 0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61,
	0x6e, 0x63, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x22, 0xae, 0x03, 0x0a, 0x0b, 0x53, 0x63, 0x61, 0x6e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x63, 0x61,
	0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x2b, 0x0a, 0x11, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x5f,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x62,
	0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x2e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x53, 0x63, 0x61,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x4f, 0x0a, 0x14, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x74, 0x5f, 0x62, 0x65, 0x6e,
	0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c,
	0x69, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x13, 0x63, 0x6f, 0x6d,
	0x70, 0x6c, 0x69, 0x61, 0x6e, 0x74, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x73,
	0x12, 0x56, 0x0a, 0x18, 0x6e, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e,
	0x74, 0x5f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x18, 0x07, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e,
	0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x52, 0x16, 0x6e, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x74, 0x42, 0x65,
	0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x22, 0xb0, 0x01, 0x0a, 0x0a, 0x53, 0x63, 0x61,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x3d, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74,
	0x6f, 0x61, 0x73, 0x74, 0x2e, 0x53, 0x63, 0x61, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e,
	0x53, 0x63, 0x61, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x45, 0x6e, 0x75, 0x6d, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72,
	0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x3c, 0x0a,
	0x0e, 0x53, 0x63, 0x61, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x45, 0x6e, 0x75, 0x6d, 0x12,
	0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09,
	0x53, 0x55, 0x43, 0x43, 0x45, 0x45, 0x44, 0x45, 0x44, 0x10, 0x02, 0x22, 0x7f, 0x0a, 0x10, 0x43,
	0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x55, 0x0a, 0x15, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x6f, 0x63,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x67, 0x72, 0x61, 0x66, 0x65, 0x61, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x4f, 0x63, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x65,
	0x52, 0x14, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x4f, 0x63, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x22, 0x56, 0x0a, 0x0a,
	0x44, 0x69, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x15,
	0x0a, 0x06, 0x69, 0x73, 0x5f, 0x64, 0x69, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05,
	0x69, 0x73, 0x44, 0x69, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x73, 0x79, 0x6d, 0x6c,
	0x69, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x53, 0x79, 0x6d,
	0x6c, 0x69, 0x6e, 0x6b, 0x22, 0x87, 0x01, 0x0a, 0x10, 0x50, 0x6f, 0x73, 0x69, 0x78, 0x50, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x65, 0x72,
	0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0d, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x4e, 0x75, 0x6d,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x75,
	0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x03, 0x67, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x72,
	0x0a, 0x14, 0x50, 0x65, 0x72, 0x4f, 0x73, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x37, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x72, 0x61, 0x66, 0x65, 0x61,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x21, 0x0a, 0x0c, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b,
	0x49, 0x64, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61,
	0x73, 0x74, 0x2f, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x6c, 0x69, 0x62, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x67, 0x6f, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_scannerlib_proto_api_proto_rawDescOnce sync.Once
	file_scannerlib_proto_api_proto_rawDescData = file_scannerlib_proto_api_proto_rawDesc
)

func file_scannerlib_proto_api_proto_rawDescGZIP() []byte {
	file_scannerlib_proto_api_proto_rawDescOnce.Do(func() {
		file_scannerlib_proto_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_scannerlib_proto_api_proto_rawDescData)
	})
	return file_scannerlib_proto_api_proto_rawDescData
}

var file_scannerlib_proto_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_scannerlib_proto_api_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_scannerlib_proto_api_proto_goTypes = []interface{}{
	(ScanStatus_ScanStatusEnum)(0),                   // 0: localtoast.ScanStatus.ScanStatusEnum
	(*ScanConfig)(nil),                               // 1: localtoast.ScanConfig
	(*OptOutConfig)(nil),                             // 2: localtoast.OptOutConfig
	(*BenchmarkConfig)(nil),                          // 3: localtoast.BenchmarkConfig
	(*ScanResults)(nil),                              // 4: localtoast.ScanResults
	(*ScanStatus)(nil),                               // 5: localtoast.ScanStatus
	(*ComplianceResult)(nil),                         // 6: localtoast.ComplianceResult
	(*DirContent)(nil),                               // 7: localtoast.DirContent
	(*PosixPermissions)(nil),                         // 8: localtoast.PosixPermissions
	(*PerOsBenchmarkConfig)(nil),                     // 9: localtoast.PerOsBenchmarkConfig
	(*durationpb.Duration)(nil),                      // 10: google.protobuf.Duration
	(*compliance_go_proto.ComplianceNote)(nil),       // 11: grafeas.v1.ComplianceNote
	(*timestamppb.Timestamp)(nil),                    // 12: google.protobuf.Timestamp
	(*compliance_go_proto.ComplianceOccurrence)(nil), // 13: grafeas.v1.ComplianceOccurrence
	(*compliance_go_proto.ComplianceVersion)(nil),    // 14: grafeas.v1.ComplianceVersion
}
var file_scannerlib_proto_api_proto_depIdxs = []int32{
	10, // 0: localtoast.ScanConfig.scan_timeout:type_name -> google.protobuf.Duration
	10, // 1: localtoast.ScanConfig.benchmark_check_timeout:type_name -> google.protobuf.Duration
	2,  // 2: localtoast.ScanConfig.opt_out_config:type_name -> localtoast.OptOutConfig
	3,  // 3: localtoast.ScanConfig.benchmark_configs:type_name -> localtoast.BenchmarkConfig
	11, // 4: localtoast.BenchmarkConfig.compliance_note:type_name -> grafeas.v1.ComplianceNote
	12, // 5: localtoast.ScanResults.start_time:type_name -> google.protobuf.Timestamp
	12, // 6: localtoast.ScanResults.end_time:type_name -> google.protobuf.Timestamp
	5,  // 7: localtoast.ScanResults.status:type_name -> localtoast.ScanStatus
	6,  // 8: localtoast.ScanResults.compliant_benchmarks:type_name -> localtoast.ComplianceResult
	6,  // 9: localtoast.ScanResults.non_compliant_benchmarks:type_name -> localtoast.ComplianceResult
	0,  // 10: localtoast.ScanStatus.status:type_name -> localtoast.ScanStatus.ScanStatusEnum
	13, // 11: localtoast.ComplianceResult.compliance_occurrence:type_name -> grafeas.v1.ComplianceOccurrence
	14, // 12: localtoast.PerOsBenchmarkConfig.version:type_name -> grafeas.v1.ComplianceVersion
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_scannerlib_proto_api_proto_init() }
func file_scannerlib_proto_api_proto_init() {
	if File_scannerlib_proto_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_scannerlib_proto_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scannerlib_proto_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OptOutConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scannerlib_proto_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BenchmarkConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scannerlib_proto_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanResults); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scannerlib_proto_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scannerlib_proto_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ComplianceResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scannerlib_proto_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DirContent); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scannerlib_proto_api_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PosixPermissions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scannerlib_proto_api_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PerOsBenchmarkConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_scannerlib_proto_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_scannerlib_proto_api_proto_goTypes,
		DependencyIndexes: file_scannerlib_proto_api_proto_depIdxs,
		EnumInfos:         file_scannerlib_proto_api_proto_enumTypes,
		MessageInfos:      file_scannerlib_proto_api_proto_msgTypes,
	}.Build()
	File_scannerlib_proto_api_proto = out.File
	file_scannerlib_proto_api_proto_rawDesc = nil
	file_scannerlib_proto_api_proto_goTypes = nil
	file_scannerlib_proto_api_proto_depIdxs = nil
}
