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

// Generated from scannerlib/proto/scan_instructions.proto using "bazel build"
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.14.0

package scan_instructions_go_proto

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RepeatConfig_RepeatType int32

const (
	RepeatConfig_ONCE                            RepeatConfig_RepeatType = 0
	RepeatConfig_FOR_EACH_USER                   RepeatConfig_RepeatType = 5
	RepeatConfig_FOR_EACH_USER_WITH_LOGIN        RepeatConfig_RepeatType = 1
	RepeatConfig_FOR_EACH_SYSTEM_USER_WITH_LOGIN RepeatConfig_RepeatType = 2
	RepeatConfig_FOR_EACH_OPEN_IPV4_PORT         RepeatConfig_RepeatType = 3
	RepeatConfig_FOR_EACH_OPEN_IPV6_PORT         RepeatConfig_RepeatType = 4
)

// Enum value maps for RepeatConfig_RepeatType.
var (
	RepeatConfig_RepeatType_name = map[int32]string{
		0: "ONCE",
		5: "FOR_EACH_USER",
		1: "FOR_EACH_USER_WITH_LOGIN",
		2: "FOR_EACH_SYSTEM_USER_WITH_LOGIN",
		3: "FOR_EACH_OPEN_IPV4_PORT",
		4: "FOR_EACH_OPEN_IPV6_PORT",
	}
	RepeatConfig_RepeatType_value = map[string]int32{
		"ONCE":                            0,
		"FOR_EACH_USER":                   5,
		"FOR_EACH_USER_WITH_LOGIN":        1,
		"FOR_EACH_SYSTEM_USER_WITH_LOGIN": 2,
		"FOR_EACH_OPEN_IPV4_PORT":         3,
		"FOR_EACH_OPEN_IPV6_PORT":         4,
	}
)

func (x RepeatConfig_RepeatType) Enum() *RepeatConfig_RepeatType {
	p := new(RepeatConfig_RepeatType)
	*p = x
	return p
}

func (x RepeatConfig_RepeatType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RepeatConfig_RepeatType) Descriptor() protoreflect.EnumDescriptor {
	return file_scannerlib_proto_scan_instructions_proto_enumTypes[0].Descriptor()
}

func (RepeatConfig_RepeatType) Type() protoreflect.EnumType {
	return &file_scannerlib_proto_scan_instructions_proto_enumTypes[0]
}

func (x RepeatConfig_RepeatType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RepeatConfig_RepeatType.Descriptor instead.
func (RepeatConfig_RepeatType) EnumDescriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{5, 0}
}

type PermissionCheck_BitMatchCriterion int32

const (
	PermissionCheck_NOT_SET             PermissionCheck_BitMatchCriterion = 0
	PermissionCheck_BOTH_SET_AND_CLEAR  PermissionCheck_BitMatchCriterion = 1
	PermissionCheck_EITHER_SET_OR_CLEAR PermissionCheck_BitMatchCriterion = 2
)

// Enum value maps for PermissionCheck_BitMatchCriterion.
var (
	PermissionCheck_BitMatchCriterion_name = map[int32]string{
		0: "NOT_SET",
		1: "BOTH_SET_AND_CLEAR",
		2: "EITHER_SET_OR_CLEAR",
	}
	PermissionCheck_BitMatchCriterion_value = map[string]int32{
		"NOT_SET":             0,
		"BOTH_SET_AND_CLEAR":  1,
		"EITHER_SET_OR_CLEAR": 2,
	}
)

func (x PermissionCheck_BitMatchCriterion) Enum() *PermissionCheck_BitMatchCriterion {
	p := new(PermissionCheck_BitMatchCriterion)
	*p = x
	return p
}

func (x PermissionCheck_BitMatchCriterion) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PermissionCheck_BitMatchCriterion) Descriptor() protoreflect.EnumDescriptor {
	return file_scannerlib_proto_scan_instructions_proto_enumTypes[1].Descriptor()
}

func (PermissionCheck_BitMatchCriterion) Type() protoreflect.EnumType {
	return &file_scannerlib_proto_scan_instructions_proto_enumTypes[1]
}

func (x PermissionCheck_BitMatchCriterion) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PermissionCheck_BitMatchCriterion.Descriptor instead.
func (PermissionCheck_BitMatchCriterion) EnumDescriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{7, 0}
}

type ContentEntryCheck_MatchType int32

const (
	ContentEntryCheck_NONE_MATCH             ContentEntryCheck_MatchType = 0
	ContentEntryCheck_ALL_MATCH_STRICT_ORDER ContentEntryCheck_MatchType = 1
	ContentEntryCheck_ALL_MATCH_ANY_ORDER    ContentEntryCheck_MatchType = 2
)

// Enum value maps for ContentEntryCheck_MatchType.
var (
	ContentEntryCheck_MatchType_name = map[int32]string{
		0: "NONE_MATCH",
		1: "ALL_MATCH_STRICT_ORDER",
		2: "ALL_MATCH_ANY_ORDER",
	}
	ContentEntryCheck_MatchType_value = map[string]int32{
		"NONE_MATCH":             0,
		"ALL_MATCH_STRICT_ORDER": 1,
		"ALL_MATCH_ANY_ORDER":    2,
	}
)

func (x ContentEntryCheck_MatchType) Enum() *ContentEntryCheck_MatchType {
	p := new(ContentEntryCheck_MatchType)
	*p = x
	return p
}

func (x ContentEntryCheck_MatchType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ContentEntryCheck_MatchType) Descriptor() protoreflect.EnumDescriptor {
	return file_scannerlib_proto_scan_instructions_proto_enumTypes[2].Descriptor()
}

func (ContentEntryCheck_MatchType) Type() protoreflect.EnumType {
	return &file_scannerlib_proto_scan_instructions_proto_enumTypes[2]
}

func (x ContentEntryCheck_MatchType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ContentEntryCheck_MatchType.Descriptor instead.
func (ContentEntryCheck_MatchType) EnumDescriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{9, 0}
}

type GroupCriterion_Type int32

const (
	GroupCriterion_LESS_THAN                 GroupCriterion_Type = 0
	GroupCriterion_GREATER_THAN              GroupCriterion_Type = 1
	GroupCriterion_NO_LESS_RESTRICTIVE_UMASK GroupCriterion_Type = 2
	GroupCriterion_UNIQUE                    GroupCriterion_Type = 3
)

// Enum value maps for GroupCriterion_Type.
var (
	GroupCriterion_Type_name = map[int32]string{
		0: "LESS_THAN",
		1: "GREATER_THAN",
		2: "NO_LESS_RESTRICTIVE_UMASK",
		3: "UNIQUE",
	}
	GroupCriterion_Type_value = map[string]int32{
		"LESS_THAN":                 0,
		"GREATER_THAN":              1,
		"NO_LESS_RESTRICTIVE_UMASK": 2,
		"UNIQUE":                    3,
	}
)

func (x GroupCriterion_Type) Enum() *GroupCriterion_Type {
	p := new(GroupCriterion_Type)
	*p = x
	return p
}

func (x GroupCriterion_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GroupCriterion_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_scannerlib_proto_scan_instructions_proto_enumTypes[3].Descriptor()
}

func (GroupCriterion_Type) Type() protoreflect.EnumType {
	return &file_scannerlib_proto_scan_instructions_proto_enumTypes[3]
}

func (x GroupCriterion_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GroupCriterion_Type.Descriptor instead.
func (GroupCriterion_Type) EnumDescriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{11, 0}
}

type SQLCheck_SQLDatabase int32

const (
	SQLCheck_DB_UNSPECIFIED SQLCheck_SQLDatabase = 0
	SQLCheck_DB_MYSQL       SQLCheck_SQLDatabase = 1
)

// Enum value maps for SQLCheck_SQLDatabase.
var (
	SQLCheck_SQLDatabase_name = map[int32]string{
		0: "DB_UNSPECIFIED",
		1: "DB_MYSQL",
	}
	SQLCheck_SQLDatabase_value = map[string]int32{
		"DB_UNSPECIFIED": 0,
		"DB_MYSQL":       1,
	}
)

func (x SQLCheck_SQLDatabase) Enum() *SQLCheck_SQLDatabase {
	p := new(SQLCheck_SQLDatabase)
	*p = x
	return p
}

