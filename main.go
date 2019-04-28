package main

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
	"log"
)

const (
	endpoint   = "oss-cn-shenzhen.aliyuncs.com"
	ak         = "LTAIeqrJGFbG7IBh"
	sk         = "3Vk6Cg7iKLov7MlcliSqYncLtEvk92"
	bucketName = "drone-build-cache"
	repo       = ""
	debug      = true
)

func main() {
	lg := initLogger()

	client, err := oss.New(endpoint, ak, sk)
	if err != nil {
		lg.Fatal("ossClient create error", zap.Error(err))
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		lg.Fatal("get bucketName error", zap.String("bucket", bucketName), zap.Error(err))
	}

	lg.Info("ok", zap.String("bucket", bucket.BucketName))
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
