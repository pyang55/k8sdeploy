package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"k8sdeploy/pkg/operator"
)

var kubeConfig string
var nameSpace string
var ns []string
var set string
var clusterName string
var token string
var region string
var releaseName string
var chartPath string

var kcfg = &cobra.Command{
	Use:   "kubeconfig",
	Short: "deploy using kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {
		today := time.Now()
		namespaces := strings.Split(nameSpace, ",")
		//assumes chartPath includes .tgz file
		if _, err := os.Stat(chartPath); err == nil {
			for _, z := range namespaces {
				go operator.HelmDeploy(kcfgActionConfig(kubeConfig), chartPath, releaseName, z, set)
			}
			client, _ := operator.GetClient(kubeConfig)
			operator.DeployLogger(client, releaseName, namespaces, today)
		} else {
			fmt.Printf("%s", err)
		}
	},
}

var eks = &cobra.Command{
	Use:   "eks",
	Short: "deploying to eks without a kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {
		today := time.Now()
		namespaces := strings.Split(nameSpace, ",")
		//assumes chartPath includes .tgz file
		if _, err := os.Stat(chartPath); err == nil {
			for _, z := range namespaces {
				go operator.HelmDeploy(eksActionConfig(region, clusterName), chartPath, releaseName, z, set)
			}
			eksClient, _, _ := operator.NewClientset(operator.GetEksCluster(region, clusterName))
			operator.DeployLogger(eksClient, releaseName, namespaces, today)
		} else {
			fmt.Printf("%s", err)
		}
	},
}

var gke = &cobra.Command{
	Use:   "gke",
	Short: "deploying to gke without a kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {
		//
	},
}

var aks = &cobra.Command{
	Use:   "aks",
	Short: "deploying to aks without a kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {
		//
	},
}

func eksActionConfig(region string, clusterName string) *action.Configuration {
	cluster := operator.GetEksCluster(region, clusterName)
	_, configs, err := operator.NewClientset(cluster)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}
	actionConfig, err := operator.GetActionConfig(nameSpace, configs)
	if err != nil {
		log.Println(err)
	}
	return actionConfig
}

func kcfgActionConfig(kubeConfig string) *action.Configuration {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(kube.GetConfig(kubeConfig, "", nameSpace), nameSpace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		panic(err)
	}
	return actionConfig
}

func init() {
	kcfg.Flags().StringVar(&chartPath, "chartdir", "", "Local path to the directory containing chart(s) to upload (Defaults to cwd)")
	kcfg.Flags().StringVar(&kubeConfig, "configpath", "", "full path to kubeconfig file")
	kcfg.Flags().StringVar(&nameSpace, "namespace", "", "namespace to deploy to, can enter a comma separated list for multiple namespaces. Defaults to all pods in the cluster if nothing is entered")
	kcfg.Flags().StringVar(&releaseName, "releasename", "", "the release name of the deployment")
	kcfg.Flags().StringVar(&set, "set", "", "Specify each parameter using the `--set key=value[,key=value]`")

	eks.Flags().StringVar(&clusterName, "clustername", "", "If using --platform eks you must specify clustername")
	eks.Flags().StringVar(&region, "region", "", "If using --platform eks you must specify region")
	eks.Flags().StringVar(&chartPath, "chartdir", "", "Local path to the directory containing chart(s) to upload (Defaults to cwd)")
	eks.Flags().StringVar(&nameSpace, "namespace", "", "namespace to deploy to, can enter a comma separated list for multiple namespaces. Defaults to all pods in the cluster if nothing is entered")
	eks.Flags().StringVar(&releaseName, "releasename", "", "the release name of the deployment")
	eks.Flags().StringVar(&set, "set", "", "Specify each parameter using the `--set key=value[,key=value]`")

	deploy.AddCommand(kcfg)
	deploy.AddCommand(eks)
	deploy.AddCommand(gke)
	deploy.AddCommand(aks)

	if err := kcfg.MarkFlagRequired("releasename"); err != nil {
		log.Fatalf("error marking flag as required: %v", err)
	}
	if err := kcfg.MarkFlagRequired("configpath"); err != nil {
		log.Fatalf("error marking flag as required: %v", err)
	}
	if err := kcfg.MarkFlagRequired("chartdir"); err != nil {
		log.Fatalf("error marking flag as required: %v", err)
	}
	if err := kcfg.MarkFlagRequired("namespace"); err != nil {
		log.Fatalf("error marking flag as required: %v", err)
	}
	if err := eks.MarkFlagRequired("clustername"); err != nil {
		log.Fatalf("error marking flag as required: %v", err)
	}
	if err := eks.MarkFlagRequired("region"); err != nil {
		log.Fatalf("error marking flag as required: %v", err)
	}
}
