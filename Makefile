.PHONY: all build unit-test clean

all: build 
test: unit-test

PLATFORM=local

build:
	@echo "Let's build it"
	@echo $(git rev-parse HEAD|cut -c1-7)
	@export APP_BUILD_INFO=$(git rev-parse HEAD|cut -c1-7)
	@docker build --build-arg APP_BUILD_INFO=$APP_BUILD_INFO \
	--target build -o bin/ . \
	--platform ${PLATFORM}

unit-test:
	@echo "Run tests here..."
	@docker build --target unit-test .

lint:
	@echo "Run lint here..."
	@docker build --target lint
clean:
	@echo "Cleaning up..."
	rm -rf bin