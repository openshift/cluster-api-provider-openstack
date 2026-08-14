#!/bin/bash

# This script is run by rebasebot as a post-rebase hook during the rebase of
# openshift/cluster-api-provider-openstack.
# It replaces the merge-bot's --run-make flag which ran `make merge-bot`.

set -e
set -o pipefail

if [[ -n "$(git status --porcelain)" ]]; then
    echo "post-rebase hook requires a clean worktree" >&2
    exit 1
fi

# Rebase replays old downstream commits that predate cluster-capi-operator.
# Conflict resolution can leave stale openshift/ artifacts, removed API version
# directories, and main.go/suite_test.go scheme registrations.
"$(dirname "$0")/cleanup-stale-rebase-artifacts.sh"

make merge-bot

if [[ -z "$REBASEBOT_GIT_USERNAME" || -z "$REBASEBOT_GIT_EMAIL" ]]; then
    author_flag=()
else
    author_flag=(--author="$REBASEBOT_GIT_USERNAME <$REBASEBOT_GIT_EMAIL>")
fi

if [[ -n $(git status --porcelain) ]]; then
    git add -A
    git commit "${author_flag[@]}" -q -m "UPSTREAM: <drop>: Run make merge-bot"
fi
