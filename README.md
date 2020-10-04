# Template Arhat Extension Go

[![CI](https://github.com/arhat-dev/template-arhat-ext-go/workflows/CI/badge.svg)](https://github.com/arhat-dev/template-arhat-ext-go/actions?query=workflow%3ACI)
[![Build](https://github.com/arhat-dev/template-arhat-ext-go/workflows/Build/badge.svg)](https://github.com/arhat-dev/template-arhat-ext-go/actions?query=workflow%3ABuild)
[![PkgGoDev](https://pkg.go.dev/badge/arhat.dev/template-arhat-ext-go)](https://pkg.go.dev/arhat.dev/template-arhat-ext-go)
[![GoReportCard](https://goreportcard.com/badge/arhat.dev/template-arhat-ext-go)](https://goreportcard.com/report/arhat.dev/template-arhat-ext-go)
[![codecov](https://codecov.io/gh/arhat-dev/template-arhat-ext-go/branch/master/graph/badge.svg)](https://codecov.io/gh/arhat-dev/template-arhat-ext-go)

Template repo for extensions integrating with [`arhat`](https://github.com/arhat-dev/arhat) via [`arhat-proto`](https://github.com/arhat-dev/arhat-proto) in Go

## Usage

1. Create a new repo using this template
2. Rename `template-arhat-ext-go` (most importantly, the module name in `go.mod` file) according to your preference
3. Update application configuration definition in [`pkg/conf`](./pkg/conf/)
4. Implement your extension by updating code in [`pkg/device`](./pkg/device/)
5. Update deployment charts in [`cicd/deploy/charts/`](./cicd/deploy/charts/template-arhat-ext-go/)
6. Document supported device params and metrics (and their params) in [`docs/device.md`](./docs/device.md)
7. Deploy to somewhere it can communicate with `arhat` (host or container)

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
