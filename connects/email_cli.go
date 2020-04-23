package connects

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	gomail "gopkg.in/gomail.v2"
)

// ErrEmailProxyNotInited Email代理未初始化错误
var ErrEmailProxyNotInited = errors.New("email proxy not inited yet")

// ErrEmailURLSchemaWrong Email代理解析配置URL时Schema错误
var ErrEmailURLSchemaWrong = errors.New("schema must be email")

// ErrEmailURLUsernameEmpty Email代理解析配置URL时没有username
var ErrEmailURLUsernameEmpty = errors.New("email need to set a username")

// ErrEmailURLPasswordEmpty Email代理解析配置URL时没有password
var ErrEmailURLPasswordEmpty = errors.New("email need to set a password")

// ErrEmailURLHostEmpty Email代理解析配置URL时没有host
var ErrEmailURLHostEmpty = errors.New("email need to set host")

// ErrEmailURLPortWrong Email代理解析配置URL时port错误
var ErrEmailURLPortWrong = errors.New("email need to set int port")

// EmailProxyCallback 邮件服务代理操作的回调函数
type emailProxyCallback func(emailCli *gomail.Dialer) error

// EmailOptions 邮件服务的配置参数
type emailOptions struct {
	Host     string
	Port     int
	Username string
	Password string
	SSL      bool
	TLS      bool
}

// EmailProxy 邮件服务客户端的代理
type emailProxy struct {
	Ok        bool
	Options   *emailOptions
	Cli       *gomail.Dialer
	callBacks []emailProxyCallback
}

// EmailMessageTemplate 邮件服务发送邮件的模板信息
type EmailMessageTemplate struct {
	Subject       string
	PlainTemplate string
	HTMLTemplate  string
	Embed         []string
	Cc            []string
	Attach        []string
}

// NewEmailProxy 创建新的EmailProxy实例
func NewEmailProxy() *emailProxy {
	proxy := new(emailProxy)
	proxy.Ok = false
	return proxy
}

// 将url解析为email的初始化参数
func parseEmailURL(address string) (*emailOptions, error) {
	result := &emailOptions{}
	u, err := url.Parse(address)
	if err != nil {
		return result, err
	}
	if u.Scheme != "email" {
		return result, ErrEmailURLSchemaWrong
	}

	user := u.User.Username()
	if user == "" {
		return result, ErrEmailURLUsernameEmpty
	}
	result.Username = user

	password, has := u.User.Password()
	if has == false {
		return result, ErrEmailURLPasswordEmpty
	}
	result.Password = password
	hostinfo := strings.Split(u.Host, ":")
	if len(hostinfo) != 2 {
		return result, ErrEmailURLHostEmpty
	}
	result.Host = hostinfo[0]
	number, err := strconv.Atoi(hostinfo[1])
	if err != nil {
		return result, ErrEmailURLPortWrong
	}
	result.Port = number

	if u.RawQuery != "" {
		v, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			log.Fatal(err)
			return result, nil
		}
		ssl, sslok := v["ssl"]
		if sslok {
			bv := strings.ToLower(ssl[0])
			switch bv {
			case "true":
				result.SSL = true
			default:
				result.SSL = false

			}
		} else {
			result.SSL = false
		}

		tls, tlsok := v["tls"]
		if tlsok {
			bv := strings.ToLower(tls[0])
			switch bv {
			case "true":
				result.TLS = true
			default:
				result.TLS = false
			}
		} else {
			result.TLS = false
		}
	} else {
		result.SSL = false
		result.TLS = false
	}
	return result, nil
}

// Init 使用配置给代理赋值客户端实例
func (proxy *emailProxy) Init(options *emailOptions) error {
	d := gomail.NewDialer(
		options.Host,
		options.Port,
		options.Username,
		options.Password)

	if options.SSL {
		d.SSL = options.SSL
	}
	if options.TLS {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	proxy.Cli = d
	proxy.Options = options
	for _, cb := range proxy.callBacks {
		//go cb(proxy.Cli)
		cb(proxy.Cli)
	}
	proxy.Ok = true
	return nil
}

// InitFromURL 使用配置给代理赋值客户端实例
func (proxy *emailProxy) InitFromURL(address string) error {
	options, err := parseEmailURL(address)
	if err != nil {
		return err
	}
	err = proxy.Init(options)
	return err
}

// Regist 注册回调函数,在init执行后执行回调函数
func (proxy *emailProxy) Regist(cb emailProxyCallback) {
	proxy.callBacks = append(proxy.callBacks, cb)
}

// ExecuteTemplate 用模板构造字符串
func ExecuteTemplate(tp string, kwargs map[string]interface{}) (string, error) {
	t, err := template.New("").Parse(tp)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, kwargs)
	if err != nil {
		return "", err
	}
	res := buf.String()
	fmt.Println(res)
	return res, nil
}

// Send 发送邮件
func (proxy *emailProxy) Send(tp *EmailMessageTemplate, to string, kwargs map[string]interface{}) error {
	if proxy.Ok == false {
		return ErrEmailProxyNotInited
	}
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(proxy.Options.Username, "[Info]"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", tp.Subject)
	if len(tp.Cc) != 0 {
		m.SetHeader("Cc", tp.Cc...)
	}
	if tp.PlainTemplate != "" {
		s, err := ExecuteTemplate(tp.PlainTemplate, kwargs)
		if err != nil {
			return err
		}
		m.SetBody("text/plain", s)
	}
	if tp.HTMLTemplate != "" {
		s, err := ExecuteTemplate(tp.HTMLTemplate, kwargs)
		if err != nil {
			return err
		}
		m.SetBody("text/html", s)
		if len(tp.Embed) != 0 {
			for _, emb := range tp.Embed {
				m.Embed(emb)
			}
		}
	}
	if len(tp.Attach) != 0 {
		for _, attach := range tp.Attach {
			m.Attach(attach)
		}
	}
	err := proxy.Cli.DialAndSend(m)
	return err
}

// SendTest 尝试发送邮件
func (proxy *emailProxy) SendTest(to string) error {

	tp := &EmailMessageTemplate{
		Subject: "测试邮件",
		//PlainTemplate: `hello here is {{ .username }}`,
		HTMLTemplate: `<h2> hello here is  <b>{{ .username }}</b> </h2>`,
	}
	kwargs := map[string]interface{}{
		"username": proxy.Options.Username}
	err := proxy.Send(tp, to, kwargs)
	return err
}

// Publish 向固定群体广播发送邮件
func (proxy *emailProxy) Publish(msg *EmailMessageTemplate, togroup map[string]map[string]interface{}) {
	for to, kwargs := range togroup {
		go proxy.Send(msg, to, kwargs)
	}
}

// Email 默认的etcd代理对象
var Email = NewEmailProxy()
