#!/bin/bash
user_home=$(eval echo ~${SUDO_USER:-$USER})
rm -rf "$user_home/.qutine" "$user_home/.local/bin/qtn"
echo "qutine uninstalled"
