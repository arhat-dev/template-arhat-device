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

# native
template-arhat-device:
	sh scripts/build/build.sh $@

# linux
template-arhat-device.linux.x86:
	sh scripts/build/build.sh $@

template-arhat-device.linux.amd64:
	sh scripts/build/build.sh $@

template-arhat-device.linux.armv5:
	sh scripts/build/build.sh $@

template-arhat-device.linux.armv6:
	sh scripts/build/build.sh $@

template-arhat-device.linux.armv7:
	sh scripts/build/build.sh $@

template-arhat-device.linux.arm64:
	sh scripts/build/build.sh $@

template-arhat-device.linux.mips:
	sh scripts/build/build.sh $@

template-arhat-device.linux.mipshf:
	sh scripts/build/build.sh $@

template-arhat-device.linux.mipsle:
	sh scripts/build/build.sh $@

template-arhat-device.linux.mipslehf:
	sh scripts/build/build.sh $@

template-arhat-device.linux.mips64:
	sh scripts/build/build.sh $@

template-arhat-device.linux.mips64hf:
	sh scripts/build/build.sh $@

template-arhat-device.linux.mips64le:
	sh scripts/build/build.sh $@

template-arhat-device.linux.mips64lehf:
	sh scripts/build/build.sh $@

template-arhat-device.linux.ppc64:
	sh scripts/build/build.sh $@

template-arhat-device.linux.ppc64le:
	sh scripts/build/build.sh $@

template-arhat-device.linux.s390x:
	sh scripts/build/build.sh $@

template-arhat-device.linux.riscv64:
	sh scripts/build/build.sh $@

template-arhat-device.linux.all: \
	template-arhat-device.linux.x86 \
	template-arhat-device.linux.amd64 \
	template-arhat-device.linux.armv5 \
	template-arhat-device.linux.armv6 \
	template-arhat-device.linux.armv7 \
	template-arhat-device.linux.arm64 \
	template-arhat-device.linux.mips \
	template-arhat-device.linux.mipshf \
	template-arhat-device.linux.mipsle \
	template-arhat-device.linux.mipslehf \
	template-arhat-device.linux.mips64 \
	template-arhat-device.linux.mips64hf \
	template-arhat-device.linux.mips64le \
	template-arhat-device.linux.mips64lehf \
	template-arhat-device.linux.ppc64 \
	template-arhat-device.linux.ppc64le \
	template-arhat-device.linux.s390x \
	template-arhat-device.linux.riscv64

template-arhat-device.darwin.amd64:
	sh scripts/build/build.sh $@

# # currently darwin/arm64 build will fail due to golang link error
# template-arhat-device.darwin.arm64:
# 	sh scripts/build/build.sh $@

template-arhat-device.darwin.all: \
	template-arhat-device.darwin.amd64

template-arhat-device.windows.x86:
	sh scripts/build/build.sh $@

template-arhat-device.windows.amd64:
	sh scripts/build/build.sh $@

template-arhat-device.windows.armv5:
	sh scripts/build/build.sh $@

template-arhat-device.windows.armv6:
	sh scripts/build/build.sh $@

template-arhat-device.windows.armv7:
	sh scripts/build/build.sh $@

# # currently no support for windows/arm64
# template-arhat-device.windows.arm64:
# 	sh scripts/build/build.sh $@

template-arhat-device.windows.all: \
	template-arhat-device.windows.x86 \
	template-arhat-device.windows.amd64

# # android build requires android sdk
# template-arhat-device.android.amd64:
# 	sh scripts/build/build.sh $@

# template-arhat-device.android.x86:
# 	sh scripts/build/build.sh $@

# template-arhat-device.android.armv5:
# 	sh scripts/build/build.sh $@

# template-arhat-device.android.armv6:
# 	sh scripts/build/build.sh $@

# template-arhat-device.android.armv7:
# 	sh scripts/build/build.sh $@

