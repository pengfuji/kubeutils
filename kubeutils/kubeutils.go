/*
 * @Time : 2024/7/29 12:38
 * @Author : diehao.yuan
 * @Email : diehao.yuan@outlook.com
 * @File : kubeutils.go
 */
package kubeutils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

// 定义kubeutils接口
type KubeUtilser interface {
	// namespace, item
	Create(string) error
	// namespace, name, gracePeriodSeconds
	Delete(string, string, *int64) error
	// namespace, nameList, gracePeriodSeconds
	DeleteList(string, []string, *int64) error
	// namespace, item
	Update(string) error
	// namespace, labelSelector, fieldSelector
	List(string, string, string) (interface{}, error)
	// namespace name
	Get(string, string) (interface{}, error)
}

type ResourceInstance struct {
	Kubeconfig string
	Clientset  *kubernetes.Clientset
}

func (c *ResourceInstance) Init(kubeconfig string) {
	c.Kubeconfig = kubeconfig

	// 生成Clientset
	restConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(c.Kubeconfig))
	if err != nil {
		msg := "解析kubeconfig错误: " + err.Error()
		panic(msg)
	}

	// 设置超时时间
	restConfig.Timeout = 15 * time.Second
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		msg := "创建clientset失败: " + err.Error()
		panic(msg)
	}
	c.Clientset = clientSet
}
