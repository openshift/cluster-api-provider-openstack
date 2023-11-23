package e2e

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	configv1 "github.com/openshift/api/config/v1"
	mapiv1alpha1 "github.com/openshift/api/machine/v1alpha1"
	mapiv1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/cluster-capi-operator/e2e/framework"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	yaml "sigs.k8s.io/yaml"

	"github.com/openshift/cluster-api-provider-openstack/openshift/pkg/infraclustercontroller"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
)

const (
	openStackMachineTemplateName = "openstack-machine-template"
)

var _ = Describe("Cluster API OpenStack MachineSet", Ordered, func() {
	var openStackMachineTemplate *infrav1.OpenStackMachineTemplate
	var machineSet *clusterv1.MachineSet
	var mapiMachineSpec *mapiv1alpha1.OpenstackProviderSpec

	BeforeAll(func() {
		if platform != configv1.OpenStackPlatformType {
			Skip("Skipping OpenStack E2E tests")
		}
		framework.CreateCoreCluster(cl, clusterName, "OpenStackCluster")
		mapiMachineSpec = getOpenStackMAPIProviderSpec(cl)
	})

	AfterEach(func() {
		if platform != configv1.OpenStackPlatformType {
			// Because AfterEach always runs, even when tests are skipped, we have to
			// explicitly skip it here for other platforms.
			Skip("Skipping OpenStack E2E tests")
		}
		framework.DeleteMachineSets(cl, machineSet)
		framework.WaitForMachineSetsDeleted(cl, machineSet)
		framework.DeleteObjects(cl, openStackMachineTemplate)
	})

	It("should be able to run a machine", func() {
		openStackMachineTemplate = createOpenStackMachineTemplate(cl, mapiMachineSpec)

		machineSet = framework.CreateMachineSet(cl, framework.NewMachineSetParams(
			"openstack-machineset",
			clusterName,
			"",
			1,
			corev1.ObjectReference{
				Kind:       "OpenStackMachineTemplate",
				APIVersion: infraAPIVersion,
				Name:       openStackMachineTemplateName,
			},
		))

		framework.WaitForMachineSet(cl, machineSet.Name)
	})
})

func getOpenStackMAPIProviderSpec(cl client.Client) *mapiv1alpha1.OpenstackProviderSpec {
	machineSetList := &mapiv1beta1.MachineSetList{}
	Expect(cl.List(ctx, machineSetList, client.InNamespace(framework.MAPINamespace))).To(Succeed())

	Expect(machineSetList.Items).ToNot(HaveLen(0))
	machineSet := machineSetList.Items[0]
	Expect(machineSet.Spec.Template.Spec.ProviderSpec.Value).ToNot(BeNil())

	providerSpec := &mapiv1alpha1.OpenstackProviderSpec{}
	Expect(yaml.Unmarshal(machineSet.Spec.Template.Spec.ProviderSpec.Value.Raw, providerSpec)).To(Succeed())

	return providerSpec
}

func createOpenStackMachineTemplate(cl client.Client, mapiProviderSpec *mapiv1alpha1.OpenstackProviderSpec) *infrav1.OpenStackMachineTemplate {
	By("Creating OpenStack machine template")

	Expect(mapiProviderSpec).ToNot(BeNil())
	Expect(mapiProviderSpec.Flavor).ToNot(BeEmpty())
	// NOTE(stephenfin): Installer does not populate ps.Image when ps.RootVolume is set and will
	// instead populate ps.RootVolume.SourceUUID. Moreover, according to the ClusterOSImage option
	// definition this is always the name of the image and never the UUID. We should allow UUID
	// at some point and this will need an update.
	if mapiProviderSpec.RootVolume != nil {
		Expect(mapiProviderSpec.RootVolume.SourceUUID).ToNot(BeEmpty())
	} else {
		Expect(mapiProviderSpec.Image).ToNot(BeEmpty())
	}
	Expect(len(mapiProviderSpec.Networks)).To(BeNumerically(">", 0))
	Expect(len(mapiProviderSpec.Networks[0].Subnets)).To(BeNumerically(">", 0))
	Expect(mapiProviderSpec.Tags).ToNot(BeNil())
	Expect(len(mapiProviderSpec.Tags)).To(BeNumerically(">", 0))

	var image string
	var rootVolume *infrav1.RootVolume

	if mapiProviderSpec.RootVolume != nil {
		rootVolume = &infrav1.RootVolume{
			Size:             mapiProviderSpec.RootVolume.Size,
			VolumeType:       mapiProviderSpec.RootVolume.VolumeType,
			AvailabilityZone: mapiProviderSpec.RootVolume.Zone,
		}
	} else {
		image = mapiProviderSpec.Image
	}

	// NOTE(stephenfin): We intentionally ignore additional networks for now since we don't care
	// about e.g. Manila shares. We can re-evaluate this if necessary.
	ports := []infrav1.PortOpts{}
	for _, subnet := range mapiProviderSpec.Networks[0].Subnets {
		port := infrav1.PortOpts{
			FixedIPs: []infrav1.FixedIP{
				{
					Subnet: &infrav1.SubnetFilter{
						// NOTE(stephenfin): Only one of name or ID will be set.
						ID:   subnet.Filter.ID,
						Name: subnet.Filter.Name,
						Tags: subnet.Filter.Tags,
					},
				},
			},
		}
		ports = append(ports, port)
	}
	port := infrav1.PortOpts{
		Network: &infrav1.NetworkFilter{
			// The installer still sets NetworkParam.Filter.ID rather than NetworkParam.ID,
			// at least as of 4.15.
			ID:   mapiProviderSpec.Networks[0].Filter.ID, //nolint:staticcheck
			Name: mapiProviderSpec.Networks[0].Filter.Name,
		},
	}
	ports = append(ports, port)

	// NOTE(stephenfin): Ditto for security groups
	securityGroups := []infrav1.SecurityGroupFilter{
		{
			Name: mapiProviderSpec.SecurityGroups[0].Name,
			ID:   mapiProviderSpec.SecurityGroups[0].UUID,
		},
	}

	openStackMachineSpec := infrav1.OpenStackMachineSpec{
		CloudName: infraclustercontroller.CloudName,
		Flavor:    mapiProviderSpec.Flavor,
		IdentityRef: &infrav1.OpenStackIdentityReference{
			Kind: "Secret",
			Name: infraclustercontroller.CredentialsSecretName,
		},
		Image:          image,
		Ports:          ports,
		RootVolume:     rootVolume,
		SecurityGroups: securityGroups,
		Tags:           mapiProviderSpec.Tags,
	}

	openStackMachineTemplate := &infrav1.OpenStackMachineTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      openStackMachineTemplateName,
			Namespace: framework.CAPINamespace,
		},
		Spec: infrav1.OpenStackMachineTemplateSpec{
			Template: infrav1.OpenStackMachineTemplateResource{
				Spec: openStackMachineSpec,
			},
		},
	}

	if err := cl.Create(ctx, openStackMachineTemplate); err != nil && !apierrors.IsAlreadyExists(err) {
		Expect(err).ToNot(HaveOccurred())
	}

	return openStackMachineTemplate
}