func (x SQLCheck_SQLDatabase) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SQLCheck_SQLDatabase) Descriptor() protoreflect.EnumDescriptor {
	return file_scannerlib_proto_scan_instructions_proto_enumTypes[4].Descriptor()
}

func (SQLCheck_SQLDatabase) Type() protoreflect.EnumType {
	return &file_scannerlib_proto_scan_instructions_proto_enumTypes[4]
}

func (x SQLCheck_SQLDatabase) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SQLCheck_SQLDatabase.Descriptor instead.
func (SQLCheck_SQLDatabase) EnumDescriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{13, 0}
}

type BenchmarkScanInstructionDef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Instructions:
	//	*BenchmarkScanInstructionDef_Generic
	//	*BenchmarkScanInstructionDef_ScanTypeSpecific
	Instructions isBenchmarkScanInstructionDef_Instructions `protobuf_oneof:"instructions"`
}

func (x *BenchmarkScanInstructionDef) Reset() {
	*x = BenchmarkScanInstructionDef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BenchmarkScanInstructionDef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BenchmarkScanInstructionDef) ProtoMessage() {}

func (x *BenchmarkScanInstructionDef) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BenchmarkScanInstructionDef.ProtoReflect.Descriptor instead.
func (*BenchmarkScanInstructionDef) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{0}
}

func (m *BenchmarkScanInstructionDef) GetInstructions() isBenchmarkScanInstructionDef_Instructions {
	if m != nil {
		return m.Instructions
	}
	return nil
}

func (x *BenchmarkScanInstructionDef) GetGeneric() *BenchmarkScanInstruction {
	if x, ok := x.GetInstructions().(*BenchmarkScanInstructionDef_Generic); ok {
		return x.Generic
	}
	return nil
}

func (x *BenchmarkScanInstructionDef) GetScanTypeSpecific() *ScanTypeSpecificInstruction {
	if x, ok := x.GetInstructions().(*BenchmarkScanInstructionDef_ScanTypeSpecific); ok {
		return x.ScanTypeSpecific
	}
	return nil
}

type isBenchmarkScanInstructionDef_Instructions interface {
	isBenchmarkScanInstructionDef_Instructions()
}

type BenchmarkScanInstructionDef_Generic struct {
	Generic *BenchmarkScanInstruction `protobuf:"bytes,1,opt,name=generic,proto3,oneof"`
}

type BenchmarkScanInstructionDef_ScanTypeSpecific struct {
	ScanTypeSpecific *ScanTypeSpecificInstruction `protobuf:"bytes,2,opt,name=scan_type_specific,json=scanTypeSpecific,proto3,oneof"`
}

func (*BenchmarkScanInstructionDef_Generic) isBenchmarkScanInstructionDef_Instructions() {}

func (*BenchmarkScanInstructionDef_ScanTypeSpecific) isBenchmarkScanInstructionDef_Instructions() {}

type ScanTypeSpecificInstruction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InstanceScanning *BenchmarkScanInstruction `protobuf:"bytes,1,opt,name=instance_scanning,json=instanceScanning,proto3" json:"instance_scanning,omitempty"`
	ImageScanning    *BenchmarkScanInstruction `protobuf:"bytes,2,opt,name=image_scanning,json=imageScanning,proto3" json:"image_scanning,omitempty"`
}

func (x *ScanTypeSpecificInstruction) Reset() {
	*x = ScanTypeSpecificInstruction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanTypeSpecificInstruction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanTypeSpecificInstruction) ProtoMessage() {}

func (x *ScanTypeSpecificInstruction) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanTypeSpecificInstruction.ProtoReflect.Descriptor instead.
func (*ScanTypeSpecificInstruction) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{1}
}

func (x *ScanTypeSpecificInstruction) GetInstanceScanning() *BenchmarkScanInstruction {
	if x != nil {
		return x.InstanceScanning
	}
	return nil
}

func (x *ScanTypeSpecificInstruction) GetImageScanning() *BenchmarkScanInstruction {
	if x != nil {
		return x.ImageScanning
	}
	return nil
}

type BenchmarkScanInstruction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CheckAlternatives []*CheckAlternative `protobuf:"bytes,1,rep,name=check_alternatives,json=checkAlternatives,proto3" json:"check_alternatives,omitempty"`
}

func (x *BenchmarkScanInstruction) Reset() {
	*x = BenchmarkScanInstruction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BenchmarkScanInstruction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BenchmarkScanInstruction) ProtoMessage() {}

func (x *BenchmarkScanInstruction) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BenchmarkScanInstruction.ProtoReflect.Descriptor instead.
func (*BenchmarkScanInstruction) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{2}
}

func (x *BenchmarkScanInstruction) GetCheckAlternatives() []*CheckAlternative {
	if x != nil {
		return x.CheckAlternatives
	}
	return nil
}

type CheckAlternative struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileChecks []*FileCheck `protobuf:"bytes,1,rep,name=file_checks,json=fileChecks,proto3" json:"file_checks,omitempty"`
	SqlChecks  []*SQLCheck  `protobuf:"bytes,2,rep,name=sql_checks,json=sqlChecks,proto3" json:"sql_checks,omitempty"`
}

func (x *CheckAlternative) Reset() {
	*x = CheckAlternative{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckAlternative) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckAlternative) ProtoMessage() {}

func (x *CheckAlternative) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckAlternative.ProtoReflect.Descriptor instead.
func (*CheckAlternative) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{3}
}

func (x *CheckAlternative) GetFileChecks() []*FileCheck {
	if x != nil {
		return x.FileChecks
	}
	return nil
}

func (x *CheckAlternative) GetSqlChecks() []*SQLCheck {
	if x != nil {
		return x.SqlChecks
	}
	return nil
}

type FileCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FilesToCheck []*FileSet `protobuf:"bytes,1,rep,name=files_to_check,json=filesToCheck,proto3" json:"files_to_check,omitempty"`
	// Types that are assignable to CheckType:
	//	*FileCheck_Existence
	//	*FileCheck_Permission
	//	*FileCheck_Content
	//	*FileCheck_ContentEntry
	CheckType          isFileCheck_CheckType `protobuf_oneof:"check_type"`
	NonComplianceMsg   string                `protobuf:"bytes,6,opt,name=non_compliance_msg,json=nonComplianceMsg,proto3" json:"non_compliance_msg,omitempty"`
	FileDisplayCommand string                `protobuf:"bytes,7,opt,name=file_display_command,json=fileDisplayCommand,proto3" json:"file_display_command,omitempty"`
	RepeatConfig       *RepeatConfig         `protobuf:"bytes,9,opt,name=repeat_config,json=repeatConfig,proto3" json:"repeat_config,omitempty"`
}

func (x *FileCheck) Reset() {
	*x = FileCheck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileCheck) ProtoMessage() {}

func (x *FileCheck) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileCheck.ProtoReflect.Descriptor instead.
func (*FileCheck) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{4}
}

func (x *FileCheck) GetFilesToCheck() []*FileSet {
	if x != nil {
		return x.FilesToCheck
	}
	return nil
}

func (m *FileCheck) GetCheckType() isFileCheck_CheckType {
	if m != nil {
		return m.CheckType
	}
	return nil
}

func (x *FileCheck) GetExistence() *ExistenceCheck {
	if x, ok := x.GetCheckType().(*FileCheck_Existence); ok {
		return x.Existence
	}
	return nil
}

func (x *FileCheck) GetPermission() *PermissionCheck {
	if x, ok := x.GetCheckType().(*FileCheck_Permission); ok {
		return x.Permission
	}
	return nil
}

func (x *FileCheck) GetContent() *ContentCheck {
	if x, ok := x.GetCheckType().(*FileCheck_Content); ok {
		return x.Content
	}
	return nil
}

func (x *FileCheck) GetContentEntry() *ContentEntryCheck {
	if x, ok := x.GetCheckType().(*FileCheck_ContentEntry); ok {
		return x.ContentEntry
	}
	return nil
}

func (x *FileCheck) GetNonComplianceMsg() string {
	if x != nil {
		return x.NonComplianceMsg
	}
	return ""
}

func (x *FileCheck) GetFileDisplayCommand() string {
	if x != nil {
		return x.FileDisplayCommand
	}
	return ""
}

func (x *FileCheck) GetRepeatConfig() *RepeatConfig {
	if x != nil {
		return x.RepeatConfig
	}
	return nil
}

