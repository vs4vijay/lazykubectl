package k8s

import (
	"k8s.io/client-go/util/homedir"
)

func Home() string {
	return homedir.HomeDir()
}
