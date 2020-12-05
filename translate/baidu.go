package translate

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//申请的信息
//var appID = ""
//var apiSecret = ""

//百度翻译api接口
var Url = "http://api.fanyi.baidu.com/api/trans/vip/translate"

type apiParam struct {
	Q     string
	From  string
	To    string
	Appid string
	Salt  string
	Sign  string
	Tts   string
	Dict  string
}

func newTranslateModeler(param TransParam) apiParam {
	if param.From == "" {
		param.From = "auto"
	}
	if param.To == "" {
		param.To = "zh"
	}
	tran := apiParam{
		Q:    param.Query,
		From: param.From,
		To:   param.To,
		Tts:  "0",
		Dict: "0",
	}
	tran.Appid = param.Appid
	tran.Salt = strconv.Itoa(time.Now().Second())
	content := param.Appid + param.Query + tran.Salt + param.ApiSecret
	sign := sumString(content) //计算sign值
	tran.Sign = sign
	return tran
}

func (tran apiParam) toValues() url.Values {
	values := url.Values{
		"q":     {tran.Q},
		"from":  {tran.From},
		"to":    {tran.To},
		"appid": {tran.Appid},
		"salt":  {tran.Salt},
		"sign":  {tran.Sign},
		"tts":   {tran.Tts},
		"dict":  {tran.Dict},
	}
	return values
}

//计算文本的md5值
func sumString(content string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(content))
	bys := md5Ctx.Sum(nil)
	value := hex.EncodeToString(bys)
	return value
}

func Baidu(param TransParam) (transResult TransResult, err error) {
	tran := newTranslateModeler(param)
	values := tran.toValues()
	resp, err := http.PostForm(Url, values)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// txt := string(body)
	// fmt.Println(txt)
	var baiduResult baiduResult
	if err = json.Unmarshal(body, &baiduResult); err != nil {
		return
	}
	if baiduResult.ErrorCode != "" && baiduResult.ErrorCode != "52000" {
		err = errors.New(baiduResult.ErrorCode + "-" + baiduResult.ErrorMsg)
		return
	}
	if len(baiduResult.TransResult) == 0 {
		err = errors.New("No trans_result found")
		return
	}
	transResult = TransResult{
		From:   baiduResult.From,
		To:     baiduResult.To,
		Src:    baiduResult.TransResult[0].Src,
		Dst:    baiduResult.TransResult[0].Dst,
		SrcTts: baiduResult.TransResult[0].SrcTts,
		DstTts: baiduResult.TransResult[0].DstTts,
	}
	var bdDict baiduDict
	if baiduResult.TransResult[0].Dict != "" {
		if err = json.Unmarshal([]byte(baiduResult.TransResult[0].Dict), &bdDict); err != nil {
			// 格式不对，翻译失败
			//fmt.Println("Baidu translate err0:", err)
			err = nil
		}
		if bdDict.Lang == "1" {
			var enDict baiduEnDict
			if err = json.Unmarshal([]byte(baiduResult.TransResult[0].Dict), &enDict); err != nil {
				// 格式不对，翻译失败
				//fmt.Println("Baidu translate err1:", err)
				err = nil
			}
			transResult.MeansEn = enDict.WordResult.SimpleMeans
		} else if bdDict.Lang == "0" {
			var zhDict baiduZhDict
			if err = json.Unmarshal([]byte(baiduResult.TransResult[0].Dict), &zhDict); err != nil {
				// 格式不对，翻译失败
				//fmt.Println("Baidu translate err2:", err)
				err = nil
			}
			transResult.MeansZh = zhDict.WordResult.SimpleMeans
		}
	}
	if bdDict.Lang == "" {
		if transResult.From == "zh" {
			transResult.Lang = "0"
		} else {
			transResult.Lang = "1"
		}
	} else {
		transResult.Lang = bdDict.Lang
	}

	return
}

type TransParam struct {
	Appid     string
	ApiSecret string
	Query     string
	From      string
	To        string
}

type baiduResult struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src    string `json:"src"`
		Dst    string `json:"dst"`
		SrcTts string `json:"src_tts"`
		DstTts string `json:"dst_tts"`
		Dict   string `json:"dict"`
	} `json:"trans_result"`
	ErrorCode string `json:"error_code"` // 52000：成功，其他：失败
	ErrorMsg  string `json:"error_msg"`
}

type baiduDict struct {
	Lang string `json:"lang"`
}

type baiduEnDict struct {
	Lang       string `json:"lang"`
	WordResult struct {
		SimpleMeans SimpleMeansEn `json:"simple_means"`
	} `json:"word_result"`
}

type baiduZhDict struct {
	Lang       string `json:"lang"`
	WordResult struct {
		SimpleMeans SimpleMeansZh `json:"simple_means"`
	} `json:"word_result"`
}

type SimpleMeansEn struct {
	WordName  string   `json:"word_name"`
	From      string   `json:"from"`
	WordMeans []string `json:"word_means,omitempty"`
	Exchange  struct {
		WordThird []string `json:"word_third,omitempty"` // 第三人称单数
		WordIng   []string `json:"word_ing,omitempty"`   // 进行时态
		WordDone  []string `json:"word_done,omitempty"`  // 完成时态
		WordPast  []string `json:"word_past,omitempty"`  // 过去时态
		WordPl    []string `json:"word_pl,omitempty"`    // 复数形式
	} `json:"exchange,omitempty"`
	Tags struct {
		Core  []string `json:"core,omitempty"`
		Other []string `json:"other,omitempty"`
	} `json:"tags,omitempty"`
	Symbols []struct {
		PhEn  string `json:"ph_en,omitempty"`
		PhAm  string `json:"ph_am,omitempty"`
		Parts []struct {
			Part  string   `json:"part,omitempty"`
			Means []string `json:"means,omitempty"`
		} `json:"parts,omitempty"`
	} `json:"symbols,omitempty"`
}

type SimpleMeansZh struct {
	WordName  string   `json:"word_name"`
	From      string   `json:"from"`
	WordMeans []string `json:"word_means,omitempty"`
	Symbols   []struct {
		WordSymbol string `json:"word_symbol,omitempty"`
		Parts      []struct {
			PartName string `json:"part_name,omitempty"`
			Means    []struct {
				Text     string   `json:"text,omitempty"`
				Part     string   `json:"part,omitempty"`
				WordMean string   `json:"word_mean,omitempty"`
				Means    []string `json:"means,omitempty"`
			} `json:"means,omitempty"`
		} `json:"parts,omitempty"`
	} `json:"symbols,omitempty"`
}

type TransResult struct {
	From    string        `json:"from"`
	To      string        `json:"to"`
	Src     string        `json:"src"`
	Dst     string        `json:"dst"`
	SrcTts  string        `json:"src_tts"`
	DstTts  string        `json:"dst_tts"`
	Lang    string        `json:"lang"` // 0: 中文解释；1：英语解释
	MeansEn SimpleMeansEn `json:"means_en,omitempty"` // from=en时有值
	MeansZh SimpleMeansZh `json:"means_zh,omitempty"` // from=zh时有值
}
