package e2e_tests

import (
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestE2EInfobloxGoClient(t *testing.T) {
	RegisterFailHandler(Fail)
	config.DefaultReporterConfig.SlowSpecThreshold = (time.Minute * 10).Seconds()
	RunSpecs(t, "InfobloxGoClient E2E Test Suite")
}

// ConnectorFacadeE2E is an end-to-end test facade for the ibclient.Connector.
// Its purpose is to delete objects created by test, when the test is done.
type ConnectorFacadeE2E struct {
	ibclient.Connector
	deleteSet map[string]struct{}
}

func (c *ConnectorFacadeE2E) CreateObject(obj ibclient.IBObject) (ref string, err error) {
	ref, err = c.Connector.CreateObject(obj)
	if err == nil {
		c.deleteSet[ref] = struct{}{}
	}
	return ref, err
}

func (c *ConnectorFacadeE2E) DeleteObject(ref string) (refRes string, err error) {
	refRes, err = c.Connector.DeleteObject(ref)
	if err == nil {
		delete(c.deleteSet, ref)
	}

	if _, ok := err.(*ibclient.NotFoundError); ok {
		log.Printf("Object %s not found. Probably was deleted while sweeping.", ref)
		err = nil
		delete(c.deleteSet, ref)
	}

	return refRes, err
}

// SweepObjects should be executed when the test
// is done and all created objects need to be deleted
func (c *ConnectorFacadeE2E) SweepObjects() error {
	for ref, _ := range c.deleteSet {
		_, err := c.DeleteObject(ref)
		if err != nil {
			log.Printf("Failed to sweep test objects. During object %s deletion error %s was raised", ref, err)
			return err
		}
	}
	return nil
}

func GetAllReturnFields(val reflect.Value) []string {
	rf := make([]string, 0)
	for i := 0; i < val.Type().NumField(); i++ {
		switch tag := val.Type().Field(i).Tag.Get("json"); tag {
		case "-":
		case "":
			continue
		default:
			parts := strings.Split(tag, ",")
			if parts[0] != "_ref" {
				rf = append(rf, parts[0])
			}
		}
	}

	return rf
}
