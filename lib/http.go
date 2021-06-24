package lib

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

//http请求方法
func httpReq(method string, uri string, headMap map[string]string, cookies []*http.Cookie, form ...url.Values) (*http.Response, error) {
	client := &http.Client{}
	if len(cookies) > 0 {
		url, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}
		jar, _ := cookiejar.New(nil)
		jar.SetCookies(url, cookies)
		client.Jar = jar
	}
	var req *http.Request
	switch method {
	case http.MethodPost:
		resp, err := client.PostForm(uri, form[0])
		if err != nil {
			return nil, err
		}
		return resp, nil
	case http.MethodGet:
		req, _ = http.NewRequest(method, uri, nil)
		for k, v := range headMap {
			req.Header.Add(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, nil
}

//unicode转中文
func UnescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
