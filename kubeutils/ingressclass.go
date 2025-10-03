/*
 * @Time : 2024/7/29 14:37
 * @Author : diehao.yuan
 * @Email : diehao.yuan@outlook.com
 * @File : ingressclass.go
 */
package kubeutils

import (
	"context"
	"kubeutils/utils/log"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
)

// 定义结构体
type IngressClass struct {
	InstanceInterface typedv1.NetworkingV1Interface
	Item              *networkingv1.IngressClass
}

// New函数可以设置一些默认值
func NewIngressClass(kubeconfig string, item *networkingv1.IngressClass) *IngressClass {
	// 首先调用instance的init函数，生成一个ResourceInstance的实例，并配置默认值和生成clientset
	instance := ResourceInstance{}
	instance.Init(kubeconfig)

	// 定义一个IngressClass实例
	resource := IngressClass{}
	resource.InstanceInterface = instance.Clientset.NetworkingV1()
	resource.Item = item
	return &resource
}

// 创建资源
func (c *IngressClass) Create(namespace string) error {
	log.Infof("Name: ", c.Item.Name, "Namespace: ", namespace, "Create IngressClass!")
	_, err := c.InstanceInterface.IngressClasses().Create(context.TODO(), c.Item, metav1.CreateOptions{})
	return err
}

// 删除资源
func (c *IngressClass) Delete(namespace, name string, gracePeriodSeconds *int64) error {
	log.Warnf("Namespace: ", namespace, "Name: ", name, "Delete IngressClass!")
	deleteOptions := metav1.DeleteOptions{}

	//  gracePeriodSeconds可配置，如果为0代表是强制删除
	if gracePeriodSeconds != nil {
		// 说明传递了 gracePeriodSeconds
		deleteOptions.GracePeriodSeconds = gracePeriodSeconds
	}
	err := c.InstanceInterface.IngressClasses().Delete(context.TODO(), name, deleteOptions)
	return err
}

// 删除多个资源
func (c *IngressClass) DeleteList(namespace string, nameList []string, gracePeriodSeconds *int64) error {
	// 删除多个时，结构体会接收一个nameList的切片，循环该切片，然后调用Delete函数即可
	for _, name := range nameList {
		// 调用删除函数
		c.Delete("", name, gracePeriodSeconds)
	}
	// 忽略错误
	return nil
}

// 更新资源
func (c *IngressClass) Update(namespace string) error {
	log.Warnf("Namespace: ", namespace, "Name: ", c.Item.Name, "Update IngressClass!")
	_, err := c.InstanceInterface.IngressClasses().Update(context.TODO(), c.Item, metav1.UpdateOptions{})
	return err
}

// 获取资源列表
func (c *IngressClass) List(namespace, labelSelector, fieldSelector string) (items interface{}, err error) {
	log.Infof("Get IngressClass List!")
	// 有可能是根据条件进行查询
	listOptions := metav1.ListOptions{
		FieldSelector: fieldSelector,
		LabelSelector: labelSelector,
	}
	list, err := c.InstanceInterface.IngressClasses().List(context.TODO(), listOptions)
	items = list.Items
	return items, err
}

// 获取资源详情
func (c *IngressClass) Get(namespace, name string) (item interface{}, err error) {
	log.Infof("Name: ", name, "Get IngressClass Info!")
	i, err := c.InstanceInterface.IngressClasses().Get(context.TODO(), name, metav1.GetOptions{})
	i.APIVersion = "networking.k8s.io/v1"
	i.Kind = "IngressClass"
	item = i
	return item, err
}
