manifests_dir ?= ../manifests
manifests_prefix ?= 0000_30_cluster-api-provider-openstack_

TOOLS_DIR=../hack/tools
KUSTOMIZE=$(TOOLS_DIR)/bin/kustomize

define manifest_name
    $(addsuffix ".yaml",$(addprefix $(manifests_dir)/$(manifests_prefix),$(1)))
endef

manifest_names = 02_crds 03_rbac 04_infrastructure-components
source_manifest = manifests/source-manifest.yaml

.PHONY: all_manifests
all_manifests: $(foreach m,$(manifest_names),$(call manifest_name,$(m)))

$(call manifest_name,02_crds): $(KUSTOMIZE) ALWAYS | $(manifests_dir)
	$(KUSTOMIZE) build crd > $@

$(call manifest_name,03_rbac): $(KUSTOMIZE) ALWAYS | $(manifests_dir)
	$(KUSTOMIZE) build rbac > $@

$(source_manifest): $(KUSTOMIZE) ALWAYS | $(manifests_dir)
	$(KUSTOMIZE) build infrastructure-components > $@

$(call manifest_name,04_infrastructure-components): $(KUSTOMIZE) $(source_manifest) manifests/kustomization.yaml | $(manifests_dir)
	$(KUSTOMIZE) build manifests > $@


$(manifests_dir):
	mkdir -p $@

$(KUSTOMIZE):
	$(MAKE) -C $(TOOLS_DIR) bin/kustomize

.PHONY: ALWAYS
ALWAYS:
