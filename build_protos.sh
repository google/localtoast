#!/bin/sh
rm -rf scannerlib/proto/*_go_proto

# Install and prepare Grafeas.
if [ ! -e scannerlib/proto/v1 ]; then
  wget https://github.com/grafeas/grafeas/archive/53c88dd9ce08ce959b308b41597cbfce5d999ede.tar.gz
  tar -xf 53c88dd9ce08ce959b308b41597cbfce5d999ede.tar.gz
  mv grafeas-53c88dd9ce08ce959b308b41597cbfce5d999ede/proto/v1 scannerlib/proto
  rm -r *53c88dd9ce08ce959b308b41597cbfce5d999ede*
fi

sed -i 's\option go_package = ".*";\option go_package = "github.com/google/localtoast/scannerlib/proto/compliance_go_proto";\g' scannerlib/proto/v1/compliance.proto
sed -i 's\option go_package = ".*";\option go_package = "github.com/google/localtoast/scannerlib/proto/severity_go_proto";\g' scannerlib/proto/v1/severity.proto

# Compile protos.
protoc -I=scannerlib --go_out=scannerlib/proto scannerlib/proto/*.proto scannerlib/proto/v1/compliance.proto scannerlib/proto/v1/severity.proto

# Clean up.
mv scannerlib/proto/github.com/google/localtoast/scannerlib/proto/* scannerlib/proto/
rm -r scannerlib/proto/github.com
