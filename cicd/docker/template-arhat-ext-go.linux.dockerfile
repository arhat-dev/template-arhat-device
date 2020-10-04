ARG ARCH=amd64

FROM arhatdev/builder-go:alpine as builder
FROM arhatdev/go:alpine-${ARCH}
ARG APP=template-arhat-ext-go

ENTRYPOINT [ "/template-arhat-ext-go" ]
