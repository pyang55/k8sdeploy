# k8sdeploy
> A tool to deploy to multiple namespaces in a kubernetes cluster. There is added feature set to deploy to different types of clusters like GKE and EKS without the use of a kubeconfig (currently only EKS is supported)

## Installation

OS X & Linux golang:

```sh
go build -o k8sdeploy
```

Docker (creating binary in docker and exporting it):

```sh
docker run --rm -i hadolint/hadolint hadolint --ignore DL3000 --ignore DL3008 --ignore DL3018 --ignore DL3019 --ignore DL4000 - < ./Dockerfile
docker build --tag k8sdeploy:latest .
docker create --name temp k8sdeploy:latest
docker cp temp:/go/src/k8sdeploy/k8sdeploy .
docker rm temp
```

## Development setup
* If you are running this locally or as part of your CI/CD pipeline, you must have the proper permissions into your environment
* Setting deploy_timestamp to now in deployment.yaml. There is an K8s event watcher that looks for all update events after deployment timestamp is set.
```sh
metadata:
  annotations:
    deploy_timestamp: {{ now }}
```

## Usage example

Deploying directly to EKS cluster without a kubeconfig:

```sh
k8sdeploy deploy eks --clustername <eks-cluster-name> --releasename <name-of-release> --region <cluster-region> --namespace <namespace1,namespace2,namespace3> --chartdir <full-path-to-tgz-chart-file> --set <set-string-values>
```

Deploying with a kubeconfig:

```sh
k8sdeploy deploy kubeconfig --configpath <full-path-to-kubeconfig> --releasename <name-of-release> --namespace <namespace1,namespace2,namespace3> --chartdir <full-path-to-tgz-chart-file> --set <set-string-values>
```

## Sample output

```sh
build	29-Jul-2020 19:23:20	Starting deployment in namespace=name-space-1 for app=customapp at 2020-07-29 19:23:20 -0700 PDT
build	29-Jul-2020 19:23:20	Waiting for deployment  rollout to finish: 0 of 2 updated replicas are available...
build	29-Jul-2020 19:23:20	Waiting for deployment  rollout to finish: 0 of 2 updated replicas are available...
build	29-Jul-2020 19:23:20	Starting deployment in namespace=name-space-2 for app=customapp at 2020-07-29 19:23:20 -0700 PDT
build	29-Jul-2020 19:23:20	Waiting for deployment  rollout to finish: 0 of 2 updated replicas are available...
build	29-Jul-2020 19:23:35	Waiting for deployment  rollout to finish: 1 of 2 updated replicas are available...
build	29-Jul-2020 19:23:35	Waiting for deployment  rollout to finish: 1 of 2 updated replicas are available...
build	29-Jul-2020 19:23:49	Waiting for deployment  rollout to finish: 1 of 2 updated replicas are available...
build	29-Jul-2020 19:23:56	Waiting for deployment  rollout to finish: 2 of 2 updated replicas are available...
build	29-Jul-2020 19:23:56	Successful Deployment of customapp on name-space-2
build	29-Jul-2020 19:23:58	Waiting for deployment  rollout to finish: 2 of 2 updated replicas are available...
build	29-Jul-2020 19:23:58	Successful Deployment of customapp on name-space-2
build	29-Jul-2020 19:24:10	All deployments finished, sutting down watcher gracefully
build	29-Jul-2020 19:24:10	+----------------+--------------+---------+
build	29-Jul-2020 19:24:10	| APP            | NAMESPACE    | STATUS  |
build	29-Jul-2020 19:24:10	+----------------+--------------+---------+
build	29-Jul-2020 19:24:10	| customapp      | name-space-1 | Success |
build	29-Jul-2020 19:24:10	| customapp      | name-space-2 | Success |
build	29-Jul-2020 19:24:10	+----------------+--------------+---------+
```

## Notes
* This is also written with support for helm3 only. There is only support to deploy with kubeconfig for helm2
* This currently works for me. This is what we use to deploy to multiple microservices to all of production. If this works for you, that is fantastic...if it doesn't...well, it still works for me.
