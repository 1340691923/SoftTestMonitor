package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/1340691923/SoftTestMonitor/lib"
	"github.com/go-playground/validator"
	"github.com/spf13/cobra"
)

//停止执行的开关
var stopC = make(chan int,1)

//监控软考平台cmd
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "监控软考平台，出成绩后将第一时间查询出成绩并通过邮件发送给您",
	Long: `监控软考平台，出成绩后将第一时间查询出成绩并通过邮件发送给您 config.json 配置文件中所需参数： 0.year->考试时间(建议查询之前先执行showExaminationTimeCmd 查看所有考试时间),
1.name->姓名，2.id_card->身份证,3.appsecret->六派数据的appsecret信息(用于解析验证码,一个手机号码可领取10次免费解析机会,六派数据网站地址为:https://www.6api.net/my/),解析验证码申请地址 https://www.6api.net/api/captcha/
4.send_user_163_mail_address->邮件发送者163用户账号,5.send_user_163_mail_auth_code->邮件发送者163邮箱授权码,6.receive_user_163_mail_address->邮件接收者邮箱地址,7.time_interval->轮询间隔时间,单位为分钟,最小为1(定时去软考网检测分数是否可查,一旦可查立马发邮件给您)`,
	Run: runMonitorCmd,
}

func init() {
	monitorCmd.Flags().StringVarP(&configFile, "configFile", "c", "config.json", "参数配置文件")
	monitorCmd.MarkFlagRequired("configFile")
	rootCmd.AddCommand(monitorCmd)
}

//详细注释见运行手动查询分数命令
func runMonitorCmd(cmd *cobra.Command, args []string) {
	var config lib.Config
	err := lib.LoadJSONConfig(configFile,&config)
	if err!=nil{
		fmt.Println(err)
		return
	}

	v := validator.New()
	err = v.Struct(config)
	if err!=nil{
		fmt.Println(err)
		return
	}
	if config.TimeInterval <1{
		fmt.Println("轮询时间最小为1分钟")
		return
	}

	parseCaptcha,err := lib.ParseCaptchaCommand(lib.LiuPaiCommand, map[string]interface{}{"appsecret":config.Appsecret})

	if err!=nil{
		fmt.Println(err)
		return
	}

	lib.SMTP_MAIL_USER = config.SendUser163MailAddress
	lib.SMTP_MAIL_PWD = config.SendUser163MailAuthCode

	go func() {
		select {
		case <-stopC:
			panic(errors.New(fmt.Sprintf("成绩已发送至您的邮箱!请检查邮件:%v",config.ReceiveUser163MailAddress)))
		}
	}()

	fmt.Println(fmt.Sprintf("开始轮询,轮询时间间隔为%v分钟",config.TimeInterval))
	for {
		time.Sleep(time.Minute * time.Duration(config.TimeInterval))

		ruanKao := new(lib.RuanKao)
		err := ruanKao.SetWelComeCookie()
		if err!=nil{
			fmt.Println(err)
			continue
		}

		list,err := ruanKao.GetExaminationTimeList()
		if err!=nil{
			fmt.Println(err)
			continue
		}
		fmt.Println("软考网站公开的所有考试时间如下:")
		for _,t := range list {
			fmt.Println(t)
		}
		if !lib.InArrayStr(config.Year,list){
			fmt.Println("还没到成绩公开时间,开始继续下一次的轮询")
			continue
		}
		captchaCode,err := ruanKao.GetCaptchaCode(parseCaptcha)
		if err!=nil{
			fmt.Println(err)
			continue
		}
		fmt.Println("验证码获取成功",captchaCode)
		err = ruanKao.VerifyCaptchaUri(captchaCode)
		if err!=nil{
			fmt.Println(err)
			continue
		}
		fmt.Println("验证码验证成功",captchaCode)
		scoreInfo ,err :=ruanKao.GetScore(lib.RuanKaoUser{
			Name:            config.Name,
			Idcard:         config.IdCard,
			ExaminationTime:config.Year,
			CaptchaCode:     captchaCode,
		})
		if err!=nil{
			fmt.Println(err)
			continue
		}

		go func() {
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
			stopC <- 1
		}()
		fmt.Println("分数获取成功----------->")
		fmt.Println("---------------------------------------")
		fmt.Println("考试时间:",scoreInfo.KSSJ)
		fmt.Println("---------------------------------------")
		fmt.Println("资格名称:",scoreInfo.ZGMC)
		fmt.Println("---------------------------------------")
		fmt.Println("准考证号:",scoreInfo.ZKZH)
		fmt.Println("---------------------------------------")
		fmt.Println("证件号:",scoreInfo.ZJH)
		fmt.Println("---------------------------------------")
		fmt.Println("姓名:",scoreInfo.XM)
		fmt.Println("---------------------------------------")
		fmt.Println("上午成绩:",scoreInfo.SWCJ)
		fmt.Println("---------------------------------------")
		fmt.Println("下午成绩:",scoreInfo.XWCJ)
		fmt.Println("---------------------------------------")
	}
}
