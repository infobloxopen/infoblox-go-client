package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing VDiscoveryTaskOperations", func() {
	Context("VDiscoveryTask object", func() {
		name := "task3"
		fqdnOrIp := "test.amazonaws.com"
		memberV := "infoblox.localdomain"
		port := uint(443)
		protocol := "HTTPS"
		rv := NewVDiscoveryTask(VDiscoveryTask{
			Name: name, FqdnOrIp: fqdnOrIp,
			MemberV:  memberV,
			Port:     port,
			Protocol: protocol,
		})

		It("should set fields correctly", func() {
			Expect(rv.Name).To(Equal(name))
			Expect(rv.FqdnOrIp).To(Equal(fqdnOrIp))
			Expect(rv.MemberV).To(Equal(memberV))
			Expect(rv.Port).To(Equal(port))
			Expect(rv.Protocol).To(Equal(protocol))
		})

		It("should set base fields correctly", func() {
			Expect(rv.ObjectType()).To(Equal("vdiscoverytask"))
			Expect(rv.ReturnFields()).To(ConsistOf("name", "driver_type", "fqdn_or_ip", "username", "member",
				"port", "protocol", "auto_consolidate_cloud_ea", "auto_consolidate_managed_tenant",
				"auto_consolidate_managed_vm", "merge_data", "private_network_view", "private_network_view_mapping_policy",
				"public_network_view", "public_network_view_mapping_policy", "allow_unsecured_connection",
				"auto_create_dns_hostname_template", "auto_create_dns_record", "auto_create_dns_record_type",
				"comment", "credentials_type", "dns_view_private_ip", "dns_view_public_ip", "domain_name",
				"enabled", "identity_version", "state", "state_msg",
				"update_dns_view_private_ip", "update_dns_view_public_ip", "update_metadata"))
		})
	})

	Describe("Allocate specific VDiscoveryTask ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"

		vDis := VDiscoveryTask{
			Name: "testTask", DriverType: "AWS", FqdnOrIp: "ec2.eu-west1.amazonaws.com",
			MemberV:                         "test.localdomain",
			Port:                            443,
			Protocol:                        "HTTPS",
			AutoConsolidateCloudEa:          true,
			AutoConsolidateManagedTenant:    true,
			AutoConsolidateManagedVm:        true,
			MergeData:                       true,
			PrivateNetworkView:              "default",
			PrivateNetworkViewMappingPolicy: "DIRECT",
			PublicNetworkView:               "default",
			PublicNetworkViewMappingPolicy:  "DIRECT",
			UpdateMetadata:                  true,
			Username:                        "test",
			Password:                        "test",
		}
		fakeRefReturn := fmt.Sprintf("vdiscoverytask/ZG5zLmJpbmRfY25h:%s/%20%20", vDis.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewVDiscoveryTask(vDis),
			resultObject: NewVDiscoveryTask(VDiscoveryTask{
				Ref: fakeRefReturn,
				Name: vDis.Name, DriverType: vDis.DriverType, FqdnOrIp: vDis.FqdnOrIp,
				MemberV:                         vDis.MemberV,
				Port:                            vDis.Port,
				Protocol:                        vDis.Protocol,
				AutoConsolidateCloudEa:          vDis.AutoConsolidateCloudEa,
				AutoConsolidateManagedTenant:    vDis.AutoConsolidateManagedTenant,
				AutoConsolidateManagedVm:        vDis.AutoConsolidateManagedVm,
				MergeData:                       vDis.MergeData,
				PrivateNetworkView:              vDis.PrivateNetworkView,
				PrivateNetworkViewMappingPolicy: vDis.PublicNetworkViewMappingPolicy,
				PublicNetworkView:               vDis.PublicNetworkView,
				PublicNetworkViewMappingPolicy:  vDis.PublicNetworkViewMappingPolicy,
				UpdateMetadata:                  vDis.UpdateMetadata,
				Username:                        vDis.Username,
				Password:                        vDis.Password,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualTask *VDiscoveryTask
		var err error
		It("should pass expected vDiscovery Task Object to CreateObject", func() {
			actualTask, err = objMgr.CreateVDiscoveryTask(vDis)
		})

		It("should return expected vDiscovery Task Object", func() {
			Expect(actualTask).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get vDiscovery Task by Reference", func() {
		name := "test"
		fakeRefReturn := fmt.Sprintf("vdiscoverytask/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		vDis := VDiscoveryTask{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewVDiscoveryTask(vDis),
			getObjectRef: fakeRefReturn,
			resultObject: []VDiscoveryTask{*NewVDiscoveryTask(VDiscoveryTask{Name: vDis.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualTask *[]VDiscoveryTask
		var err error
		It("should pass expected vDiscovery Task Object to GetObject", func() {
			actualTask, err = objMgr.GetVDiscoveryTask(vDis)

		})

		It("should return expected vDiscovery Task Object", func() {
			Expect(*actualTask).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete vDiscovery Task by Reference", func() {

		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("vdiscoverytask/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		vDis := VDiscoveryTask{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			deleteObjectRef: vDis.Ref,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected vDiscovery Task Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteVDiscoveryTask(vDis)
		})
		It("should return expected vDiscovery Task Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
