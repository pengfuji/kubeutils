/*
 * @Time : 2024/7/29 14:41
 * @Author : diehao.yuan
 * @Email : diehao.yuan@outlook.com
 * @File : tools.go
 */
package kubeutils

import (
	"context"
	"errors"
	"fmt"
	"github.com/dotbalo/kubeutils/utils/logs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
	"time"
)

// 定义结构体
type Tools struct {
	ClusterId     string
	DynamicClient *dynamic.DynamicClient
}

func NewClientSet(kubeconfig string, timeout int) (clientset *kubernetes.Clientset, err error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		msg := "解析kubeconfig错误: " + err.Error()
		return nil, errors.New(msg)
	}

	// 设置超时时间
	config.Timeout = time.Duration(timeout) * time.Second
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		msg := "创建clientset失败: " + err.Error()
		return nil, errors.New(msg)
	}
	return clientset, nil
}

// New函数可以用于配置一些默认值
func NewTools(kubeconfig string) (tools *Tools, err error) {
	// 加载kubeconfig文件，并创建restConfig
	config, _ := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	// 创建一个discovery客户端，这个客户端用于发现k8s集群可以支持的资源类型，同时一些自定义资源也可以使用该客户端进行发现
	// discoveryClient, _ := discovery.NewDiscoveryClientForConfig(config)
	// 创建dynamic client，用于创建k8s非结构化的数据，也就是Unstructured类型的数据
	// k8s自带的核心资源比如deployment、service，都是结构化数据，这些结构化数据都实现了统一的接口，也就是Object.runtime，这些类型可以使用clientset创建
	// 但是非结构化的数据需要使用dynamic client创建
	dynamicClient, err := dynamic.NewForConfig(config)
	// 创建一个非结构化数据对象，用于接受解析的yaml文件内容
	// obj := &unstructured.Unstructured{}
	// // GVK：Group Version Kind
	// _, gvk, err := serializeryaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode([]byte(yaml), nil, obj)
	// if err != nil {
	// 	return nil, err
	// }
	// // 获取非结构化数据的namespace
	// namespace := obj.GetNamespace()
	// if namespace == "" {
	// 	namespace = "default"
	// }
	// name := obj.GetName()
	// // 创建一个gvr：Group Version Resource  Resource=Kind + s example deployments
	// resource := gvk.Kind + "s"
	// gvr := schema.GroupVersionResource{
	// 	Group:    gvk.Group,
	// 	Version:  gvk.Version,
	// 	Resource: strings.ToLower(resource),
	// }
	// // 创建dynamic资源接口
	// dynamicResourceInterface := dynamicClient.Resource(gvr).Namespace(namespace)
	// tools.DynamicClient = dynamicResourceInterface
	// tools.Obj = obj
	// tools.Name = name
	tools = &Tools{}
	tools.DynamicClient = dynamicClient
	return tools, err
}

func createOrUpdate(dynamicClient *dynamic.DynamicClient, yamlContent, method string) (string, error) {
	// 拆分yaml
	var errMsgList []string
	methodMsg := "创建"
	yamlList := strings.Split(yamlContent, "---")
	// 循环yaml列表
	for k, v := range yamlList {
		// 去除空行和\n
		if v == "" || v == "\n" {
			continue
		}
		index := k + 1
		logs.Debug(map[string]interface{}{"yaml": v, "index": index}, "基于yaml创建或更新")
		// 创建一个非结构化数据对象，用于接受解析的yaml文件内容
		obj := &unstructured.Unstructured{}
		// GVK：Group Version Kind
		_, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode([]byte(v), nil, obj)
		if err != nil {
			// 如果此处失败，说明这一列的yaml内容有问题，直接下一个即可
			msg := fmt.Sprintf("第%d项yaml数据序列化: %s \n", index, err.Error())
			errMsgList = append(errMsgList, msg)
			continue
		}
		// yaml序列化成功继续往下执行
		// 获取非结构化数据的namespace
		namespace := obj.GetNamespace()
		if namespace == "" {
			namespace = "default"
		}
		// 创建一个gvr：Group Version Resource  Resource=Kind + s example deployments
		resource := gvk.Kind + "s"
		gvr := schema.GroupVersionResource{
			Group:    gvk.Group,
			Version:  gvk.Version,
			Resource: strings.ToLower(resource),
		}
		// 创建dynamic资源接口
		dynamicResourceInterface := dynamicClient.Resource(gvr).Namespace(namespace)
		switch method {
		case "Create":
			_, err = dynamicResourceInterface.Create(context.TODO(), obj, metav1.CreateOptions{})
		case "Update":
			methodMsg = "更新"
			_, err = dynamicResourceInterface.Update(context.TODO(), obj, metav1.UpdateOptions{})
		case "Apply":
			methodMsg = "应用"
			name := obj.GetName()
			_, err = dynamicResourceInterface.Apply(context.TODO(), name, obj, metav1.ApplyOptions{})
		case "Delete":
			methodMsg = "删除"
			name := obj.GetName()
			err = dynamicResourceInterface.Delete(context.TODO(), name, metav1.DeleteOptions{})
		}
		if err != nil {
			msg := fmt.Sprintf("第%d项yaml数据%s失败: %s", index, methodMsg, err.Error())
			errMsgList = append(errMsgList, msg)
			continue
		}
	}
	if len(errMsgList) == 0 {
		return "", nil
	} else {
		errMsg := fmt.Sprintf("%s失败", methodMsg)
		return strings.Join(errMsgList, "\n"), errors.New(errMsg)
	}

}

// 创建资源
func (c *Tools) Create(yamlContent string) (msg string, err error) {
	msg, err = createOrUpdate(c.DynamicClient, yamlContent, "Create")
	return
}

func (c *Tools) Update(yamlContent string) (msg string, err error) {
	msg, err = createOrUpdate(c.DynamicClient, yamlContent, "Update")
	return
}

func (c *Tools) Apply(yamlContent string) (msg string, err error) {
	msg, err = createOrUpdate(c.DynamicClient, yamlContent, "Apply")
	return
}

func (c *Tools) Delete(yamlContent string) (msg string, err error) {
	msg, err = createOrUpdate(c.DynamicClient, yamlContent, "Delete")
	return
}

// 创建资源
// func (c *Tools) Create() error {
// 	_, err := c.DynamicClient.Create(context.TODO(), c.Obj, metav1.CreateOptions{})
// 	return err
// }

// // 删除资源
// func (c *Tools) Delete() error {
// 	err := c.DynamicClient.Delete(context.TODO(), c.Name, metav1.DeleteOptions{})
// 	return err
// }

// // 更新资源
// func (c *Tools) Update() error {
// 	_, err := c.DynamicClient.Update(context.TODO(), c.Obj, metav1.UpdateOptions{})
// 	return err
// }

// // Apply资源
// func (c *Tools) Apply() error {
// 	_, err := c.DynamicClient.Apply(context.TODO(), c.Name, c.Obj, metav1.ApplyOptions{})
// 	return err
// }
