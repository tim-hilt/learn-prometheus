#! /usr/bin/env sh

NAME="ml-simulation"

docker build -t $NAME . && \
    kind create cluster --config=kind-config.yaml && \
    kind load docker-image $NAME && \
    helm dependency build ./deploy && \
    helm install $NAME ./deploy
