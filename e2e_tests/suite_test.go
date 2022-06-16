package e2e_tests

import (
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
	"testing"
	"time"
)

func TestE2EInfobloxGoClient(t *testing.T) {
	RegisterFailHandler(Fail)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.SlowSpecThreshold = time.Minute * 10
	RunSpecs(t, "InfobloxGoClient E2E Test Suite", reporterConfig)
}

// ConnectorFacadeE2E is an end-to-end test facade for the ibclient.Connector.
// Its purpose is to delete objects created by test, when the test is done.
type ConnectorFacadeE2E struct {
	ibclient.Connector
	deleteSet []string
}

func (c *ConnectorFacadeE2E) CreateObject(obj ibclient.IBObject) (ref string, err error) {
	ref, err = c.Connector.CreateObject(obj)
	if err == nil {
		c.addDeleteRef(ref)
	}
	return ref, err
}

func (c *ConnectorFacadeE2E) DeleteObject(ref string) (refRes string, err error) {
	refRes, err = c.Connector.DeleteObject(ref)
	if err == nil {
		c.removeDeleteRef(ref)
	}

	if _, ok := err.(*ibclient.NotFoundError); ok {
		log.Printf("Object %s not found. Probably was deleted while sweeping.", ref)
		err = nil
		c.removeDeleteRef(ref)
	}

	return refRes, err
}

// SweepObjects should be executed when the test
// is done and all created objects need to be deleted
func (c *ConnectorFacadeE2E) SweepObjects() error {
	for i := len(c.deleteSet) - 1; i >= 0; i-- {
		_, err := c.DeleteObject(c.deleteSet[i])
		if err != nil {
			log.Printf("Failed to sweep test objects. During object %s deletion error %s was raised", c.deleteSet[i], err)
			return err
		}
	}
	return nil
}

func (c *ConnectorFacadeE2E) addDeleteRef(ref string) {
	isInTheSet := false
	for _, r := range c.deleteSet {
		if r == ref {
			isInTheSet = true
		}
	}
	if !isInTheSet {
		c.deleteSet = append(c.deleteSet, ref)
	}
}

func (c *ConnectorFacadeE2E) removeDeleteRef(ref string) {
	for i := range c.deleteSet {
		if c.deleteSet[i] == ref {
			c.deleteSet = append(c.deleteSet[:i], c.deleteSet[i+1:]...)
		}
	}
}
