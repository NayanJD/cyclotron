#!/bin/bash

set -ex

# source lib.sh
fetch_droplet_tags() {
  # Endpoint for droplet metadata
  METADATA_URL="http://169.254.169.254/metadata/v1/tags"

  # Fetch the tags from the metadata endpoint
  tags=$(curl -s $METADATA_URL)

  # Check if the curl command was successful
  if [ $? -ne 0 ]; then
    echo "Failed to retrieve metadata. Are you sure you're running this on a DigitalOcean droplet?"
    exit 1
  fi

  # Declare an associative array (dictionary)
  declare -gA tag_dict

  # Check if any tags were returned
  if [ -z "$tags" ]; then
    echo "No tags found for this droplet."
  else
    echo "Processing tags and converting them into key-value pairs..."

    # Loop through each tag and split it into key-value pairs
    for tag in $tags; do
      if [[ "$tag" == *":"* ]]; then
        # Split the tag by the colon
        key=$(echo "$tag" | cut -d':' -f1)
        value=$(echo "$tag" | cut -d':' -f2)

        echo "Key: $key, Value: $value"

        # Add the key-value pair to the dictionary
        tag_dict["$key"]="$value"
      else
        # If the tag doesn't contain a colon, store it with an empty value
        tag_dict["$tag"]=""
      fi
    done

    # Output the contents of the dictionary
    echo "Tags in dictionary format:"
    for key in "${!tag_dict[@]}"; do
      echo "Key: $key, Value: ${tag_dict[$key]}"
    done
  fi
}

fetch_droplet_tags

echo "Volume: ${tag_dict["VOLUME_NAME"]}"

# make postgres dir
mkdir -p /var/lib/postgresql

# mount volume to
disk_path="/dev/disk/by-id/scsi-0DO_Volume_${tag_dict["VOLUME_NAME"]}"
echo $disk_path

mkfs.ext4 $disk_path

mount -o discard,defaults,noatime "${disk_path}" /var/lib/postgresql

# add entry to fstab
echo "/dev/disk/by-id/scsi-0DO_Volume_${tag_dict["VOLUME_NAME"]} /var/lib/postgresql ext4 defaults,nofail,discard 0 0" | sudo tee -a /etc/fstab

# Install postgresql
apt install postgresql -y

sudo -u postgres psql -c "ALTER USER postgres with encrypted password 'somesecret';"
sudo -u postgres psql -c "CREATE DATABASE cyclotron"

# Allow trust from all IPs
sudo echo "host     all             all        all                          md5" >> /etc/postgresql/16/main/pg_hba.conf

systemctl restart postgresql.service
 

