package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	jsoniter "github.com/json-iterator/go"
)

//配置
type Config struct {
	Year                      string `json:"year" validate:"required"`
	Name                      string `json:"name" validate:"required"`
	IdCard                    string `json:"id_card" validate:"required"`
	Appsecret                 string `json:"appsecret" validate:"required"`
	SendUser163MailAddress    string `json:"send_user_163_mail_address" validate:"required"`
	SendUser163MailAuthCode   string `json:"send_user_163_mail_auth_code" validate:"required"`
	ReceiveUser163MailAddress string `json:"receive_user_163_mail_address" validate:"required"`
	TimeInterval              int    `json:"time_interval"`
}

// LoadJSONConfig 读取配置文件 json格式
func LoadJSONConfig(filename string, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

//下载配置文件
func DownloadConfigFile(fname string)(err error){
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var config Config
	filePtr, err := os.Create(fname)
	if err != nil {
		return errors.New(fmt.Sprintf("创建配置文件异常:%s", err.Error()))
	}
	defer filePtr.Close()
	// 带JSON缩进格式写文件
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return errors.New(fmt.Sprintf("创建配置文件异常:%s", err.Error()))
	}
	_,err = filePtr.Write(data)
	return
}
