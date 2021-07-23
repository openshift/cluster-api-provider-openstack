/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clients

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/common/extensions"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/schedulerhints"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	netext "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsbinding"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/utils/openstack/clientconfig"
	azutils "github.com/gophercloud/utils/openstack/compute/v2/availabilityzones"
	flavorutils "github.com/gophercloud/utils/openstack/compute/v2/flavors"
	imageutils "github.com/gophercloud/utils/openstack/imageservice/v2/images"
	configclient "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"github.com/openshift/machine-api-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	openstackconfigv1 "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
)

const (
	CloudsSecretKey = "clouds.yaml"

	TimeoutTrunkDelete       = 3 * time.Minute
	RetryIntervalTrunkDelete = 5 * time.Second

	TimeoutPortDelete       = 3 * time.Minute
	RetryIntervalPortDelete = 5 * time.Second

	// Maximum port name length supported by Neutron
	PortNameMaxSize = 255

	// MachineRegionLabelName as annotation name for a machine region
	MachineRegionLabelName = "machine.openshift.io/region"

	// MachineAZLabelName as annotation name for a machine AZ
	MachineAZLabelName = "machine.openshift.io/zone"

	// MachineInstanceTypeLabelName as annotation name for a machine instance type
	MachineInstanceTypeLabelName = "machine.openshift.io/instance-type"
)

type InstanceService struct {
	provider       *gophercloud.ProviderClient
	computeClient  *gophercloud.ServiceClient
	identityClient *gophercloud.ServiceClient
	networkClient  *gophercloud.ServiceClient
	imagesClient   *gophercloud.ServiceClient
	volumeClient   *gophercloud.ServiceClient

	regionName string
}

type Instance struct {
	servers.Server
}

type ServerNetwork struct {
	networkID    string
	subnetID     string
	portTags     []string
	vnicType     string
	portSecurity *bool
}

// for updating the state of ports with port security
var portWithPortSecurityExtensions struct {
	ports.Port
	portsecurity.PortSecurityExt
}

type InstanceListOpts struct {
	// Name of the image in URL format.
	Image string `q:"image"`

	// Name of the flavor in URL format.
	Flavor string `q:"flavor"`

	// Name of the server as a string; can be queried with regular expressions.
	// Realize that ?name=bob returns both bob and bobb. If you need to match bob
	// only, you can use a regular expression matching the syntax of the
	// underlying database server implemented for Compute.
	Name string `q:"name"`
}

type serverMetadata struct {
	// AZ contains name of the server's availability zone
	AZ string `json:"OS-EXT-AZ:availability_zone"`

	// Flavor refers to a JSON object, which itself indicates the hardware
	// configuration of the deployed server.
	Flavor map[string]interface{} `json:"flavor"`

	// Status contains the current operational status of the server,
	// such as IN_PROGRESS or ACTIVE.
	Status string `json:"status"`
}

func GetCloudFromSecret(kubeClient kubernetes.Interface, namespace string, secretName string, cloudName string) (clientconfig.Cloud, error) {
	emptyCloud := clientconfig.Cloud{}

	if secretName == "" {
		return emptyCloud, nil
	}

	if secretName != "" && cloudName == "" {
		return emptyCloud, fmt.Errorf("Secret name set to %v but no cloud was specified. Please set cloud_name in your machine spec.", secretName)
	}

	secret, err := kubeClient.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return emptyCloud, fmt.Errorf("Failed to get secrets from kubernetes api: %v", err)
	}

	content, ok := secret.Data[CloudsSecretKey]
	if !ok {
		return emptyCloud, fmt.Errorf("OpenStack credentials secret %v did not contain key %v",
			secretName, CloudsSecretKey)
	}
	var clouds clientconfig.Clouds
	err = yaml.Unmarshal(content, &clouds)
	if err != nil {
		return emptyCloud, fmt.Errorf("failed to unmarshal clouds credentials stored in secret %v: %v", secretName, err)
	}

	return clouds.Clouds[cloudName], nil
}

// TODO: Eventually we'll have a NewInstanceServiceFromCluster too
func NewInstanceServiceFromMachine(kubeClient kubernetes.Interface, machine *machinev1.Machine) (*InstanceService, error) {
	cloud, err := GetCloud(kubeClient, machine)
	if err != nil {
		return nil, err
	}

	return NewInstanceServiceFromCloud(cloud, GetCACertificate(kubeClient))
}

