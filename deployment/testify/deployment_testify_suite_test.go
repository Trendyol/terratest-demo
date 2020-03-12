package deployment

import (
	"fmt"
	"strings"
	"terratest-demo/util"
	"testing"
	"time"

	httpHelper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	coreV1 "k8s.io/api/core/v1"
)

const (
	/*  Namespace for testing purpose */
	TestNamespace = "terratest"
)

type DeploymentTestSuite struct {
	suite.Suite
}

func (s *DeploymentTestSuite) SetupSuite() {
	/* we setup our KUBECONFIG environment variable to communicate our cluster via kubectl */
	util.SetupK8sConfig()
}

func (s *DeploymentTestSuite) TestItShouldGiveOutputAsHelloWorld() {

	/* This our manifest to deploy deployment & service resources */
	k8sManifestPath := "../deployment.yaml"

	/* Preparing kubectl options */
	options := k8s.NewKubectlOptions("", "", TestNamespace)

	/* Clean our cluster after this method execution */
	defer k8s.KubectlDelete(s.T(), options, k8sManifestPath)

	/* Deploy our manifest files to cluster */
	k8s.KubectlApply(s.T(), options, k8sManifestPath)

	/* Wait until our application become available */
	k8s.WaitUntilServiceAvailable(s.T(), options, "hello-world-service", 10, 1*time.Second)

	/* Get deployed service to access our application */
	service := k8s.GetService(s.T(), options, "hello-world-service")

	/* 51 to 68 , Retrieve url to our service's url */
	serviceEndpoint := k8s.GetServiceEndpoint(s.T(), options, service,
		5000)
	serviceEndpointSplitted := strings.Split(serviceEndpoint, ":")
	nodes := k8s.GetNodes(s.T(), options)

	var node coreV1.Node
	for _, n := range nodes {
		if n.Name == serviceEndpointSplitted[0] {
			node = n
			break
		}
	}

	address := node.Status.Addresses[0].Address
	serviceURL := fmt.Sprintf("http://%s", address+":"+serviceEndpointSplitted[1])

	/* Sending HTTP request to our service via it's url and validate the response for expected status code and the body */
	err := httpHelper.HttpGetWithRetryE(s.T(), serviceURL, nil, 200, "Hello world!", 20, 1*time.Second)

	/* Check there is no error */
	assert.NoError(s.T(), err, "it should'nt get error")
}

func TestDeploymentSuite(t *testing.T) {
	suite.Run(t, new(DeploymentTestSuite))
}
