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

// Package protofilehandler provides a utility for reading and writing protos
// into zipped or unzipped textproto or binproto files.
package protofilehandler

import (
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type fileType struct {
	isGZipped  bool
	isBinProto bool
}

func fileTypeForPath(filePath string) (*fileType, error) {
	parts := strings.Split(path.Base(filePath), ".")
	if len(parts) < 2 { // No extension
		return nil, errors.New("invalid filename: Doesn't have an extension")
	}
	isGZipped := false
	extension := parts[len(parts)-1]
	if extension == "gz" {
		isGZipped = true
		if len(parts) < 3 {
			return nil, errors.New("invalid filename: Gzipped file doesn't have an extension")
		}
		extension = parts[len(parts)-2]
	}
	isBinProto := false
	switch extension {
	case "binproto":
		isBinProto = true
	case "textproto":
		isBinProto = false
	default:
		return nil, errors.New("invalid filename: not a .textproto or .binproto")
	}
	return &fileType{isGZipped: isGZipped, isBinProto: isBinProto}, nil
}

// ReadProtoFromFile reads a proto message from a .textproto or .binproto file.
// If the file is gzipped, it's unzipped first.
func ReadProtoFromFile(filePath string, inputProto proto.Message) error {
	ft, err := fileTypeForPath(filePath)
	if err != nil {
		return err
	}

	var protoTxt []byte
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	if ft.isGZipped {
		reader, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer reader.Close()
		if protoTxt, err = io.ReadAll(reader); err != nil {
			return err
		}
	} else if protoTxt, err = io.ReadAll(f); err != nil {
		return err
	}

	if ft.isBinProto {
		if err := proto.Unmarshal(protoTxt, inputProto); err != nil {
			return err
		}
	} else if err := prototext.Unmarshal(protoTxt, inputProto); err != nil {
		return err
	}
	return nil
}

// WriteProtoToFile writes a proto message from a .textproto or .binproto file.
// If the file name additionally has the .gz suffix, it's zipped before writing.
func WriteProtoToFile(filePath string, outputProto proto.Message) error {
	ft, err := fileTypeForPath(filePath)
	if err != nil {
		return err
	}
	var protoTxt []byte
	if ft.isBinProto {
		if protoTxt, err = proto.Marshal(outputProto); err != nil {
			return err
		}
	} else {
		if protoTxt, err = (prototext.MarshalOptions{Multiline: true}.Marshal(outputProto)); err != nil {
			return err
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	if ft.isGZipped {
		writer := gzip.NewWriter(f)
		if _, err := writer.Write(protoTxt); err != nil {
			return err
		}
		if err := writer.Close(); err != nil {
			return err
		}
	} else if _, err := io.WriteString(f, string(protoTxt)); err != nil {
		return err
	}
	return nil
}
