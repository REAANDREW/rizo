package tcp_test

import (
	// "io/ioutil"
	// "log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRizoTCP(t *testing.T) {
	// log.SetOutput(ioutil.Discard)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rizo TCP Suite")
}
