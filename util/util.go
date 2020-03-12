package util

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func SetupK8sConfig() {
	_, filename, _, _ := runtime.Caller(0)
	k8sConfigFilesPath := path.Join(path.Join(path.Dir(filename), ".."), "kubeconfigs")

	fmt.Print("k8sConfigFilesPath", k8sConfigFilesPath)
	var KUBECONFIG string

	filepath.Walk(k8sConfigFilesPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".k8s" {
			if KUBECONFIG != "" {
				KUBECONFIG += ":" + info.Name()
			} else {
				KUBECONFIG += info.Name()
			}
		}
		return nil
	})

	os.Setenv("KUBECONFIG", k8sConfigFilesPath+"/"+KUBECONFIG)
}
