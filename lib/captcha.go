package lib

import "errors"

//解析验证码接口
type ParseCaptchaInterface interface {
	//通过验证码图片的base64得到验证码code
	GetCaptchaRes(captchaBase64 string)(captchaRes string,err error)
}

const LiuPaiCommand  = 1

//命令列表
var commandMap = map[int]func(config map[string]interface{}) *Liupai{
	LiuPaiCommand:NewLiuPai,
}
//根据命令释放对应的对象
func ParseCaptchaCommand(command int,config map[string]interface{})(ParseCaptchaInterface,error){
	var f func(config map[string]interface{}) *Liupai
	var has bool
	if f,has=commandMap[command];!has{
		return nil,errors.New("没找到该解析验证码命令")
	}
	return f(config),nil
}
