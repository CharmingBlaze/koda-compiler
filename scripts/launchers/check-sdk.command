#!/bin/bash
cd "$(dirname "$0")"
if [[ ! -x "./koda" ]]; then
  osascript -e 'display alert "koda not found in this folder"'
  exit 1
fi
./koda doctor
read -r -p "Press Enter to close…"
