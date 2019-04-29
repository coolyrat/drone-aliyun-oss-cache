package oss

import (
	"drone-alicloud-oss/log"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/drone/drone-cache-lib/storage"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
)

type Config struct {
	AK       string
	SK       string
	Endpoint string
	Bucket   string
}

type ossStorage struct {
	bucket *oss.Bucket
}

func NewStorage(cfg *Config) storage.Storage {
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
	log.Logger.Info("oss getting", zap.String("path", p))
	body, err := s.bucket.GetObject(p)
	if err != nil {
		log.Logger.Error("error get object", zap.String("path", p), zap.Error(err))
		return err
	}
	defer body.Close()

	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Logger.Error("error read object", zap.String("path", p), zap.Error(err))
		return err
	}

	if _, err := dst.Write(b); err != nil {
		log.Logger.Error("error write object", zap.String("path", p), zap.Error(err))
		return err
	}

	return nil
}

func (s *ossStorage) Put(p string, src io.Reader) error {
	log.Logger.Info("oss putting", zap.String("path", p))
	if err := s.bucket.PutObject(p, src); err != nil {
		log.Logger.Error("error put object", zap.String("path", p), zap.Error(err))
		return err
	}
	return nil
}

func (s *ossStorage) List(p string) ([]storage.FileEntry, error) {
	log.Logger.Info("oss listing", zap.String("path", p))
	b, err := s.bucket.ListObjects(oss.Prefix(p))
	if err != nil {
		log.Logger.Error("error list object", zap.String("path", p), zap.Error(err))
		return nil, err
	}

	ee := make([]storage.FileEntry, b.MaxKeys)
	for i, o := range b.Objects {
		ee[i] = storage.FileEntry{
			Path:         o.Key,
			Size:         o.Size,
			LastModified: o.LastModified,
		}
	}
	log.Logger.Info("oss list objects", zap.Reflect("objects", ee))

	return ee, nil
}

func (s *ossStorage) Delete(p string) error {
	log.Logger.Info("oss deleting", zap.String("path", p))
	if err := s.bucket.DeleteObject(p); err != nil {
		log.Logger.Error("error delete object", zap.String("path", p), zap.Error(err))
		return err
	}
	return nil
}
