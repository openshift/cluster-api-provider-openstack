manifests_dir ?= ./manifests
manifests_prefix ?= 0000_30_cluster-api-provider-openstack_

define manifest_name
    $(addsuffix ".yaml",$(addprefix $(manifests_dir)/$(manifests_prefix),$(1)))
endef

.PHONY: manifests
all_manifests: $(foreach m,02_crds 03_rbac,$(call manifest_name,$(m)))

$(call manifest_name,02_crds): ALWAYS | $(manifests_dir)
	kustomize build config/crd > $@

$(call manifest_name,03_rbac): ALWAYS | $(manifests_dir)
	kustomize build config/rbac > $@

$(manifests_dir):
	mkdir -p $@

.PHONY: ALWAYS
ALWAYS:
