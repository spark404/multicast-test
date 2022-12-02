#!/bin/zsh

docker build -t receiver -f cmd/receiver/Dockerfile .
docker run --rm -v `pwd`:/hostos receiver cp /usr/bin/receiver /hostos/receiver
docker tag receiver registry.vagrant-cluster.laptop/receiver
docker push registry.vagrant-cluster.laptop/receiver

docker build -t sender -f cmd/sender/Dockerfile .
docker run --rm -v `pwd`:/hostos sender cp /usr/bin/sender /hostos/sender
docker tag sender registry.vagrant-cluster.laptop/sender
docker push registry.vagrant-cluster.laptop/sender

kubectl create -f k3s/receiver.yaml
