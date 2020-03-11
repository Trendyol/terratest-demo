package helm

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/suite"
	coreV1 "k8s.io/api/core/v1"
	"testing"
)

type HelmTestSuite struct {
	suite.Suite
}

func (s *HelmTestSuite) TestPodTemplateRendersContainerImage() {

	helmChartPath := "./charts/minimal-pod"

	options := &helm.Options{
		SetValues: map[string]string{
			"image": "nginx:1.15.8",
		},
	}

	output := helm.RenderTemplate(s.T(), options, helmChartPath, "nginx", []string{"templates/pod.yaml"})

	var pod coreV1.Pod

	helm.UnmarshalK8SYaml(s.T(), output, &pod)

	expectedContainerImage := "nginx:1.15.8"
	podFirstContainerImage := pod.Spec.Containers[0].Image

	if expectedContainerImage != podFirstContainerImage {
		s.T().Fatalf("Rendered container image (%s) is not expected (%s)", podFirstContainerImage, expectedContainerImage)
	}
}

func TestRunHelmTestSuite(t *testing.T) {
	suite.Run(t, new(HelmTestSuite))
}
