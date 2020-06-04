$env:IMG_VERSION="dev-0.0.1"
docker build -t 47.98.42.103:8333/service/business_config_center:$env:IMG_VERSION .
docker push 47.98.42.103:8333/service/business_config_center:$env:IMG_VERSION
