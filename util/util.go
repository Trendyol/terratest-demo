package util

import (
	"os"
	"path/filepath"
)

const KubeConfigFilesPath = "../../kubeconfigs"

func SetupKubeconfigEnv() {
	var KUBECONFIG string

	filepath.Walk(KubeConfigFilesPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".k8s" {
			if KUBECONFIG != "" {
				KUBECONFIG += ":" + info.Name()
			} else {
				KUBECONFIG += info.Name()
			}
		}
		return nil
	})

	os.Setenv("KUBECONFIG", KubeConfigFilesPath+"/"+KUBECONFIG)
}
