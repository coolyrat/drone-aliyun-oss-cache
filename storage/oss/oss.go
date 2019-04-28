package oss

import (
	"drone-alicloud-oss/log"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/drone/drone-cache-lib/storage"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
)

type config struct {
	AK       string
	SK       string
	Endpoint string
	Bucket   string
}

type ossStorage struct {
	bucket *oss.Bucket
}

func NewStorage(cfg *config) storage.Storage {
	client, err := oss.New(cfg.Endpoint, cfg.AK, cfg.SK)
	if err != nil {
		log.Logger.Fatal("fatal new oss client", zap.Error(err))
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		log.Logger.Fatal("fatal use bucket", zap.String("bucket", cfg.Bucket), zap.Error(err))
	}

	return &ossStorage{
		bucket: bucket,
	}
}

func (s *ossStorage) Get(p string, dst io.Writer) error {
	body, err := s.bucket.GetObject(p)
	if err != nil {
		log.Logger.Error("error get object", zap.String("object", p), zap.Error(err))
		return err
	}
	defer body.Close()

	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Logger.Error("error read object", zap.String("object", p), zap.Error(err))
		return err
	}

	if _, err := dst.Write(b); err != nil {
		log.Logger.Error("error write object", zap.String("object", p), zap.Error(err))
		return err
	}

	return nil
}

func (s *ossStorage) Put(p string, src io.Reader) error {
	if err := s.bucket.PutObject(p, src); err != nil {
		log.Logger.Error("error put object", zap.String("object", p), zap.Error(err))
		return err
	}
	return nil
}

func (s *ossStorage) List(p string) ([]storage.FileEntry, error) {
	panic("implement me")
}

func (s *ossStorage) Delete(p string) error {
	if err := s.bucket.DeleteObject(p); err != nil {
		log.Logger.Error("error delete object", zap.String("object", p), zap.Error(err))
		return err
	}
	return nil
}
