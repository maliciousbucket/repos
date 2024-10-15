#!/usr/bin/env bash

set -eux
set -o pipefail

env_file="${2}"
source ./utils.sh

read_env "${env_file}"

check_env() {
  if [ "${GITHUB_TOKEN}" == "" ]; then
      echo -e "GITHUB_TOKEN env variable is empty"
     exit 1
  fi

  if [ "${GITHUB_REPOSITORY}" == "" ]; then
      echo -e "GITHUB_REPOSITORY env variable is empty"
     exit 1
  fi

}

run_workflow() {
  check_env
  local TOKEN=$GITHUB_TOKEN
  local REPO_NAME=$GITHUB_REPOSITORY
  local ACTION_HOOK=$1

  RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" \
       -H "Accept: application/vnd.github+json" \
       -H "Authorization: token ${TOKEN}" \
       -H "Content-Type: application/json" \
       -d '{"event_type": "'${ACTION_HOOK}'"}' \
       https://api.github.com/repos/"${REPO_NAME}"/dispatches)

  echo "Response status: ${RESPONSE}"
}

action="${1}"


case "${action}" in
  "deployment" )
    run_workflow "deployment"
  ;;
  "test"| "test-run" )
  run_workflow test_run
  ;;
  "ci" | "ci-run" )
  run_workflow "trigger_ci"
  ;;
  *)
  echo -e "Unknown workflow type: ${action}"
  exit 1
  ;;
esac