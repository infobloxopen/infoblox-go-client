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
		ra := NewAWSRte53TaskGroup(AWSRte53TaskGroup{
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

		aws := AWSRte53TaskGroup{
			Name: "test_sync",
			NetworkView: "test_netview",
			GridMember: "test.localdomain",
			NetworkViewMappingPolicy: "DIRECT",
		}
		fakeRefReturn := fmt.Sprintf("awsrte53taskgroup/ZG5zLmJpbmRfY25h:%s/%20%20", aws.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewAWSRte53TaskGroup(aws),
			resultObject: NewAWSRte53TaskGroup(AWSRte53TaskGroup{
				Ref: fakeRefReturn,
				NetworkView: aws.NetworkView,
				Name: aws.Name,
				GridMember: aws.GridMember,
				NetworkViewMappingPolicy: "DIRECT",

			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualTask *AWSRte53TaskGroup
		var err error
		It("should pass expected AWSRte53Group Task Object to CreateObject", func() {
			actualTask, err = objMgr.CreateAWSRte53TaskGroup(aws)
		})

		It("should return expected AWSRte53Group Task Object", func() {
			Expect(actualTask).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get AWSRte53Group Task by Reference", func() {
		name := "test"
		fakeRefReturn := fmt.Sprintf("awsrte53taskgroup/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aws := AWSRte53TaskGroup{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewAWSRte53TaskGroup(aws),
			getObjectRef: fakeRefReturn,
			resultObject: []AWSRte53TaskGroup{*NewAWSRte53TaskGroup(AWSRte53TaskGroup{Name: aws.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualTask *[]AWSRte53TaskGroup
		var err error
		It("should pass expected AWSRte53Group Task Object to GetObject", func() {
			actualTask, err = objMgr.GetAWSRte53TaskGroup(aws)

		})

		It("should return expected AWSRte53Group Task Object", func() {
			Expect(*actualTask).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete AWSRte53Group Task by Reference", func() {

		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("awsrte53taskgroup/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aws := AWSRte53TaskGroup{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			deleteObjectRef: aws.Ref,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected AWSRte53Group Task Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteAWSRte53TaskGroup(aws)
		})
		It("should return expected AWSRte53Group Task Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
