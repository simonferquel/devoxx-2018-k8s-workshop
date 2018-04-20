FROM golang:1.10-alpine3.7 AS builder-base
RUN apk add --no-cache \
  coreutils \
  make
WORKDIR /go/src/github.com/simonferquel/devoxx-2018-k8s-workshop
COPY pkg pkg
COPY vendor vendor
RUN go build -o /dev/null ./pkg/... 
COPY . .

FROM alpine:3.7 AS runtime-base
RUN apk add ca-certificates --no-cache

FROM builder-base AS controller-builder
RUN make bin/etcdaas-controller

FROM runtime-base AS controller
COPY --from=controller-builder /go/src/github.com/simonferquel/devoxx-2018-k8s-workshop/bin/etcdaas-controller /etcdaas-controller
ENTRYPOINT [ "/etcdaas-controller" ]
