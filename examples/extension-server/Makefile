tools.dir = tools
tools.bindir = tools/bin
tools.srcdir = tools/src

tools/controller-gen     = $(tools.bindir)/controller-gen
$(tools.bindir)/%: $(tools.srcdir)/%/pin.go $(tools.srcdir)/%/go.mod
	cd $(<D) && GOOS= GOARCH= go build -o $(abspath $@) $$(sed -En 's,^import _ "(.*)".*,\1,p' pin.go)

clean: tools.clean
	rm -fr bin

#1. make build to build the extension server locally.
build: generate manifests
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/extension-server ./cmd/extension-server

#2. make image to tag the imagea and push the locally build image to docker hub. 
image: build
	docker build -t extension-server:latest1 -f tools/docker/extension-server/Dockerfile .
	docker tag extension-server:latest1 avidocker692/extension-server:latest1
	docker push avidocker692/extension-server:latest1

generate: $(tools/controller-gen)
	$(tools/controller-gen) object:headerFile="$(tools.dir)/boilerplate.generatego.txt",year=2024 paths="{./api/...}"

# Use manifests-combined as the default manifests target
# manifests: manifests-combined

manifests: manifests-v1alpha1

# manifests: manifests-v1alpha2

# manifests: manifests-v1alpha3

# manifests: $(tools/controller-gen)
# 	$(tools/controller-gen) crd:allowDangerousTypes=true paths="./..." output:crd:artifacts:config=charts/extension-server/crds/generated

#3. make manifests commands to generate CRDs
# Generate a single CRD only for v1alpha1 version
manifests-v1alpha1: $(tools/controller-gen)
	$(tools/controller-gen) crd:allowDangerousTypes=true paths="./api/v1alpha1" output:crd:artifacts:config=charts/extension-server/crds/generated
	# Rename the generated CRD to dp.wso2.com_apis_v1alpha1.yaml
	mv charts/extension-server/crds/generated/dp.wso2.com_apis.yaml charts/extension-server/crds/generated/dp.wso2.com_apis_v1alpha1.yaml

# Generate a single CRD only for v1alpha2 version
manifests-v1alpha2: $(tools/controller-gen)
	$(tools/controller-gen) crd:allowDangerousTypes=true paths="./api/v1alpha2" output:crd:artifacts:config=charts/extension-server/crds/generated
	# Rename the generated CRD to dp.wso2.com_apis_v1alpha2.yaml
	mv charts/extension-server/crds/generated/dp.wso2.com_apis.yaml charts/extension-server/crds/generated/dp.wso2.com_apis_v1alpha2.yaml

# Generate a single CRD only for v1alpha3 version
manifests-v1alpha3: $(tools/controller-gen)
	$(tools/controller-gen) crd:allowDangerousTypes=true paths="./api/v1alpha3" output:crd:artifacts:config=charts/extension-server/crds/generated
	# Rename the generated CRD to dp.wso2.com_apis_v1alpha3.yaml
	mv charts/extension-server/crds/generated/dp.wso2.com_apis.yaml charts/extension-server/crds/generated/dp.wso2.com_apis_v1alpha3.yaml

# Generate a single CRD for v1alpha1,v1alpha2 and v1alpha3
manifests-combined: $(tools/controller-gen)
	$(tools/controller-gen) crd:allowDangerousTypes=true paths="./api/v1alpha1;./api/v1alpha2;./api/v1alpha3" output:crd:artifacts:config=charts/extension-server/crds/generated



# generate-v1alpha1: $(tools/controller-gen)
# 	$(tools/controller-gen) object:headerFile="$(tools.dir)/boilerplate.generatego.txt",year=2024 paths="{./api/v1alpha1/}"	

# generate-v1alpha2: $(tools/controller-gen)
# 	$(tools/controller-gen) object:headerFile="$(tools.dir)/boilerplate.generatego.txt",year=2024 paths="{./api/v1alpha2/}"

# generate-v1alpha3: $(tools/controller-gen)
# 	$(tools/controller-gen) object:headerFile="$(tools.dir)/boilerplate.generatego.txt",year=2024 paths="{./api/v1alpha3/}"

tools.clean: # Remove all tools
	@$(LOG_TARGET)
	rm -rf $(tools.bindir)

settings: 
	-/home/avindu/Desktop/envoy-gateway/examples/extension-server/main.sh || true
	-kubectl apply -f /home/avindu/Desktop/envoy-gateway-2/examples/extension-server/settings.yaml || true
	-kubectl rollout restart deployment/envoy-gateway -n envoy-gateway-system || true

apis:
	kubectl apply -f /home/avindu/Desktop/envoy-gateway-2/examples/extension-server/tapi.yaml

delete-apis:
	kubectl delete -f /home/avindu/Desktop/envoy-gateway-2/examples/extension-server/tapi.yaml

qs:
	kubectl apply -f /home/avindu/Desktop/envoy-gateway-2/examples/extension-server/qs.yaml

qs-api:
	kubectl apply -f /home/avindu/Desktop/envoy-gateway-2/examples/extension-server/qs-api.yaml

install-envoy:
	helm install eg oci://docker.io/envoyproxy/gateway-helm --version v1.1.4 -n envoy-gateway-system --create-namespace --values /home/avindu/Desktop/envoy-gateway-2/gateway-helm/values.yaml

uninstall-envoy:
	helm uninstall eg  -n envoy-gateway-system

delete-qs:
	kubectl delete -f /home/avindu/Desktop/envoy-gateway-2/examples/extension-server/qs.yaml

test-qs:
	GATEWAY_HOST=$$(kubectl get gateway/eg -o jsonpath='{.status.addresses[0].value}') && echo host: $$GATEWAY_HOST  && \
	curl --verbose --header "Host: www.example.com" http://$${GATEWAY_HOST}/hello

test-auth:
	GATEWAY_HOST=$$(kubectl get gateway/eg -o jsonpath='{.status.addresses[0].value}') && echo host: $$GATEWAY_HOST  && \
	curl --verbose --header "Host: www.example.com" http://$${GATEWAY_HOST}/hello  --user 'user:p@ssw0rd'

test-api:
	GATEWAY_HOST=$$(kubectl get gateway/eg -o jsonpath='{.status.addresses[0].value}') && echo host: $$GATEWAY_HOST  && \
	curl --verbose --header "Host: www.example.com" http://$${GATEWAY_HOST}/example-context-1/hello  --user 'user:p@ssw0rd'


install-extension-server:
	helm install -n envoy-gateway-system extension-server /home/avindu/Desktop/envoy-gateway-2/examples/extension-server/charts/extension-server

uninstall-extension-server:
	helm uninstall extension-server -n envoy-gateway-system

clean-all:
	-make delete-apis || true
	-make delete-qs || true
	-make uninstall-envoy || true
	-make uninstall-extension-server || true

install-all:
	minikube addons enable metallb
	kubectl apply -f /home/avindu/Desktop/envoy-gateway-2/examples/extension-server/metallb.yaml
	make install-envoy
	make install-extension-server
	make settings
	make qs
	# make apis
	# make qs-api

.PHONY: build image manifests generate tools.clean clean
