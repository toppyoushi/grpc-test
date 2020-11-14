package main

import (
	"fmt"

	"net"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	pb "github.com/toppyoushi/grpc-apis/pkg/helloworld"
	imp "github.com/toppyoushi/grpc-test/helloworld"
	"google.golang.org/grpc"
)

var (
	configPath string
)

func setLogFormatter(format string) {
	var formatter log.Formatter

	switch format {
	case "json":
		formatter = &log.JSONFormatter{}
	case "text":
		formatter = &log.TextFormatter{}
	default:
		formatter = &log.JSONFormatter{}
	}

	log.WithField("action", "set log format").Debug(format)
	log.SetFormatter(formatter)
}

func SetLog(e fsnotify.Event) {
	entry := log.WithField("action", fmt.Sprintf("config file %s %s", e.Name, e.Op))
	level, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		entry.Error(err)
		return
	}

	setLogFormatter(viper.GetString("log.format"))
	entry.WithField("action", "set log level").Debug(level)
	log.SetLevel(level)
}

func init() {
	flag.StringVarP(&configPath, "config", "c", "./conf/conf.yaml", "specify a config file")
	flag.Parse()

	viper.SetConfigFile(configPath)
	viper.WatchConfig()
	viper.OnConfigChange(SetLog)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	log.SetReportCaller(true)
	level, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)
	setLogFormatter(viper.GetString("log.format"))

}

func main() {
	server := grpc.NewServer()

	pb.RegisterGreeterServer(server, new(imp.GreeterServerImp))

	host, port := viper.GetString("server.host"), viper.GetInt("server.port")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		log.Fatal(err)
	}

	go func(e error) {
		err = server.Serve(lis)
	}(err)

	if err != nil {
		log.Fatal(err)
	}
}
