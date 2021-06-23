package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	jsoniter "github.com/json-iterator/go"
)

//六派
type Liupai struct {
	appsecret string
}

func NewLiuPai(config map[string]interface{}) *Liupai {
	liupai := new(Liupai)
	if _,has:=config["appsecret"];has{
		liupai.appsecret = config["appsecret"].(string)
	}
	return liupai
}

func(this *Liupai) GetCaptchaRes(captchaBase64 string)(captchaRes string,err error){
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	const LiupaiUri = "http://open.liupai.net/captcha/query"
	liupaiForm := url.Values{}

	liupaiForm.Add("appkey", this.appsecret)
	liupaiForm.Add("typeid", "1")
	liupaiForm.Add("minlen", "4")
	liupaiForm.Add("maxlen", "4")
	liupaiForm.Add("pic", captchaBase64)

	liupaiUriResp, err := http.PostForm(LiupaiUri, liupaiForm)

	if err != nil {
		return "",errors.New(fmt.Sprintf("六派请求失败:%v",err.Error()))
	}

	defer liupaiUriResp.Body.Close()

	liupaiUriRespB, err := ioutil.ReadAll(liupaiUriResp.Body)

	if err != nil {
		return "",errors.New(fmt.Sprintf("六派请求失败:%v",err.Error()))
	}

	type LiupaiResModel struct {
		Status int         `json:"status"`
		Msg    string      `json:"msg"`
		Result interface{} `json:"result"`
	}

	var liupaiResModel LiupaiResModel

	err = json.Unmarshal(liupaiUriRespB, &liupaiResModel)

	if err != nil {
		return "",errors.New(fmt.Sprintf("JSON 解析六派返回结果失败:%v,返回结果为%v",err.Error(),string(liupaiUriRespB)))
	}

	if liupaiResModel.Status != 200 {
		switch liupaiResModel.Result.(type) {
		case string:
			return "",errors.New(fmt.Sprintf(" 六派解析验证码失败 (Status:%v,Msg:%c,result:%c )", liupaiResModel.Status, liupaiResModel.Msg, liupaiResModel.Result.(string)))
		case []interface{}:
			return "",errors.New(fmt.Sprintf(" 六派解析验证码失败 (Status:%v,Msg:%c )", liupaiResModel.Status, liupaiResModel.Msg))
		}
		return
	}

	return liupaiResModel.Result.(map[string]interface{})["val"].(string),nil
}