type isFileCheck_CheckType interface {
	isFileCheck_CheckType()
}

type FileCheck_Existence struct {
	Existence *ExistenceCheck `protobuf:"bytes,2,opt,name=existence,proto3,oneof"`
}

type FileCheck_Permission struct {
	Permission *PermissionCheck `protobuf:"bytes,3,opt,name=permission,proto3,oneof"`
}

type FileCheck_Content struct {
	Content *ContentCheck `protobuf:"bytes,4,opt,name=content,proto3,oneof"`
}

type FileCheck_ContentEntry struct {
	ContentEntry *ContentEntryCheck `protobuf:"bytes,5,opt,name=content_entry,json=contentEntry,proto3,oneof"`
}

func (*FileCheck_Existence) isFileCheck_CheckType() {}

func (*FileCheck_Permission) isFileCheck_CheckType() {}

func (*FileCheck_Content) isFileCheck_CheckType() {}

func (*FileCheck_ContentEntry) isFileCheck_CheckType() {}

type RepeatConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type   RepeatConfig_RepeatType            `protobuf:"varint,1,opt,name=type,proto3,enum=localtoast.RepeatConfig_RepeatType" json:"type,omitempty"`
	OptOut []*RepeatConfig_OptOutSubstitution `protobuf:"bytes,2,rep,name=opt_out,json=optOut,proto3" json:"opt_out,omitempty"`
}

func (x *RepeatConfig) Reset() {
	*x = RepeatConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepeatConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepeatConfig) ProtoMessage() {}

func (x *RepeatConfig) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepeatConfig.ProtoReflect.Descriptor instead.
func (*RepeatConfig) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{5}
}

func (x *RepeatConfig) GetType() RepeatConfig_RepeatType {
	if x != nil {
		return x.Type
	}
	return RepeatConfig_ONCE
}

func (x *RepeatConfig) GetOptOut() []*RepeatConfig_OptOutSubstitution {
	if x != nil {
		return x.OptOut
	}
	return nil
}

type ExistenceCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShouldExist bool `protobuf:"varint,1,opt,name=should_exist,json=shouldExist,proto3" json:"should_exist,omitempty"`
}

func (x *ExistenceCheck) Reset() {
	*x = ExistenceCheck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExistenceCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExistenceCheck) ProtoMessage() {}

func (x *ExistenceCheck) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExistenceCheck.ProtoReflect.Descriptor instead.
func (*ExistenceCheck) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{6}
}

func (x *ExistenceCheck) GetShouldExist() bool {
	if x != nil {
		return x.ShouldExist
	}
	return false
}

type PermissionCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SetBits         int32                             `protobuf:"varint,1,opt,name=set_bits,json=setBits,proto3" json:"set_bits,omitempty"`
	ClearBits       int32                             `protobuf:"varint,2,opt,name=clear_bits,json=clearBits,proto3" json:"clear_bits,omitempty"`
	BitsShouldMatch PermissionCheck_BitMatchCriterion `protobuf:"varint,3,opt,name=bits_should_match,json=bitsShouldMatch,proto3,enum=localtoast.PermissionCheck_BitMatchCriterion" json:"bits_should_match,omitempty"`
	User            *PermissionCheck_OwnerCheck       `protobuf:"bytes,4,opt,name=user,proto3" json:"user,omitempty"`
	Group           *PermissionCheck_OwnerCheck       `protobuf:"bytes,5,opt,name=group,proto3" json:"group,omitempty"`
}

func (x *PermissionCheck) Reset() {
	*x = PermissionCheck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PermissionCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PermissionCheck) ProtoMessage() {}

func (x *PermissionCheck) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PermissionCheck.ProtoReflect.Descriptor instead.
func (*PermissionCheck) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{7}
}

func (x *PermissionCheck) GetSetBits() int32 {
	if x != nil {
		return x.SetBits
	}
	return 0
}

func (x *PermissionCheck) GetClearBits() int32 {
	if x != nil {
		return x.ClearBits
	}
	return 0
}

func (x *PermissionCheck) GetBitsShouldMatch() PermissionCheck_BitMatchCriterion {
	if x != nil {
		return x.BitsShouldMatch
	}
	return PermissionCheck_NOT_SET
}

func (x *PermissionCheck) GetUser() *PermissionCheck_OwnerCheck {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *PermissionCheck) GetGroup() *PermissionCheck_OwnerCheck {
	if x != nil {
		return x.Group
	}
	return nil
}

type ContentCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ContentCheck) Reset() {
	*x = ContentCheck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContentCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContentCheck) ProtoMessage() {}

func (x *ContentCheck) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContentCheck.ProtoReflect.Descriptor instead.
func (*ContentCheck) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{8}
}

func (x *ContentCheck) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type ContentEntryCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Delimiter     []byte                      `protobuf:"bytes,1,opt,name=delimiter,proto3" json:"delimiter,omitempty"`
	MatchType     ContentEntryCheck_MatchType `protobuf:"varint,2,opt,name=match_type,json=matchType,proto3,enum=localtoast.ContentEntryCheck_MatchType" json:"match_type,omitempty"`
	MatchCriteria []*MatchCriterion           `protobuf:"bytes,3,rep,name=match_criteria,json=matchCriteria,proto3" json:"match_criteria,omitempty"`
}

func (x *ContentEntryCheck) Reset() {
	*x = ContentEntryCheck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContentEntryCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContentEntryCheck) ProtoMessage() {}

func (x *ContentEntryCheck) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContentEntryCheck.ProtoReflect.Descriptor instead.
func (*ContentEntryCheck) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{9}
}

func (x *ContentEntryCheck) GetDelimiter() []byte {
	if x != nil {
		return x.Delimiter
	}
	return nil
}

func (x *ContentEntryCheck) GetMatchType() ContentEntryCheck_MatchType {
	if x != nil {
		return x.MatchType
	}
	return ContentEntryCheck_NONE_MATCH
}

func (x *ContentEntryCheck) GetMatchCriteria() []*MatchCriterion {
	if x != nil {
		return x.MatchCriteria
	}
	return nil
}

type MatchCriterion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FilterRegex   string            `protobuf:"bytes,1,opt,name=filter_regex,json=filterRegex,proto3" json:"filter_regex,omitempty"`
	ExpectedRegex string            `protobuf:"bytes,2,opt,name=expected_regex,json=expectedRegex,proto3" json:"expected_regex,omitempty"`
	GroupCriteria []*GroupCriterion `protobuf:"bytes,3,rep,name=group_criteria,json=groupCriteria,proto3" json:"group_criteria,omitempty"`
}

func (x *MatchCriterion) Reset() {
	*x = MatchCriterion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MatchCriterion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MatchCriterion) ProtoMessage() {}

func (x *MatchCriterion) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MatchCriterion.ProtoReflect.Descriptor instead.
func (*MatchCriterion) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{10}
}

func (x *MatchCriterion) GetFilterRegex() string {
	if x != nil {
		return x.FilterRegex
	}
	return ""
}

func (x *MatchCriterion) GetExpectedRegex() string {
	if x != nil {
		return x.ExpectedRegex
	}
	return ""
}

func (x *MatchCriterion) GetGroupCriteria() []*GroupCriterion {
	if x != nil {
		return x.GroupCriteria
	}
	return nil
}

type GroupCriterion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupIndex int32               `protobuf:"varint,1,opt,name=group_index,json=groupIndex,proto3" json:"group_index,omitempty"`
	Type       GroupCriterion_Type `protobuf:"varint,2,opt,name=type,proto3,enum=localtoast.GroupCriterion_Type" json:"type,omitempty"`
	// Types that are assignable to ComparisonValue:
	//	*GroupCriterion_Const
	//	*GroupCriterion_Today_
	ComparisonValue isGroupCriterion_ComparisonValue `protobuf_oneof:"comparison_value"`
}

func (x *GroupCriterion) Reset() {
	*x = GroupCriterion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupCriterion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupCriterion) ProtoMessage() {}

func (x *GroupCriterion) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupCriterion.ProtoReflect.Descriptor instead.
func (*GroupCriterion) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{11}
}

func (x *GroupCriterion) GetGroupIndex() int32 {
	if x != nil {
		return x.GroupIndex
	}
	return 0
}

