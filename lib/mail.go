package lib

import (
	"fmt"
	"net/smtp"
	"strings"
)

const (
	// 邮件服务器地址
	SMTP_MAIL_HOST = "smtp.126.com"
	// 端口
	SMTP_MAIL_PORT = "25"
)

var (
	// 发送邮件用户账号 126邮箱
	SMTP_MAIL_USER = ""
	// 授权密码
	SMTP_MAIL_PWD = ""
)

//发送邮件
func SendSMTPMail(mailAddress string, subject string, body string) error {
	// 通常身份应该是空字符串，填充用户名.
	auth := smtp.PlainAuth("", SMTP_MAIL_USER, SMTP_MAIL_PWD, SMTP_MAIL_HOST)
	// (text/plain)纯文本 , (text/html)html
	contentType := "Content-Type: text/html; charset=UTF-8"
	nickname := "肖文龙(邮箱：1340691923@qq.com)"
	msg := []byte("To: " + mailAddress + "\r\nFrom: " + nickname + "<" + SMTP_MAIL_USER + ">\r\nSubject: " + subject +
		"\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(mailAddress, ",")
	err := smtp.SendMail(fmt.Sprintf("%s:%s", SMTP_MAIL_HOST, SMTP_MAIL_PORT), auth, SMTP_MAIL_USER, sendTo, msg)
	return err
}
