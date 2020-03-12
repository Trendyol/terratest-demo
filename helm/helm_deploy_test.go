package helm

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	httpHelper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"terratest-demo/util"
	"testing"
	"time"
)

type HelmDeploymentSuite struct {
	suite.Suite
}

func (s *HelmDeploymentSuite) SetupTest(){
	util.SetupK8sConfig()
}

func (s *HelmDeploymentSuite) TestPodDeployment() {

	helmChartPath := "./charts/minimal-pod"

	kubectlOptions := k8s.NewKubectlOptions("", "", "terratest")

	helmOptions := &helm.Options{
		SetValues: map[string]string{
			"image": "nginx:1.15.8",
		},
	}

	releaseName := fmt.Sprintf("nginx-%s", strings.ToLower(random.UniqueId()))

	defer helm.Delete(s.T(), helmOptions, releaseName, false)

	helm.Install(s.T(), helmOptions, helmChartPath, releaseName)

	podName := fmt.Sprintf("%s-minimal-pod", releaseName)

	k8s.WaitUntilPodAvailable(s.T(), kubectlOptions, podName, 20, 1*time.Second)

	tunnel := k8s.NewTunnel(kubectlOptions, k8s.ResourceTypePod, podName, 0, 80)

	defer tunnel.Close()

	tunnel.ForwardPort(s.T())

	endpoint := fmt.Sprintf("http://%s", tunnel.Endpoint())

	err := httpHelper.HttpGetWithRetryWithCustomValidationE(s.T(), endpoint,nil,20,1 * time.Second, func(statusCode int, body string) bool {
		return statusCode == 200 && strings.Contains(body, "Welcome to nginx")
	})

	assert.NoError(s.T(), err, "it should'nt get error")
}

func TestHelmDeploy(t *testing.T) {
	suite.Run(t, new(HelmDeploymentSuite))
}
