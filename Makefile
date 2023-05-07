
ifndef IMAGE_REPO
	IMAGE_REPO=theotw
endif
BASE_VERSION := $(shell cat 'version.txt')
BUILD_DATE := $(shell date '+%Y%m%d%H%M')

BUILD_VERSION := ${BASE_VERSION}.${BUILD_DATE}
ifndef IMAGE_TAG
	IMAGE_TAG := ${BUILD_VERSION}
endif



printversion:
	echo Base: ${BASE_VERSION}
	echo Date: ${BUILD_DATE}
	echo Image: ${IMAGE_TAG}
	echo Build: ${BUILD_VERSION}

echoenv:
	echo "PATH ${PATH}"
	echo "REPO ${IMAGE_REPO}"
	echo "TAG ${IMAGE_TAG}"

incontainergenerate:
	#Generate the x509 Certs
	./cicd_gen_certs.sh

buildall: buildlinux buildmac

buildmac: export GOOS=darwin
buildmac: export GOARCH=amd64
buildmac: export CGO_ENABLED=0
buildmac: basebuild


buildlinux:	export GOOS=linux
buildlinux: export GOARCH=amd64
buildlinux: export CGO_ENABLED=0
buildlinux: basebuild

buildarm: export GOOS=linux
buildarm: export GOARCH=arm
buildarm: export CGO_ENABLED=0
buildarm: basebuild

build: basebuild

basebuild: export GO111MODULE=on
basebuild: export GOPROXY=${GOPROXY_ENV}
basebuild: export GOSUM=${GOSUM_ENV}
basebuild: export LDFLAGS=-ldflags "-X github.com/theotw/k8srelay/pkg.VERSION=${BUILD_VERSION}"
basebuild:
	mkdir -p out
	rm -f  out/bridgeserver_x64_linux
	go mod tidy
	go build ${LDFLAGS} -v -o out/k8srelaylet_${GOARCH}_${GOOS} apps/k8s-relaylet.go
	go build ${LDFLAGS} -v -o out/k8srelayserver_${GOARCH}_${GOOS} apps/k8s-relay-server.go


buildtest: export GOOS=linux
buildtest: export GOARCH=amd64
buildtest: export CGO_ENABLED=0
buildtest: export GO111MODULE=on
buildtest: export LDFLAGS=-ldflags "-X github.com/theotw/k8srelay/pkg.VERSION=${BUILD_VERSION}"
buildtest: export COVERPKG=-coverpkg=./pkg/...
buildtest:

clean:
	rm -r -f tmp
	rm -r -f pkg/bridgemodel/generated/v1
	rm -r -f out
	rm go.sum

serverimage:
	DOCKER_BUILDKIT=1 docker build --no-cache --build-arg IMAGE_REPO=${IMAGE_REPO} --build-arg IMAGE_TAG=${IMAGE_TAG} --tag ${IMAGE_REPO}/k8srelayserver:${IMAGE_TAG} --target k8srelayserver .
relayletimage:
	DOCKER_BUILDKIT=1 docker build --no-cache --build-arg IMAGE_REPO=${IMAGE_REPO} --build-arg IMAGE_TAG=${IMAGE_TAG} --tag ${IMAGE_REPO}/k8srelaylet:${IMAGE_TAG} --target k8srelaylet .


allimages: serverimage relayletimage
buildAndTag: allimages tag
tag:
	docker tag ${IMAGE_REPO}/k8srelayserver:${IMAGE_TAG} ${IMAGE_REPO}/k8srelayserver:latest
	docker tag ${IMAGE_REPO}/k8srelayserver:${IMAGE_TAG} ${IMAGE_REPO}/k8srelayserver:${BASE_VERSION}
	docker tag ${IMAGE_REPO}/k8srelaylet:${IMAGE_TAG} ${IMAGE_REPO}/k8srelaylet:latest
	docker tag ${IMAGE_REPO}/k8srelaylet:${IMAGE_TAG} ${IMAGE_REPO}/k8srelaylet:${BASE_VERSION}


push:
	docker push ${IMAGE_REPO}/k8srelayserver:${IMAGE_TAG}
	docker push ${IMAGE_REPO}/k8srelayserver:latest
	docker push ${IMAGE_REPO}/k8srelayserver:${BASE_VERSION}
	docker push ${IMAGE_REPO}/k8srelaylet:${IMAGE_TAG}
	docker push ${IMAGE_REPO}/k8srelaylet:latest
	docker push ${IMAGE_REPO}/k8srelaylet:${BASE_VERSION}



l1:
	SUCCESS=0; \
	go install github.com/jstemmer/go-junit-report@latest; \
	mkdir -p out; \
	go test -v -coverpkg=github.com/theotw/k8srelay/pkg/... -coverprofile=out/unit_coverage.out github.com/theotw/k8srelay/pkg/... > out/l1_out.txt 2>&1 || SUCCESS=1; \
	cat out/l1_out.txt | go-junit-report > out/l1_report.xml || echo "Failure generating report xml"; \
	cat out/l1_out.txt; \
	exit $$SUCCESS;


