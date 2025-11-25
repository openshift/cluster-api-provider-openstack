# Copyright 2023 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

manifests_dir ?= ./manifests
manifests_prefix ?= 0000_30_cluster-api-provider-openstack_

TOOLS_DIR=../hack/tools
KUSTOMIZE=$(TOOLS_DIR)/bin/kustomize
CONTROLLER_GEN=$(TOOLS_DIR)/bin/controller-gen
GOLANGCI_LINT=$(TOOLS_DIR)/bin/golangci-lint

define manifest_name
    $(addsuffix ".yaml",$(addprefix $(manifests_dir)/$(manifests_prefix),$(1)))
endef

manifest_names = 00_credentials-request 04_infrastructure-components
infrastructure_components = kustomize/cluster-capi-configmap/infrastructure-components.yaml

.PHONY: generate
generate: $(foreach m,$(manifest_names),$(call manifest_name,$(m)))

$(call manifest_name,00_credentials-request): $(KUSTOMIZE) ALWAYS | $(manifests_dir)
	$(KUSTOMIZE) build kustomize/credentials-request > $@

$(infrastructure_components): $(KUSTOMIZE) ALWAYS
	$(KUSTOMIZE) build kustomize/infrastructure-components > $@

$(call manifest_name,04_infrastructure-components): $(KUSTOMIZE) $(infrastructure_components) ALWAYS | $(manifests_dir)
	$(KUSTOMIZE) build kustomize/cluster-capi-configmap > $@

$(manifests_dir):
	mkdir -p $@

$(KUSTOMIZE):
	$(MAKE) -C $(TOOLS_DIR) bin/kustomize

$(CONTROLLER_GEN):
	$(MAKE) -C $(TOOLS_DIR) bin/controller-gen

$(GOLANGCI_LINT):
	$(MAKE) -C $(TOOLS_DIR) bin/golangci-lint

.PHONY: verify
verify: lint modules generate
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "generated files are out of date, run make generate"; exit 1; \
	fi

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint codebase
	$(GOLANGCI_LINT) run -v --fast=false

.PHONY: lint-update
lint-update: $(GOLANGCI_LINT) ## Lint codebase
	$(GOLANGCI_LINT) run -v --fast=false --fix

.PHONY: modules
modules:
	go mod tidy
	go mod vendor

.PHONY: test
test:
	go test $(shell go list ./... | grep -v /e2e)

.PHONY: e2e
e2e:
	./hack/test.sh "./e2e/..." 30m

.PHONY: ALWAYS
ALWAYS:
