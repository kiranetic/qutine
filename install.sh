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

go build -o qtn ./cmd/qtn
mkdir -p "/home/$username/.qutine"
mv qtn "/home/$username/.local/bin/"
echo "password:$password" > "/home/$username/.qutine/config" # Placeholder; will hash later
chown -R "$username:$username" "/home/$username/.qutine"
echo "qutine installed for UID $uid"
