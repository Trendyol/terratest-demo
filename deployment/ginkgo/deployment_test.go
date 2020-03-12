package ginkgo

import (
	"fmt"
	httpHelper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coreV1 "k8s.io/api/core/v1"
	"strings"
	. "terratest-demo/util"
	"time"
)

type TestingT struct {
	GinkgoTInterface
	desc GinkgoTestDescription
}

func NewTestingT() TestingT {
	return TestingT{GinkgoT(), CurrentGinkgoTestDescription()}
}

func (i TestingT) Helper() {

}
func (i TestingT) Name() string {
	return i.desc.FullTestText
}

var _ = Describe("When I Deploy Hello-World Service", func() {
	var serviceURL string

	BeforeAll(func() {
		SetupK8sConfig()
		k8sManifestPath := "../deployment.yaml"

		options := k8s.NewKubectlOptions("", "", "terratest")

		//defer k8s.KubectlDelete(NewTestingT(), options, k8sManifestPath)

		k8s.KubectlApply(NewTestingT(), options, k8sManifestPath)

		k8s.WaitUntilServiceAvailable(NewTestingT(), options, "hello-world-service", 10, 1*time.Second)

		service := k8s.GetService(NewTestingT(), options, "hello-world-service")

		serviceEndpoint := k8s.GetServiceEndpoint(NewTestingT(), options, service,
			5000)

		serviceEndpointSplitted := strings.Split(serviceEndpoint, ":")

		nodes := k8s.GetNodes(NewTestingT(), options)

		var node coreV1.Node

		for _, n := range nodes {
			if n.Name == serviceEndpointSplitted[0] {
				node = n
				break
			}
		}

		address := node.Status.Addresses[0].Address

		serviceURL = fmt.Sprintf("http://%s", address+":"+serviceEndpointSplitted[1])

	})
	Describe("And I send HTTP GET request to Hello-World Service", func() {
		var payload string
		var err error
		BeforeAll(func() {
			payload, err = httpHelper.HTTPDoWithRetryE(NewTestingT(), "GET", serviceURL, nil, nil, 200, 5, 1*time.Second, nil)
		})

		It("Should Return \" Hello world!\" in payload", func() {
			Expect(payload).To(Equal("Hello world!"))
		})

		It("Should not return errors", func() {
			Expect(err).To(BeNil())
		})
	})
})
