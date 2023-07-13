package storage

import (
	"errors"
	"io"
	"time"

	"github.com/minio/minio-go/v6"
)

// New creates a new stager for putting payloads in object storage.
func New(c CompletedConfig) Stager {
    return &S3Stager{
        c,
    }
}

// Stager provides the mechanism to storage a payload
type Stager interface {
	Stage(*Input) (string, error)
    GetURL(string) (string, error)
}

// Stager provides the mechanism to stage a payload via AWS S3
type S3Stager struct {
    CompletedConfig
}

// Input contains data and metadata to be staged
type Input struct {
	Payload io.ReadCloser
	Key     string
	Account string
	OrgId   string
	Size    int64
}

// Close closes the underlying ReadCloser as long as it isn't nil
func (i *Input) Close() {
	if i.Payload != nil {
		i.Payload.Close()
	}
}

// Stage stores the file in s3 compatible storage and returns a presigned url
func (s *S3Stager) Stage(in *Input) (string, error) {
	bucketName := s.StageBucket
	objectName := in.Key
	object := in.Payload
	contentType := "application/gzip"

	_, err := s.Client.PutObject(bucketName,
		objectName,
		object,
		in.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
			UserMetadata: map[string]string{
				"requestID": in.Key,
				"account":   in.Account,
				"org":       in.OrgId,
			},
		},
	)
	if err != nil {
		return "", errors.New("Failed to upload to storage" + err.Error())
	}
	return s.GetURL(in.Key)
}

// GetURL retrieves a presigned url from s3 compatible storage
func (s *S3Stager) GetURL(requestID string) (string, error) {
	url, err := s.Client.PresignedGetObject(s.StageBucket, requestID, time.Second*24*60*60, nil)
	if err != nil {
		return "", errors.New("Failed to generate presigned url: " + err.Error())
	}
	return url.String(), nil
}
