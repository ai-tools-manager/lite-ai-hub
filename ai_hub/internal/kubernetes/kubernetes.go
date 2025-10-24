package kubernetes

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	kedahttpv1alpha1 "github.com/kedacore/keda-http-add-on/api/v1alpha1"
	kedahttpclientset "github.com/kedacore/keda-http-add-on/pkg/generated/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset        *kubernetes.Clientset
	namespace        string
	logger           *zap.SugaredLogger
	sidecarImage     string
}

func NewClient(namespace string, sidecarImage string, logger *zap.SugaredLogger) (*Client, error) {
	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()
	if err != nil {
		kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build kubernetes config: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}



	return &Client{
		clientset:        clientset,
		namespace:        namespace,
		logger:           logger,
		sidecarImage:     sidecarImage,
	}, nil
}

func (c *Client) GetClientset() *kubernetes.Clientset {
	return c.clientset
}

func (c *Client) GetNamespace() string {
	return c.namespace
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.clientset.CoreV1().Namespaces().Get(ctx, c.namespace, metav1.GetOptions{})
	return err
}

func (c *Client) DeployWorker(ctx context.Context, name, image string, envs map[string]string) (string, error) {
	c.logger.Infof("Deploying worker %s with image %s", name, image)

