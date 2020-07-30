package operator

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/strvals"
)

//run a helm deploy from a chart
func HelmDeploy(actionConfig *action.Configuration, chrtPath string, deployName string, namespace string, opts ...string) error {
	chart, err := loader.Load(chrtPath)
	if err != nil {
		return err
	}

	client := action.NewGet(actionConfig)
	set, rawVals, _ := checkSets(opts)
	_, err = client.Run(deployName)
	if err != nil {
		//if brand new release, we install new chart
		fmt.Printf("Can not find release %s...installing\n", deployName)
		if set {
			install := action.NewInstall(actionConfig)
			install.Namespace = namespace
			install.ReleaseName = deployName
			_, err := install.Run(chart, rawVals)
			if err != nil {
				panic(err)
			}
		} else {
			install := action.NewInstall(actionConfig)
			install.Namespace = namespace
			install.ReleaseName = deployName
			_, err := install.Run(chart, nil)
			if err != nil {
				panic(err)
			}
		}

		// we're running this again after inital installation
		deployCheck, err := client.Run(deployName)
		if err != nil {
			log.Fatalf("%s", err)
			os.Exit(3)
		}
		if deployCheck.Info.Status.String() != "deployed" {
			fmt.Printf("%s\n", deployCheck.Info.Description)
			os.Exit(3)
		}
		// checks if existing release had been deployed, if so do rolling update
		// this will also ensure that any previously failed deployments will be upgraded
	} else {
		fmt.Printf("Found existing deployment for %s...updating\n", deployName)
		if set {
			upgrade := action.NewUpgrade(actionConfig)
			upgrade.Namespace = namespace
			upgrade.ResetValues = true
			upgrade.Atomic = true // this will roll back on failure
			_, err := upgrade.Run(deployName, chart, rawVals)
			if err != nil {
				panic(err)
			}
		} else {
			upgrade := action.NewUpgrade(actionConfig)
			upgrade.Namespace = namespace
			upgrade.ResetValues = true
			upgrade.Atomic = true // this will roll back on failure
			_, err := upgrade.Run(deployName, chart, nil)
			if err != nil {
				panic(err)
			}
		}
	}
	return nil
}

// parses the extra variables entered
func vals(values []string) (map[string]interface{}, error) {
	base := map[string]interface{}{}
	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			return nil, fmt.Errorf("failed parsing --set data: %s", err)
		}
	}
	return base, nil
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

//checks if theres extra variables set
func checkSets(sets []string) (bool, map[string]interface{}, error) {
	if len(sets) > 0 && sets[0] != "" {
		raw := strings.Split(sets[0], ",")
		rawVals, err := vals(raw)
		if err != nil {
			return false, nil, err
		}
		return true, rawVals, nil
	}
	return false, nil, nil
}
