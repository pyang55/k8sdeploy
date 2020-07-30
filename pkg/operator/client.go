package operator

import (
	"fmt"
	"log"
	"os"

	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

// GetClient generates a k8s client based on kubeconfig
func GetClient(kfg string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kfg)
	if err != nil {
		panic(err.Error())
	}

	return kubernetes.NewForConfig(config)
}

// this is in case user wants to connect to EKS cluster with kubeconfigs
func NewClientset(cluster *eks.Cluster) (*kubernetes.Clientset, *rest.Config, error) {
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, nil, err
	}
	opts := &token.GetTokenOptions{
		ClusterID: aws.StringValue(cluster.Name),
	}
	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		return nil, nil, err
	}
	ca, err := base64.StdEncoding.DecodeString(aws.StringValue(cluster.CertificateAuthority.Data))
	if err != nil {
		return nil, nil, err
	}

	restConfigs := &rest.Config{
		Host:        aws.StringValue(cluster.Endpoint),
		BearerToken: tok.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: ca,
		},
	}

	clientset, err := kubernetes.NewForConfig(restConfigs)
	if err != nil {
		return nil, nil, err
	}
	return clientset, restConfigs, nil
}

func GetActionConfig(namespace string, config *rest.Config) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	var kubeConfig *genericclioptions.ConfigFlags
	// Create the ConfigFlags struct instance with initialized values from ServiceAccount

	insecure := true
	kubeConfig = genericclioptions.NewConfigFlags(false)
	kubeConfig.APIServer = &config.Host
	kubeConfig.BearerToken = &config.BearerToken
	kubeConfig.Namespace = &namespace
	kubeConfig.Insecure = &insecure
	if err := actionConfig.Init(kubeConfig, namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		return nil, err
	}
	return actionConfig, nil
}

// GetEksCluster returns a cluster name in a aws region
func GetEksCluster(region string, clustername string) *eks.Cluster {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	eksSvc := eks.New(sess)
	input := &eks.DescribeClusterInput{
		Name: aws.String(clustername),
	}
	result, err := eksSvc.DescribeCluster(input)
	if err != nil {
		log.Fatalf("Error calling DescribeCluster: %v", err)
	}
	return result.Cluster
}
