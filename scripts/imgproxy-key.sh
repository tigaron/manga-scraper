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

# Display the results
echo "Key: $key"
echo "Salt: $salt"
