FROM golang:1.10
RUN go get -u k8s.io/code-generator/cmd/client-gen \
&&  go get -u k8s.io/code-generator/cmd/conversion-gen \
&&  go get -u k8s.io/code-generator/cmd/deepcopy-gen \
&&  go get -u k8s.io/code-generator/cmd/defaulter-gen \
&&  go get -u k8s.io/code-generator/cmd/go-to-protobuf \
&&  go get -u k8s.io/code-generator/cmd/import-boss \
&&  go get -u k8s.io/code-generator/cmd/informer-gen \
&&  go get -u k8s.io/code-generator/cmd/lister-gen \
&&  go get -u k8s.io/code-generator/cmd/openapi-gen \
&&  go get -u k8s.io/code-generator/cmd/set-gen
WORKDIR /go/src/github.com/simonferquel/devoxx-2018-k8s-workshop