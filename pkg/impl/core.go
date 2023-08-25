package impl

import "KubeVulpes/pkg/db"

type CoreV1Interface interface {
	UserGetter
}

type KubeVulpes struct {
	factory db.ShareDaoFactory
}

func (k *KubeVulpes) User() UserInterface {
	return newUser(k)
}
