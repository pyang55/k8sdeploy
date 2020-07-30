module k8sdeploy

go 1.13

replace (
	github.com/Sirupsen/logrus v1.0.5 => github.com/sirupsen/logrus v1.0.5
	github.com/Sirupsen/logrus v1.3.0 => github.com/Sirupsen/logrus v1.0.6
	github.com/Sirupsen/logrus v1.4.0 => github.com/sirupsen/logrus v1.0.6
	github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20190717161051-705d9623b7c1+incompatible
)

require (
	github.com/aws/aws-sdk-go v1.30.0
	github.com/fatih/color v1.9.0 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	helm.sh/helm/v3 v3.2.4
	k8s.io/api v0.18.6
	k8s.io/cli-runtime v0.18.6
	k8s.io/client-go v0.18.6
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/aws-iam-authenticator v0.5.1
)