func (x *GroupCriterion) GetType() GroupCriterion_Type {
	if x != nil {
		return x.Type
	}
	return GroupCriterion_LESS_THAN
}

func (m *GroupCriterion) GetComparisonValue() isGroupCriterion_ComparisonValue {
	if m != nil {
		return m.ComparisonValue
	}
	return nil
}

func (x *GroupCriterion) GetConst() int32 {
	if x, ok := x.GetComparisonValue().(*GroupCriterion_Const); ok {
		return x.Const
	}
	return 0
}

func (x *GroupCriterion) GetToday() *GroupCriterion_Today {
	if x, ok := x.GetComparisonValue().(*GroupCriterion_Today_); ok {
		return x.Today
	}
	return nil
}

type isGroupCriterion_ComparisonValue interface {
	isGroupCriterion_ComparisonValue()
}

type GroupCriterion_Const struct {
	Const int32 `protobuf:"varint,3,opt,name=const,proto3,oneof"`
}

type GroupCriterion_Today_ struct {
	Today *GroupCriterion_Today `protobuf:"bytes,4,opt,name=today,proto3,oneof"`
}

func (*GroupCriterion_Const) isGroupCriterion_ComparisonValue() {}

func (*GroupCriterion_Today_) isGroupCriterion_ComparisonValue() {}

type FileSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to FilePath:
	//	*FileSet_SingleFile_
	//	*FileSet_FilesInDir_
	//	*FileSet_ProcessPath_
	//	*FileSet_UnixEnvVarPaths_
	FilePath isFileSet_FilePath `protobuf_oneof:"file_path"`
}

func (x *FileSet) Reset() {
	*x = FileSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileSet) ProtoMessage() {}

func (x *FileSet) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileSet.ProtoReflect.Descriptor instead.
func (*FileSet) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{12}
}

func (m *FileSet) GetFilePath() isFileSet_FilePath {
	if m != nil {
		return m.FilePath
	}
	return nil
}

func (x *FileSet) GetSingleFile() *FileSet_SingleFile {
	if x, ok := x.GetFilePath().(*FileSet_SingleFile_); ok {
		return x.SingleFile
	}
	return nil
}

func (x *FileSet) GetFilesInDir() *FileSet_FilesInDir {
	if x, ok := x.GetFilePath().(*FileSet_FilesInDir_); ok {
		return x.FilesInDir
	}
	return nil
}

func (x *FileSet) GetProcessPath() *FileSet_ProcessPath {
	if x, ok := x.GetFilePath().(*FileSet_ProcessPath_); ok {
		return x.ProcessPath
	}
	return nil
}

func (x *FileSet) GetUnixEnvVarPaths() *FileSet_UnixEnvVarPaths {
	if x, ok := x.GetFilePath().(*FileSet_UnixEnvVarPaths_); ok {
		return x.UnixEnvVarPaths
	}
	return nil
}

type isFileSet_FilePath interface {
	isFileSet_FilePath()
}

type FileSet_SingleFile_ struct {
	SingleFile *FileSet_SingleFile `protobuf:"bytes,1,opt,name=single_file,json=singleFile,proto3,oneof"`
}

type FileSet_FilesInDir_ struct {
	FilesInDir *FileSet_FilesInDir `protobuf:"bytes,2,opt,name=files_in_dir,json=filesInDir,proto3,oneof"`
}

type FileSet_ProcessPath_ struct {
	ProcessPath *FileSet_ProcessPath `protobuf:"bytes,3,opt,name=process_path,json=processPath,proto3,oneof"`
}

type FileSet_UnixEnvVarPaths_ struct {
	UnixEnvVarPaths *FileSet_UnixEnvVarPaths `protobuf:"bytes,4,opt,name=unix_env_var_paths,json=unixEnvVarPaths,proto3,oneof"`
}

func (*FileSet_SingleFile_) isFileSet_FilePath() {}

func (*FileSet_FilesInDir_) isFileSet_FilePath() {}

func (*FileSet_ProcessPath_) isFileSet_FilePath() {}

func (*FileSet_UnixEnvVarPaths_) isFileSet_FilePath() {}

type SQLCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetDatabase   SQLCheck_SQLDatabase `protobuf:"varint,1,opt,name=target_database,json=targetDatabase,proto3,enum=localtoast.SQLCheck_SQLDatabase" json:"target_database,omitempty"`
	Query            string               `protobuf:"bytes,2,opt,name=query,proto3" json:"query,omitempty"`
	ExpectResults    bool                 `protobuf:"varint,3,opt,name=expect_results,json=expectResults,proto3" json:"expect_results,omitempty"`
	NonComplianceMsg string               `protobuf:"bytes,4,opt,name=non_compliance_msg,json=nonComplianceMsg,proto3" json:"non_compliance_msg,omitempty"`
}

func (x *SQLCheck) Reset() {
	*x = SQLCheck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SQLCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SQLCheck) ProtoMessage() {}

func (x *SQLCheck) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SQLCheck.ProtoReflect.Descriptor instead.
func (*SQLCheck) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{13}
}

func (x *SQLCheck) GetTargetDatabase() SQLCheck_SQLDatabase {
	if x != nil {
		return x.TargetDatabase
	}
	return SQLCheck_DB_UNSPECIFIED
}

func (x *SQLCheck) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *SQLCheck) GetExpectResults() bool {
	if x != nil {
		return x.ExpectResults
	}
	return false
}

func (x *SQLCheck) GetNonComplianceMsg() string {
	if x != nil {
		return x.NonComplianceMsg
	}
	return ""
}

type RepeatConfig_OptOutSubstitution struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Wildcard string `protobuf:"bytes,1,opt,name=wildcard,proto3" json:"wildcard,omitempty"`
	Value    string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *RepeatConfig_OptOutSubstitution) Reset() {
	*x = RepeatConfig_OptOutSubstitution{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepeatConfig_OptOutSubstitution) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepeatConfig_OptOutSubstitution) ProtoMessage() {}

func (x *RepeatConfig_OptOutSubstitution) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepeatConfig_OptOutSubstitution.ProtoReflect.Descriptor instead.
func (*RepeatConfig_OptOutSubstitution) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{5, 0}
}

func (x *RepeatConfig_OptOutSubstitution) GetWildcard() string {
	if x != nil {
		return x.Wildcard
	}
	return ""
}

func (x *RepeatConfig_OptOutSubstitution) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type PermissionCheck_OwnerCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ShouldOwn bool   `protobuf:"varint,2,opt,name=should_own,json=shouldOwn,proto3" json:"should_own,omitempty"`
}

func (x *PermissionCheck_OwnerCheck) Reset() {
	*x = PermissionCheck_OwnerCheck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[15]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PermissionCheck_OwnerCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PermissionCheck_OwnerCheck) ProtoMessage() {}

func (x *PermissionCheck_OwnerCheck) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[15]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PermissionCheck_OwnerCheck.ProtoReflect.Descriptor instead.
func (*PermissionCheck_OwnerCheck) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{7, 0}
}

func (x *PermissionCheck_OwnerCheck) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PermissionCheck_OwnerCheck) GetShouldOwn() bool {
	if x != nil {
		return x.ShouldOwn
	}
	return false
}

type GroupCriterion_Today struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GroupCriterion_Today) Reset() {
	*x = GroupCriterion_Today{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[16]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupCriterion_Today) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupCriterion_Today) ProtoMessage() {}

func (x *GroupCriterion_Today) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[16]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupCriterion_Today.ProtoReflect.Descriptor instead.
func (*GroupCriterion_Today) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{11, 0}
}

type FileSet_SingleFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *FileSet_SingleFile) Reset() {
	*x = FileSet_SingleFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[17]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileSet_SingleFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileSet_SingleFile) ProtoMessage() {}

func (x *FileSet_SingleFile) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[17]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileSet_SingleFile.ProtoReflect.Descriptor instead.
func (*FileSet_SingleFile) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{12, 0}
}

