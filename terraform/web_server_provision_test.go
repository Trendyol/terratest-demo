package main

import (
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
	"time"
)

func TestWebServer(t *testing.T) {
	tfOptions := &terraform.Options{
		TerraformDir: "./web-server",
	}

	defer terraform.Destroy(t, tfOptions)

	terraform.InitAndApply(t, tfOptions)

	tfOutput := terraform.Output(t, tfOptions, "url")

	http_helper.HttpGetWithRetry(t, tfOutput, nil, 200, "Hello, World", 15, 1*time.Second)
}
