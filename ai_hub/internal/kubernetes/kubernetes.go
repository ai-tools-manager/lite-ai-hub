package kubernetes

import (
	"context"
	"errors" // ДОБАВЛЕНО
	"fmt"

	kedaclientset "github.com/kedacore/keda/v2/pkg/generated/clientset/versioned"
	"go.uber.org/zap"

	// Kubernetes
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// int32Ptr - вспомогательная функция для получения указателя на int32.
func int32Ptr(i int32) *int32 { return &i }

// pathTypePtr - вспомогательная функция для получения указателя на PathType (для Ingress).
func pathTypePtr(pt networkingv1.PathType) *networkingv1.PathType { return &pt }

type Client struct {
	kubeClient    kubernetes.Interface
	kedaClient    kedaclientset.Interface
	dynamicClient dynamic.Interface
	namespace     string
	logger        *zap.SugaredLogger
}

func NewClient(namespace string, logger *zap.SugaredLogger) (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Warnf("Не удалось получить in-cluster config, используется kubeconfig: %v", err)

		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}

		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		config, err = kubeConfig.ClientConfig()
		if err != nil {
			logger.Errorf("Ошибка сборки конфига из kubeconfig (по правилам kubectl): %v", err)
			return nil, err
		}
		logger.Info("Конфиг успешно загружен по правилам kubectl (env KECONFIG или ~/.kube/config)")
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Errorf("Ошибка создания K8s клиента: %v", err)
		return nil, err
	}

	kedaClient, err := kedaclientset.NewForConfig(config)
	if err != nil {
		logger.Errorf("Ошибка создания KEDA Core клиента: %v", err)
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		logger.Errorf("Ошибка создания Dynamic клиента: %v", err)
		return nil, err
	}

	client := &Client{
		kubeClient:    kubeClient,
		kedaClient:    kedaClient,
		dynamicClient: dynamicClient,
		namespace:     namespace,
		logger:        logger,
	}

	if err := client.checkDependencies(); err != nil {
		return nil, err
	}

	return client, nil
}

// Метод для проверки установки keda
func (c *Client) checkDependencies() error {
	c.logger.Info("Проверка наличия KEDA API...")
	apiGroups, err := c.kubeClient.Discovery().ServerGroups()
	if err != nil {
		c.logger.Warnf("Не удалось получить список API групп (это может быть проблемой RBAC): %v", err)
		// Не возвращаем ошибку, т.к. KEDA может быть установлена, а прав на Discovery не быть
		// Попытаемся проверить ресурсы напрямую (резервный метод)
		return c.checkDependenciesFallback(context.Background())
	}

	foundKedaCore := false
	foundKedaHttp := false

	for _, group := range apiGroups.Groups {
		if group.Name == "keda.sh" {
			foundKedaCore = true
		}
		if group.Name == "http.keda.sh" {
			foundKedaHttp = true
		}
	}

	if !foundKedaCore {
		c.logger.Error("API группа 'keda.sh' не найдена. KEDA Core не установлен или неисправен.")
		return errors.New("KEDA Core (keda.sh) API group not found. Please install KEDA")
	}

	if !foundKedaHttp {
		c.logger.Error("API группа 'http.keda.sh' не найдена. KEDA HTTP Add-on не установлен или неисправен.")
		return errors.New("KEDA HTTP Add-on (http.keda.sh) API group not found. Please install the KEDA HTTP Add-on")
	}

	c.logger.Info("KEDA Core и KEDA HTTP Add-on успешно найдены.")
	return nil
}

