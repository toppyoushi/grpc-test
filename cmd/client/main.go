package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	pb "github.com/toppyoushi/grpc-apis/pkg/helloworld"
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

func setLog(e fsnotify.Event) {
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
	flag.StringVarP(&configPath, "config", "c", "conf/conf.yaml", "specify a config file")
	flag.Parse()

	viper.SetConfigFile(configPath)
	viper.WatchConfig()
	viper.OnConfigChange(setLog)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	log.SetReportCaller(true)
	level, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		panic(err)
	}
	setLogFormatter(viper.GetString("log.format"))
	log.SetLevel(level)

}

func main() {
	ctx := context.Background()

	// requestID := uuid.New()
	host := viper.GetString("server.host")
	port := viper.GetInt("server.port")
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		log.Panic(err)
	}

	client := pb.NewGreeterClient(conn)

	req := &pb.HelloReq{
		Msg: "Hello,I am a grpc client",
	}

	sayHelloCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rsp, err := client.SayHello(sayHelloCtx, req)

	if err != nil {
		log.Fatal(err)
	}

	log.Info(rsp)
}
