/*
 * @Time : 2024/7/29 12:39
 * @Author : diehao.yuan
 * @Email : diehao.yuan@outlook.com
 * @File : pod.go
 */
package kubeutils

import (
	"context"
	"kubeutils/utils/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// 定义结构体
type Pod struct {
	InstanceInterface typedv1.CoreV1Interface
	Item              *corev1.Pod
}

// New函数可以用于配置一些默认值
func NewPod(kubeconfig string, item *corev1.Pod) *Pod {
	// 首先调用instance的init函数，生成一个ResourceInstance实例，并配置默认值和生成clientset
	instance := ResourceInstance{}
	instance.Init(kubeconfig)

	// 定义一个Pod实例
	resource := Pod{}
	resource.InstanceInterface = instance.Clientset.CoreV1()
	resource.Item = item
	return &resource
}

// 创建资源
func (c *Pod) Create(namespace string) error {
	log.Infof("Name: ", c.Item.Name, "Namespace: ", namespace, "Create Pod!")
	_, err := c.InstanceInterface.Pods(namespace).Create(context.TODO(), c.Item, metav1.CreateOptions{})
	return err
}

// 删除资源
func (c *Pod) Delete(namespace, name string, gracePeriodSeconds *int64) error {
	log.Warnf("Namespace: ", namespace, "Name: ", name, "Delete Pod!")
	deleteOptions := metav1.DeleteOptions{}

	// gracePeriodSeconds 可配置，如果为0代表是强制删除
	if gracePeriodSeconds != nil {
		// 说明传递了gracePeriodSeconds
		deleteOptions.GracePeriodSeconds = gracePeriodSeconds
	}
	err := c.InstanceInterface.Pods(namespace).Delete(context.TODO(), name, deleteOptions)
	return err
}

// 删除多个资源
func (c *Pod) DeleteList(namespace string, nameList []string, gracePeriodSeconds *int64) error {
	// 删除多个时，结构体会接收一个nameList的切片，循环该切片，然后调用Delete函数即可
	for _, name := range nameList {
		// 调用删除函数
		c.Delete(namespace, name, gracePeriodSeconds)
	}
	// 忽略错误
	return nil
}

// 更新资源
func (c *Pod) Update(namespace string) error {
	log.Warnf("Namespace: ", namespace, "Name: ", c.Item.Name, "Update Pod!")
	_, err := c.InstanceInterface.Pods(namespace).Update(context.TODO(), c.Item, metav1.UpdateOptions{})
	return err
}

// 获取资源列表
func (c *Pod) List(namespace, labelSelector, fieldSelector string) (items interface{}, err error) {
	log.Infof("Namespace: ", namespace, "Get Pod List!")
	// 有可能是根据查询条件进行查询
	listOptions := metav1.ListOptions{
		FieldSelector: fieldSelector,
		LabelSelector: labelSelector,
	}
	list, err := c.InstanceInterface.Pods(namespace).List(context.TODO(), listOptions)
	items = list.Items
	return items, err
}

// 获取资源配置
func (c *Pod) Get(namespace, name string) (item *corev1.Pod, err error) {
	log.Infof("Namespace: ", namespace, "Name: ", name, "Get Pod Info!")
	item, err = c.InstanceInterface.Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	item.APIVersion = "v1"
	item.Kind = "Pod"
	return item, err
}
