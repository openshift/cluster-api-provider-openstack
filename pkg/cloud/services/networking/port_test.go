/*
Copyright 2021 The Kubernetes Authors.

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

package networking

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsbinding"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients/mock"
)

func Test_GetOrCreatePort(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Arbitrary GUIDs used in the tests
	netID := "7fd24ceb-788a-441f-ad0a-d8e2f5d31a1d"
	subnetID1 := "d9c88a6d-0b8c-48ff-8f0e-8d85a078c194"
	subnetID2 := "d9c2346d-05gc-48er-9ut4-ig83ayt8c7h4"
	portID1 := "50214c48-c09e-4a54-914f-97b40fd22802"
	portID2 := "4c096384-f0a5-466d-9534-06a7ed281a79"
	hostID := "825c1b11-3dca-4bfe-a2d8-a3cc1964c8d5"
	trunkID := "eb7541fa-5e2a-4cca-b2c3-dfa409b917ce"
	portSecurityGroupID := "f51d1206-fc5a-4f7a-a5c0-2e03e44e4dc0"

	// Other arbitrary variables passed in to the tests
	instanceSecurityGroups := []string{"instance-secgroup"}
	securityGroupUUIDs := []string{portSecurityGroupID}
	portSecurityGroupFilters := []infrav1.SecurityGroupFilter{{ID: portSecurityGroupID, Name: "port-secgroup"}}
	valueSpecs := map[string]string{"key": "value"}

	pointerToTrue := pointerTo(true)
	pointerToFalse := pointerTo(false)

	tests := []struct {
		name                   string
		portName               string
		port                   infrav1.PortOpts
		instanceSecurityGroups []string
		tags                   []string
		expect                 func(m *mock.MockNetworkClientMockRecorder)
		// Note the 'wanted' port isn't so important, since it will be whatever we tell ListPort or CreatePort to return.
		// Mostly in this test suite, we're checking that ListPort/CreatePort is called with the expected port opts.
		want    *ports.Port
		wantErr bool
	}{
		{
			"gets and returns existing port if name matches",
			"foo-port-1",
			infrav1.PortOpts{
				Network: &infrav1.NetworkFilter{
					ID: netID,
				},
			},
			nil,
			[]string{},
			func(m *mock.MockNetworkClientMockRecorder) {
				m.
					ListPort(ports.ListOpts{
						Name:      "foo-port-1",
						NetworkID: netID,
					}).Return([]ports.Port{{
					ID: portID1,
				}}, nil)
			},
			&ports.Port{
				ID: portID1,
			},
			false,
		},
		{
			"errors if multiple matching ports are found",
			"foo-port-1",
			infrav1.PortOpts{
				Network: &infrav1.NetworkFilter{
					ID: netID,
				},
			},
			nil,
			[]string{},
			func(m *mock.MockNetworkClientMockRecorder) {
				m.
					ListPort(ports.ListOpts{
						Name:      "foo-port-1",
						NetworkID: netID,
					}).Return([]ports.Port{
					{
						ID:        portID1,
						NetworkID: netID,
						Name:      "foo-port-1",
					},
					{
						ID:        portID2,
						NetworkID: netID,
						Name:      "foo-port-2",
					},
				}, nil)
			},
			nil,
			true,
		},
		{
			"creates port with defaults (description and secgroups) if not specified in portOpts",
			"foo-port-1",
			infrav1.PortOpts{
				Network: &infrav1.NetworkFilter{
					ID: netID,
				},
			},
			instanceSecurityGroups,
			[]string{},
			func(m *mock.MockNetworkClientMockRecorder) {
				// No ports found
				m.
					ListPort(ports.ListOpts{
						Name:      "foo-port-1",
						NetworkID: netID,
					}).Return([]ports.Port{}, nil)
				m.
					CreatePort(portsbinding.CreateOptsExt{
						CreateOptsBuilder: ports.CreateOpts{
							Name:                "foo-port-1",
							Description:         "Created by cluster-api-provider-openstack cluster test-cluster",
							SecurityGroups:      &instanceSecurityGroups,
							NetworkID:           netID,
							AllowedAddressPairs: []ports.AddressPair{},
						},
					}).Return(&ports.Port{ID: portID1}, nil)
			},
			&ports.Port{ID: portID1},
			false,
		},
		{
			"creates port with specified portOpts if no matching port exists",
			"foo-port-bar",
			infrav1.PortOpts{
				Network: &infrav1.NetworkFilter{
					ID: netID,
				},
				NameSuffix:   "bar",
				Description:  "this is a test port",
				MACAddress:   "fe:fe:fe:fe:fe:fe",
				AdminStateUp: pointerToTrue,
				FixedIPs: []infrav1.FixedIP{{
					Subnet: &infrav1.SubnetFilter{
						Name: "subnetFoo",
					},
					IPAddress: "192.168.0.50",
				}, {IPAddress: "192.168.1.50"}},
				SecurityGroupFilters: portSecurityGroupFilters,
				AllowedAddressPairs: []infrav1.AddressPair{{
					IPAddress:  "10.10.10.10",
					MACAddress: "f1:f1:f1:f1:f1:f1",
				}},
				HostID:   hostID,
				VNICType: "direct",
				Profile: infrav1.BindingProfile{
					OVSHWOffload: true,
					TrustedVF:    true,
				},
				DisablePortSecurity: pointerToFalse,
				Tags:                []string{"my-port-tag"},
			},
			nil,
			nil,
			func(m *mock.MockNetworkClientMockRecorder) {
				portCreateOpts := ports.CreateOpts{
					NetworkID:    netID,
					Name:         "foo-port-bar",
					Description:  "this is a test port",
					AdminStateUp: pointerToTrue,
					MACAddress:   "fe:fe:fe:fe:fe:fe",
					FixedIPs: []ports.IP{
						{
							SubnetID:  subnetID1,
							IPAddress: "192.168.0.50",
						}, {
							IPAddress: "192.168.1.50",
						},
					},
					SecurityGroups: &securityGroupUUIDs,
					AllowedAddressPairs: []ports.AddressPair{{
						IPAddress:  "10.10.10.10",
						MACAddress: "f1:f1:f1:f1:f1:f1",
					}},
				}
				portsecurityCreateOptsExt := portsecurity.PortCreateOptsExt{
					CreateOptsBuilder:   portCreateOpts,
					PortSecurityEnabled: pointerToTrue,
				}
				portbindingCreateOptsExt := portsbinding.CreateOptsExt{
					// Note for the test matching, the order in which the builders are composed
					// must be the same as in the function we are testing.
					CreateOptsBuilder: portsecurityCreateOptsExt,
					HostID:            hostID,
					VNICType:          "direct",
					Profile: map[string]interface{}{
						"capabilities": []string{"switchdev"},
						"trusted":      true,
					},
				}
				m.
					ListPort(ports.ListOpts{
						Name:      "foo-port-bar",
						NetworkID: netID,
					}).Return([]ports.Port{}, nil)
				m.
					CreatePort(portbindingCreateOptsExt).
					Return(&ports.Port{
						ID: portID1,
					}, nil)
				m.ReplaceAllAttributesTags("ports", portID1, attributestags.ReplaceAllOpts{Tags: []string{"my-port-tag"}}).Return([]string{"my-port-tag"}, nil)
				m.
					ListSubnet(subnets.ListOpts{
						Name:      "subnetFoo",
						NetworkID: netID,
					}).Return([]subnets.Subnet{
					{
						ID:        subnetID1,
						Name:      "subnetFoo",
						NetworkID: netID,
					},
				}, nil)
			},
			&ports.Port{
				ID: portID1,
			},
			false,
		},
		{
			"fails to create port with specified portOpts if subnet query returns more than one subnet",
			"foo-port-bar",
			infrav1.PortOpts{
				Network: &infrav1.NetworkFilter{
					ID: netID,
				},
				NameSuffix:  "foo-port-bar",
				Description: "this is a test port",
				FixedIPs: []infrav1.FixedIP{{
					Subnet: &infrav1.SubnetFilter{
						Tags: "Foo",
					},
					IPAddress: "192.168.0.50",
				}},
			},
			nil,
			nil,
			func(m *mock.MockNetworkClientMockRecorder) {
				m.
					ListPort(ports.ListOpts{
						Name:      "foo-port-bar",
						NetworkID: netID,
					}).Return([]ports.Port{}, nil)
				m.
					ListSubnet(subnets.ListOpts{
						Tags:      "Foo",
						NetworkID: netID,
					}).Return([]subnets.Subnet{
					{
						ID:        subnetID1,
						NetworkID: netID,
						Name:      "subnetFoo",
					},
					{
						ID:        subnetID2,
						NetworkID: netID,
						Name:      "subnetBar",
					},
				}, nil)
			},
			nil,
			true,
		},
		{
			"overrides default (instance) security groups if port security groups are specified",
			"foo-port-1",
			infrav1.PortOpts{
				Network: &infrav1.NetworkFilter{
					ID: netID,
				},
				SecurityGroupFilters: portSecurityGroupFilters,
			},
			instanceSecurityGroups,
			[]string{},
			func(m *mock.MockNetworkClientMockRecorder) {
				// No ports found
				m.
					ListPort(ports.ListOpts{
						Name:      "foo-port-1",
						NetworkID: netID,
					}).Return([]ports.Port{}, nil)
				m.
					CreatePort(portsbinding.CreateOptsExt{
						CreateOptsBuilder: ports.CreateOpts{
							Name:                "foo-port-1",
							Description:         "Created by cluster-api-provider-openstack cluster test-cluster",
							SecurityGroups:      &securityGroupUUIDs,
							NetworkID:           netID,
							AllowedAddressPairs: []ports.AddressPair{},
						},
					},
					).Return(&ports.Port{ID: portID1}, nil)
			},
			&ports.Port{ID: portID1},
			false,
		},
		{
			"creates port with instance tags when port tags aren't specified",
			"foo-port-1",
			infrav1.PortOpts{
				Network: &infrav1.NetworkFilter{
					ID: netID,
				},
			},
			nil,
			[]string{"my-instance-tag"},
			func(m *mock.MockNetworkClientMockRecorder) {
				// No ports found
				m.
					ListPort(ports.ListOpts{
						Name:      "foo-port-1",
						NetworkID: netID,
					}).Return([]ports.Port{}, nil)
				m.CreatePort(portsbinding.CreateOptsExt{
					CreateOptsBuilder: ports.CreateOpts{
						Name:                "foo-port-1",
						Description:         "Created by cluster-api-provider-openstack cluster test-cluster",
						NetworkID:           netID,
						AllowedAddressPairs: []ports.AddressPair{},
					},
				}).Return(&ports.Port{ID: portID1}, nil)
				m.ReplaceAllAttributesTags("ports", portID1, attributestags.ReplaceAllOpts{Tags: []string{"my-instance-tag"}}).Return([]string{"my-instance-tag"}, nil)
			},
			&ports.Port{ID: portID1},
			false,
		},
		{
			"creates port with port specific tags appending to instance tags",
			"foo-port-1",
			infrav1.PortOpts{
				Network: &infrav1.NetworkFilter{
					ID: netID,
				},
				Tags: []string{"my-port-tag"},
			},
			nil,
			[]string{"my-instance-tag"},
			func(m *mock.MockNetworkClientMockRecorder) {
				// No ports found
				m.
					ListPort(ports.ListOpts{
						Name:      "foo-port-1",
						NetworkID: netID,
					}).Return([]ports.Port{}, nil)
				m.CreatePort(portsbinding.CreateOptsExt{
					CreateOptsBuilder: ports.CreateOpts{
						Name:                "foo-port-1",
						Description:         "Created by cluster-api-provider-openstack cluster test-cluster",
						NetworkID:           netID,
						AllowedAddressPairs: []ports.AddressPair{},
					},
				}).Return(&ports.Port{ID: portID1}, nil)
				m.
					ReplaceAllAttributesTags("ports", portID1, attributestags.ReplaceAllOpts{Tags: []string{"my-instance-tag", "my-port-tag"}}).
					Return([]string{"my-instance-tag", "my-port-tag"}, nil)
			},
			&ports.Port{ID: portID1},
			false,
		},
		{
			"creates port and trunk (with tags) if they aren't found",
			"foo-port-1",
			infrav1.PortOpts{
				Network: &infrav1.NetworkFilter{
					ID: netID,
				},
				Trunk: pointerToTrue,
			},
			nil,
			[]string{"my-tag"},
			func(m *mock.MockNetworkClientMockRecorder) {
				// No ports found
				m.
					ListPort(ports.ListOpts{
						Name:      "foo-port-1",
						NetworkID: netID,
					}).Return([]ports.Port{}, nil)
				m.
					CreatePort(portsbinding.CreateOptsExt{
						CreateOptsBuilder: ports.CreateOpts{
							Name:                "foo-port-1",
							Description:         "Created by cluster-api-provider-openstack cluster test-cluster",
							NetworkID:           netID,
							AllowedAddressPairs: []ports.AddressPair{},
						},
					}).Return(&ports.Port{Name: "foo-port-1", ID: portID1}, nil)
				m.
					ListTrunk(trunks.ListOpts{
						Name:   "foo-port-1",
						PortID: portID1,
					}).Return([]trunks.Trunk{}, nil)
				m.
					CreateTrunk(trunks.CreateOpts{
						Name:        "foo-port-1",
						PortID:      portID1,
						Description: "Created by cluster-api-provider-openstack cluster test-cluster",
					}).Return(&trunks.Trunk{ID: trunkID}, nil)

				m.ReplaceAllAttributesTags("ports", portID1, attributestags.ReplaceAllOpts{Tags: []string{"my-tag"}}).Return([]string{"my-tag"}, nil)
				m.ReplaceAllAttributesTags("trunks", trunkID, attributestags.ReplaceAllOpts{Tags: []string{"my-tag"}}).Return([]string{"my-tag"}, nil)
			},
			&ports.Port{Name: "foo-port-1", ID: portID1},
			false,
		},
		{
			"creates port with value_specs",
			"foo-port-1",
			infrav1.PortOpts{
				Network: &infrav1.NetworkFilter{
					ID: netID,
				},
				ValueSpecs: []infrav1.ValueSpec{
					{
						Name:  "Not important",
						Key:   "key",
						Value: "value",
					},
				},
			},
			nil,
			nil,
			func(m *mock.MockNetworkClientMockRecorder) {
				// No ports found
				m.
					ListPort(ports.ListOpts{
						Name:      "foo-port-1",
						NetworkID: netID,
					}).Return([]ports.Port{}, nil)
				m.
					CreatePort(portsbinding.CreateOptsExt{
						CreateOptsBuilder: ports.CreateOpts{
							Name:                "foo-port-1",
							Description:         "Created by cluster-api-provider-openstack cluster test-cluster",
							NetworkID:           netID,
							AllowedAddressPairs: []ports.AddressPair{},
							ValueSpecs:          &valueSpecs,
						},
					}).Return(&ports.Port{ID: portID1}, nil)
			},
			&ports.Port{ID: portID1},
			false,
		},
		{
			"creates port with propagate uplink status",
			"foo-port-1",
			infrav1.PortOpts{
				Network: &infrav1.NetworkFilter{
					ID: netID,
				},
				PropagateUplinkStatus: pointerToTrue,
			},
			instanceSecurityGroups,
			[]string{},
			func(m *mock.MockNetworkClientMockRecorder) {
				// No ports found
				m.
					ListPort(ports.ListOpts{
						Name:      "foo-port-1",
						NetworkID: netID,
					}).Return([]ports.Port{}, nil)
				m.
					CreatePort(portsbinding.CreateOptsExt{
						CreateOptsBuilder: ports.CreateOpts{
							Name:                  "foo-port-1",
							Description:           "Created by cluster-api-provider-openstack cluster test-cluster",
							SecurityGroups:        &instanceSecurityGroups,
							NetworkID:             netID,
							AllowedAddressPairs:   []ports.AddressPair{},
							PropagateUplinkStatus: pointerToTrue,
						},
					}).Return(&ports.Port{ID: portID1, PropagateUplinkStatus: *pointerToTrue}, nil)
			},
			&ports.Port{ID: portID1, PropagateUplinkStatus: *pointerToTrue},
			false,
		},
	}

	eventObject := &infrav1.OpenStackMachine{}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			mockClient := mock.NewMockNetworkClient(mockCtrl)
			tt.expect(mockClient.EXPECT())
			s := Service{
				client: mockClient,
			}
			got, err := s.GetOrCreatePort(
				eventObject,
				"test-cluster",
				tt.portName,
				&tt.port,
				tt.instanceSecurityGroups,
				tt.tags,
			)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func Test_GarbageCollectErrorInstancesPort(t *testing.T) {
	const (
		instanceName = "foo"
		portID1      = "dc6e0ae3-dad6-4240-a9cb-e541916f20d3"
		portID2      = "a38ab1cb-c2cc-4c1b-9d1d-696ec73356d2"
		trunkID1     = "6ecfc5c2-7ee7-41b3-bba7-874b9b1423b7"
	)
	portName1 := GetPortName(instanceName, nil, 0)
	portName2 := GetPortName(instanceName, nil, 1)

	tests := []struct {
		// man is the name of the test.
		name string
		// expect allows definition of any expected calls to the mock.
		expect func(m *mock.MockNetworkClientMockRecorder)
		// portOpts defines the instance ports as defined in the OSM spec.
		portOpts []infrav1.PortOpts
		// trunkSupported indicates whether we should check for trunk ports
		trunkSupported bool
		// wantErr defines whether the test is supposed to fail.
		wantErr bool
	}{
		{
			name: "garbage collects all ports for an instance",
			expect: func(m *mock.MockNetworkClientMockRecorder) {
				o1 := ports.ListOpts{
					Name: portName1,
				}
				p1 := []ports.Port{
					{
						ID:   portID1,
						Name: portName1,
					},
				}
				m.ListPort(o1).Return(p1, nil)
				m.DeletePort(portID1)
				o2 := ports.ListOpts{
					Name: portName2,
				}
				p2 := []ports.Port{
					{
						ID:   portID2,
						Name: portName2,
					},
				}

				m.ListPort(o2).Return(p2, nil)
				m.DeletePort(portID2)
			},
			portOpts: []infrav1.PortOpts{
				{},
				{},
			},
			wantErr: false,
		},
		{
			name: "succeed if there are no ports to be cleaned up",
			expect: func(m *mock.MockNetworkClientMockRecorder) {
				o1 := ports.ListOpts{
					Name: portName1,
				}
				m.ListPort(o1).Return([]ports.Port{}, nil)
			},
			portOpts: []infrav1.PortOpts{
				{},
			},
			wantErr: false,
		},
		{
			name: "cleanup trunk port",
			expect: func(m *mock.MockNetworkClientMockRecorder) {
				o1 := ports.ListOpts{
					Name: portName1,
				}
				p1 := []ports.Port{
					{
						ID:   portID1,
						Name: portName1,
					},
				}
				m.ListPort(o1).Return(p1, nil)
				m.ListTrunk(trunks.ListOpts{
					PortID: portID1,
				}).Return([]trunks.Trunk{
					{
						ID: trunkID1,
					},
				}, nil)
				m.DeleteTrunk(trunkID1)
				m.ListTrunkSubports(trunkID1).Return([]trunks.Subport{}, nil)
				m.DeletePort(portID1)
			},
			portOpts: []infrav1.PortOpts{
				{},
			},
			trunkSupported: true,
			wantErr:        false,
		},
	}

	eventObject := &infrav1.OpenStackMachine{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			g := NewWithT(t)
			mockClient := mock.NewMockNetworkClient(mockCtrl)
			tt.expect(mockClient.EXPECT())
			s := Service{
				client: mockClient,
			}
			err := s.GarbageCollectErrorInstancesPort(
				eventObject,
				instanceName,
				tt.portOpts,
				tt.trunkSupported,
			)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
		})
	}
}

func pointerTo(b bool) *bool {
	return &b
}

// newEventObject returns a minimal runtime.Object for the eventObject argument.
func newEventObject() *corev1.Node {
	return &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "test-node"}}
}

// portWithIP builds a ports.Port fixture with one FixedIP and optional
// pre-existing AllowedAddressPairs.
func portWithIP(id, networkID, ip string, existingPairs ...ports.AddressPair) ports.Port {
	return ports.Port{
		ID:                  id,
		NetworkID:           networkID,
		MACAddress:          "fa:16:3e:00:00:" + id[len(id)-2:],
		FixedIPs:            []ports.IP{{IPAddress: ip, SubnetID: "subnet-" + id}},
		AllowedAddressPairs: existingPairs,
	}
}

// minimalPortOpts returns a PortOpts that will cause GetOrCreatePort to create
// a new port. It sets only the fields required by the function under test.
func minimalPortOpts(networkID string, symmetric bool, allowedPairs ...infrav1.AddressPair) *infrav1.PortOpts {
	return &infrav1.PortOpts{
		Network: &infrav1.NetworkFilter{
			ID: networkID,
		},
		AllowedAddressPairs:          allowedPairs,
		SymmetricAllowedAddressPairs: symmetric,
	}
}

// newService constructs a networking.Service wired to the provided mock client.
func newService(t *testing.T, mc *mock.MockNetworkClient) *Service {
	t.Helper()
	return &Service{client: mc}
}

// TestGetOrCreatePort_SymmetricDisabled checks that when
// SymmetricAllowedAddressPairs=false the VIP discovery ListPort and UpdatePort
// are never called, even when AllowedAddressPairs is populated.
func TestGetOrCreatePort_SymmetricDisabled(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mc := mock.NewMockNetworkClient(mockCtrl)
	svc := newService(t, mc)

	const networkID = "net-1"
	const portName = "test-port"
	created := portWithIP("port-1", networkID, "10.0.0.100")

	// First ListPort: check for pre-existing port → none.
	mc.EXPECT().
		ListPort(ports.ListOpts{Name: portName, NetworkID: networkID}).
		Return([]ports.Port{}, nil)
	// CreatePort: return the new machine port.
	mc.EXPECT().
		CreatePort(gomock.Any()).
		Return(&created, nil)

	// VIP discovery ListPort and UpdatePort must NOT be called.
	// gomock will fail the test automatically if they are called unexpectedly.

	opts := minimalPortOpts(networkID, false, infrav1.AddressPair{IPAddress: "10.0.0.5"})
	port, err := svc.GetOrCreatePort(newEventObject(), "cluster", portName, opts, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port.ID != created.ID {
		t.Errorf("got port ID %q, want %q", port.ID, created.ID)
	}
}

// TestGetOrCreatePort_SymmetricEnabled_UpdatesMatchingVIP is the core happy
// path: symmetric=true, one VIP IP in AllowedAddressPairs, the corresponding
// VIP port exists in the network → one UpdatePort call that adds the machine
// IP to the VIP's AllowedAddressPairs.
func TestGetOrCreatePort_SymmetricEnabled_UpdatesMatchingVIP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mc := mock.NewMockNetworkClient(mockCtrl)
	svc := newService(t, mc)

	const networkID = "net-2"
	const portName = "test-port"
	const machineIP = "10.0.0.100"
	const vipIP = "10.0.0.5"

	machinePort := portWithIP("machine-port", networkID, machineIP)
	vipPort := portWithIP("vip-port", networkID, vipIP)

	// 1. Check for pre-existing port.
	mc.EXPECT().
		ListPort(ports.ListOpts{Name: portName, NetworkID: networkID}).
		Return([]ports.Port{}, nil)

	// 2. Create the machine port.
	mc.EXPECT().
		CreatePort(gomock.Any()).
		Return(&machinePort, nil)

	// 3. VIP discovery: return all ports on the network.
	mc.EXPECT().
		ListPort(ports.ListOpts{NetworkID: networkID}).
		Return([]ports.Port{machinePort, vipPort}, nil)

	// 4. UpdatePort: VIP port must be updated with the machine IP.
	mc.EXPECT().
		UpdatePort("vip-port", gomock.Any()).
		DoAndReturn(func(id string, opts ports.UpdateOpts) (*ports.Port, error) {
			if opts.AllowedAddressPairs == nil {
				t.Errorf("UpdatePort called with nil AllowedAddressPairs")
				return nil, fmt.Errorf("nil pairs")
			}
			found := false
			for _, p := range *opts.AllowedAddressPairs {
				if p.IPAddress == machineIP {
					found = true
				}
			}
			if !found {
				t.Errorf("machine IP %q not found in UpdatePort opts %+v", machineIP, *opts.AllowedAddressPairs)
			}
			updated := vipPort
			updated.AllowedAddressPairs = *opts.AllowedAddressPairs
			return &updated, nil
		})

	opts := minimalPortOpts(networkID, true, infrav1.AddressPair{IPAddress: vipIP})
	port, err := svc.GetOrCreatePort(newEventObject(), "cluster", portName, opts, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port.ID != machinePort.ID {
		t.Errorf("got port ID %q, want %q", port.ID, machinePort.ID)
	}
}

// TestGetOrCreatePort_SymmetricEnabled_NoAddressPairs verifies that when
// symmetric=true but AllowedAddressPairs is empty, the len(addressPairs) > 0
// guard in GetOrCreatePort short-circuits: no VIP discovery ListPort and no
// UpdatePort are called.
func TestGetOrCreatePort_SymmetricEnabled_NoAddressPairs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mc := mock.NewMockNetworkClient(mockCtrl)
	svc := newService(t, mc)

	const networkID = "net-3"
	const portName = "test-port"
	created := portWithIP("port-3", networkID, "10.0.0.6")

	mc.EXPECT().
		ListPort(ports.ListOpts{Name: portName, NetworkID: networkID}).
		Return([]ports.Port{}, nil)
	mc.EXPECT().
		CreatePort(gomock.Any()).
		Return(&created, nil)

	// No second ListPort or UpdatePort expected.

	opts := minimalPortOpts(networkID, true /* no AllowedAddressPairs */)
	_, err := svc.GetOrCreatePort(newEventObject(), "cluster", portName, opts, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestGetOrCreatePort_SymmetricEnabled_VIPNotInNetwork verifies that when no
// port in the network has an IP matching an allowed address pair, UpdatePort
// is never called.
func TestGetOrCreatePort_SymmetricEnabled_VIPNotInNetwork(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mc := mock.NewMockNetworkClient(mockCtrl)
	svc := newService(t, mc)

	const networkID = "net-4"
	const portName = "test-port"
	machinePort := portWithIP("port-4", networkID, "10.0.0.100")

	mc.EXPECT().
		ListPort(ports.ListOpts{Name: portName, NetworkID: networkID}).
		Return([]ports.Port{}, nil)
	mc.EXPECT().
		CreatePort(gomock.Any()).
		Return(&machinePort, nil)
	// VIP discovery: only the machine port is on this network;
	// no port has IP 10.0.0.5 so no match is found.
	mc.EXPECT().
		ListPort(ports.ListOpts{NetworkID: networkID}).
		Return([]ports.Port{machinePort}, nil)

	// UpdatePort must NOT be called.

	opts := minimalPortOpts(networkID, true, infrav1.AddressPair{IPAddress: "10.0.0.5"})
	_, err := svc.GetOrCreatePort(newEventObject(), "cluster", portName, opts, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestGetOrCreatePort_SymmetricEnabled_MachineIPAlreadyInVIP verifies
// idempotency: updateVipWithAllowedAddressPairs skips UpdatePort when the
// machine IP is already present in the VIP's AllowedAddressPairs.
func TestGetOrCreatePort_SymmetricEnabled_MachineIPAlreadyInVIP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mc := mock.NewMockNetworkClient(mockCtrl)
	svc := newService(t, mc)

	const networkID = "net-5"
	const portName = "test-port"
	const machineIP = "10.0.0.100"
	const vipIP = "10.0.0.5"

	machinePort := portWithIP("machine-port-5", networkID, machineIP)
	// VIP already carries the machine IP — no update should be needed.
	vipPort := portWithIP("vip-port-5", networkID, vipIP,
		ports.AddressPair{IPAddress: machineIP, MACAddress: machinePort.MACAddress},
	)

	mc.EXPECT().
		ListPort(ports.ListOpts{Name: portName, NetworkID: networkID}).
		Return([]ports.Port{}, nil)
	mc.EXPECT().
		CreatePort(gomock.Any()).
		Return(&machinePort, nil)
	mc.EXPECT().
		ListPort(ports.ListOpts{NetworkID: networkID}).
		Return([]ports.Port{machinePort, vipPort}, nil)

	// Machine IP already present → UpdatePort must NOT be called.

	opts := minimalPortOpts(networkID, true, infrav1.AddressPair{IPAddress: vipIP})
	_, err := svc.GetOrCreatePort(newEventObject(), "cluster", portName, opts, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestGetOrCreatePort_SymmetricEnabled_PreservesExistingVIPPairs verifies
// that pre-existing pairs on the VIP port are kept when the machine IP is
// appended.
func TestGetOrCreatePort_SymmetricEnabled_PreservesExistingVIPPairs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mc := mock.NewMockNetworkClient(mockCtrl)
	svc := newService(t, mc)

	const networkID = "net-6"
	const portName = "test-port"
	const machineIP = "10.0.0.100"
	const vipIP = "10.0.0.5"
	const existingPairIP = "192.168.1.50"

	machinePort := portWithIP("machine-port-6", networkID, machineIP)
	vipPort := portWithIP("vip-port-6", networkID, vipIP,
		ports.AddressPair{IPAddress: existingPairIP},
	)

	mc.EXPECT().
		ListPort(ports.ListOpts{Name: portName, NetworkID: networkID}).
		Return([]ports.Port{}, nil)
	mc.EXPECT().
		CreatePort(gomock.Any()).
		Return(&machinePort, nil)
	mc.EXPECT().
		ListPort(ports.ListOpts{NetworkID: networkID}).
		Return([]ports.Port{machinePort, vipPort}, nil)

	mc.EXPECT().
		UpdatePort("vip-port-6", gomock.Any()).
		DoAndReturn(func(id string, opts ports.UpdateOpts) (*ports.Port, error) {
			pairIPs := make(map[string]bool)
			for _, p := range *opts.AllowedAddressPairs {
				pairIPs[p.IPAddress] = true
			}
			if !pairIPs[existingPairIP] {
				t.Errorf("pre-existing pair IP %q was dropped from UpdatePort opts", existingPairIP)
			}
			if !pairIPs[machineIP] {
				t.Errorf("machine IP %q was not added to UpdatePort opts", machineIP)
			}
			updated := vipPort
			updated.AllowedAddressPairs = *opts.AllowedAddressPairs
			return &updated, nil
		})

	opts := minimalPortOpts(networkID, true, infrav1.AddressPair{IPAddress: vipIP})
	_, err := svc.GetOrCreatePort(newEventObject(), "cluster", portName, opts, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestGetOrCreatePort_SymmetricEnabled_MultipleVIPs verifies that when
// AllowedAddressPairs lists multiple VIP IPs, each matching port on the
// network gets its own UpdatePort call.
func TestGetOrCreatePort_SymmetricEnabled_MultipleVIPs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mc := mock.NewMockNetworkClient(mockCtrl)
	svc := newService(t, mc)

	const networkID = "net-7"
	const portName = "test-port"
	const machineIP = "10.0.0.100"

	machinePort := portWithIP("machine-port-7", networkID, machineIP)
	vipPort1 := portWithIP("vip-port-7a", networkID, "10.0.0.5")
	vipPort2 := portWithIP("vip-port-7b", networkID, "10.0.0.7")

	mc.EXPECT().
		ListPort(ports.ListOpts{Name: portName, NetworkID: networkID}).
		Return([]ports.Port{}, nil)
	mc.EXPECT().
		CreatePort(gomock.Any()).
		Return(&machinePort, nil)
	mc.EXPECT().
		ListPort(ports.ListOpts{NetworkID: networkID}).
		Return([]ports.Port{machinePort, vipPort1, vipPort2}, nil)

	// Both VIP ports must receive exactly one UpdatePort call each.
	mc.EXPECT().UpdatePort("vip-port-7a", gomock.Any()).Return(&vipPort1, nil)
	mc.EXPECT().UpdatePort("vip-port-7b", gomock.Any()).Return(&vipPort2, nil)

	opts := minimalPortOpts(networkID, true,
		infrav1.AddressPair{IPAddress: "10.0.0.5"},
		infrav1.AddressPair{IPAddress: "10.0.0.7"},
	)
	_, err := svc.GetOrCreatePort(newEventObject(), "cluster", portName, opts, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestGetOrCreatePort_SymmetricEnabled_VIPDiscoveryListPortError verifies
// that when the VIP discovery ListPort call fails, the error is propagated
// and UpdatePort is never attempted.
func TestGetOrCreatePort_SymmetricEnabled_VIPDiscoveryListPortError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mc := mock.NewMockNetworkClient(mockCtrl)
	svc := newService(t, mc)

	const networkID = "net-8"
	const portName = "test-port"
	created := portWithIP("port-8", networkID, "10.0.0.11")

	mc.EXPECT().
		ListPort(ports.ListOpts{Name: portName, NetworkID: networkID}).
		Return([]ports.Port{}, nil)
	mc.EXPECT().
		CreatePort(gomock.Any()).
		Return(&created, nil)
	mc.EXPECT().
		ListPort(ports.ListOpts{NetworkID: networkID}).
		Return(nil, fmt.Errorf("neutron: internal server error"))

	opts := minimalPortOpts(networkID, true, infrav1.AddressPair{IPAddress: "10.0.0.5"})
	_, err := svc.GetOrCreatePort(newEventObject(), "cluster", portName, opts, nil, nil)
	if err == nil {
		t.Fatal("expected an error from VIP discovery ListPort failure, got nil")
	}
}

// TestGetOrCreatePort_SymmetricEnabled_UpdatePortError verifies that a
// Neutron error from UpdatePort is propagated back through GetOrCreatePort.
func TestGetOrCreatePort_SymmetricEnabled_UpdatePortError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mc := mock.NewMockNetworkClient(mockCtrl)
	svc := newService(t, mc)

	const networkID = "net-9"
	const portName = "test-port"
	const machineIP = "10.0.0.100"
	const vipIP = "10.0.0.5"

	machinePort := portWithIP("machine-port-9", networkID, machineIP)
	vipPort := portWithIP("vip-port-9", networkID, vipIP)

	mc.EXPECT().
		ListPort(ports.ListOpts{Name: portName, NetworkID: networkID}).
		Return([]ports.Port{}, nil)
	mc.EXPECT().
		CreatePort(gomock.Any()).
		Return(&machinePort, nil)
	mc.EXPECT().
		ListPort(ports.ListOpts{NetworkID: networkID}).
		Return([]ports.Port{machinePort, vipPort}, nil)
	mc.EXPECT().
		UpdatePort("vip-port-9", gomock.Any()).
		Return(nil, fmt.Errorf("neutron: quota exceeded"))

	opts := minimalPortOpts(networkID, true, infrav1.AddressPair{IPAddress: vipIP})
	_, err := svc.GetOrCreatePort(newEventObject(), "cluster", portName, opts, nil, nil)
	if err == nil {
		t.Fatal("expected an error from UpdatePort failure, got nil")
	}
}

// TestGetOrCreatePort_ReturnsExistingPort verifies the short-circuit path:
// when ListPort finds exactly one pre-existing port it is returned directly
// with no CreatePort, VIP discovery, or UpdatePort calls.
func TestGetOrCreatePort_ReturnsExistingPort(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mc := mock.NewMockNetworkClient(mockCtrl)
	svc := newService(t, mc)

	const networkID = "net-11"
	const portName = "test-port"
	existing := portWithIP("existing-port", networkID, "10.0.0.20")

	mc.EXPECT().
		ListPort(ports.ListOpts{Name: portName, NetworkID: networkID}).
		Return([]ports.Port{existing}, nil)

	// None of these must be called.

	opts := minimalPortOpts(networkID, true, infrav1.AddressPair{IPAddress: "10.0.0.5"})
	port, err := svc.GetOrCreatePort(newEventObject(), "cluster", portName, opts, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port.ID != existing.ID {
		t.Errorf("got port ID %q, want %q", port.ID, existing.ID)
	}
}
