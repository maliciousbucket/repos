#!/usr/bin/env sh

set -eux
set -o pipefail


source tools/utils.sh

read_env "${ENV}"

check_env() {
  if [  -z  "${GITHUB_TOKEN}" ]; then
      echo -e "GITHUB_TOKEN env variable is empty"
     exit 1
  fi

  if [ -z  "${GITHUB_REPOSITORY}" ]; then
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
       -d '{"event_type": "'${ACTION_HOOK}'", "client_payload": {"delay": "'${DELAY}'"}}' \
       https://api.github.com/repos/"${REPO_NAME}"/dispatches)

  echo "Response status: ${RESPONSE}"
}

action="${1}"


case "${action}" in
  "deployment" )
    run_workflow deployment
  ;;
  "test"| "test-run" )
  run_workflow test_run
  ;;
  "ci" | "ci-run" )
  run_workflow trigger_ci
  ;;
  *)
  echo  "Unknown workflow type: ${action}"
  exit 1
  ;;
esac