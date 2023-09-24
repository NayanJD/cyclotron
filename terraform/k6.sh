#!/bin/bash

set -x
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version
GITHUB_TOKEN=$(aws secretsmanager get-secret-value --secret-id cyclotron/commons --region ap-south-1 --query SecretString --output text | jq -r '.GITHUB_TOKEN')

echo $GITHUB_TOKEN

cd /home/ubuntu
git clone https://$GITHUB_TOKEN:@github.com/nayanjd/cyclotron && cd cyclotron
git checkout feat/benchmarking-infra

sudo mkfs -t ext4 /dev/xvdj
sudo mkfs -t ext4 /dev/xvdp

mkdir /var/lib/jaeger
mount /dev/xvdj /var/lib/jaeger

mkdir /var/lib/prometheus
mount /dev/xvdp /var/lib/prometheus

cp /etc/fstab /etc/fstab.orig
echo "/dev/xvdj  /var/lib/jaeger  ext4 defaults,nofail  0  2" >> /etc/fstab
echo "/dev/xvdp  /var/lib/prometheus  ext4 defaults,nofail  0  2" >> /etc/fstab

mount -a

docker-compose -f docker-compose-ec2.yml up -d

