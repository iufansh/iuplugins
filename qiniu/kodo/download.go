package kodo

import (
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func DownloadPublic(domain, key string) string {
	return storage.MakePublicURL(domain, key)
}

func DownloadPrivate(accessKey, secretKey, domain, key string, expireSec int64) string {
	mac := qbox.NewMac(accessKey, secretKey)
	deadline := time.Now().Add(time.Second * time.Duration(expireSec)).Unix()
	return storage.MakePrivateURL(mac, domain, key, deadline)
}
