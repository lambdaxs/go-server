package oss

import (
	"bytes"
	"context"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

type QiniuClient struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

func NewQiniuClient(accessKey, secretKey string) *QiniuClient {
	return &QiniuClient{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

func (q *QiniuClient) GetToken(bucket string) (token string) {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	putPolicy.Expires = 3600
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	token = putPolicy.UploadToken(mac)
	return
}

func (q *QiniuClient) UploadFile(bucket string, key string, mimeType string, body []byte) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		MimeType: mimeType,
	}
	if err := formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(body), int64(len(body)), &putExtra); err != nil {
		return "", err
	}
	return ret.Key, nil
}
