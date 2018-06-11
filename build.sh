 #! /bin/bash

APP_NAME=${1:-front}
APP_VERSION=${2:-v0.1.5}
NEW_FEATURE=$3
APP_DB="root:@tcp(db)/demo"


docker build --rm \
--build-arg APP_NAME=$APP_NAME \
--build-arg NEW_FEATURE=$NEW_FEATURE \
--build-arg APP_VERSION=$APP_VERSION \
--build-arg APP_DB=$APP_DB \
--build-arg APP_BUILD_INFO=$(git rev-parse HEAD|cut -c1-7) \
-t denvasyliev/$APP_NAME:$APP_VERSION-$(git rev-parse HEAD|cut -c1-7) .

docker push denvasyliev/$APP_NAME:$APP_VERSION-$(git rev-parse HEAD|cut -c1-7)