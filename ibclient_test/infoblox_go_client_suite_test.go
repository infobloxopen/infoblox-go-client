package ibclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInfobloxGoClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "InfobloxGoClient Suite")
}
