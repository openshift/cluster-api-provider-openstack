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

BIN_DIR := bin
TOOLS_DIR := tools

$(RELEASE_DIR):
	mkdir -p $(RELEASE_DIR)/

MANIFESTS_GEN := go run ./vendor/github.com/openshift/cluster-capi-operator/manifests-gen/

.PHONY: check-env
check-env:
ifndef PROVIDER_VERSION
	$(error PROVIDER_VERSION is undefined)
endif

.PHONY: update-manifests-gen
update-manifests-gen:
	cd tools && go get github.com/openshift/cluster-capi-operator/manifests-gen && go mod tidy && go mod vendor

.PHONY: ocp-manifests
ocp-manifests: $(RELEASE_DIR) check-env ## Builds openshift specific manifests.
	# Generate provider manifests.
	cd tools && $(MANIFESTS_GEN) --provider-name "openstack" --provider-type "InfrastructureProvider" --provider-version "${PROVIDER_VERSION}" --base-path "../../" --manifests-path "../manifests"

# FIXME(stephenfin): Delete everything from here down once we are fully migrated to manifests-gen

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

.PHONY: verify
verify: generate
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "generated files are out of date, run make generate"; exit 1; \
	fi

.PHONY: ALWAYS
ALWAYS:
