# CI Fix Plan

## Summary
The CI verification job `pull-ci-openshift-cluster-api-provider-openstack-main-verify` failed because vendored files in `hack/tools/vendor/` are out of date compared to the current code. Running `make vendor` updates the vendor directories and resolves the drift.

## Fixable Failures

### ci/prow/verify
**Category**: codegen-outdated
**Root Cause**: Vendored files are out of date. Specifically, `hack/tools/vendor/sigs.k8s.io/cluster-api-provider-openstack/controllers/openstackmachine_controller.go` is not synchronized with the root `controllers/openstackmachine_controller.go`.
**Affected Files**:
- `hack/tools/vendor/sigs.k8s.io/cluster-api-provider-openstack/controllers/openstackmachine_controller.go`
- Any other out-of-sync files in `hack/tools/vendor/`
**Fix**:
1. Run `make vendor` (which runs `go mod vendor` and `cd hack/tools; go mod vendor`).
2. Run `make verify-vendoring` (or `make verify`) to confirm the files match.

## Skipped Failures
*None*
