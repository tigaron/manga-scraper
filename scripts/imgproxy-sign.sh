#!/usr/bin/env bash

# Check if key and salt are provided
if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Usage: $0 <key> <salt>"
  exit 1
fi

# Function to convert string to hex and then to uppercase
string_to_hex() {
  local input="$1"
  # Convert string to hex using xxd and then to uppercase using tr
  echo -n "$input" | xxd -p | tr 'a-f' 'A-F'
}

# Assign the key and salt to variables
key="$1"
salt="$2"

# Convert key and salt to hex
key_hex=$(string_to_hex "$key")
salt_hex=$(string_to_hex "$salt")

# Prompt user to input the unsigned URL
echo "Enter the unsigned URL:"
read unsigned_url

# https://imgproxy.hostinger.fourleaves.studio/insecure/resize:fill:300:400:0/plain/https://asuratoon.com/wp-content/uploads/2023/12/ReincarnatorCover01.png
# Check if the input contains the specified string
if [[ $unsigned_url != *"https://imgproxy.hostinger.fourleaves.studio/insecure/"* ]]; then
  echo "Error: The input must contain 'https://imgproxy.hostinger.fourleaves.studio/insecure/'"
  exit 1
fi

# Extract the path after 'insecure/'
path_to_sign=$(echo "$unsigned_url" | sed 's|https://imgproxy.hostinger.fourleaves.studio/insecure/||')

# Add the salt to the beginning of the path
string_to_sign="$salt/$path_to_sign"

# Calculate the HMAC digest using SHA256 and encode it with URL-safe Base64
hmac_signature=$(echo -n "$string_to_sign" | openssl dgst -sha256 -mac HMAC -macopt hexkey:$key_hex -binary | openssl base64 | tr '+/' '-_' | tr -d '=')

# Construct the signed URL
signed_url="https://imgproxy.hostinger.fourleaves.studio/$hmac_signature/$path_to_sign"

# Display the signed URL
echo "Signed URL: $signed_url"
