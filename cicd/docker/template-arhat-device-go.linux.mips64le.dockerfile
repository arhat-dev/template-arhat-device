ARG ARCH=amd64

# we do not use cgo so we can build on alpine and copy it to debian
FROM arhatdev/builder-go:alpine as builder
FROM arhatdev/go:debian-${ARCH}
ARG APP=template-arhat-device-go

ENTRYPOINT [ "/template-arhat-device-go" ]