// ДОБАВЛЕНО: Резервный метод проверки, если Discovery API недоступен
func (c *Client) checkDependenciesFallback(ctx context.Context) error {
	c.logger.Info("Резервная проверка KEDA: пытаемся получить CRD...")

	scaledObjectGVR := schema.GroupVersionResource{
		Group:    "keda.sh",
		Version:  "v1alpha1",
		Resource: "scaledobjects",
	}
	// Мы пытаемся получить список с лимитом 1, чтобы просто проверить, существует ли CRD
	_, err := c.dynamicClient.Resource(scaledObjectGVR).Namespace(c.namespace).List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		c.logger.Errorf("Резервная проверка KEDA Core не удалась: %v", err)
		return errors.New("KEDA Core (keda.sh) CRD not accessible. Please install KEDA and check RBAC")
	}

	// Проверка KEDA HTTP Add-on (HTTPScaledObjects)
	httpScaledObjectGVR := schema.GroupVersionResource{
		Group:    "http.keda.sh",
		Version:  "v1alpha1",
		Resource: "httpscaledobjects",
	}
	_, err = c.dynamicClient.Resource(httpScaledObjectGVR).Namespace(c.namespace).List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		c.logger.Errorf("Резервная проверка KEDA HTTP Add-on не удалась: %v", err)
		return errors.New("KEDA HTTP Add-on (http.keda.sh) CRD not accessible. Please install the Add-on and check RBAC")
	}

	c.logger.Info("Резервная проверка KEDA прошла успешно.")
	return nil
}

// ... GetKubeClient, GetNamespace, Ping ... (без изменений)
func (c *Client) GetKubeClient() kubernetes.Interface {
	return c.kubeClient
}

func (c *Client) GetNamespace() string {
	return c.namespace
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.kubeClient.CoreV1().Namespaces().Get(ctx, c.namespace, metav1.GetOptions{})
	return err
}

// Деплоит воркер с автомасштабированием по запросам
func (c *Client) DeployWorker(ctx context.Context, name, image string, minReplicas int64) (string, error) {
	hostName := fmt.Sprintf("%s.worker.local", name)

	c.logger.Infof("Развертывание воркера %s с образом %s для хоста %s", name, image, hostName)

	labels := map[string]string{"app": name}

	// 1. Определение Deployment (без изменений)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: c.namespace},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(int32(minReplicas)),
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "worker",
						Image: image,
						Ports: []corev1.ContainerPort{{ContainerPort: 8080}},
					}},
				},
			},
		},
	}
	c.logger.Infof("Создание Deployment: %s...", name)
	_, err := c.kubeClient.AppsV1().Deployments(c.namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		c.logger.Errorf("Ошибка создания Deployment: %v", err)
		return "", err
	}

	// 2. Определение "Реального" Service (без изменений)
	// KEDA будет направлять трафик на этот сервис ПОСЛЕ "пробуждения"
	serviceName := name + "-svc"
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: serviceName, Namespace: c.namespace},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       80,
				TargetPort: intstr.FromInt32(8080),
			}},
		},
	}
	c.logger.Infof("Создание 'Реального' Service: %s...", serviceName)
	_, err = c.kubeClient.CoreV1().Services(c.namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		c.logger.Errorf("Ошибка создания 'Реального' Service: %v", err)
		return "", err
	}

	// 3. Определение HTTPScaledObject
	c.logger.Infof("Создание HTTPScaledObject: %s...", name)
	httpScaledObjectGVR := schema.GroupVersionResource{
		Group:    "http.keda.sh",
		Version:  "v1alpha1",
		Resource: "httpscaledobjects",
	}

	httpScaledObject := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "http.keda.sh/v1alpha1",
			"kind":       "HTTPScaledObject",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": c.namespace,
				"annotations": map[string]interface{}{
					"httpscaledobject.keda.sh/skip-scaledobject-creation": "false",
				},
			},
			"spec": map[string]interface{}{
				"hosts":        []string{hostName},
				"pathPrefixes": []string{"/"},
				"scaleTargetRef": map[string]interface{}{
					"name":       name,
					"kind":       "Deployment",
					"apiVersion": "apps/v1",
					"service":    serviceName,
					"port":       int64(80),
				},
				"replicas": map[string]interface{}{
					"min": minReplicas,
					"max": int64(20),
				},
				"scaledownPeriod": int64(120),

				"targetPendingRequests": int64(200),
			},
		},
	}
	_, err = c.dynamicClient.Resource(httpScaledObjectGVR).Namespace(c.namespace).Create(ctx, httpScaledObject, metav1.CreateOptions{})
	if err != nil {
		c.logger.Errorf("Ошибка создания HTTPScaledObject: %v", err)
		return "", err
	}

	c.logger.Infof("Ресурсы для HTTP воркера '%s' созданы.", name)
	return hostName, nil
}
