 # Localtoast
Localtoast is a scanner for running security-related configuration checks such as [CIS benchmarks](https://www.cisecurity.org/cis-benchmarks) in an easily configurable manner.

The scanner can either be used as a standalone binary to scan the local machine or as a library with a custom wrapper to perform scans on e.g. container images or remote hosts.

## How to use

### As a standalone binary:

1. `make`
2. `./localtoast --config=configs/example.textproto --result=scan-result.textproto`

### As a library:
1. Import `github.com/google/localtoast/scannerlib` into your Go project
2. Write a custom implementation for the `ScanAPIProvider` interface
3. Call `scannerlib.Scanner{}.Scan()` with the appropriate config and the implementation

See the [scan config](scannerlib/proto/api.proto) and [result](scannerlib/proto/scan_instructions.proto) protos for details on the input+output format.

## Defining custom checks
To add your own checks to a scan config,

1. Define the check in one of the appropriate [config files](configs/cos_97/instance_scanning.textproto)
  * [Example](https://github.com/google/localtoast/commit/8aecb12ead5408c26905282ee6b441026ea264f0)
  * See the [instruction proto](scannerlib/proto/scan_instructions.proto) for details on the instruction syntax
2. Use the modified config file in your scans, e.g. `sudo ./localtoast --config=configs/cos_97/instance_scanning.textproto --result=scan-result.textproto`

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
