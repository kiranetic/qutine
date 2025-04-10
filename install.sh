#!/bin/bash

echo "Enter username for qutine installation:"
read username
uid=$(id -u "$username")
if [ -z "$uid" ]; then
  echo "User not found"
  exit 1
fi

echo "Enter password for qutine:"
read -s password
echo

user_home=$(getent passwd "$username" | cut -d: -f6)
bin_dir="$user_home/.local/bin"
qtn_path="$bin_dir/qtn"
config_dir="$user_home/.qutine"

go build -o qtn ./cmd/qtn || {
  echo "Failed to build qtn"
  exit 1
}

mkdir -p "$bin_dir" "$config_dir" || {
  echo "Failed to create directories"
  exit 1
}

mv qtn "$qtn_path" || {
  echo "Failed to move qtn to $bin_dir"
  exit 1
}

"$qtn_path" hash-password "$password" > "$config_dir/config" || {
  echo "Failed to hash password"
  exit 1
}

chmod 600 "$config_dir/config"
chown -R "$username:$username" "$bin_dir" "$config_dir" || {
  echo "Failed to set ownership"
  exit 1
}

echo "qutine installed for UID $uid"