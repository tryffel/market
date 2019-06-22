package storage

import (
	"errors"
	"fmt"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/market/config"
	"github.com/tryffel/market/modules/Error"
	"io"
)

type Minio struct {
	client    *minio.Client
	endpoint  string
	accessKey string
	secretKey string
	bucket    string
}

// Create new Minio instance
func NewMinio(config *config.Config) *Minio {

	m := &Minio{
		endpoint:  config.Minio.Url,
		accessKey: config.Minio.AccessKey,
		secretKey: config.Minio.SecretKey,
		bucket:    config.Minio.Bucket,
	}

	client, err := minio.New(m.endpoint, m.accessKey, m.secretKey, false)
	if err != nil {
		logrus.Fatal("Failed to initialize minio client: ", err)
		panic(err)
	}

	m.client = client
	bucket, err := m.client.BucketExists(m.bucket)
	if err != nil {
		logrus.Fatal("Failed to list minio buckets ", err)
		panic(err)
	}

	if bucket == false {
		logrus.Debug("Creating new minio bucket")
		err = m.client.MakeBucket(m.bucket, "")
		if err != nil {
			logrus.Fatal("Failed to create minio bucket", err)
			panic(err)
		}
	}

	return m
}

// Get photo from bucket
func (m *Minio) GetFile(name string) (io.Reader, error) {
	object, err := m.client.GetObject(m.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		err = Error.Wrap(&err, "failed to retrieve file from minio")
	}
	return object, err
}

// Put photo in bucket, return true if successful
func (m *Minio) PutFile(name string, file io.Reader, size int64) error {
	written, err := m.client.PutObject(m.bucket, name, file, size, minio.PutObjectOptions{})
	if err != nil {
		err = Error.Wrap(&err, "Failed to put file to minio")
	}
	if written != size {
		err = errors.New("minio didn't return same size as input!")
	}
	return err
}

func (m *Minio) FileExists(name string) bool {
	return true
}

func fullPath(file string) string {
	r := []rune(file)
	top := string(r[0:2])
	bottom := string(r[2:4])
	path := fmt.Sprintf("%s/%s/%s", top, bottom, file)
	return path
}
