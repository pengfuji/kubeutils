/*
 * @Time : 2024/7/29 14:37
 * @Author : diehao.yuan
 * @Email : diehao.yuan@outlook.com
 * @File : deployment.go
 */
package kubeutils

import (
	"context"
	"kubeutils/utils/log"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

// 定义结构体
type Deployment struct {
	InstanceInterface typedv1.AppsV1Interface
	Item              *appsv1.Deployment
}

// New函数用于配置一些默认值
func NewDeployment(kubeconfig string, items *appsv1.Deployment) *Deployment {
	// 首先调用instance的init函数，生成一个ResourceInstance的实例，并配置默认值和生成clientset
	instance := ResourceInstance{}
	instance.Init(kubeconfig)

	// 定义一个Deployment函数
	resource := Deployment{}
	resource.InstanceInterface = instance.Clientset.AppsV1()
	resource.Item = items
	return &resource
}

// 创建资源
func (c *Deployment) Create(namespace string) error {
	log.Infof("Name: ", c.Item.Name, "Namespace: ", namespace, "Create Deployment!")
	_, err := c.InstanceInterface.Deployments(namespace).Create(context.TODO(), c.Item, metav1.CreateOptions{})
	return err
}

// 删除资源
func (c *Deployment) Delete(namespace, name string, gracePeriodSeconds *int64) error {
	log.Warnf("Namespace: ", namespace, "Delete Deployment!")
	deleteOptions := metav1.DeleteOptions{}

	// gracePeriodSeconds可配置，如果为0代表是强制删除
	if gracePeriodSeconds != nil {
		// 说明传递了gracePeriodSeconds
		deleteOptions.GracePeriodSeconds = gracePeriodSeconds
	}
	err := c.InstanceInterface.Deployments(namespace).Delete(context.TODO(), name, deleteOptions)
	return err
}

// 删除多个资源
func (c *Deployment) DeleteList(namespace string, nameList []string, gracePeriodSeconds *int64) error {
	// 删除多个时，结构体会接收一个nameList的切片，循环该切片，然后调用Delete函数即可
	for _, name := range nameList {
		// 调用删除函数
		c.Delete("", name, gracePeriodSeconds)
	}
	return nil
}

// 更新资源
func (c *Deployment) Update(namespace string) error {
	log.Warnf("Namespace: ", namespace, "Name: ", c.Item.Name, "Update Deployment!")
	_, err := c.InstanceInterface.Deployments(namespace).Update(context.TODO(), c.Item, metav1.UpdateOptions{})
	return err
}

// 获取资源列表
func (c *Deployment) List(namespace, labelSelector, fieldSelector string) (items interface{}, err error) {
	log.Infof("Get Deployment List!")
	// 有可能是根据查询条件进行查询
	listOptions := metav1.ListOptions{
		FieldSelector: fieldSelector,
		LabelSelector: labelSelector,
	}
	list, err := c.InstanceInterface.Deployments(namespace).List(context.TODO(), listOptions)
	items = list.Items
	return items, err
}

// 获取资源详情
func (c *Deployment) Get(namespace, name string) (item interface{}, err error) {
	log.Infof("Name: ", name, "Get Deployment Info!")
	i, err := c.InstanceInterface.Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	i.APIVersion = "apps/v1"
	i.Kind = "Deployment"
	item = i
	return item, err
}
