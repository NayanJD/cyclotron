#!/bin/bash

set -x

sudo apt update
sudo apt install jq unzip -y

curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

curl -L "https://go.dev/dl/go1.20.11.linux-amd64.tar.gz" -O
tar -C /usr/local -xzf go1.20.11.linux-amd64.tar.gz

sudo ln -s /usr/local/go/bin/go /usr/bin

GITHUB_TOKEN=$(aws secretsmanager get-secret-value --secret-id cyclotron/commons --region ap-south-1 --query SecretString --output text | jq -r '.GITHUB_TOKEN')

cd /home/ubuntu

curl "https://raw.githubusercontent.com/12moons/ec2-tags-env/master/import-tags.sh" -o import-tags.sh
sed -i -e 's/\r$//' import-tags.sh
chmod +x import-tags.sh

. ./import-tags.sh

git clone https://$GITHUB_TOKEN:@github.com/nayanjd/cyclotron && cd cyclotron
git checkout feat/benchmarking-infra-final

cd /home/ubuntu/cyclotron/user
go mod tidy

PGPASSWORD=$(aws secretsmanager get-secret-value --secret-id cyclotron/commons --region ap-south-1 --query SecretString --output text | jq -r '.PGPASSWORD')

go run cmd/main.go --postgres-conn-url postgres://postgres:$PGPASSWORD@$POSTGRES_URL/cyclotron?sslmode=disable --otel-grpc-url $JAEGER_URL:4317 --log-level DEBUG
