/*
 * @Time : 2024/7/29 14:37
 * @Author : diehao.yuan
 * @Email : diehao.yuan@outlook.com
 * @File : ingress.go
 */
package kubeutils

import (
	"context"
	"github.com/YuanDieHao/kubeutils/utils/log"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
)

// 定义结构体
type Ingress struct {
	InstanceInterface typedv1.NetworkingV1Interface
	Item              *networkingv1.Ingress
}

// New函数用于配置一些默认信息
func NewIngerss(kubeconfig string, item *networkingv1.Ingress) *Ingress {
	// 首先调用instance的init函数，生成一个ResourceInstance的实例，并配置默认值和生成clientset
	instance := ResourceInstance{}
	instance.Init(kubeconfig)

	// 定义一个ConfigMap实例
	resource := Ingress{}
	resource.InstanceInterface = instance.Clientset.NetworkingV1()
	resource.Item = item
	return &resource
}

// 创建资源
func (c *Ingress) Create(namespace string) error {
	log.Infof("Name: ", c.Item.Name, "Namespace: ", namespace, "Create Ingress!")
	_, err := c.InstanceInterface.Ingresses(namespace).Create(context.TODO(), c.Item, metav1.CreateOptions{})
	return err
}

// 删除资源
func (c *Ingress) Delete(namespace, name string, gracePeriodSeconds *int64) error {
	log.Warnf("Namespace: ", namespace, "Delete Ingress!")
	deleteOptions := metav1.DeleteOptions{}

	// gracePeriodSeconds可配置，如果为0代表是强制删除
	if gracePeriodSeconds != nil {
		// 说明传递了gracePeriodSeconds
		deleteOptions.GracePeriodSeconds = gracePeriodSeconds
	}
	err := c.InstanceInterface.Ingresses(namespace).Delete(context.TODO(), name, deleteOptions)
	return err
}

// 删除多个资源
func (c *Ingress) DeleteList(namespace string, nameList []string, gracePeriodSeconds *int64) error {
	// 删除多个时，结构体会接收一个nameList的切片，循环该切片，然后调用Delete函数即可
	for _, name := range nameList {
		// 调用删除函数
		c.Delete("", name, gracePeriodSeconds)
	}
	// 忽略错误
	return nil
}

// 更新资源
func (c *Ingress) Update(namespace string) error {
	log.Warnf("Namespace: ", namespace, "Name: ", c.Item.Name, "Update Ingress!")
	_, err := c.InstanceInterface.Ingresses(namespace).Update(context.TODO(), c.Item, metav1.UpdateOptions{})
	return err
}

// 获取资源列表
func (c *Ingress) List(namespace, labelSelector, fieldSelector string) (items interface{}, err error) {
	log.Infof("Get ConfigMap List!")
	// 有可能是根据查询条件进行查询
	listOptions := metav1.ListOptions{
		FieldSelector: fieldSelector,
		LabelSelector: labelSelector,
	}
	list, err := c.InstanceInterface.Ingresses(namespace).List(context.TODO(), listOptions)
	items = list.Items
	return items, err
}

// 获取资源详情
func (c *Ingress) Get(namespace, name string) (item interface{}, err error) {
	log.Infof("Name: ", name, "Get Ingress Info!")
	i, err := c.InstanceInterface.Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	i.APIVersion = "networking.k8s.io/v1"
	i.Kind = "Ingress"
	item = i
	return item, err
}
