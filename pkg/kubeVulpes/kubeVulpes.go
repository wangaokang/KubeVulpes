package kubeVulpes

import (
	"KubeVulpes/pkg/db"
	"KubeVulpes/pkg/impl"
)

var CoreV1 impl.CoreV1Interface

type KubeVulpes struct {
	factory db.ShareDaoFactory
}

// var了一个变量要给这个变量复制才不会报空指针的错误
func Setup() {
	CoreV1 = impl.New()

}
