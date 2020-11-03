# Template Go

[![CI](https://github.com/arhat-ext/template-go/workflows/CI/badge.svg)](https://github.com/arhat-ext/template-go/actions?query=workflow%3ACI)
[![Build](https://github.com/arhat-ext/template-go/workflows/Build/badge.svg)](https://github.com/arhat-ext/template-go/actions?query=workflow%3ABuild)
[![PkgGoDev](https://pkg.go.dev/badge/ext.arhat.dev/template-go)](https://pkg.go.dev/ext.arhat.dev/template-go)
[![GoReportCard](https://goreportcard.com/badge/ext.arhat.dev/template-go)](https://goreportcard.com/report/ext.arhat.dev/template-go)
[![codecov](https://codecov.io/gh/arhat-ext/template-go/branch/master/graph/badge.svg)](https://codecov.io/gh/arhat-ext/template-go)

Template repo for extensions integrating with [`arhat`](https://github.com/arhat-dev/arhat) via [`arhat-proto`](https://github.com/arhat-dev/arhat-proto) using [`arhat.dev/libext`](https://github.com/arhat-dev/libext-go) in Go

## Usage

1. Create a new repo using this template
2. Rename `template-go` (most importantly, the module name in `go.mod` file) according to your preference
3. Update application configuration definition in [`pkg/conf`](./pkg/conf/)
4. Update extension name `my-extension-name` in [`pkg/cmd`](./pkg/cmd/template-go.go)
5. Implement your extension by updating code in `pkg/<some name>`
   - e.g. peripheral extension: in [`pkg/peripheral`](./pkg/peripheral/)
6. Update deployment charts in [`cicd/deploy/charts/`](./cicd/deploy/charts/)
7. Document supported features [`docs`](./docs)
   - e.g. peripheral extension: in [`docs/peripheral.md`](./docs/peripheral.md)
8. Deploy to somewhere it can communicate with `arhat` (host or container)

## LICENSE

```text
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
