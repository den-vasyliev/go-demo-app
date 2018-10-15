 #! /bin/bash
GIT_REPO=msrn/services
APP_NAME=${1}
APP_VERSION=${2:-v3.5.1}


docker build --rm \
--build-arg APP_NAME=$APP_NAME \
--build-arg APP_VERSION=$APP_VERSION \
--build-arg APP_BUILD_INFO=$(git rev-parse HEAD|cut -c1-7) \
-t $GIT_REPO:$APP_NAME-$APP_VERSION-$(git rev-parse HEAD|cut -c1-7) .

docker push $GIT_REPO:$APP_NAME-$APP_VERSION-$(git rev-parse HEAD|cut -c1-7)