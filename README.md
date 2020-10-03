# Template Arhat Device Go

[![CI](https://github.com/arhat-dev/template-arhat-device/workflows/CI/badge.svg)](https://github.com/arhat-dev/template-arhat-device/actions?query=workflow%3ACI)
[![Build](https://github.com/arhat-dev/template-arhat-device/workflows/Build/badge.svg)](https://github.com/arhat-dev/template-arhat-device/actions?query=workflow%3ABuild)
[![PkgGoDev](https://pkg.go.dev/badge/arhat.dev/template-arhat-device)](https://pkg.go.dev/arhat.dev/template-arhat-device)
[![GoReportCard](https://goreportcard.com/badge/arhat.dev/template-arhat-device)](https://goreportcard.com/report/arhat.dev/template-arhat-device)
[![codecov](https://codecov.io/gh/arhat-dev/template-arhat-device/branch/master/graph/badge.svg)](https://codecov.io/gh/arhat-dev/template-arhat-device)

Template repo for devices integrating with [`arhat`](https://github.com/arhat-dev/arhat) via [`arhat-proto`](https://github.com/arhat-dev/arhat-proto) in Go

## Usage

1. Create a new repo using this template
2. Rename `template-arhat-device` (most importantly, the module name in `go.mod` file) according to your preference
3. Update application configuration definition in [`pkg/conf`](./pkg/conf/)
4. Implement your device by updating code in [`pkg/device`](./pkg/device/)
5. Update deployment charts in [`cicd/deploy/charts/`](./cicd)
6. Document supported device params and metrics (and their params) in `docs/device.md`
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