func NewInstanceService() (*InstanceService, error) {
	cloud := clientconfig.Cloud{}
	return NewInstanceServiceFromCloud(cloud, nil)
}

func NewInstanceServiceFromCloud(cloud clientconfig.Cloud, cert []byte) (*InstanceService, error) {
	provider, err := GetProviderClient(cloud, cert)
	if err != nil {
		return nil, err
	}

	identityClient, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{
		Region: "",
	})
	if err != nil {
		return nil, fmt.Errorf("Create identityClient err: %v", err)
	}
	serverClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: cloud.RegionName,
	})

	if err != nil {
		return nil, fmt.Errorf("Create serviceClient err: %v", err)
	}

	networkingClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: cloud.RegionName,
	})
	if err != nil {
		return nil, fmt.Errorf("Create networkingClient err: %v", err)
	}

	imagesClient, err := openstack.NewImageServiceV2(provider, gophercloud.EndpointOpts{
		Region: cloud.RegionName,
	})
	if err != nil {
		return nil, fmt.Errorf("Create ImageClient err: %v", err)
	}

	volumeClient, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{
		Region: cloud.RegionName,
	})
	if err != nil {
		return nil, fmt.Errorf("Create VolumeClient err: %v", err)
	}

	return &InstanceService{
		provider:       provider,
		identityClient: identityClient,
		computeClient:  serverClient,
		networkClient:  networkingClient,
		imagesClient:   imagesClient,
		volumeClient:   volumeClient,
		regionName:     cloud.RegionName,
	}, nil
}

// UpdateToken to update token if need.
func (is *InstanceService) UpdateToken() error {
	token := is.provider.Token()
	result, err := tokens.Validate(is.identityClient, token)
	if err != nil {
		return fmt.Errorf("Validate token err: %v", err)
	}
	if result {
		return nil
	}
	klog.V(2).Infof("Token is out of date, getting new token.")
	reAuthFunction := is.provider.ReauthFunc
	if reAuthFunction() != nil {
		return fmt.Errorf("reAuth err: %v", err)
	}
	return nil
}

func (is *InstanceService) AssociateFloatingIP(instanceID, floatingIP string) error {
	opts := floatingips.AssociateOpts{
		FloatingIP: floatingIP,
	}
	return floatingips.AssociateInstance(is.computeClient, instanceID, opts).ExtractErr()
}

func (is *InstanceService) GetAcceptableFloatingIP() (string, error) {
	page, err := floatingips.List(is.computeClient).AllPages()
	if err != nil {
		return "", fmt.Errorf("Get floating IP list failed: %v", err)
	}
	list, err := floatingips.ExtractFloatingIPs(page)
	if err != nil {
		return "", err
	}
	for _, floatingIP := range list {
		if floatingIP.FixedIP == "" {
			return floatingIP.IP, nil
		}
	}
	return "", fmt.Errorf("Don't have acceptable floating IP")
}

// A function for getting the id of a network by querying openstack with filters
func getNetworkIDsByFilter(is *InstanceService, opts *networks.ListOpts) ([]string, error) {
	if opts == nil {
		return []string{}, fmt.Errorf("No Filters were passed")
	}
	pager := networks.List(is.networkClient, opts)
	var uuids []string
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)
		if err != nil {
			return false, err
		} else if len(networkList) == 0 {
			return false, fmt.Errorf("No networks could be found with the filters provided")
		}
		for _, network := range networkList {
			uuids = append(uuids, network.ID)
		}
		return true, nil
	})
	if err != nil {
		return []string{}, err
	}
	return uuids, nil
}

// A function for getting the id of a subnet by querying openstack with filters
func getSubnetsByFilter(is *InstanceService, opts *subnets.ListOpts) ([]subnets.Subnet, error) {
	if opts == nil {
		return []subnets.Subnet{}, fmt.Errorf("No Filters were passed")
	}
	pager := subnets.List(is.networkClient, opts)
	var snets []subnets.Subnet
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		subnetList, err := subnets.ExtractSubnets(page)
		if err != nil {
			return false, err
		} else if len(subnetList) == 0 {
			return false, fmt.Errorf("No subnets could be found with the filters provided")
		}
		for _, subnet := range subnetList {
			snets = append(snets, subnet)
		}
		return true, nil
	})
	if err != nil {
		return []subnets.Subnet{}, err
	}
	return snets, nil
}

