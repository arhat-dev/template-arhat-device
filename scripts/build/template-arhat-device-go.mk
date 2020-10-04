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
template-arhat-device-go:
	sh scripts/build/build.sh $@

# linux
template-arhat-device-go.linux.x86:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.amd64:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.armv5:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.armv6:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.armv7:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.arm64:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.mips:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.mipshf:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.mipsle:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.mipslehf:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.mips64:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.mips64hf:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.mips64le:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.mips64lehf:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.ppc64:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.ppc64le:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.s390x:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.riscv64:
	sh scripts/build/build.sh $@

template-arhat-device-go.linux.all: \
	template-arhat-device-go.linux.x86 \
	template-arhat-device-go.linux.amd64 \
	template-arhat-device-go.linux.armv5 \
	template-arhat-device-go.linux.armv6 \
	template-arhat-device-go.linux.armv7 \
	template-arhat-device-go.linux.arm64 \
	template-arhat-device-go.linux.mips \
	template-arhat-device-go.linux.mipshf \
	template-arhat-device-go.linux.mipsle \
	template-arhat-device-go.linux.mipslehf \
	template-arhat-device-go.linux.mips64 \
	template-arhat-device-go.linux.mips64hf \
	template-arhat-device-go.linux.mips64le \
	template-arhat-device-go.linux.mips64lehf \
	template-arhat-device-go.linux.ppc64 \
	template-arhat-device-go.linux.ppc64le \
	template-arhat-device-go.linux.s390x \
	template-arhat-device-go.linux.riscv64

template-arhat-device-go.darwin.amd64:
	sh scripts/build/build.sh $@

# # currently darwin/arm64 build will fail due to golang link error
# template-arhat-device-go.darwin.arm64:
# 	sh scripts/build/build.sh $@

template-arhat-device-go.darwin.all: \
	template-arhat-device-go.darwin.amd64

template-arhat-device-go.windows.x86:
	sh scripts/build/build.sh $@

template-arhat-device-go.windows.amd64:
	sh scripts/build/build.sh $@

template-arhat-device-go.windows.armv5:
	sh scripts/build/build.sh $@

template-arhat-device-go.windows.armv6:
	sh scripts/build/build.sh $@

template-arhat-device-go.windows.armv7:
	sh scripts/build/build.sh $@

# # currently no support for windows/arm64
# template-arhat-device-go.windows.arm64:
# 	sh scripts/build/build.sh $@

template-arhat-device-go.windows.all: \
	template-arhat-device-go.windows.x86 \
	template-arhat-device-go.windows.amd64

# # android build requires android sdk
# template-arhat-device-go.android.amd64:
# 	sh scripts/build/build.sh $@

# template-arhat-device-go.android.x86:
# 	sh scripts/build/build.sh $@

# template-arhat-device-go.android.armv5:
# 	sh scripts/build/build.sh $@

# template-arhat-device-go.android.armv6:
# 	sh scripts/build/build.sh $@

# template-arhat-device-go.android.armv7:
# 	sh scripts/build/build.sh $@

# template-arhat-device-go.android.arm64:
# 	sh scripts/build/build.sh $@

# template-arhat-device-go.android.all: \
# 	template-arhat-device-go.android.amd64 \
# 	template-arhat-device-go.android.arm64 \
# 	template-arhat-device-go.android.x86 \
# 	template-arhat-device-go.android.armv7 \
# 	template-arhat-device-go.android.armv5 \
# 	template-arhat-device-go.android.armv6

template-arhat-device-go.freebsd.amd64:
	sh scripts/build/build.sh $@

template-arhat-device-go.freebsd.x86:
	sh scripts/build/build.sh $@

template-arhat-device-go.freebsd.armv5:
	sh scripts/build/build.sh $@

template-arhat-device-go.freebsd.armv6:
	sh scripts/build/build.sh $@

template-arhat-device-go.freebsd.armv7:
	sh scripts/build/build.sh $@

template-arhat-device-go.freebsd.arm64:
	sh scripts/build/build.sh $@

template-arhat-device-go.freebsd.all: \
	template-arhat-device-go.freebsd.amd64 \
	template-arhat-device-go.freebsd.arm64 \
	template-arhat-device-go.freebsd.armv7 \
	template-arhat-device-go.freebsd.x86 \
	template-arhat-device-go.freebsd.armv5 \
	template-arhat-device-go.freebsd.armv6

template-arhat-device-go.netbsd.amd64:
	sh scripts/build/build.sh $@

template-arhat-device-go.netbsd.x86:
	sh scripts/build/build.sh $@

template-arhat-device-go.netbsd.armv5:
	sh scripts/build/build.sh $@

template-arhat-device-go.netbsd.armv6:
	sh scripts/build/build.sh $@

template-arhat-device-go.netbsd.armv7:
	sh scripts/build/build.sh $@

template-arhat-device-go.netbsd.arm64:
	sh scripts/build/build.sh $@

template-arhat-device-go.netbsd.all: \
	template-arhat-device-go.netbsd.amd64 \
	template-arhat-device-go.netbsd.arm64 \
	template-arhat-device-go.netbsd.armv7 \
	template-arhat-device-go.netbsd.x86 \
	template-arhat-device-go.netbsd.armv5 \
	template-arhat-device-go.netbsd.armv6

template-arhat-device-go.openbsd.amd64:
	sh scripts/build/build.sh $@

template-arhat-device-go.openbsd.x86:
	sh scripts/build/build.sh $@

template-arhat-device-go.openbsd.armv5:
	sh scripts/build/build.sh $@

template-arhat-device-go.openbsd.armv6:
	sh scripts/build/build.sh $@

template-arhat-device-go.openbsd.armv7:
	sh scripts/build/build.sh $@

template-arhat-device-go.openbsd.arm64:
	sh scripts/build/build.sh $@

template-arhat-device-go.openbsd.all: \
	template-arhat-device-go.openbsd.amd64 \
	template-arhat-device-go.openbsd.arm64 \
	template-arhat-device-go.openbsd.armv7 \
	template-arhat-device-go.openbsd.x86 \
	template-arhat-device-go.openbsd.armv5 \
	template-arhat-device-go.openbsd.armv6

template-arhat-device-go.solaris.amd64:
	sh scripts/build/build.sh $@

template-arhat-device-go.aix.ppc64:
	sh scripts/build/build.sh $@

template-arhat-device-go.dragonfly.amd64:
	sh scripts/build/build.sh $@

template-arhat-device-go.plan9.amd64:
	sh scripts/build/build.sh $@

template-arhat-device-go.plan9.x86:
	sh scripts/build/build.sh $@

template-arhat-device-go.plan9.armv5:
	sh scripts/build/build.sh $@

template-arhat-device-go.plan9.armv6:
	sh scripts/build/build.sh $@

template-arhat-device-go.plan9.armv7:
	sh scripts/build/build.sh $@

template-arhat-device-go.plan9.all: \
	template-arhat-device-go.plan9.amd64 \
	template-arhat-device-go.plan9.armv7 \
	template-arhat-device-go.plan9.x86 \
	template-arhat-device-go.plan9.armv5 \
	template-arhat-device-go.plan9.armv6
