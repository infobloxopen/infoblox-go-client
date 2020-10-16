package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing AWSRoute53TaskGroupOperations", func() {
	Context("AWSRoute53TaskGroup object", func() {
		name := "task3"
		networkView := "test_netview"
		gridMember :="test.localdomain"
		ra := NewAWSRoute53TaskGroup(AWSRoute53TaskGroup{
			Name: name,
			NetworkView: networkView,
			GridMember: gridMember,
		})

		It("should set fields correctly", func() {
			Expect(ra.Name).To(Equal(name))
			Expect(ra.NetworkView).To(Equal(networkView))
			Expect(ra.GridMember).To(Equal(gridMember))
		})

		It("should set base fields correctly", func() {
			Expect(ra.ObjectType()).To(Equal("awsrte53taskgroup"))
			Expect(ra.ReturnFields()).To(ConsistOf("account_id", "comment", "consolidate_zones", "consolidated_view", "disabled",
				"grid_member", "name", "network_view", "network_view_mapping_policy", "sync_status"))
		})
	})

	Describe("Allocate specific AWS Route53 Task ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"

		aws := AWSRoute53TaskGroup{
			Name: "test_sync",
			NetworkView: "test_netview",
			GridMember: "test.localdomain",
			NetworkViewMappingPolicy: "DIRECT",
		}
		fakeRefReturn := fmt.Sprintf("awsrte53taskgroup/ZG5zLmJpbmRfY25h:%s/%20%20", aws.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewAWSRoute53TaskGroup(aws),
			resultObject: NewAWSRoute53TaskGroup(AWSRoute53TaskGroup{
				Ref: fakeRefReturn,
				NetworkView: aws.NetworkView,
				Name: aws.Name,
				GridMember: aws.GridMember,
				NetworkViewMappingPolicy: "DIRECT",

			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualTask *AWSRoute53TaskGroup
		var err error
		It("should pass expected AWSRoute53Group Task Object to CreateObject", func() {
			actualTask, err = objMgr.CreateAWSRoute53TaskGroup(aws)
		})

		It("should return expected AWSRoute53Group Task Object", func() {
			Expect(actualTask).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get AWSRoute53Group Task by Reference", func() {
		name := "test"
		fakeRefReturn := fmt.Sprintf("awsrte53taskgroup/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aws := AWSRoute53TaskGroup{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewAWSRoute53TaskGroup(aws),
			getObjectRef: fakeRefReturn,
			resultObject: []AWSRoute53TaskGroup{*NewAWSRoute53TaskGroup(AWSRoute53TaskGroup{Name: aws.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualTask *[]AWSRoute53TaskGroup
		var err error
		It("should pass expected AWSRoute53Group Task Object to GetObject", func() {
			actualTask, err = objMgr.GetAWSRoute53TaskGroup(aws)

		})

		It("should return expected AWSRoute53Group Task Object", func() {
			Expect(*actualTask).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete AWSRoute53Group Task by Reference", func() {

		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("awsrte53taskgroup/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aws := AWSRoute53TaskGroup{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			deleteObjectRef: aws.Ref,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected AWSRoute53Group Task Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteAWSRoute53TaskGroup(aws)
		})
		It("should return expected AWSRoute53Group Task Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
