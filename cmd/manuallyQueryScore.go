package cmd

import (
	"bytes"
	"fmt"
	"sync"
	"github.com/go-playground/validator"
	"github.com/1340691923/SoftTestMonitor/lib"
	"github.com/spf13/cobra"
)

//用户输入的配置文件路径
var configFile string

//手动查询分数命令
var manuallyQueryScoreCmd = &cobra.Command{
	Use:   "manuallyQueryScore",
	Short: "手动查询分数",
	Long: `手动查询分数功能  
config.json 配置文件中所需参数： 
0.year->考试时间(建议查询之前先执行showExaminationTimeCmd 查看所有考试时间)
1.name->姓名
2.id_card->身份证
3.appsecret->六派数据的appsecret信息(用于解析验证码,一个手机号码可领取10次免费解析机会,六派数据网站地址为:https://www.6api.net/my/),解析验证码申请地址 https://www.6api.net/api/captcha/
4.send_user_163_mail_address->邮件发送者126用户账号
5.send_user_163_mail_auth_code->邮件发送者126邮箱授权码
6.receive_user_163_mail_address->邮件接收者邮箱地址
`,
	Run: runManuallyQueryScoreCmd,
}

func init() {
	manuallyQueryScoreCmd.Flags().StringVarP(&configFile, "configFile", "c", "config.json", "参数配置文件")
	manuallyQueryScoreCmd.MarkFlagRequired("configFile")
	rootCmd.AddCommand(manuallyQueryScoreCmd)
}

//运行手动查询分数命令
func runManuallyQueryScoreCmd(cmd *cobra.Command, args []string){
	//读取配置文件
	var config lib.Config
	err := lib.LoadJSONConfig(configFile,&config)
	if err!=nil{
		fmt.Println(err)
		return
	}
	//验证文件中是否有空
	v := validator.New()
	err = v.Struct(config)
	if err!=nil{
		fmt.Println(err)
		return
	}
	//设置软考网的cookie
	ruanKao := new(lib.RuanKao)
	err = ruanKao.SetWelComeCookie()
	if err!=nil{
		fmt.Println(err)
		return
	}
	//解析验证码 用六派命令 以后可能有别的命令
	fmt.Println("cookie获取成功")
	parseCaptcha,err := lib.ParseCaptchaCommand(lib.LiuPaiCommand, map[string]interface{}{"appsecret":config.Appsecret})

	if err!=nil{
		fmt.Println(err)
		return
	}
	//获取验证码
	captchaCode,err := ruanKao.GetCaptchaCode(parseCaptcha)
	if err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Println("验证码获取成功",captchaCode)
	//校验验证码
	err = ruanKao.VerifyCaptchaUri(captchaCode)
	if err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Println("验证码验证成功",captchaCode)
	//获取成绩
	scoreInfo ,err :=ruanKao.GetScore(lib.RuanKaoUser{
		Name:            config.Name,
		Idcard:         config.IdCard,
		ExaminationTime:config.Year,
		CaptchaCode:     captchaCode,
	})

	if err!=nil{
		fmt.Println(err)
		return
	}
	//发送邮件
	lib.SMTP_MAIL_USER = config.SendUser163MailAddress
	lib.SMTP_MAIL_PWD = config.SendUser163MailAuthCode
	//并发执行 输出与发送
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		var buffer bytes.Buffer
		buffer.WriteString("分数获取成功->")
		buffer.WriteString("考试时间：")

		buffer.WriteString("资格名称:")
		buffer.WriteString(scoreInfo.ZGMC)
		buffer.WriteString("准考证号:")
		buffer.WriteString(scoreInfo.ZKZH)
		buffer.WriteString("证件号:")
		buffer.WriteString(scoreInfo.ZJH)
		buffer.WriteString("姓名:")
		buffer.WriteString(scoreInfo.XM)
		buffer.WriteString("上午成绩:")
		buffer.WriteString(scoreInfo.SWCJ)
		buffer.WriteString("下午成绩:")
		buffer.WriteString(scoreInfo.XWCJ)
		err = lib.SendSMTPMail(config.ReceiveUser163MailAddress,"软考成绩已出",buffer.String())
		if err!=nil{
			fmt.Println(fmt.Sprintf("邮箱发送异常:%v",err.Error()))
			return
		}
	}()

	fmt.Println("分数获取成功----------->")
	fmt.Println("-------------------------")
	fmt.Println("考试时间:",scoreInfo.KSSJ)
	fmt.Println("-------------------------")
	fmt.Println("资格名称:",scoreInfo.ZGMC)
	fmt.Println("-------------------------")
	fmt.Println("准考证号:",scoreInfo.ZKZH)
	fmt.Println("-------------------------")
	fmt.Println("证件号:",scoreInfo.ZJH)
	fmt.Println("-------------------------")
	fmt.Println("姓名:",scoreInfo.XM)
	fmt.Println("-------------------------")
	fmt.Println("上午成绩:",scoreInfo.SWCJ)
	fmt.Println("-------------------------")
	fmt.Println("下午成绩:",scoreInfo.XWCJ)
	fmt.Println("-------------------------")
	wg.Wait()
}
