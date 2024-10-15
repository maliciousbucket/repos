#!/usr/bin/env bash

read_env() {
  local filePath="${1:-.env}"
  if [ ! -f "$filePath" ]; then
    echo "missing ${filePath}"
    exit 1
  fi

  echo "Reading $filePath"
  while read -r LINE; do
    CLEANED_LINE=$(echo "$LINE" | awk '{$1=$1};1' | tr -d '\r')

    if [[ $CLEANED_LINE != '#'* ]] && [[ $CLEANED_LINE == *'='* ]]; then
      VAR_NAME=$(echo "$CLEANED_LINE" | cut -d '=' -f 1)
      VAR_VALUE=$(echo "$CLEANED_LINE" | cut -d '=' -f 2- | sed 's/^"//;s/"$//')
      export "$VAR_NAME=$VAR_VALUE"
    fi
  done < "$filePath"
}