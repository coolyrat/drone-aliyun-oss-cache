package main

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/urfave/cli.v1"
	"log"
)

const (
	endpoint   = "oss-cn-shenzhen.aliyuncs.com"
	bucketName = "drone-build-cache"
	repo       = ""
	debug      = true
)

var (
	gitVersion = "unknown"
	goVersion  = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "oss cache plugin"
	app.Version = gitVersion
	app.Flags = []cli.Flag{
		cli.StringFlag{},
	}

	//lg := initLogger()

	//client, err := oss.New(endpoint, ak, sk)
	//if err != nil {
	//	lg.Fatal("ossClient create error", zap.Error(err))
	//}
	//
	//bucket, err := client.Bucket(bucketName)
	//if err != nil {
	//	lg.Fatal("get bucketName error", zap.String("bucket", bucketName), zap.Error(err))
	//}

	//lg.Info("ok", zap.String("bucket", bucket.BucketName))
	fmt.Println("gitVersion: ", gitVersion, " goVersion: ", goVersion)
}

func initLogger() *zap.Logger {
	var lg *zap.Logger
	var err error
	if debug {
		lg, err = zap.NewDevelopment()
	} else {
		lg, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("new Logger failed: %v", err)
	}
	return lg
}
