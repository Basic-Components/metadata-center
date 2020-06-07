$env:IMG_VERSION="dev-0.0.1"
docker build -t hsz1273327/components_manager:$env:IMG_VERSION .
docker push hsz1273327/components_manager:$env:IMG_VERSION
