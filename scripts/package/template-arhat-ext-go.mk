# Copyright 2020 The arhat.dev Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# linux
#
package.template-arhat-ext-go.deb.amd64:
	sh scripts/package/package.sh $@

package.template-arhat-ext-go.deb.armv6:
	sh scripts/package/package.sh $@

package.template-arhat-ext-go.deb.armv7:
	sh scripts/package/package.sh $@

package.template-arhat-ext-go.deb.arm64:
	sh scripts/package/package.sh $@

package.template-arhat-ext-go.deb.all: \
	package.template-arhat-ext-go.deb.amd64 \
	package.template-arhat-ext-go.deb.armv6 \
	package.template-arhat-ext-go.deb.armv7 \
	package.template-arhat-ext-go.deb.arm64

package.template-arhat-ext-go.rpm.amd64:
	sh scripts/package/package.sh $@

package.template-arhat-ext-go.rpm.armv7:
	sh scripts/package/package.sh $@

package.template-arhat-ext-go.rpm.arm64:
	sh scripts/package/package.sh $@

package.template-arhat-ext-go.rpm.all: \
	package.template-arhat-ext-go.rpm.amd64 \
	package.template-arhat-ext-go.rpm.armv7 \
	package.template-arhat-ext-go.rpm.arm64

package.template-arhat-ext-go.linux.all: \
	package.template-arhat-ext-go.deb.all \
	package.template-arhat-ext-go.rpm.all

#
# windows
#

package.template-arhat-ext-go.msi.amd64:
	sh scripts/package/package.sh $@

package.template-arhat-ext-go.msi.arm64:
	sh scripts/package/package.sh $@

package.template-arhat-ext-go.msi.all: \
	package.template-arhat-ext-go.msi.amd64 \
	package.template-arhat-ext-go.msi.arm64

package.template-arhat-ext-go.windows.all: \
	package.template-arhat-ext-go.msi.all

#
# darwin
#

package.template-arhat-ext-go.pkg.amd64:
	sh scripts/package/package.sh $@

package.template-arhat-ext-go.pkg.arm64:
	sh scripts/package/package.sh $@

package.template-arhat-ext-go.pkg.all: \
	package.template-arhat-ext-go.pkg.amd64 \
	package.template-arhat-ext-go.pkg.arm64

package.template-arhat-ext-go.darwin.all: \
	package.template-arhat-ext-go.pkg.all
