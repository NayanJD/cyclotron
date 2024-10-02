#!/bin/bash

set -x

# sudo mkfs -t ext4 /dev/xvdj

# # mkdir -p /var/lib/postgres/data
# mount /dev/xvdj /bitnami/postgresql

# cp /etc/fstab /etc/fstab.orig
# echo "/dev/xvdj /bitnami/postgresql  ext4 defaults,nofail  0  2" >> /etc/fstab

# mkdir -p /bitnami/postgresql/data

mkdir /tmp/ssm
cd /tmp/ssm
wget https://s3.amazonaws.com/ec2-downloads-windows/SSMAgent/latest/debian_amd64/amazon-ssm-agent.deb
sudo dpkg -i amazon-ssm-agent.deb


sudo apt update
sudo apt install jq git -y

GITHUB_TOKEN=$(aws secretsmanager get-secret-value --secret-id cyclotron/commons --region ap-south-1 --query SecretString --output text | jq -r '.GITHUB_TOKEN')

sudo cd /home/binami

git clone https://$GITHUB_TOKEN:@github.com/nayanjd/cyclotron && cd cyclotron
git checkout feat/benchmarking-infra-final

CYCLOTRON_SECRET=$(aws secretsmanager get-secret-value --secret-id cyclotron/commons --region ap-south-1 --query SecretString --output text)

PGPASSWORD=$(sudo grep -oP "\'.*\'" /home/bitnami/bitnami_credentials | tr -d \'\" | awk '{print $3}')

UPDATED_CYCLOTRON_SECRET=$(echo $CYCLOTRON_SECRET | jq --arg PGPASSWORD $PGPASSWORD '.PGPASSWORD = $PGPASSWORD' -c)

aws secretsmanager update-secret --secret-id cyclotron/commons --region ap-south-1 --secret-string $UPDATED_CYCLOTRON_SECRET --output text

sudo echo "host     all             all        all                          md5" >> /opt/bitnami/postgresql/conf/pg_hba.conf 

sudo echo "listen_addresses = '*'" >> /opt/bitnami/postgresql/conf/postgresql.conf

sudo systemctl restart bitnami.service

PGPASSWORD=$PGPASSWORD psql -U postgres -c "CREATE DATABASE cyclotron"





