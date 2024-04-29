#!/bin/sh
rm -rf scannerlib/proto/*_go_proto

# Install and prepare Grafeas.
if [ ! -e scannerlib/proto/v1 ]; then
  wget --no-verbose https://github.com/grafeas/grafeas/archive/220ed72376f81d0dd5233839d22c5627eb8d9494.tar.gz
  tar -xf 220ed72376f81d0dd5233839d22c5627eb8d9494.tar.gz
  mv grafeas-220ed72376f81d0dd5233839d22c5627eb8d9494/proto/v1 scannerlib/proto
  rm -r *220ed72376f81d0dd5233839d22c5627eb8d9494*
fi

sed -i 's\option go_package = ".*";\option go_package = "github.com/google/localtoast/scannerlib/proto/compliance_go_proto";\g' scannerlib/proto/v1/compliance.proto
sed -i 's\option go_package = ".*";\option go_package = "github.com/google/localtoast/scannerlib/proto/severity_go_proto";\g' scannerlib/proto/v1/severity.proto

# Compile protos.
protoc -I=scannerlib --go_out=scannerlib/proto scannerlib/proto/*.proto scannerlib/proto/v1/compliance.proto scannerlib/proto/v1/severity.proto

# Clean up.
mv scannerlib/proto/github.com/google/localtoast/scannerlib/proto/* scannerlib/proto/
rm -r scannerlib/proto/github.com
