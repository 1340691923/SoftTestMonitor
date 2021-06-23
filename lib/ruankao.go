package lib

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
)

//一些头信息
var headMap = map[string]string{
	"Referer":                   "https://bm.ruankao.org.cn/",
	"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36",
	"Upgrade-Insecure-Requests": "1",
	"Sec-Fetch-User":            "?1",
	"Sec-Fetch-Site":            "same-site",
	"Sec-Fetch-Mode":            "navigate",
	"Sec-Fetch-Dest":            "document",
	"sec-ch-ua-mobile":          "?0",
	"sec-ch-ua":                 `" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`,
}

//软考结构体
type RuanKao struct {
	cookies   []*http.Cookie
	cookieStr string
}
//软考考生信息结构体
type RuanKaoUser struct {
	Name            string `json:"name"`
	Idcard          string `json:"idcard"`
	ExaminationTime string `json:"examination_time"`
	CaptchaCode     string `json:"captcha_code"`
}

//设置cookie
func (this *RuanKao) SetWelComeCookie() (err error) {

	const WelcomeUri = "https://bm.ruankao.org.cn/sign/welcome"

	welcomeRes, err := httpReq(http.MethodGet, WelcomeUri, headMap, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("请求软考主站时异常,err:%v", err.Error()))
	}

	defer welcomeRes.Body.Close()

	this.cookies = welcomeRes.Cookies()
	sArr := []string{}
	for _, cookieStr := range welcomeRes.Header["Set-Cookie"] {
		tmpArr := strings.Split(cookieStr, ";")
		tmp := tmpArr[0]
		sArr = append(sArr, tmp)
	}
	this.cookieStr = strings.Join(sArr, ";")
	return nil
}

//获取所有考试时间
func (this *RuanKao) GetExaminationTimeList() (list []string, err error) {

	const ScorePageuri = "https://query.ruankao.org.cn/score"

	scorePageuriRes, err := httpReq(http.MethodGet, ScorePageuri, headMap, this.cookies)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取查询成绩页面时异常:%v", err.Error()))
	}

	defer scorePageuriRes.Body.Close()

	doc, err := goquery.NewDocumentFromReader(scorePageuriRes.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("解析查询成绩页面DOM时异常:%v", err.Error()))
	}
	doc.Find(".select > li").Each(func(i int, selection *goquery.Selection) {
		list = append(list, selection.Text())
	})
	return list, nil
}

//获取验证码
func (this *RuanKao) GetCaptchaCode(parseCaptchaInterface ParseCaptchaInterface) (captchaCode string, err error) {

	getCaptchaUri := fmt.Sprintf("https://query.ruankao.org.cn/score/captcha")

	captchaRes, err := httpReq(http.MethodGet, getCaptchaUri, nil, this.cookies)

	if err != nil {
		return "", errors.New(fmt.Sprintf("获取软考成绩查询页面验证码时异常:%v", err.Error()))
	}

	defer captchaRes.Body.Close()

	captchaResb, err := ioutil.ReadAll(captchaRes.Body)

	if err != nil {
		return "", errors.New(fmt.Sprintf("解析软考成绩查询页面验证码时异常:%v", err.Error()))
	}

	captchaBase64 := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(captchaResb))

	return parseCaptchaInterface.GetCaptchaRes(captchaBase64)
}

//验证验证码
func (this *RuanKao) VerifyCaptchaUri(captcha string) (err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	const VerifyCaptchaUri = "https://query.ruankao.org.cn/score/VerifyCaptcha"
	type ResModel struct {
		Flag   int           `json:"flag"`
		Msg    string        `json:"msg"`
		Data   []interface{} `json:"data"`
		Status int           `json:"status"`
	}
	var resModel ResModel
	values := url.Values{}
	values.Set("captcha", captcha)

	res, err := httpReq(http.MethodPost, VerifyCaptchaUri, nil, this.cookies, values)
	if err != nil {
		return errors.New(fmt.Sprintf("请求验证验证码接口时异常:%v", err.Error()))
	}
	defer res.Body.Close()
	s, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("解析验证验证码接口返回值时异常:%v", err.Error()))
	}

	err = json.Unmarshal(s, &resModel)

	if err != nil {
		return errors.New(fmt.Sprintf("JSON 解析验证验证码接口返回值时异常:%v", err.Error()))
	}

	if resModel.Flag != 1 {
		return errors.New(fmt.Sprintf("验证验证码时异常:flag:%v:err:%c", resModel.Flag, resModel.Msg))
	}
	return nil
}

//分数信息结构体
type ScoreRes struct {
	KSSJ string
	ZGMC string
	XM   string
	ZJH  string
	ZKZH string
	XWCJ string
	SWCJ string
}

//获取分数
func (this *RuanKao) GetScore(ruanKaoUser RuanKaoUser) (scoreRes *ScoreRes, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	const ScoreResultUri = "https://query.ruankao.org.cn/score/result"

	type ResModel struct {
		Flag   int         `json:"flag"`
		Data   interface{} `json:"data"`
		Msg    string      `json:"msg"`
		Status int         `json:"status"`
	}

	var resModel ResModel

	form := url.Values{}
	form.Set("stage", ruanKaoUser.ExaminationTime)
	form.Set("xm", ruanKaoUser.Name)
	form.Set("zjhm", ruanKaoUser.Idcard)
	form.Set("jym", ruanKaoUser.CaptchaCode)
	form.Set("select_type", "1")

	res, err := httpReq(http.MethodPost, ScoreResultUri, nil, this.cookies,form)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请求获取成绩接口时异常:%v", err.Error()))
	}
	defer res.Body.Close()
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("解析获取成绩接口返回值时异常:%v", err.Error()))
	}
	resBytes,err = UnescapeUnicode(resBytes)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("解析获取成绩接口返回值时异常:%v", err.Error()))
	}
	err = json.Unmarshal(resBytes, &resModel)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("JSON 解析获取成绩接口返回值时异常:%v", err.Error()))
	}

	if resModel.Flag != 1 {
		return nil, errors.New(fmt.Sprintf("解析获取成绩接口返回值时异常:flag:%v:err:%c", resModel.Flag, resModel.Msg))
	}

	m := resModel.Data.(map[string]interface{})

	scoreRes = &ScoreRes{
		KSSJ:  m["KSSJ"].(string),
		ZGMC: m["ZGMC"].(string),
		XM:    m["XM"].(string),
		ZJH:  m["ZJH"].(string),
		ZKZH: m["ZKZH"].(string),
		XWCJ: m["XWCJ"].(string),
		SWCJ: m["SWCJ"].(string),
	}

	return scoreRes, nil
}
