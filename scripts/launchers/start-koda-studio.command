#!/bin/bash
cd "$(dirname "$0")"
APP="$PWD/Koda Studio.app"
if [[ ! -d "$APP" ]]; then
  osascript -e 'display alert "Koda Studio not found" message "Unzip the full Koda SDK so Koda Studio.app sits next to koda and stdlib/."'
  exit 1
fi
open "$APP" --args "$@"