func getOrCreatePort(is *InstanceService, name string, portOpts openstackconfigv1.PortOpts) (*ports.Port, error) {
	var portName string
	if portOpts.NameSuffix != "" {
		portName = name + "-" + portOpts.NameSuffix
	} else {
		portName = name
	}
	if len(portName) > PortNameMaxSize {
		portName = portName[len(portName)-PortNameMaxSize:]
	}
	existingPorts, err := listPorts(is, ports.ListOpts{
		Name:      portName,
		NetworkID: portOpts.NetworkID,
	})
	if err != nil {
		return nil, err
	}
	if len(existingPorts) == 0 {
		createOpts := ports.CreateOpts{
			Name:                portName,
			NetworkID:           portOpts.NetworkID,
			Description:         portOpts.Description,
			AdminStateUp:        portOpts.AdminStateUp,
			MACAddress:          portOpts.MACAddress,
			TenantID:            portOpts.TenantID,
			ProjectID:           portOpts.ProjectID,
			SecurityGroups:      portOpts.SecurityGroups,
			AllowedAddressPairs: []ports.AddressPair{},
		}

		for _, ap := range portOpts.AllowedAddressPairs {
			createOpts.AllowedAddressPairs = append(createOpts.AllowedAddressPairs, ports.AddressPair{
				IPAddress:  ap.IPAddress,
				MACAddress: ap.MACAddress,
			})
		}
		if len(portOpts.FixedIPs) != 0 {
			fixedIPs := make([]ports.IP, len(portOpts.FixedIPs))
			for i, portOptIP := range portOpts.FixedIPs {
				fixedIPs[i].SubnetID = portOptIP.SubnetID
				fixedIPs[i].IPAddress = portOptIP.IPAddress
			}
			createOpts.FixedIPs = fixedIPs
		}
		newPort, err := ports.Create(is.networkClient, portsbinding.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			HostID:            portOpts.HostID,
			VNICType:          portOpts.VNICType,
			Profile:           nil,
		}).Extract()
		if err != nil {
			return nil, err
		}

		if portOpts.PortSecurity != nil {
			portUpdateOpts := ports.UpdateOpts{}
			if *portOpts.PortSecurity == false {
				portUpdateOpts.SecurityGroups = &[]string{}
				portUpdateOpts.AllowedAddressPairs = &[]ports.AddressPair{}
			}
			updateOpts := portsecurity.PortUpdateOptsExt{
				UpdateOptsBuilder:   portUpdateOpts,
				PortSecurityEnabled: portOpts.PortSecurity,
			}
			err = ports.Update(is.networkClient, newPort.ID, updateOpts).ExtractInto(&portWithPortSecurityExtensions)
			if err != nil {
				return nil, fmt.Errorf("Failed to update port security on port %s: %v", newPort.ID, err)
			}
		}

		return newPort, nil
	} else if len(existingPorts) == 1 {
		return &existingPorts[0], nil
	}

	return nil, fmt.Errorf("multiple ports found with name \"%s\"", portName)
}

func listPorts(is *InstanceService, opts ports.ListOpts) ([]ports.Port, error) {
	allPages, err := ports.List(is.networkClient, opts).AllPages()
	if err != nil {
		return []ports.Port{}, err
	}

	portList, err := ports.ExtractPorts(allPages)
	if err != nil {
		return []ports.Port{}, err
	}

	return portList, nil
}

func isDuplicate(list []string, name string) bool {
	if list == nil || len(list) == 0 {
		return false
	}
	for _, element := range list {
		if element == name {
			return true
		}
	}
	return false
}

func GetSecurityGroups(is *InstanceService, sg_param []openstackconfigv1.SecurityGroupParam) ([]string, error) {
	var sgIDs []string
	for _, sg := range sg_param {
		listOpts := groups.ListOpts(sg.Filter)
		listOpts.Name = sg.Name
		listOpts.ID = sg.UUID
		pages, err := groups.List(is.networkClient, listOpts).AllPages()
		if err != nil {
			return nil, err
		}

		SGList, err := groups.ExtractGroups(pages)
		if err != nil {
			return nil, err
		}

		for _, group := range SGList {
			if isDuplicate(sgIDs, group.ID) {
				continue
			}
			sgIDs = append(sgIDs, group.ID)
		}
	}
	return sgIDs, nil
}

