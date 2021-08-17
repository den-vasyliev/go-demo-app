.PHONY: all build unit-test clean

all: build 
test: unit-test

PLATFORM=local

BUILDER = img

build:
	@echo "Let's build it"
	#@echo $(git rev-parse HEAD|cut -c1-7)
	#@export APP_BUILD_INFO=$(git rev-parse HEAD|cut -c1-7)
	@${BUILDER} build \
	--target build -o bin/ . \
	--platform ${PLATFORM}

unit-test:
	@echo "Run tests here..."
	@${BUILDER} build --target unit-test .

lint:
	@echo "Run lint here..."
	@${BUILDER} build --target lint .

clean:
	@echo "Cleaning up..."
	rm -rf bin