func (x *FileSet_SingleFile) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type FileSet_FilesInDir struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DirPath           string   `protobuf:"bytes,1,opt,name=dir_path,json=dirPath,proto3" json:"dir_path,omitempty"`
	Recursive         bool     `protobuf:"varint,2,opt,name=recursive,proto3" json:"recursive,omitempty"`
	FilesOnly         bool     `protobuf:"varint,3,opt,name=files_only,json=filesOnly,proto3" json:"files_only,omitempty"`
	DirsOnly          bool     `protobuf:"varint,4,opt,name=dirs_only,json=dirsOnly,proto3" json:"dirs_only,omitempty"`
	SkipSymlinks      bool     `protobuf:"varint,7,opt,name=skip_symlinks,json=skipSymlinks,proto3" json:"skip_symlinks,omitempty"`
	FilenameRegex     string   `protobuf:"bytes,5,opt,name=filename_regex,json=filenameRegex,proto3" json:"filename_regex,omitempty"`
	OptOutPathRegexes []string `protobuf:"bytes,6,rep,name=opt_out_path_regexes,json=optOutPathRegexes,proto3" json:"opt_out_path_regexes,omitempty"`
}

func (x *FileSet_FilesInDir) Reset() {
	*x = FileSet_FilesInDir{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[18]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileSet_FilesInDir) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileSet_FilesInDir) ProtoMessage() {}

func (x *FileSet_FilesInDir) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[18]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileSet_FilesInDir.ProtoReflect.Descriptor instead.
func (*FileSet_FilesInDir) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{12, 1}
}

func (x *FileSet_FilesInDir) GetDirPath() string {
	if x != nil {
		return x.DirPath
	}
	return ""
}

func (x *FileSet_FilesInDir) GetRecursive() bool {
	if x != nil {
		return x.Recursive
	}
	return false
}

func (x *FileSet_FilesInDir) GetFilesOnly() bool {
	if x != nil {
		return x.FilesOnly
	}
	return false
}

func (x *FileSet_FilesInDir) GetDirsOnly() bool {
	if x != nil {
		return x.DirsOnly
	}
	return false
}

func (x *FileSet_FilesInDir) GetSkipSymlinks() bool {
	if x != nil {
		return x.SkipSymlinks
	}
	return false
}

func (x *FileSet_FilesInDir) GetFilenameRegex() string {
	if x != nil {
		return x.FilenameRegex
	}
	return ""
}

func (x *FileSet_FilesInDir) GetOptOutPathRegexes() []string {
	if x != nil {
		return x.OptOutPathRegexes
	}
	return nil
}

type FileSet_ProcessPath struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcName string `protobuf:"bytes,1,opt,name=proc_name,json=procName,proto3" json:"proc_name,omitempty"`
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
}

func (x *FileSet_ProcessPath) Reset() {
	*x = FileSet_ProcessPath{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[19]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileSet_ProcessPath) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileSet_ProcessPath) ProtoMessage() {}

func (x *FileSet_ProcessPath) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[19]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileSet_ProcessPath.ProtoReflect.Descriptor instead.
func (*FileSet_ProcessPath) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{12, 2}
}

func (x *FileSet_ProcessPath) GetProcName() string {
	if x != nil {
		return x.ProcName
	}
	return ""
}

func (x *FileSet_ProcessPath) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

type FileSet_UnixEnvVarPaths struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VarName   string `protobuf:"bytes,1,opt,name=var_name,json=varName,proto3" json:"var_name,omitempty"`
	FilesOnly bool   `protobuf:"varint,2,opt,name=files_only,json=filesOnly,proto3" json:"files_only,omitempty"`
	DirsOnly  bool   `protobuf:"varint,3,opt,name=dirs_only,json=dirsOnly,proto3" json:"dirs_only,omitempty"`
}

func (x *FileSet_UnixEnvVarPaths) Reset() {
	*x = FileSet_UnixEnvVarPaths{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[20]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileSet_UnixEnvVarPaths) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileSet_UnixEnvVarPaths) ProtoMessage() {}

