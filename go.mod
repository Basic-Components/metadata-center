module github.com/Basic-Components/components_manager

go 1.12

require (
	github.com/Basic-Components/connectproxy v0.0.0-20200428173545-00b34ef9d92e
	github.com/delicb/gstring v1.0.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/static v0.0.0-20191128031702-f81c604d8ac2
	github.com/gin-gonic/gin v1.6.2
	github.com/go-pg/pg/v9 v9.1.6
	github.com/json-iterator/go v1.1.9
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/sirupsen/logrus v1.5.0
	github.com/small-tk/pathlib v0.0.0-20190601032836-742166d9b695
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/toorop/gin-logrus v0.0.0-20190701131413-6c374ad36b67
	github.com/xeipuuv/gojsonschema v1.2.0
	go.etcd.io/etcd v3.3.20+incompatible
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0 // indirect
