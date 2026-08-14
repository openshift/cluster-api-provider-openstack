#!/bin/bash

# Remove artifacts that old downstream carry commits can reintroduce during a
# rebasebot rebase. Without this cleanup, go mod tidy fails on stale cluster-api
# import paths (e.g. sigs.k8s.io/cluster-api/api/v1beta1) and removed API
# versions (v1alpha5/v1alpha6/v1alpha7).

set -euo pipefail

rm -rf openshift/vendor openshift/go.mod openshift/e2e openshift/pkg
rm -rf api/v1alpha5 api/v1alpha6 api/v1alpha7

restore_ref=""
if [[ -n "${REBASEBOT_SOURCE:-}" ]] && git show "source/${REBASEBOT_SOURCE}:main.go" >/dev/null 2>&1; then
    restore_ref="source/${REBASEBOT_SOURCE}"
elif git show dest/main:main.go >/dev/null 2>&1; then
    restore_ref="dest/main"
fi

if [[ -n "$restore_ref" ]]; then
    for f in main.go test/e2e/suites/apivalidations/suite_test.go; do
        if git show "${restore_ref}:${f}" >/dev/null 2>&1; then
            git show "${restore_ref}:${f}" > "$f"
        fi
    done
fi
