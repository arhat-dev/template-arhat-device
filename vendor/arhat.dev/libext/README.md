# libext

[![CI](https://github.com/arhat-dev/libext-go/workflows/CI/badge.svg)](https://github.com/arhat-dev/libext-go/actions?query=workflow%3ACI)
[![PkgGoDev](https://pkg.go.dev/badge/arhat.dev/libext)](https://pkg.go.dev/arhat.dev/libext)
[![GoReportCard](https://goreportcard.com/badge/arhat.dev/libext)](https://goreportcard.com/report/arhat.dev/libext)
[![codecov](https://codecov.io/gh/arhat-dev/libext-go/branch/master/graph/badge.svg)](https://codecov.io/gh/arhat-dev/libext-go)

Library for building [`arhat`](https://github.com/arhat-dev/arhat) extensions in Go

## Support Matrix

- Protocols: `tcp`, `tcp-tls`, `udp`, `udp-dtls`, `unix`, `unix-tls`, `pipe`, `pipe-tls`
- Codec: `json`, `protobuf`

please refer to [Benchmark](#benchmark) for performance evaluation of different combinations

## Usage

### Controller for Extension (Client)

__TL;DR__: create a project using this template: [`arhat-ext/template-go`](https://github.com/arhat-ext/template-go) or have a look at [`examples/client_example_test.go`](./examples/client_example_test.go)

### Extension Hub (Server)

__TL;DR__: have a look at [`examples/server_example_test.go`](./examples/server_example_test.go)

### Benchmark

You can find reference benchmark results in CI log, to run it locally:

```bash
make test.benchmark
```

## LICENSE

```txt
Copyright 2020 The arhat.dev Authors.

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
