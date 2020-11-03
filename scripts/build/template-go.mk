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
template-go:
	sh scripts/build/build.sh $@

# linux
template-go.linux.x86:
	sh scripts/build/build.sh $@

template-go.linux.amd64:
	sh scripts/build/build.sh $@

template-go.linux.armv5:
	sh scripts/build/build.sh $@

template-go.linux.armv6:
	sh scripts/build/build.sh $@

template-go.linux.armv7:
	sh scripts/build/build.sh $@

template-go.linux.arm64:
	sh scripts/build/build.sh $@

template-go.linux.mips:
	sh scripts/build/build.sh $@

template-go.linux.mipshf:
	sh scripts/build/build.sh $@

template-go.linux.mipsle:
	sh scripts/build/build.sh $@

template-go.linux.mipslehf:
	sh scripts/build/build.sh $@

template-go.linux.mips64:
	sh scripts/build/build.sh $@

template-go.linux.mips64hf:
	sh scripts/build/build.sh $@

template-go.linux.mips64le:
	sh scripts/build/build.sh $@

template-go.linux.mips64lehf:
	sh scripts/build/build.sh $@

template-go.linux.ppc64:
	sh scripts/build/build.sh $@

template-go.linux.ppc64le:
	sh scripts/build/build.sh $@

template-go.linux.s390x:
	sh scripts/build/build.sh $@

template-go.linux.riscv64:
	sh scripts/build/build.sh $@

template-go.linux.all: \
	template-go.linux.x86 \
	template-go.linux.amd64 \
	template-go.linux.armv5 \
	template-go.linux.armv6 \
	template-go.linux.armv7 \
	template-go.linux.arm64 \
	template-go.linux.mips \
	template-go.linux.mipshf \
	template-go.linux.mipsle \
	template-go.linux.mipslehf \
	template-go.linux.mips64 \
	template-go.linux.mips64hf \
	template-go.linux.mips64le \
	template-go.linux.mips64lehf \
	template-go.linux.ppc64 \
	template-go.linux.ppc64le \
	template-go.linux.s390x \
	template-go.linux.riscv64

template-go.darwin.amd64:
	sh scripts/build/build.sh $@

# # currently darwin/arm64 build will fail due to golang link error
# template-go.darwin.arm64:
# 	sh scripts/build/build.sh $@

template-go.darwin.all: \
	template-go.darwin.amd64

template-go.windows.x86:
	sh scripts/build/build.sh $@

template-go.windows.amd64:
	sh scripts/build/build.sh $@

template-go.windows.armv5:
	sh scripts/build/build.sh $@

template-go.windows.armv6:
	sh scripts/build/build.sh $@

template-go.windows.armv7:
	sh scripts/build/build.sh $@

# # currently no support for windows/arm64
# template-go.windows.arm64:
# 	sh scripts/build/build.sh $@

template-go.windows.all: \
	template-go.windows.amd64 \
	template-go.windows.armv7 \
	template-go.windows.x86 \
	template-go.windows.armv5 \
	template-go.windows.armv6

# # android build requires android sdk
# template-go.android.amd64:
# 	sh scripts/build/build.sh $@

# template-go.android.x86:
# 	sh scripts/build/build.sh $@

# template-go.android.armv5:
# 	sh scripts/build/build.sh $@

# template-go.android.armv6:
# 	sh scripts/build/build.sh $@

# template-go.android.armv7:
# 	sh scripts/build/build.sh $@

# template-go.android.arm64:
# 	sh scripts/build/build.sh $@

# template-go.android.all: \
# 	template-go.android.amd64 \
# 	template-go.android.arm64 \
# 	template-go.android.x86 \
# 	template-go.android.armv7 \
# 	template-go.android.armv5 \
# 	template-go.android.armv6

template-go.freebsd.amd64:
	sh scripts/build/build.sh $@

template-go.freebsd.x86:
	sh scripts/build/build.sh $@

template-go.freebsd.armv5:
	sh scripts/build/build.sh $@

template-go.freebsd.armv6:
	sh scripts/build/build.sh $@

template-go.freebsd.armv7:
	sh scripts/build/build.sh $@

template-go.freebsd.arm64:
	sh scripts/build/build.sh $@

template-go.freebsd.all: \
	template-go.freebsd.amd64 \
	template-go.freebsd.arm64 \
	template-go.freebsd.armv7 \
	template-go.freebsd.x86 \
	template-go.freebsd.armv5 \
	template-go.freebsd.armv6

template-go.netbsd.amd64:
	sh scripts/build/build.sh $@

template-go.netbsd.x86:
	sh scripts/build/build.sh $@

template-go.netbsd.armv5:
	sh scripts/build/build.sh $@

template-go.netbsd.armv6:
	sh scripts/build/build.sh $@

template-go.netbsd.armv7:
	sh scripts/build/build.sh $@

template-go.netbsd.arm64:
	sh scripts/build/build.sh $@

template-go.netbsd.all: \
	template-go.netbsd.amd64 \
	template-go.netbsd.arm64 \
	template-go.netbsd.armv7 \
	template-go.netbsd.x86 \
	template-go.netbsd.armv5 \
	template-go.netbsd.armv6

template-go.openbsd.amd64:
	sh scripts/build/build.sh $@

template-go.openbsd.x86:
	sh scripts/build/build.sh $@

template-go.openbsd.armv5:
	sh scripts/build/build.sh $@

template-go.openbsd.armv6:
	sh scripts/build/build.sh $@

template-go.openbsd.armv7:
	sh scripts/build/build.sh $@

template-go.openbsd.arm64:
	sh scripts/build/build.sh $@

template-go.openbsd.all: \
	template-go.openbsd.amd64 \
	template-go.openbsd.arm64 \
	template-go.openbsd.armv7 \
	template-go.openbsd.x86 \
	template-go.openbsd.armv5 \
	template-go.openbsd.armv6

template-go.solaris.amd64:
	sh scripts/build/build.sh $@

template-go.aix.ppc64:
	sh scripts/build/build.sh $@

template-go.dragonfly.amd64:
	sh scripts/build/build.sh $@

template-go.plan9.amd64:
	sh scripts/build/build.sh $@

template-go.plan9.x86:
	sh scripts/build/build.sh $@

template-go.plan9.armv5:
	sh scripts/build/build.sh $@

template-go.plan9.armv6:
	sh scripts/build/build.sh $@

template-go.plan9.armv7:
	sh scripts/build/build.sh $@

template-go.plan9.all: \
	template-go.plan9.amd64 \
	template-go.plan9.armv7 \
	template-go.plan9.x86 \
	template-go.plan9.armv5 \
	template-go.plan9.armv6
