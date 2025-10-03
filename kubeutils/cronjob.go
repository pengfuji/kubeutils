/*
 * @Time : 2024/7/29 14:36
 * @Author : diehao.yuan
 * @Email : diehao.yuan@outlook.com
 * @File : cronjob.go
 */
package kubeutils

import (
	"context"
	"kubeutils/utils/log"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

// 定义结构体
type CronJob struct {
	InstanceInterface typedv1.BatchV1Interface
	Item              *batchv1.CronJob
}

// New函数用于配置一些默认值
func NewCronJob(kubeconfig string, item *batchv1.CronJob) *CronJob {
	// 首先调用instance的init函数，生成一个ResourceInstance的实例，并配置默认值和生成clientset
	instance := ResourceInstance{}
	instance.Init(kubeconfig)

	// 定义一个CronJob实例
	resource := CronJob{}
	resource.InstanceInterface = instance.Clientset.BatchV1()
	resource.Item = item
	return &resource
}

// 创建资源
func (c *CronJob) Create(namespace string) error {
	log.Infof("Name: ", c.Item.Name, "Namespace: ", namespace, "Create CronJob!")
	_, err := c.InstanceInterface.CronJobs(namespace).Create(context.TODO(), c.Item, metav1.CreateOptions{})
	return err
}

// 删除资源
func (c *CronJob) Delete(namespace, name string, gracePeriodSeconds *int64) error {
	log.Warnf("Namespace: ", namespace, "Name: ", name, "Delete CronJob!")
	deleteOptions := metav1.DeleteOptions{}

	// gracePeriodSeconds可配置，如果为0代表是强制删除
	if gracePeriodSeconds != nil {
		// 说明传递了gracePeriodSeconds
		deleteOptions.GracePeriodSeconds = gracePeriodSeconds
	}
	err := c.InstanceInterface.CronJobs(namespace).Delete(context.TODO(), name, deleteOptions)
	return err
}

// 删除多个资源
func (c *CronJob) DeleteList(namespace string, nameList []string, gracePeriodSeconds *int64) error {
	// 删除多个时，结构体会接受一个nameList的切片，循环该切片，然后调用Delete函数即可
	for _, name := range nameList {
		// 调用删除函数
		c.Delete("", name, gracePeriodSeconds)
	}
	// 忽略错误
	return nil
}

// 更新资源
func (c *CronJob) Update(namespace string) error {
	log.Warnf("Namespace: ", namespace, "Name: ", c.Item.Name, "Update CronJob!")
	_, err := c.InstanceInterface.CronJobs(namespace).Update(context.TODO(), c.Item, metav1.UpdateOptions{})
	return err
}

// 获取资源列表
func (c *CronJob) List(namespace, labelSelector, fieldSelector string) (items interface{}, err error) {
	log.Infof("Get CronJob List!")
	// 有可能是根据查询条件进行查询
	listOptions := metav1.ListOptions{
		FieldSelector: fieldSelector,
		LabelSelector: labelSelector,
	}
	list, err := c.InstanceInterface.CronJobs(namespace).List(context.TODO(), listOptions)
	items = list.Items
	return items, err
}

// 获取资源详情
func (c *CronJob) Get(namespace, name string) (item interface{}, err error) {
	log.Infof("Name: ", name, "Get CronJob Info!")
	i, err := c.InstanceInterface.CronJobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	i.APIVersion = "batch/v1"
	i.Kind = "CronJob"
	item = i
	return item, err
}
