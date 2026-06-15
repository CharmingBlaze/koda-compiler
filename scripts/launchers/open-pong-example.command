#!/bin/bash
cd "$(dirname "$0")"
APP="./Koda Studio.app"
if [[ ! -d "$APP" ]]; then
  osascript -e 'display alert "Koda Studio not found"'
  exit 1
fi
if [[ -d "./examples/games/pong" ]]; then
  open "$PWD/Koda Studio.app" --args "$PWD/examples/games/pong"
else
  open "$PWD/Koda Studio.app"
fi
