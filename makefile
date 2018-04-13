EXEC_EXT :=
ifeq ($(OS),Windows_NT)
	EXEC_EXT := .exe
	HOME := ${HOMEDRIVE}${HOMEPATH}
endif

.PHONY: generator-image generate-all test bin/etcdaas-cli bin/etcdaas-controller bin/etcdaas-api image/etcdaas-controller image/etcdaas-api images

bin/etcdaas-cli:
	go build -i -o ./bin/etcdaas-cli$(EXEC_EXT) ./cmd/etcdaas-cli
	
bin/etcdaas-controller:
	go build -i -o ./bin/etcdaas-controller$(EXEC_EXT) ./cmd/etcdaas-controller

bin/etcdaas-api:
	go build -i -o ./bin/etcdaas-api$(EXEC_EXT) ./cmd/etcdaas-api

image/etcdaas-controller:
	docker build -t etcdaas-controller --target controller .

image/etcdaas-api:
	docker build -t etcdaas-api --target api .

images: image/etcdaas-controller image/etcdaas-api

generator-image:
	docker build -t k8s-generators ./tools/generators

generate-all: generator-image
	docker run -v $(CURDIR):/go/src/github.com/simonferquel/devoxx-2018-k8s-workshop --rm k8s-generators bash ./generate-groups.sh all \
	github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/client \
	github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/apis \
	"etcdaas:v1alpha1"

test:
	go test -i ./...

start-dev-etcd:
	docker start devoxx-etcd || docker run --name devoxx-etcd -d -p 9898:2379 quay.io/coreos/etcd:v3.2.9 //usr/local/bin/etcd -advertise-client-urls=http://0.0.0.0:2379 -listen-client-urls=http://0.0.0.0:2379

start-dev-api: bin/etcdaas-api start-dev-etcd
	kubectl apply -f ./k8s-assets/apiservice-dev.yml
	./bin/etcdaas-api$(EXEC_EXT) --etcd-servers=http://127.0.0.1:9898 --kubeconfig=$(HOME)/.kube/config --authorization-kubeconfig=$(HOME)/.kube/config --authentication-kubeconfig=$(HOME)/.kube/config --secure-port 9443

teardown-dev-api:
	kubectl delete -f ./k8s-assets/apiservice-dev.yml

start-dev-controller: bin/etcdaas-controller
	./bin/etcdaas-controller$(EXEC_EXT) -kubeconfig $(HOME)/.kube/config

deploy-crd:
	kubectl apply -f ./k8s-assets/crd.yml

uninstall-crd:
	kubectl delete -f ./k8s-assets/crd.yml || echo

deploy-api: image/etcdaas-api
	kubectl apply -f ./k8s-assets/sa.yml
	kubectl apply -f ./k8s-assets/api-deployment.yml

uninstall-api:
	kubectl delete -f ./k8s-assets/api-deployment.yml || echo

deploy-controller: image/etcdaas-controller
	kubectl apply -f ./k8s-assets/sa.yml
	kubectl apply -f ./k8s-assets/controller-deployment.yml
	
uninstall-controller:
	kubectl delete -f ./k8s-assets/controller-deployment.yml || echo

deploy-all: deploy-api deploy-controller

uninstall-all: uninstall-api uninstall-controller uninstall-crd