# template-arhat-device.android.arm64:
# 	sh scripts/build/build.sh $@

# template-arhat-device.android.all: \
# 	template-arhat-device.android.amd64 \
# 	template-arhat-device.android.arm64 \
# 	template-arhat-device.android.x86 \
# 	template-arhat-device.android.armv7 \
# 	template-arhat-device.android.armv5 \
# 	template-arhat-device.android.armv6

template-arhat-device.freebsd.amd64:
	sh scripts/build/build.sh $@

template-arhat-device.freebsd.x86:
	sh scripts/build/build.sh $@

template-arhat-device.freebsd.armv5:
	sh scripts/build/build.sh $@

template-arhat-device.freebsd.armv6:
	sh scripts/build/build.sh $@

template-arhat-device.freebsd.armv7:
	sh scripts/build/build.sh $@

template-arhat-device.freebsd.arm64:
	sh scripts/build/build.sh $@

template-arhat-device.freebsd.all: \
	template-arhat-device.freebsd.amd64 \
	template-arhat-device.freebsd.arm64 \
	template-arhat-device.freebsd.armv7 \
	template-arhat-device.freebsd.x86 \
	template-arhat-device.freebsd.armv5 \
	template-arhat-device.freebsd.armv6

template-arhat-device.netbsd.amd64:
	sh scripts/build/build.sh $@

template-arhat-device.netbsd.x86:
	sh scripts/build/build.sh $@

template-arhat-device.netbsd.armv5:
	sh scripts/build/build.sh $@

template-arhat-device.netbsd.armv6:
	sh scripts/build/build.sh $@

template-arhat-device.netbsd.armv7:
	sh scripts/build/build.sh $@

template-arhat-device.netbsd.arm64:
	sh scripts/build/build.sh $@

template-arhat-device.netbsd.all: \
	template-arhat-device.netbsd.amd64 \
	template-arhat-device.netbsd.arm64 \
	template-arhat-device.netbsd.armv7 \
	template-arhat-device.netbsd.x86 \
	template-arhat-device.netbsd.armv5 \
	template-arhat-device.netbsd.armv6

template-arhat-device.openbsd.amd64:
	sh scripts/build/build.sh $@

template-arhat-device.openbsd.x86:
	sh scripts/build/build.sh $@

template-arhat-device.openbsd.armv5:
	sh scripts/build/build.sh $@

template-arhat-device.openbsd.armv6:
	sh scripts/build/build.sh $@

template-arhat-device.openbsd.armv7:
	sh scripts/build/build.sh $@

template-arhat-device.openbsd.arm64:
	sh scripts/build/build.sh $@

template-arhat-device.openbsd.all: \
	template-arhat-device.openbsd.amd64 \
	template-arhat-device.openbsd.arm64 \
	template-arhat-device.openbsd.armv7 \
	template-arhat-device.openbsd.x86 \
	template-arhat-device.openbsd.armv5 \
	template-arhat-device.openbsd.armv6

template-arhat-device.solaris.amd64:
	sh scripts/build/build.sh $@

template-arhat-device.aix.ppc64:
	sh scripts/build/build.sh $@

template-arhat-device.dragonfly.amd64:
	sh scripts/build/build.sh $@

template-arhat-device.plan9.amd64:
	sh scripts/build/build.sh $@

template-arhat-device.plan9.x86:
	sh scripts/build/build.sh $@

template-arhat-device.plan9.armv5:
	sh scripts/build/build.sh $@

template-arhat-device.plan9.armv6:
	sh scripts/build/build.sh $@

template-arhat-device.plan9.armv7:
	sh scripts/build/build.sh $@

template-arhat-device.plan9.all: \
	template-arhat-device.plan9.amd64 \
	template-arhat-device.plan9.armv7 \
	template-arhat-device.plan9.x86 \
	template-arhat-device.plan9.armv5 \
	template-arhat-device.plan9.armv6
