package kodo

import (
	"bytes"
	"context"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type QiniuKodoUpload struct {
	Zone      string // 华东-ZoneHuadong;华北-ZoneHuabei;华南-ZoneHuanan;北美-ZoneBeimei;新加坡-ZoneXinjiapo
	UseHttps  bool
	Bucket    string
	AccessKey string
	SecretKey string
}

type QiniuKodoPutRet struct {
	Key    string `json:"key"`
	Hash   string `json:"hash"`
	Fsize  int    `json:"fsize"`
	Bucket string `json:"bucket"`
}

func (s QiniuKodoUpload) getZone() *storage.Region {
	switch s.Zone {
	case "ZoneHuanan":
		return &storage.ZoneHuanan
	case "ZoneHuadong":
		return &storage.ZoneHuadong
	case "ZoneHuabei":
		return &storage.ZoneHuabei
	case "ZoneBeimei":
		return &storage.ZoneBeimei
	case "ZoneXinjiapo":
		return &storage.ZoneXinjiapo
	default:
		break
	}
	return nil
}

var returnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)"}`

func (s QiniuKodoUpload) ByFormFile(context context.Context, localFile string, key string, putExtraParams map[string]string, insertOnly uint16) (ret QiniuKodoPutRet, err error) {
	putPolicy := storage.PutPolicy{
		ReturnBody: returnBody,
		InsertOnly: insertOnly,
	}
	if insertOnly == 0 {
		putPolicy.Scope = s.Bucket + ":" + key
	}
	mac := qbox.NewMac(s.AccessKey, s.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = s.getZone()
	// 是否使用https域名
	cfg.UseHTTPS = s.UseHttps
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	// 可选配置
	putExtra := storage.PutExtra{}
	if putExtraParams != nil {
		params1 := make(map[string]string)
		for k, v := range putExtraParams {
			params1["x:"+k] = v
		}
		putExtra.Params = params1
	}

	err = formUploader.PutFile(context, &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		// fmt.Println(err)
		return
	}
	// fmt.Println(ret.Key, ret.Hash)
	return
}

func (s QiniuKodoUpload) ByFormBytes(context context.Context, data []byte, key string, putExtraParams map[string]string, insertOnly uint16) (ret QiniuKodoPutRet, err error) {
	putPolicy := storage.PutPolicy{
		ReturnBody: returnBody,
		InsertOnly: insertOnly,
	}
	if insertOnly == 0 {
		putPolicy.Scope = s.Bucket + ":" + key
	}
	mac := qbox.NewMac(s.AccessKey, s.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = s.getZone()
	// 是否使用https域名
	cfg.UseHTTPS = s.UseHttps
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	// 可选配置
	putExtra := storage.PutExtra{}
	if putExtraParams != nil {
		params1 := make(map[string]string)
		for k, v := range putExtraParams {
			params1["x:"+k] = v
		}
		putExtra.Params = params1
	}

	err = formUploader.Put(context, &ret, upToken, key, bytes.NewReader(data), int64(len(data)), &putExtra)
	if err != nil {
		// fmt.Println(err)
		return
	}
	// fmt.Println(ret.Key, ret.Hash)
	return
}
