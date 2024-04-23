package e2e

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	configv1 "github.com/openshift/api/config/v1"
	mapiv1alpha1 "github.com/openshift/api/machine/v1alpha1"
	mapiv1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/cluster-capi-operator/e2e/framework"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	yaml "sigs.k8s.io/yaml"

	"github.com/openshift/cluster-api-provider-openstack/openshift/pkg/infraclustercontroller"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

const (
	openStackMachineTemplateName = "openstack-machine-template"
)

var _ = Describe("Cluster API OpenStack MachineSet", Ordered, func() {
	var mapiMachineSpec *mapiv1alpha1.OpenstackProviderSpec

	BeforeAll(func() {
		if platform != configv1.OpenStackPlatformType {
			Skip("Skipping OpenStack E2E tests")
		}
		mapiMachineSpec = getOpenStackMAPIProviderSpec(cl)
	})

	It("should be able to run a machine with implicit cluster default network", func() {
		openStackMachineTemplate := createOpenStackMachineTemplate(cl, mapiMachineSpec)

		machineSet := framework.CreateMachineSet(cl, framework.NewMachineSetParams(
			"openstack-machineset",
			clusterName,
			"",
			1,
			corev1.ObjectReference{
				Kind:       "OpenStackMachineTemplate",
				APIVersion: infraAPIVersion,
				Name:       openStackMachineTemplate.Name,
			},
		))
		DeferCleanup(func() {
			By("Deleting machineset " + machineSet.Name)
			Expect(cl.Delete(ctx, machineSet)).To(Succeed())
			framework.WaitForMachineSetsDeleted(cl, machineSet)
		})

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
			SizeGiB: mapiProviderSpec.RootVolume.Size,
			BlockDeviceVolume: infrav1.BlockDeviceVolume{
				Type: mapiProviderSpec.RootVolume.VolumeType,
				AvailabilityZone: &infrav1.VolumeAvailabilityZone{
					From: infrav1.VolumeAZFromName,
					Name: ptr.To(infrav1.VolumeAZName(mapiProviderSpec.RootVolume.Zone)),
				},
			},
		}
	} else {
		image = mapiProviderSpec.Image
	}

	// NOTE(stephenfin): We intentionally ignore additional security for now.
	var securityGroupParam infrav1.SecurityGroupParam
	securityGroup := mapiProviderSpec.SecurityGroups[0]
	if securityGroup.UUID != "" {
		securityGroupParam = infrav1.SecurityGroupParam{ID: &securityGroup.UUID}
	} else {
		securityGroupParam = infrav1.SecurityGroupParam{Filter: &infrav1.SecurityGroupFilter{Name: securityGroup.Name}}
	}
	securityGroups := []infrav1.SecurityGroupParam{
		securityGroupParam,
	}

	// We intentionally omit ports so the machine will default its network
	// from the OpenStackCluster created by the infracluster controller.
	openStackMachineSpec := infrav1.OpenStackMachineSpec{
		Flavor: mapiProviderSpec.Flavor,
		IdentityRef: &infrav1.OpenStackIdentityReference{
			CloudName: infraclustercontroller.CloudName,
			Name:      infraclustercontroller.CredentialsSecretName,
		},
		Image:          infrav1.ImageParam{Filter: &infrav1.ImageFilter{Name: &image}},
		RootVolume:     rootVolume,
		SecurityGroups: securityGroups,
	}

	openStackMachineTemplate := &infrav1.OpenStackMachineTemplate{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: openStackMachineTemplateName + "-",
			Namespace:    framework.CAPINamespace,
		},
		Spec: infrav1.OpenStackMachineTemplateSpec{
			Template: infrav1.OpenStackMachineTemplateResource{
				Spec: openStackMachineSpec,
			},
		},
	}

	Expect(cl.Create(ctx, openStackMachineTemplate)).To(Succeed(), format.Object(openStackMachineTemplate, 1))
	DeferCleanup(func() error {
		return cl.Delete(ctx, openStackMachineTemplate)
	})

	return openStackMachineTemplate
}