func (x *FileSet_UnixEnvVarPaths) ProtoReflect() protoreflect.Message {
	mi := &file_scannerlib_proto_scan_instructions_proto_msgTypes[20]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileSet_UnixEnvVarPaths.ProtoReflect.Descriptor instead.
func (*FileSet_UnixEnvVarPaths) Descriptor() ([]byte, []int) {
	return file_scannerlib_proto_scan_instructions_proto_rawDescGZIP(), []int{12, 3}
}

func (x *FileSet_UnixEnvVarPaths) GetVarName() string {
	if x != nil {
		return x.VarName
	}
	return ""
}

func (x *FileSet_UnixEnvVarPaths) GetFilesOnly() bool {
	if x != nil {
		return x.FilesOnly
	}
	return false
}

func (x *FileSet_UnixEnvVarPaths) GetDirsOnly() bool {
	if x != nil {
		return x.DirsOnly
	}
	return false
}

var File_scannerlib_proto_scan_instructions_proto protoreflect.FileDescriptor

var file_scannerlib_proto_scan_instructions_proto_rawDesc = []byte{
	0x0a, 0x28, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x6c, 0x69, 0x62, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x73, 0x63, 0x61, 0x6e, 0x5f, 0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6c, 0x6f, 0x63, 0x61,
	0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x22, 0xc8, 0x01, 0x0a, 0x1b, 0x42, 0x65, 0x6e, 0x63, 0x68,
	0x6d, 0x61, 0x72, 0x6b, 0x53, 0x63, 0x61, 0x6e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x44, 0x65, 0x66, 0x12, 0x40, 0x0a, 0x07, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x69,
	0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74,
	0x6f, 0x61, 0x73, 0x74, 0x2e, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x53, 0x63,
	0x61, 0x6e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52,
	0x07, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x12, 0x57, 0x0a, 0x12, 0x73, 0x63, 0x61, 0x6e,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73,
	0x74, 0x2e, 0x53, 0x63, 0x61, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66,
	0x69, 0x63, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52,
	0x10, 0x73, 0x63, 0x61, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69,
	0x63, 0x42, 0x0e, 0x0a, 0x0c, 0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x22, 0xbd, 0x01, 0x0a, 0x1b, 0x53, 0x63, 0x61, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x53, 0x70,
	0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x51, 0x0a, 0x11, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x73, 0x63,
	0x61, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6c,
	0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d,
	0x61, 0x72, 0x6b, 0x53, 0x63, 0x61, 0x6e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x10, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x63, 0x61, 0x6e,
	0x6e, 0x69, 0x6e, 0x67, 0x12, 0x4b, 0x0a, 0x0e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x63,
	0x61, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6c,
	0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d,
	0x61, 0x72, 0x6b, 0x53, 0x63, 0x61, 0x6e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0d, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x53, 0x63, 0x61, 0x6e, 0x6e, 0x69, 0x6e,
	0x67, 0x22, 0x67, 0x0a, 0x18, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x53, 0x63,
	0x61, 0x6e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4b, 0x0a,
	0x12, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x5f, 0x61, 0x6c, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x74, 0x69,
	0x76, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6c, 0x6f, 0x63, 0x61,
	0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x41, 0x6c, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x52, 0x11, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x41, 0x6c,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x73, 0x22, 0x7f, 0x0a, 0x10, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x41, 0x6c, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x12, 0x36,
	0x0a, 0x0b, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x0a, 0x66, 0x69, 0x6c, 0x65,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x12, 0x33, 0x0a, 0x0a, 0x73, 0x71, 0x6c, 0x5f, 0x63, 0x68,
	0x65, 0x63, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6c, 0x6f, 0x63,
	0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x53, 0x51, 0x4c, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x09, 0x73, 0x71, 0x6c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x22, 0xf0, 0x03, 0x0a, 0x09,
	0x46, 0x69, 0x6c, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x39, 0x0a, 0x0e, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x5f, 0x74, 0x6f, 0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x53, 0x65, 0x74, 0x52, 0x0c, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x54, 0x6f, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x12, 0x3a, 0x0a, 0x09, 0x65, 0x78, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74,
	0x6f, 0x61, 0x73, 0x74, 0x2e, 0x45, 0x78, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x48, 0x00, 0x52, 0x09, 0x65, 0x78, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65,
	0x12, 0x3d, 0x0a, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73,
	0x74, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x34, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x18, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x00, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x44, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x5f, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6c,
	0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x00, 0x52, 0x0c, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x2c, 0x0a, 0x12, 0x6e,
	0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x6d, 0x73,
	0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x6e, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x4d, 0x73, 0x67, 0x12, 0x30, 0x0a, 0x14, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x66, 0x69, 0x6c, 0x65, 0x44, 0x69, 0x73,
	0x70, 0x6c, 0x61, 0x79, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x3d, 0x0a, 0x0d, 0x72,
	0x65, 0x70, 0x65, 0x61, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e,
	0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0c, 0x72, 0x65,
	0x70, 0x65, 0x61, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x0c, 0x0a, 0x0a, 0x63, 0x68,
	0x65, 0x63, 0x6b, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x4a, 0x04, 0x08, 0x08, 0x10, 0x09, 0x22, 0xfe,
	0x02, 0x0a, 0x0c, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x37, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x70, 0x65, 0x61,
	0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x44, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x5f,
	0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x6c, 0x6f, 0x63, 0x61,
	0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x4f, 0x70, 0x74, 0x4f, 0x75, 0x74, 0x53, 0x75, 0x62, 0x73, 0x74, 0x69,
	0x74, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x6f, 0x70, 0x74, 0x4f, 0x75, 0x74, 0x1a, 0x46,
	0x0a, 0x12, 0x4f, 0x70, 0x74, 0x4f, 0x75, 0x74, 0x53, 0x75, 0x62, 0x73, 0x74, 0x69, 0x74, 0x75,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x77, 0x69, 0x6c, 0x64, 0x63, 0x61, 0x72, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77, 0x69, 0x6c, 0x64, 0x63, 0x61, 0x72, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xa6, 0x01, 0x0a, 0x0a, 0x52, 0x65, 0x70, 0x65, 0x61,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x4f, 0x4e, 0x43, 0x45, 0x10, 0x00, 0x12,
	0x11, 0x0a, 0x0d, 0x46, 0x4f, 0x52, 0x5f, 0x45, 0x41, 0x43, 0x48, 0x5f, 0x55, 0x53, 0x45, 0x52,
	0x10, 0x05, 0x12, 0x1c, 0x0a, 0x18, 0x46, 0x4f, 0x52, 0x5f, 0x45, 0x41, 0x43, 0x48, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x5f, 0x57, 0x49, 0x54, 0x48, 0x5f, 0x4c, 0x4f, 0x47, 0x49, 0x4e, 0x10, 0x01,
	0x12, 0x23, 0x0a, 0x1f, 0x46, 0x4f, 0x52, 0x5f, 0x45, 0x41, 0x43, 0x48, 0x5f, 0x53, 0x59, 0x53,
	0x54, 0x45, 0x4d, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x57, 0x49, 0x54, 0x48, 0x5f, 0x4c, 0x4f,
	0x47, 0x49, 0x4e, 0x10, 0x02, 0x12, 0x1b, 0x0a, 0x17, 0x46, 0x4f, 0x52, 0x5f, 0x45, 0x41, 0x43,
	0x48, 0x5f, 0x4f, 0x50, 0x45, 0x4e, 0x5f, 0x49, 0x50, 0x56, 0x34, 0x5f, 0x50, 0x4f, 0x52, 0x54,
	0x10, 0x03, 0x12, 0x1b, 0x0a, 0x17, 0x46, 0x4f, 0x52, 0x5f, 0x45, 0x41, 0x43, 0x48, 0x5f, 0x4f,
	0x50, 0x45, 0x4e, 0x5f, 0x49, 0x50, 0x56, 0x36, 0x5f, 0x50, 0x4f, 0x52, 0x54, 0x10, 0x04, 0x22,
	0x33, 0x0a, 0x0e, 0x45, 0x78, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x68, 0x6f, 0x75, 0x6c, 0x64, 0x5f, 0x65, 0x78, 0x69, 0x73,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x73, 0x68, 0x6f, 0x75, 0x6c, 0x64, 0x45,
	0x78, 0x69, 0x73, 0x74, 0x22, 0xb4, 0x03, 0x0a, 0x0f, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x65, 0x74, 0x5f,
	0x62, 0x69, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x73, 0x65, 0x74, 0x42,
	0x69, 0x74, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6c, 0x65, 0x61, 0x72, 0x5f, 0x62, 0x69, 0x74,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x63, 0x6c, 0x65, 0x61, 0x72, 0x42, 0x69,
	0x74, 0x73, 0x12, 0x59, 0x0a, 0x11, 0x62, 0x69, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x75, 0x6c,
	0x64, 0x5f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2d, 0x2e,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x42, 0x69, 0x74, 0x4d, 0x61,
	0x74, 0x63, 0x68, 0x43, 0x72, 0x69, 0x74, 0x65, 0x72, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x62, 0x69,
	0x74, 0x73, 0x53, 0x68, 0x6f, 0x75, 0x6c, 0x64, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x12, 0x3a, 0x0a,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x3c, 0x0a, 0x05, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c,
	0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x1a, 0x3f, 0x0a, 0x0a, 0x4f, 0x77, 0x6e, 0x65, 0x72,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x68, 0x6f,
	0x75, 0x6c, 0x64, 0x5f, 0x6f, 0x77, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x73,
	0x68, 0x6f, 0x75, 0x6c, 0x64, 0x4f, 0x77, 0x6e, 0x22, 0x51, 0x0a, 0x11, 0x42, 0x69, 0x74, 0x4d,
	0x61, 0x74, 0x63, 0x68, 0x43, 0x72, 0x69, 0x74, 0x65, 0x72, 0x69, 0x6f, 0x6e, 0x12, 0x0b, 0x0a,
	0x07, 0x4e, 0x4f, 0x54, 0x5f, 0x53, 0x45, 0x54, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x42, 0x4f,
	0x54, 0x48, 0x5f, 0x53, 0x45, 0x54, 0x5f, 0x41, 0x4e, 0x44, 0x5f, 0x43, 0x4c, 0x45, 0x41, 0x52,
	0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x49, 0x54, 0x48, 0x45, 0x52, 0x5f, 0x53, 0x45, 0x54,
	0x5f, 0x4f, 0x52, 0x5f, 0x43, 0x4c, 0x45, 0x41, 0x52, 0x10, 0x02, 0x22, 0x28, 0x0a, 0x0c, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x8e, 0x02, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x64,
	0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x64, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x12, 0x46, 0x0a, 0x0a, 0x6d, 0x61, 0x74,
	0x63, 0x68, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x27, 0x2e,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x4d, 0x61, 0x74,
	0x63, 0x68, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x41, 0x0a, 0x0e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x63, 0x72, 0x69, 0x74, 0x65,
	0x72, 0x69, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6c, 0x6f, 0x63, 0x61,
	0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x43, 0x72, 0x69, 0x74,
	0x65, 0x72, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x43, 0x72, 0x69, 0x74,
	0x65, 0x72, 0x69, 0x61, 0x22, 0x50, 0x0a, 0x09, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x4e, 0x4f, 0x4e, 0x45, 0x5f, 0x4d, 0x41, 0x54, 0x43, 0x48, 0x10,
	0x00, 0x12, 0x1a, 0x0a, 0x16, 0x41, 0x4c, 0x4c, 0x5f, 0x4d, 0x41, 0x54, 0x43, 0x48, 0x5f, 0x53,
	0x54, 0x52, 0x49, 0x43, 0x54, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x10, 0x01, 0x12, 0x17, 0x0a,
	0x13, 0x41, 0x4c, 0x4c, 0x5f, 0x4d, 0x41, 0x54, 0x43, 0x48, 0x5f, 0x41, 0x4e, 0x59, 0x5f, 0x4f,
	0x52, 0x44, 0x45, 0x52, 0x10, 0x02, 0x22, 0x9d, 0x01, 0x0a, 0x0e, 0x4d, 0x61, 0x74, 0x63, 0x68,
	0x43, 0x72, 0x69, 0x74, 0x65, 0x72, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x67, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x65, 0x67, 0x65, 0x78, 0x12, 0x25, 0x0a, 0x0e,
	0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x67, 0x65, 0x78, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x52, 0x65,
	0x67, 0x65, 0x78, 0x12, 0x41, 0x0a, 0x0e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x63, 0x72, 0x69,
	0x74, 0x65, 0x72, 0x69, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x43, 0x72,
	0x69, 0x74, 0x65, 0x72, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x43, 0x72,
	0x69, 0x74, 0x65, 0x72, 0x69, 0x61, 0x22, 0xa9, 0x02, 0x0a, 0x0e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x43, 0x72, 0x69, 0x74, 0x65, 0x72, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x33, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c,
	0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x43, 0x72, 0x69, 0x74, 0x65,
	0x72, 0x69, 0x6f, 0x6e, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x16, 0x0a, 0x05, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00,
	0x52, 0x05, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x12, 0x38, 0x0a, 0x05, 0x74, 0x6f, 0x64, 0x61, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f,
	0x61, 0x73, 0x74, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x43, 0x72, 0x69, 0x74, 0x65, 0x72, 0x69,
	0x6f, 0x6e, 0x2e, 0x54, 0x6f, 0x64, 0x61, 0x79, 0x48, 0x00, 0x52, 0x05, 0x74, 0x6f, 0x64, 0x61,
	0x79, 0x1a, 0x07, 0x0a, 0x05, 0x54, 0x6f, 0x64, 0x61, 0x79, 0x22, 0x52, 0x0a, 0x04, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x4c, 0x45, 0x53, 0x53, 0x5f, 0x54, 0x48, 0x41, 0x4e, 0x10,
	0x00, 0x12, 0x10, 0x0a, 0x0c, 0x47, 0x52, 0x45, 0x41, 0x54, 0x45, 0x52, 0x5f, 0x54, 0x48, 0x41,
	0x4e, 0x10, 0x01, 0x12, 0x1d, 0x0a, 0x19, 0x4e, 0x4f, 0x5f, 0x4c, 0x45, 0x53, 0x53, 0x5f, 0x52,
	0x45, 0x53, 0x54, 0x52, 0x49, 0x43, 0x54, 0x49, 0x56, 0x45, 0x5f, 0x55, 0x4d, 0x41, 0x53, 0x4b,
	0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x4e, 0x49, 0x51, 0x55, 0x45, 0x10, 0x03, 0x42, 0x12,
	0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x69, 0x73, 0x6f, 0x6e, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x8d, 0x06, 0x0a, 0x07, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x74, 0x12, 0x41,
	0x0a, 0x0b, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x74, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x46,
	0x69, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x0a, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x42, 0x0a, 0x0c, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x5f, 0x69, 0x6e, 0x5f, 0x64, 0x69,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74,
	0x6f, 0x61, 0x73, 0x74, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x74, 0x2e, 0x46, 0x69, 0x6c,
	0x65, 0x73, 0x49, 0x6e, 0x44, 0x69, 0x72, 0x48, 0x00, 0x52, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x73,
	0x49, 0x6e, 0x44, 0x69, 0x72, 0x12, 0x44, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x74,
	0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x50, 0x61, 0x74, 0x68, 0x48, 0x00, 0x52, 0x0b,
	0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x50, 0x61, 0x74, 0x68, 0x12, 0x52, 0x0a, 0x12, 0x75,
	0x6e, 0x69, 0x78, 0x5f, 0x65, 0x6e, 0x76, 0x5f, 0x76, 0x61, 0x72, 0x5f, 0x70, 0x61, 0x74, 0x68,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x74,
	0x6f, 0x61, 0x73, 0x74, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x74, 0x2e, 0x55, 0x6e, 0x69,
	0x78, 0x45, 0x6e, 0x76, 0x56, 0x61, 0x72, 0x50, 0x61, 0x74, 0x68, 0x73, 0x48, 0x00, 0x52, 0x0f,
	0x75, 0x6e, 0x69, 0x78, 0x45, 0x6e, 0x76, 0x56, 0x61, 0x72, 0x50, 0x61, 0x74, 0x68, 0x73, 0x1a,
	0x20, 0x0a, 0x0a, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74,
	0x68, 0x1a, 0xfe, 0x01, 0x0a, 0x0a, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x49, 0x6e, 0x44, 0x69, 0x72,
	0x12, 0x19, 0x0a, 0x08, 0x64, 0x69, 0x72, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x64, 0x69, 0x72, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x72,
	0x65, 0x63, 0x75, 0x72, 0x73, 0x69, 0x76, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09,
	0x72, 0x65, 0x63, 0x75, 0x72, 0x73, 0x69, 0x76, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x4f, 0x6e, 0x6c, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x73,
	0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69, 0x72,
	0x73, 0x4f, 0x6e, 0x6c, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x6b, 0x69, 0x70, 0x5f, 0x73, 0x79,
	0x6d, 0x6c, 0x69, 0x6e, 0x6b, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x73, 0x6b,
	0x69, 0x70, 0x53, 0x79, 0x6d, 0x6c, 0x69, 0x6e, 0x6b, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x66, 0x69,
	0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x72, 0x65, 0x67, 0x65, 0x78, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x67, 0x65,
	0x78, 0x12, 0x2f, 0x0a, 0x14, 0x6f, 0x70, 0x74, 0x5f, 0x6f, 0x75, 0x74, 0x5f, 0x70, 0x61, 0x74,
	0x68, 0x5f, 0x72, 0x65, 0x67, 0x65, 0x78, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x11, 0x6f, 0x70, 0x74, 0x4f, 0x75, 0x74, 0x50, 0x61, 0x74, 0x68, 0x52, 0x65, 0x67, 0x65, 0x78,
	0x65, 0x73, 0x1a, 0x47, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x50, 0x61, 0x74,
	0x68, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x63, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x1a, 0x68, 0x0a, 0x0f, 0x55,
	0x6e, 0x69, 0x78, 0x45, 0x6e, 0x76, 0x56, 0x61, 0x72, 0x50, 0x61, 0x74, 0x68, 0x73, 0x12, 0x19,
	0x0a, 0x08, 0x76, 0x61, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x76, 0x61, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x4f, 0x6e, 0x6c, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x73,
	0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69, 0x72,
	0x73, 0x4f, 0x6e, 0x6c, 0x79, 0x42, 0x0b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x70, 0x61,
	0x74, 0x68, 0x22, 0xf1, 0x01, 0x0a, 0x08, 0x53, 0x51, 0x4c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12,
	0x49, 0x0a, 0x0f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61,
	0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c,
	0x74, 0x6f, 0x61, 0x73, 0x74, 0x2e, 0x53, 0x51, 0x4c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x53,
	0x51, 0x4c, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x52, 0x0e, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79,
	0x12, 0x25, 0x0a, 0x0e, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x2c, 0x0a, 0x12, 0x6e, 0x6f, 0x6e, 0x5f, 0x63,
	0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x10, 0x6e, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x61, 0x6e,
	0x63, 0x65, 0x4d, 0x73, 0x67, 0x22, 0x2f, 0x0a, 0x0b, 0x53, 0x51, 0x4c, 0x44, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x0e, 0x44, 0x42, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x42, 0x5f, 0x4d,
	0x59, 0x53, 0x51, 0x4c, 0x10, 0x01, 0x42, 0x4a, 0x5a, 0x48, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x6c, 0x6f, 0x63, 0x61,
	0x6c, 0x74, 0x6f, 0x61, 0x73, 0x74, 0x2f, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x6c, 0x69,
	0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x63, 0x61, 0x6e, 0x5f, 0x69, 0x6e, 0x73,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x67, 0x6f, 0x5f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_scannerlib_proto_scan_instructions_proto_rawDescOnce sync.Once
	file_scannerlib_proto_scan_instructions_proto_rawDescData = file_scannerlib_proto_scan_instructions_proto_rawDesc
)

func file_scannerlib_proto_scan_instructions_proto_rawDescGZIP() []byte {
	file_scannerlib_proto_scan_instructions_proto_rawDescOnce.Do(func() {
		file_scannerlib_proto_scan_instructions_proto_rawDescData = protoimpl.X.CompressGZIP(file_scannerlib_proto_scan_instructions_proto_rawDescData)
	})
	return file_scannerlib_proto_scan_instructions_proto_rawDescData
}

var file_scannerlib_proto_scan_instructions_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_scannerlib_proto_scan_instructions_proto_msgTypes = make([]protoimpl.MessageInfo, 21)
var file_scannerlib_proto_scan_instructions_proto_goTypes = []interface{}{
	(RepeatConfig_RepeatType)(0),            // 0: localtoast.RepeatConfig.RepeatType
	(PermissionCheck_BitMatchCriterion)(0),  // 1: localtoast.PermissionCheck.BitMatchCriterion
	(ContentEntryCheck_MatchType)(0),        // 2: localtoast.ContentEntryCheck.MatchType
	(GroupCriterion_Type)(0),                // 3: localtoast.GroupCriterion.Type
	(SQLCheck_SQLDatabase)(0),               // 4: localtoast.SQLCheck.SQLDatabase
	(*BenchmarkScanInstructionDef)(nil),     // 5: localtoast.BenchmarkScanInstructionDef
	(*ScanTypeSpecificInstruction)(nil),     // 6: localtoast.ScanTypeSpecificInstruction
	(*BenchmarkScanInstruction)(nil),        // 7: localtoast.BenchmarkScanInstruction
	(*CheckAlternative)(nil),                // 8: localtoast.CheckAlternative
	(*FileCheck)(nil),                       // 9: localtoast.FileCheck
	(*RepeatConfig)(nil),                    // 10: localtoast.RepeatConfig
	(*ExistenceCheck)(nil),                  // 11: localtoast.ExistenceCheck
	(*PermissionCheck)(nil),                 // 12: localtoast.PermissionCheck
	(*ContentCheck)(nil),                    // 13: localtoast.ContentCheck
	(*ContentEntryCheck)(nil),               // 14: localtoast.ContentEntryCheck
	(*MatchCriterion)(nil),                  // 15: localtoast.MatchCriterion
	(*GroupCriterion)(nil),                  // 16: localtoast.GroupCriterion
	(*FileSet)(nil),                         // 17: localtoast.FileSet
	(*SQLCheck)(nil),                        // 18: localtoast.SQLCheck
	(*RepeatConfig_OptOutSubstitution)(nil), // 19: localtoast.RepeatConfig.OptOutSubstitution
	(*PermissionCheck_OwnerCheck)(nil),      // 20: localtoast.PermissionCheck.OwnerCheck
	(*GroupCriterion_Today)(nil),            // 21: localtoast.GroupCriterion.Today
	(*FileSet_SingleFile)(nil),              // 22: localtoast.FileSet.SingleFile
	(*FileSet_FilesInDir)(nil),              // 23: localtoast.FileSet.FilesInDir
	(*FileSet_ProcessPath)(nil),             // 24: localtoast.FileSet.ProcessPath
	(*FileSet_UnixEnvVarPaths)(nil),         // 25: localtoast.FileSet.UnixEnvVarPaths
}
var file_scannerlib_proto_scan_instructions_proto_depIdxs = []int32{
	7,  // 0: localtoast.BenchmarkScanInstructionDef.generic:type_name -> localtoast.BenchmarkScanInstruction
	6,  // 1: localtoast.BenchmarkScanInstructionDef.scan_type_specific:type_name -> localtoast.ScanTypeSpecificInstruction
	7,  // 2: localtoast.ScanTypeSpecificInstruction.instance_scanning:type_name -> localtoast.BenchmarkScanInstruction
	7,  // 3: localtoast.ScanTypeSpecificInstruction.image_scanning:type_name -> localtoast.BenchmarkScanInstruction
	8,  // 4: localtoast.BenchmarkScanInstruction.check_alternatives:type_name -> localtoast.CheckAlternative
	9,  // 5: localtoast.CheckAlternative.file_checks:type_name -> localtoast.FileCheck
	18, // 6: localtoast.CheckAlternative.sql_checks:type_name -> localtoast.SQLCheck
	17, // 7: localtoast.FileCheck.files_to_check:type_name -> localtoast.FileSet
	11, // 8: localtoast.FileCheck.existence:type_name -> localtoast.ExistenceCheck
	12, // 9: localtoast.FileCheck.permission:type_name -> localtoast.PermissionCheck
	13, // 10: localtoast.FileCheck.content:type_name -> localtoast.ContentCheck
	14, // 11: localtoast.FileCheck.content_entry:type_name -> localtoast.ContentEntryCheck
	10, // 12: localtoast.FileCheck.repeat_config:type_name -> localtoast.RepeatConfig
	0,  // 13: localtoast.RepeatConfig.type:type_name -> localtoast.RepeatConfig.RepeatType
	19, // 14: localtoast.RepeatConfig.opt_out:type_name -> localtoast.RepeatConfig.OptOutSubstitution
	1,  // 15: localtoast.PermissionCheck.bits_should_match:type_name -> localtoast.PermissionCheck.BitMatchCriterion
	20, // 16: localtoast.PermissionCheck.user:type_name -> localtoast.PermissionCheck.OwnerCheck
	20, // 17: localtoast.PermissionCheck.group:type_name -> localtoast.PermissionCheck.OwnerCheck
	2,  // 18: localtoast.ContentEntryCheck.match_type:type_name -> localtoast.ContentEntryCheck.MatchType
	15, // 19: localtoast.ContentEntryCheck.match_criteria:type_name -> localtoast.MatchCriterion
	16, // 20: localtoast.MatchCriterion.group_criteria:type_name -> localtoast.GroupCriterion
	3,  // 21: localtoast.GroupCriterion.type:type_name -> localtoast.GroupCriterion.Type
	21, // 22: localtoast.GroupCriterion.today:type_name -> localtoast.GroupCriterion.Today
	22, // 23: localtoast.FileSet.single_file:type_name -> localtoast.FileSet.SingleFile
	23, // 24: localtoast.FileSet.files_in_dir:type_name -> localtoast.FileSet.FilesInDir
	24, // 25: localtoast.FileSet.process_path:type_name -> localtoast.FileSet.ProcessPath
	25, // 26: localtoast.FileSet.unix_env_var_paths:type_name -> localtoast.FileSet.UnixEnvVarPaths
	4,  // 27: localtoast.SQLCheck.target_database:type_name -> localtoast.SQLCheck.SQLDatabase
	28, // [28:28] is the sub-list for method output_type
	28, // [28:28] is the sub-list for method input_type
	28, // [28:28] is the sub-list for extension type_name
	28, // [28:28] is the sub-list for extension extendee
	0,  // [0:28] is the sub-list for field type_name
}

func init() { file_scannerlib_proto_scan_instructions_proto_init() }
func file_scannerlib_proto_scan_instructions_proto_init() {
	if File_scannerlib_proto_scan_instructions_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_scannerlib_proto_scan_instructions_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BenchmarkScanInstructionDef); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanTypeSpecificInstruction); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BenchmarkScanInstruction); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckAlternative); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileCheck); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepeatConfig); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExistenceCheck); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PermissionCheck); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContentCheck); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContentEntryCheck); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MatchCriterion); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupCriterion); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileSet); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SQLCheck); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepeatConfig_OptOutSubstitution); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[15].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PermissionCheck_OwnerCheck); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[16].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupCriterion_Today); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[17].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileSet_SingleFile); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[18].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileSet_FilesInDir); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[19].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileSet_ProcessPath); i {
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
		file_scannerlib_proto_scan_instructions_proto_msgTypes[20].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileSet_UnixEnvVarPaths); i {
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
	file_scannerlib_proto_scan_instructions_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*BenchmarkScanInstructionDef_Generic)(nil),
		(*BenchmarkScanInstructionDef_ScanTypeSpecific)(nil),
	}
	file_scannerlib_proto_scan_instructions_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*FileCheck_Existence)(nil),
		(*FileCheck_Permission)(nil),
		(*FileCheck_Content)(nil),
		(*FileCheck_ContentEntry)(nil),
	}
	file_scannerlib_proto_scan_instructions_proto_msgTypes[11].OneofWrappers = []interface{}{
		(*GroupCriterion_Const)(nil),
		(*GroupCriterion_Today_)(nil),
	}
	file_scannerlib_proto_scan_instructions_proto_msgTypes[12].OneofWrappers = []interface{}{
		(*FileSet_SingleFile_)(nil),
		(*FileSet_FilesInDir_)(nil),
		(*FileSet_ProcessPath_)(nil),
		(*FileSet_UnixEnvVarPaths_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_scannerlib_proto_scan_instructions_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   21,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_scannerlib_proto_scan_instructions_proto_goTypes,
		DependencyIndexes: file_scannerlib_proto_scan_instructions_proto_depIdxs,
		EnumInfos:         file_scannerlib_proto_scan_instructions_proto_enumTypes,
		MessageInfos:      file_scannerlib_proto_scan_instructions_proto_msgTypes,
	}.Build()
	File_scannerlib_proto_scan_instructions_proto = out.File
	file_scannerlib_proto_scan_instructions_proto_rawDesc = nil
	file_scannerlib_proto_scan_instructions_proto_goTypes = nil
	file_scannerlib_proto_scan_instructions_proto_depIdxs = nil
}