// InstanceCreate creates a compute instance.
// If ServerGroupName is nonempty and no server group exists with that name,
// then InstanceCreate creates a server group with that name.
func (is *InstanceService) InstanceCreate(clusterName string, name string, clusterSpec *openstackconfigv1.OpenstackClusterProviderSpec, config *openstackconfigv1.OpenstackProviderSpec, cmd string, keyName string, configClient configclient.ConfigV1Interface) (instance *Instance, err error) {
	// server is only non-nil in case of successful server creation.
	//
	// There are multiple preparation steps in this method, some of which
	// create resources (e.g. a volume in case of boot-from-volume with
	// "image" source), and all of which have a chance to fail.
	//
	// This variable is guaranteed to remain nil until server creation is
	// successful. Deferred cleanup functions can restore the initial state
	// in case of failure, by only acting if "server" is nil.
	var server *servers.Server

	if config == nil {
		return nil, fmt.Errorf("create Options need be specified to create instace")
	}
	if config.Trunk == true {
		trunkSupport, err := GetTrunkSupport(is)
		if err != nil {
			return nil, fmt.Errorf("There was an issue verifying whether trunk support is available, please disable it: %v", err)
		}
		if trunkSupport == false {
			return nil, fmt.Errorf("There is no trunk support. Please disable it")
		}
	}

	// Set default Tags
	machineTags := []string{
		"cluster-api-provider-openstack",
		clusterName,
	}

	// Append machine specific tags
	machineTags = append(machineTags, config.Tags...)

	// Append cluster scope tags
	if clusterSpec != nil && clusterSpec.Tags != nil {
		machineTags = append(machineTags, clusterSpec.Tags...)
	}

	// Get security groups
	securityGroups, err := GetSecurityGroups(is, config.SecurityGroups)
	if err != nil {
		return nil, err
	}
	// Get all network UUIDs
	var nets []openstackconfigv1.PortOpts
	netsWithoutAllowedAddressPairs := map[string]struct{}{}
	for _, net := range config.Networks {
		opts := networks.ListOpts(net.Filter)
		opts.ID = net.UUID
		ids, err := getNetworkIDsByFilter(is, &opts)
		if err != nil {
			return nil, err
		}
		for _, netID := range ids {
			if net.NoAllowedAddressPairs {
				netsWithoutAllowedAddressPairs[netID] = struct{}{}
			}
			if net.Subnets == nil {
				nets = append(nets, openstackconfigv1.PortOpts{
					NetworkID:    netID,
					NameSuffix:   net.UUID,
					Tags:         net.PortTags,
					VNICType:     net.VNICType,
					PortSecurity: net.PortSecurity,
				})
			}

			for _, snetParam := range net.Subnets {
				sopts := subnets.ListOpts(snetParam.Filter)
				sopts.ID = snetParam.UUID
				sopts.NetworkID = netID
				// Inherit portSecurity from network if unset on subnet
				portSecurity := net.PortSecurity
				if snetParam.PortSecurity != nil {
					portSecurity = snetParam.PortSecurity
				}

				// Query for all subnets that match filters
				snetResults, err := getSubnetsByFilter(is, &sopts)
				if err != nil {
					return nil, err
				}
				for _, snet := range snetResults {
					nets = append(nets, openstackconfigv1.PortOpts{
						NetworkID:    snet.NetworkID,
						NameSuffix:   snet.ID,
						FixedIPs:     []openstackconfigv1.FixedIPs{{SubnetID: snet.ID}},
						Tags:         append(net.PortTags, snetParam.PortTags...),
						VNICType:     net.VNICType,
						PortSecurity: portSecurity,
					})
				}
			}
		}
	}

	clusterInfra, err := configClient.Infrastructures().Get(context.TODO(), "cluster", metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve cluster Infrastructure object: %v", err)
	}

	allowedAddressPairs := []openstackconfigv1.AddressPair{}
	if clusterInfra != nil && clusterInfra.Status.PlatformStatus != nil && clusterInfra.Status.PlatformStatus.OpenStack != nil {
		clusterVips := []string{
			clusterInfra.Status.PlatformStatus.OpenStack.APIServerInternalIP,
			clusterInfra.Status.PlatformStatus.OpenStack.NodeDNSIP,
			clusterInfra.Status.PlatformStatus.OpenStack.IngressIP,
		}

		for _, vip := range clusterVips {
			if vip != "" {
				allowedAddressPairs = append(allowedAddressPairs, openstackconfigv1.AddressPair{IPAddress: vip})
			}
		}
	}

	userData := base64.StdEncoding.EncodeToString([]byte(cmd))
	var portsList []servers.Network
	for _, portOpt := range nets {
		if portOpt.NetworkID == "" {
			return nil, fmt.Errorf("A network was not found or provided for one of the networks or subnets in this machineset")
		}
		portOpt.SecurityGroups = &securityGroups
		portOpt.AllowedAddressPairs = allowedAddressPairs
		if _, ok := netsWithoutAllowedAddressPairs[portOpt.NetworkID]; ok {
			portOpt.AllowedAddressPairs = []openstackconfigv1.AddressPair{}
		}

		port, err := getOrCreatePort(is, name, portOpt)
		if err != nil {
			return nil, fmt.Errorf("Failed to create port err: %v", err)
		}

		portTags := deduplicateList(append(machineTags, portOpt.Tags...))
		_, err = attributestags.ReplaceAll(is.networkClient, "ports", port.ID, attributestags.ReplaceAllOpts{
			Tags: portTags}).Extract()
		if err != nil {
			return nil, fmt.Errorf("Tagging port for server err: %v", err)
		}
		portsList = append(portsList, servers.Network{
			Port: port.ID,
		})

		if config.Trunk == true {
			allPages, err := trunks.List(is.networkClient, trunks.ListOpts{
				Name:   name,
				PortID: port.ID,
			}).AllPages()
			if err != nil {
				return nil, fmt.Errorf("Searching for existing trunk for server err: %v", err)
			}
			trunkList, err := trunks.ExtractTrunks(allPages)
			if err != nil {
				return nil, fmt.Errorf("Searching for existing trunk for server err: %v", err)
			}
			var trunk trunks.Trunk
			if len(trunkList) == 0 {
				// create trunk with the previous port as parent
				trunkCreateOpts := trunks.CreateOpts{
					Name:   name,
					PortID: port.ID,
				}
				newTrunk, err := trunks.Create(is.networkClient, trunkCreateOpts).Extract()
				if err != nil {
					return nil, fmt.Errorf("Create trunk for server err: %v", err)
				}
				trunk = *newTrunk
			} else {
				trunk = trunkList[0]
			}

			_, err = attributestags.ReplaceAll(is.networkClient, "trunks", trunk.ID, attributestags.ReplaceAllOpts{
				Tags: machineTags}).Extract()
			if err != nil {
				return nil, fmt.Errorf("Tagging trunk for server err: %v", err)
			}
		}
	}

	for _, portCreateOpts := range config.Ports {
		port, err := getOrCreatePort(is, name, portCreateOpts)
		if err != nil {
			return nil, err
		}

		portTags := deduplicateList(append(machineTags, portCreateOpts.Tags...))
		_, err = attributestags.ReplaceAll(is.networkClient, "ports", port.ID, attributestags.ReplaceAllOpts{
			Tags: portTags}).Extract()
		if err != nil {
			return nil, fmt.Errorf("Tagging port for server err: %v", err)
		}

		portsList = append(portsList, servers.Network{
			Port: port.ID,
		})
	}

	if len(portsList) == 0 {
		return nil, fmt.Errorf("At least one network, subnet, or port must be defined as a networking interface. Please review your machineset and try again")
	}

	var serverTags []string
	if clusterSpec.DisableServerTags == false {
		serverTags = machineTags
		// NOTE(flaper87): This is the minimum required version
		// to use tags.
		is.computeClient.Microversion = "2.52"
	}

	var imageID string

	if config.RootVolume == nil {
		imageID, err = imageutils.IDFromName(is.imagesClient, config.Image)
		if err != nil {
			return nil, fmt.Errorf("Create new server err: %v", err)
		}
	}

	flavorID, err := flavorutils.IDFromName(is.computeClient, config.Flavor)
	if err != nil {
		return nil, fmt.Errorf("Create new server err: %v", err)
	}

	var serverCreateOpts servers.CreateOptsBuilder = servers.CreateOpts{
		Name:             name,
		ImageRef:         imageID,
		FlavorRef:        flavorID,
		AvailabilityZone: config.AvailabilityZone,
		Networks:         portsList,
		UserData:         []byte(userData),
		SecurityGroups:   securityGroups,
		Tags:             serverTags,
		Metadata:         config.ServerMetadata,
		ConfigDrive:      config.ConfigDrive,
	}

	// If the root volume Size is not 0, means boot from volume
	if config.RootVolume != nil && config.RootVolume.Size != 0 {
		var blocks []bootfromvolume.BlockDevice

		volumeID := config.RootVolume.SourceUUID

		// change serverCreateOpts to exclude imageRef from them
		serverCreateOpts = servers.CreateOpts{
			Name:             name,
			FlavorRef:        flavorID,
			AvailabilityZone: config.AvailabilityZone,
			Networks:         portsList,
			UserData:         []byte(userData),
			SecurityGroups:   securityGroups,
			Tags:             serverTags,
			Metadata:         config.ServerMetadata,
			ConfigDrive:      config.ConfigDrive,
		}

		if bootfromvolume.SourceType(config.RootVolume.SourceType) == bootfromvolume.SourceImage {
			// if source type is "image" then we have to create a volume from the image first
			klog.Infof("Creating a bootable volume from image %v.", config.RootVolume.SourceUUID)

			imageID, err := imageutils.IDFromName(is.imagesClient, config.RootVolume.SourceUUID)
			if err != nil {
				return nil, fmt.Errorf("Create new server err: %v", err)
			}

			// Create a volume first
			volumeCreateOpts := volumes.CreateOpts{
				Size:       config.RootVolume.Size,
				VolumeType: config.RootVolume.VolumeType,
				ImageID:    imageID,
				// The same name as the instance
				Name:             name,
				AvailabilityZone: config.RootVolume.Zone,
			}

			volume, err := volumes.Create(is.volumeClient, volumeCreateOpts).Extract()
			if err != nil {
				return nil, fmt.Errorf("Create bootable volume err: %v", err)
			}

			volumeID = volume.ID

			defer func() {
				// If the server is created in Nova, the lifetime of the associated volume will be
				// bound to it. This deferred cleanup function tackles the case where a volume was created
				// but never attached to a server. In that case, the volume must be removed manually.
				if server == nil {
					if err := volumes.Delete(is.volumeClient, volumeID, nil).ExtractErr(); err != nil {
						klog.Infof("Failed to delete stale volume %q", volumeID)
					} else {
						klog.Infof("Deleted stale volume %q", volumeID)
					}
				}
			}()

			if err := volumes.WaitForStatus(is.volumeClient, volumeID, "available", 300); err != nil {
				return nil, fmt.Errorf("bootable volume is not available: %v", err)
			}

			klog.Infof("Bootable volume %v was created successfully.", volumeID)
		}

		block := bootfromvolume.BlockDevice{
			SourceType:          bootfromvolume.SourceVolume,
			BootIndex:           0,
			UUID:                volumeID,
			DeleteOnTermination: true,
			DestinationType:     bootfromvolume.DestinationVolume,
		}
		blocks = append(blocks, block)

		serverCreateOpts = bootfromvolume.CreateOptsExt{
			CreateOptsBuilder: serverCreateOpts,
			BlockDevice:       blocks,
		}

	}

	// The Machine spec accepts both a server group ID and a server group
	// name. If both are present, assert that they are consistent, else
	// fail. If only the name is present, create the server group.
	//
	// This block validates or populates config.ServerGroupID.
	if config.ServerGroupName != "" {
		existingServerGroups, err := getServerGroupsByName(is.computeClient, config.ServerGroupName)
		if err != nil {
			return nil, fmt.Errorf("retrieving existing server groups: %v", err)
		}

		if config.ServerGroupID == "" {
			switch len(existingServerGroups) {
			case 0:
				sg, err := createServerGroup(is.computeClient, config.ServerGroupName)
				if err != nil {
					return nil, fmt.Errorf("creating the server group: %v", err)
				}
				config.ServerGroupID = sg.ID
			case 1:
				config.ServerGroupID = existingServerGroups[0].ID
			default:
				return nil, fmt.Errorf("multiple server groups found with the same ServerGroupName")
			}
		} else {
			switch len(existingServerGroups) {
			case 0:
				return nil, fmt.Errorf("incompatible ServerGroupID and ServerGroupName")
			default:
				var found bool
				for _, existingServerGroup := range existingServerGroups {
					if existingServerGroup.ID == config.ServerGroupID {
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("incompatible ServerGroupID and ServerGroupName")
				}
			}
		}
	}

	// If the spec sets a server group, then add scheduler hint
	if config.ServerGroupID != "" {
		serverCreateOpts = schedulerhints.CreateOptsExt{
			CreateOptsBuilder: serverCreateOpts,
			SchedulerHints: schedulerhints.SchedulerHints{
				Group: config.ServerGroupID,
			},
		}
	}

	server, err = servers.Create(is.computeClient, keypairs.CreateOptsExt{
		CreateOptsBuilder: serverCreateOpts,
		KeyName:           keyName,
	}).Extract()
	if err != nil {
		// Gophercloud does not guarantee "server" to be nil when "err"
		// is not nil. Resetting to nil to trigger cleanup functions.
		server = nil
		return nil, fmt.Errorf("Create new server err: %v", err)
	}

	is.computeClient.Microversion = ""
	return serverToInstance(server), nil
}

func createServerGroup(computeClient *gophercloud.ServiceClient, name string) (*servergroups.ServerGroup, error) {
	// Microversion "2.15" is the first that supports "soft"-anti-affinity.
	// Microversions starting from "2.64" accept policies as a string
	// instead of an array.
	defer func(microversion string) {
		computeClient.Microversion = microversion
	}(computeClient.Microversion)
	computeClient.Microversion = "2.15"

	return servergroups.Create(computeClient, &servergroups.CreateOpts{
		Name:     name,
		Policies: []string{"soft-anti-affinity"},
	}).Extract()
}

// deduplicateList removes all duplicate entries from a slice of strings in place
func deduplicateList(list []string) []string {
	m := map[string]bool{}
	for _, element := range list {
		if _, ok := m[element]; !ok {
			m[element] = true
		}
	}

	dedupedList := make([]string, len(m))
	i := 0
	for k := range m {
		dedupedList[i] = k
		i++
	}
	return dedupedList
}

func getServerGroupsByName(computeClient *gophercloud.ServiceClient, name string) ([]servergroups.ServerGroup, error) {
	pages, err := servergroups.List(computeClient).AllPages()
	if err != nil {
		return nil, err
	}

	allServerGroups, err := servergroups.ExtractServerGroups(pages)
	if err != nil {
		return nil, err
	}

	serverGroups := make([]servergroups.ServerGroup, 0, len(allServerGroups))
	for _, serverGroup := range allServerGroups {
		if serverGroup.Name == name {
			serverGroups = append(serverGroups, serverGroup)
		}
	}

	return serverGroups, nil
}

func (is *InstanceService) deleteInstancePorts(id string) error {
	// get instance port id
	allInterfaces, err := attachinterfaces.List(is.computeClient, id).AllPages()
	if err != nil {
		return err
	}
	instanceInterfaces, err := attachinterfaces.ExtractInterfaces(allInterfaces)
	if err != nil {
		return err
	}
	if len(instanceInterfaces) < 1 {
		return servers.Delete(is.computeClient, id).ExtractErr()
	}

	trunkSupport, err := GetTrunkSupport(is)
	if err != nil {
		return fmt.Errorf("Obtaining network extensions err: %v", err)
	}
	// get and delete trunks
	for _, port := range instanceInterfaces {
		err := attachinterfaces.Delete(is.computeClient, id, port.PortID).ExtractErr()
		if err != nil {
			return err
		}
		if trunkSupport {
			listOpts := trunks.ListOpts{
				PortID: port.PortID,
			}
			allTrunks, err := trunks.List(is.networkClient, listOpts).AllPages()
			if err != nil {
				return err
			}
			trunkInfo, err := trunks.ExtractTrunks(allTrunks)
			if err != nil {
				return err
			}
			if len(trunkInfo) == 1 {
				err = util.PollImmediate(RetryIntervalTrunkDelete, TimeoutTrunkDelete, func() (bool, error) {
					err := trunks.Delete(is.networkClient, trunkInfo[0].ID).ExtractErr()
					if err != nil {
						return false, nil
					}
					return true, nil
				})
				if err != nil {
					return fmt.Errorf("Error deleting the trunk %v", trunkInfo[0].ID)
				}
			}
		}

		// delete port
		err = util.PollImmediate(RetryIntervalPortDelete, TimeoutPortDelete, func() (bool, error) {
			err := ports.Delete(is.networkClient, port.PortID).ExtractErr()
			if err != nil {
				return false, nil
			}
			return true, nil
		})
		if err != nil {
			return fmt.Errorf("Error deleting the port %v", port.PortID)
		}
	}

	return nil
}

func (is *InstanceService) InstanceDelete(id string) error {
	err := is.deleteInstancePorts(id)
	if err != nil {
		klog.Warningf("Couldn't delete all instance %v ports: %v", id, err)
	}

	// delete instance
	return servers.Delete(is.computeClient, id).ExtractErr()
}

func GetTrunkSupport(is *InstanceService) (bool, error) {
	allPages, err := netext.List(is.networkClient).AllPages()
	if err != nil {
		return false, err
	}

	allExts, err := extensions.ExtractExtensions(allPages)
	if err != nil {
		return false, err
	}

	for _, ext := range allExts {
		if ext.Alias == "trunk" {
			return true, nil
		}
	}
	return false, nil
}

func (is *InstanceService) GetInstanceList(opts *InstanceListOpts) ([]*Instance, error) {
	var listOpts servers.ListOpts
	if opts != nil {
		listOpts = servers.ListOpts{
			// Name is a regular expression, so we need to explicitly specify a
			// whole string match. https://bugzilla.redhat.com/show_bug.cgi?id=1747270
			Name: fmt.Sprintf("^%s$", opts.Name),
		}
	} else {
		listOpts = servers.ListOpts{}
	}

	allPages, err := servers.List(is.computeClient, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("Get service list err: %v", err)
	}
	serverList, err := servers.ExtractServers(allPages)
	if err != nil {
		return nil, fmt.Errorf("Extract services list err: %v", err)
	}
	var instanceList []*Instance
	for _, server := range serverList {
		instanceList = append(instanceList, serverToInstance(&server))
	}
	return instanceList, nil
}

// DoesFlavorExist returns nil if exactly one flavor exists with the given name.
func (is *InstanceService) DoesFlavorExist(flavorName string) error {
	_, err := flavorutils.IDFromName(is.computeClient, flavorName)
	return err
}

// DoesImageExist returns nil if exactly one image exists with the given name.
func (is *InstanceService) DoesImageExist(imageName string) error {
	_, err := imageutils.IDFromName(is.imagesClient, imageName)
	return err
}

// DoesAvailabilityZoneExist return an error if AZ with the given name doesn't exist, and nil otherwise
func (is *InstanceService) DoesAvailabilityZoneExist(azName string) error {
	if azName == "" {
		return nil
	}
	zones, err := azutils.ListAvailableAvailabilityZones(is.computeClient)
	if err != nil {
		return err
	}
	if len(zones) == 0 {
		return fmt.Errorf("could not find an available compute availability zone")
	}
	for _, zoneName := range zones {
		if zoneName == azName {
			return nil
		}
	}
	return fmt.Errorf("could not find compute availability zone: %s", azName)
}

func (is *InstanceService) GetInstance(resourceId string) (instance *Instance, err error) {
	if resourceId == "" {
		return nil, fmt.Errorf("ResourceId should be specified to  get detail.")
	}
	server, err := servers.Get(is.computeClient, resourceId).Extract()
	if err != nil {
		return nil, fmt.Errorf("Get server %q detail failed: %v", resourceId, err)
	}
	return serverToInstance(server), err
}

// SetMachineLabels set labels describing the machine
func (is *InstanceService) SetMachineLabels(machine *machinev1.Machine, instanceID string) error {
	if machine.Labels[MachineRegionLabelName] != "" && machine.Labels[MachineAZLabelName] != "" && machine.Labels[MachineInstanceTypeLabelName] != "" {
		return nil
	}

	var sm serverMetadata
	err := servers.Get(is.computeClient, instanceID).ExtractInto(&sm)
	if err != nil {
		return err
	}

	if machine.Labels == nil {
		machine.Labels = make(map[string]string)
	}

	// Set the region
	machine.Labels[MachineRegionLabelName] = is.regionName

	// Set the availability zone
	machine.Labels[MachineAZLabelName] = sm.AZ

	// Set the flavor name
	flavor, err := flavors.Get(is.computeClient, sm.Flavor["id"].(string)).Extract()
	if err != nil {
		return err
	}
	machine.Labels[MachineInstanceTypeLabelName] = flavor.Name

	return nil
}

func (is *InstanceService) GetFlavorInfo(flavorID string) (flavor *flavors.Flavor, err error) {

	info, err := flavors.Get(is.computeClient, flavorID).Extract()
	if err != nil {
		return nil, fmt.Errorf("Could not find information for flavor id %s", flavorID)
	}
	return info, nil
}

func (is *InstanceService) GetFlavorID(flavorName string) (string, error) {
	return flavorutils.IDFromName(is.computeClient, flavorName)
}

func serverToInstance(server *servers.Server) *Instance {
	return &Instance{*server}
}
