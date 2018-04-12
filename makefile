EXEC_EXT :=
ifeq ($(OS),Windows_NT)
	EXEC_EXT := .exe
endif

.PHONY: generator-image generate-all test bin/etcdaas-cli bin/etcdaas-controller

bin/etcdaas-cli:
	go build -i -o ./bin/etcdaas-cli$(EXEC_EXT) ./cmd/etcdaas-cli
bin/etcdaas-controller:
	go build -i -o ./bin/etcdaas-controller$(EXEC_EXT) ./cmd/etcdaas-controller

generator-image:
	docker build -t k8s-generators ./tools/generators

generate-all: generator-image
	docker run -v $(CURDIR):/go/src/github.com/simonferquel/devoxx-2018-k8s-workshop --rm k8s-generators bash ./generate-groups.sh all \
	github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/client \
	github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/apis \
	"etcdaas:v1alpha1"

test:
	go test -i ./...