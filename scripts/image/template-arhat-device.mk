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

# build
image.build.template-arhat-device.linux.x86:
	sh scripts/image/build.sh $@

image.build.template-arhat-device.linux.amd64:
	sh scripts/image/build.sh $@

image.build.template-arhat-device.linux.armv6:
	sh scripts/image/build.sh $@

image.build.template-arhat-device.linux.armv7:
	sh scripts/image/build.sh $@

image.build.template-arhat-device.linux.arm64:
	sh scripts/image/build.sh $@

image.build.template-arhat-device.linux.ppc64le:
	sh scripts/image/build.sh $@

image.build.template-arhat-device.linux.s390x:
	sh scripts/image/build.sh $@

image.build.template-arhat-device.linux.all: \
	image.build.template-arhat-device.linux.amd64 \
	image.build.template-arhat-device.linux.arm64 \
	image.build.template-arhat-device.linux.armv7 \
	image.build.template-arhat-device.linux.armv6 \
	image.build.template-arhat-device.linux.x86 \
	image.build.template-arhat-device.linux.s390x \
	image.build.template-arhat-device.linux.ppc64le

image.build.template-arhat-device.windows.amd64:
	sh scripts/image/build.sh $@

image.build.template-arhat-device.windows.armv7:
	sh scripts/image/build.sh $@

image.build.template-arhat-device.windows.all: \
	image.build.template-arhat-device.windows.amd64 \
	image.build.template-arhat-device.windows.armv7

# push
image.push.template-arhat-device.linux.x86:
	sh scripts/image/push.sh $@

image.push.template-arhat-device.linux.amd64:
	sh scripts/image/push.sh $@

image.push.template-arhat-device.linux.armv6:
	sh scripts/image/push.sh $@

image.push.template-arhat-device.linux.armv7:
	sh scripts/image/push.sh $@

image.push.template-arhat-device.linux.arm64:
	sh scripts/image/push.sh $@

image.push.template-arhat-device.linux.ppc64le:
	sh scripts/image/push.sh $@

image.push.template-arhat-device.linux.s390x:
	sh scripts/image/push.sh $@

image.push.template-arhat-device.linux.all: \
	image.push.template-arhat-device.linux.amd64 \
	image.push.template-arhat-device.linux.arm64 \
	image.push.template-arhat-device.linux.armv7 \
	image.push.template-arhat-device.linux.armv6 \
	image.push.template-arhat-device.linux.x86 \
	image.push.template-arhat-device.linux.s390x \
	image.push.template-arhat-device.linux.ppc64le

image.push.template-arhat-device.windows.amd64:
	sh scripts/image/push.sh $@

image.push.template-arhat-device.windows.armv7:
	sh scripts/image/push.sh $@

image.push.template-arhat-device.windows.all: \
	image.push.template-arhat-device.windows.amd64 \
	image.push.template-arhat-device.windows.armv7
