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
package.template-arhat-device.deb.amd64:
	sh scripts/package/package.sh $@

package.template-arhat-device.deb.armv6:
	sh scripts/package/package.sh $@

package.template-arhat-device.deb.armv7:
	sh scripts/package/package.sh $@

package.template-arhat-device.deb.arm64:
	sh scripts/package/package.sh $@

package.template-arhat-device.deb.all: \
	package.template-arhat-device.deb.amd64 \
	package.template-arhat-device.deb.armv6 \
	package.template-arhat-device.deb.armv7 \
	package.template-arhat-device.deb.arm64

package.template-arhat-device.rpm.amd64:
	sh scripts/package/package.sh $@

package.template-arhat-device.rpm.armv7:
	sh scripts/package/package.sh $@

package.template-arhat-device.rpm.arm64:
	sh scripts/package/package.sh $@

package.template-arhat-device.rpm.all: \
	package.template-arhat-device.rpm.amd64 \
	package.template-arhat-device.rpm.armv7 \
	package.template-arhat-device.rpm.arm64

package.template-arhat-device.linux.all: \
	package.template-arhat-device.deb.all \
	package.template-arhat-device.rpm.all

#
# windows
#

package.template-arhat-device.msi.amd64:
	sh scripts/package/package.sh $@

package.template-arhat-device.msi.arm64:
	sh scripts/package/package.sh $@

package.template-arhat-device.msi.all: \
	package.template-arhat-device.msi.amd64 \
	package.template-arhat-device.msi.arm64

package.template-arhat-device.windows.all: \
	package.template-arhat-device.msi.all

#
# darwin
#

package.template-arhat-device.pkg.amd64:
	sh scripts/package/package.sh $@

package.template-arhat-device.pkg.arm64:
	sh scripts/package/package.sh $@

package.template-arhat-device.pkg.all: \
	package.template-arhat-device.pkg.amd64 \
	package.template-arhat-device.pkg.arm64

package.template-arhat-device.darwin.all: \
	package.template-arhat-device.pkg.all
