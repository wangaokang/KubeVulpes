package impl

import "KubeVulpes/pkg/db"

type CoreV1Interface interface {
	UserGetter
}

type KubeVulpes struct {
	factory db.ShareDaoFactory
}

//func NewKubeVulpes(V *KubeVulpes) UserInterface {
//	return NewUser(V)
//}

func (k *KubeVulpes) User() UserInterface {
	return newUser(k)
}

func New() CoreV1Interface {
	return &KubeVulpes{
		//factory: k.factory,
	}
}
