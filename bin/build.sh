 #! /bin/bash
APP_VERSION=${2:-v1.0.0}
GIT_REPO=denvasyliev/k8s-diy


docker build --rm \
--build-arg APP_VERSION=$APP_VERSION \
--build-arg APP_BUILD_INFO=$(git rev-parse HEAD|cut -c1-7) \
-t $GIT_REPO:$APP_VERSION-$(git rev-parse HEAD|cut -c1-7) .

docker push $GIT_REPO:$APP_VERSION-$(git rev-parse HEAD|cut -c1-7)