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
# Conflict resolution can leave openshift/ artifacts without openshift/go.mod,
# which makes them part of the main module with stale cluster-api import paths
# (e.g. sigs.k8s.io/cluster-api/api/v1beta1) and breaks go mod tidy.
rm -rf openshift/vendor openshift/go.mod openshift/e2e openshift/pkg

# Rebase conflicts can also leave apivalidations registering removed API
# versions (v1alpha6/v1alpha7). Restore the known-good file from dest/main.
if git show dest/main:test/e2e/suites/apivalidations/suite_test.go >/dev/null 2>&1; then
    git show dest/main:test/e2e/suites/apivalidations/suite_test.go \
        > test/e2e/suites/apivalidations/suite_test.go
fi

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
