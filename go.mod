module github.com/LuckyG0ldfish/balancer

go 1.14

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/free5gc/UeauCommon v1.0.0
	github.com/free5gc/aper v1.0.1
	github.com/free5gc/nas v1.0.5
	github.com/free5gc/ngap v1.0.2
	github.com/free5gc/openapi v1.0.3
	github.com/free5gc/path_util v1.0.0
	github.com/free5gc/version v1.0.0
	github.com/google/uuid v1.1.2
	github.com/ishidawataru/sctp v0.0.0-00010101000000-000000000000
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli v1.22.5
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/ishidawataru/sctp => ./lib/sctp
