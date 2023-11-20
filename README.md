 # Localtoast
Localtoast is a scanner for running security-related configuration checks such as [CIS benchmarks](https://www.cisecurity.org/cis-benchmarks) in an easily configurable manner.

The scanner can either be used as a standalone binary to scan the local machine or as a library with a custom wrapper to perform scans on e.g. container images or remote hosts.

## How to use

### As a standalone binary:

1. Install the [build deps](#build-dependencies)
2. `make`
3. `sudo ./localtoast --config=configs/example.textproto --result=scan-result.textproto`


#### Build and use OS-specific configs:
1. `make configs`
2. `sudo ./localtoast --config=configs/full/cos_97/instance_scanning.textproto --result=scan-result.textproto`

#### Build and run Localtoast with SQL scanning capabilities:
1. `make configs`
2. `make localtoast_sql`
3. `sudo localtoast_sql/localtoast_sql --config=configs/full/cassandra-cql/instance_scanning.textproto --result=scan-result.textproto --cassandra-database=localhost:9042`

### As a library:
1. Import `github.com/google/localtoast/scannerlib` and `github.com/google/localtoast/scanapi` into your Go project
2. Write a custom implementation for the `scanapi.ScanAPI` interface
3. Call `scannerlib.Scanner{}.Scan()` with the appropriate config and the implementation

See the [scan config](scannerlib/proto/api.proto) and [result](scannerlib/proto/scan_instructions.proto) protos for details on the input+output format.

## Defining custom checks
To add your own checks to a scan config,

1. Define the check in one of the [definition files](configs/defs/cos.textproto)
  * [Example](https://github.com/google/localtoast/commit/9c39a52cef30f7ad773b74a38ac9ffa7c4998ca3#diff-1350df51e73d56ca08a90aa7fc47a3032a41d85a7fe5a8b8707387000f43c0be)
  * See the [instruction proto](scannerlib/proto/scan_instructions.proto) for details on the instruction syntax
2. Add a reference to the check in [the scan config](configs/reduced/cos_97/instance_scanning.textproto) you want to extend
  * [Example](https://github.com/google/localtoast/commit/9c39a52cef30f7ad773b74a38ac9ffa7c4998ca3#diff-094e7befebe2acf9321eb3406fbb81af2880344086fe40dc97c3d4d915fe0e6e)
3. Re-build the config file with `make configs`
4. Use the re-generated config file in your scans, e.g. `sudo ./localtoast --config=configs/full/cos_97/instance_scanning.textproto --result=scan-result.textproto`

## Build dependencies
To build Localtoast, you'll need to have the following installed:
* `go`: Follow https://go.dev/doc/install
* `protoc`: Install the appropriate package, e.g. `apt install protobuf-compiler`
* `protoc-gen-go`: Run `go install google.golang.org/protobuf/cmd/protoc-gen-go`

## Contributing
Read how to [contribute to Localtoast](CONTRIBUTING.md).

## License
Localtoast is released under the [Apache 2.0 license](LICENSE).

```
Copyright 2021 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Disclaimers

Localtoast is not an official Google product.
