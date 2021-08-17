.PHONY: all build unit-test clean

all: build 
test: unit-test

PLATFORM=linux 

BUILDER = docker

TAG=denvasyliev/k8sdiy
BUILD=$$(git rev-parse HEAD|cut -c1-7)

build:
	@echo "Let's build ${BUILD}"
	@${BUILDER} build --progress plain \
	--target build . --build-arg APP_BUILD_INFO=${BUILD} \
	--platform ${PLATFORM} \
	--tag ${TAG}:build-${BUILD}

push:
	@echo "Let's push it"
	@export APP_BUILD_INFO=$(git rev-parse HEAD|cut -c1-7)
	@${BUILDER} push ${TAG}:build-${APP_BUILD_INFO}

unit-test:
	@echo "Run tests here..."
	@${BUILDER} build --target unit-test .

lint:
	@echo "Run lint here..."
	@${BUILDER} build --target lint .

clean:
	@echo "Cleaning up..."
	rm -rf bin