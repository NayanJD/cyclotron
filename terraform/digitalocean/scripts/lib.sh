#!/bin/bash

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
