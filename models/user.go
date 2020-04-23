package models

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/delicb/gstring"
	jwt "github.com/dgrijalva/jwt-go"
	pg "github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"

	"github.com/Basic-Components/components_manager/connects"
	log "github.com/Basic-Components/components_manager/logger"
	script "github.com/Basic-Components/components_manager/script"
)

// ErrUserNotMatch 用户不匹配错误
var ErrUserNotMatch = errors.New("user not match")

// ErrISSNotMatch 签发机构不匹配错误
var ErrISSNotMatch = errors.New("iss not match")

// ErrUnexpectedSigningMethod 使用未知的签名算法
var ErrUnexpectedSigningMethod = errors.New("Unexpected signing method")

// User 用户类
type User struct {
	tableName struct{} `pg:"user"`
	ID        int
	Name      string `pg:",unique,notnull"`
	Email     string `pg:",unique,notnull"`
	Password  string `pg:",notnull"`
	Right     map[string]interface{}
	Verified  bool      `pg:",notnull,default:false"`
	Ctime     time.Time `pg:"default:now()"`
	Utime     time.Time `pg:"default:now()"`
}

/// 钩子

// BeforeInsert 插入数据时将密码先做md5序列化
func (instance *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	h := md5.New()
	h.Write([]byte(instance.Password))
	instance.Password = hex.EncodeToString(h.Sum(nil))
	return ctx, nil
}

// BeforeUpdate 更新数据时记录数据的更新时间
func (instance *User) BeforeUpdate(ctx context.Context) (context.Context, error) {
	now := time.Now()
	instance.Utime = now
	return ctx, nil
}

/// 自描述

// String 用户的简单描述
func (instance User) String() string {
	res := gstring.Sprintm(
		"<Id={Id:%d};Name={Name};Email={Email:%s}>",
		map[string]interface{}{
			"Id":    instance.ID,
			"Name":  instance.Name,
			"Email": instance.Email})
	return res
}

// Info 用户的简单描述
func (instance *User) Info() map[string]interface{} {
	res := map[string]interface{}{
		"Name":   instance.Name,
		"Email":  instance.Email,
		"Right":  instance.Right,
		"Verify": instance.Verified,
	}
	return res
}

// TokenInfo 用于构造token的用户描述
func (instance *User) TokenInfo() map[string]interface{} {
	res := map[string]interface{}{
		"ID":    instance.ID,
		"Right": instance.Right}
	return res
}

/// 用户校验

// Valid 验证用户密码
func (instance *User) Valid(inputpwd string) bool {
	h := md5.New()
	h.Write([]byte(inputpwd))
	inputpwdhash := hex.EncodeToString(h.Sum(nil))
	if inputpwdhash == instance.Password {
		return true
	}
	return false

}

///SetRight 用户新增权限
func (instance *User) SetRight() {

}

///RemoveRight 用户移除权限
func (instance *User) RemoveRight() {

}

///SetAdmin 用户设为管理员
func (instance *User) SetAdmin() {

}

///RemoveAdmin 用户移除管理员权限
func (instance *User) RemoveAdmin() {

}

/// JWT

// SignJWT 创建用户的JWT
func (instance *User) SignJWT(exp int64) (string, error) {
	var claims jwt.MapClaims = instance.TokenInfo()
	now := time.Now().Unix()
	claims["exp"] = now + exp
	claims["iat"] = now
	claims["iss"] = script.Config["component_name"]

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	out, err := token.SignedString([]byte(script.Config["secret"].(string)))
	if err == nil {
		return out, nil
	}
	return "", err
}

// CheckJWT 验证用户的JWT
func (instance *User) CheckJWT(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(script.Config["secret"].(string)), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["iss"].(string) != script.Config["component_name"] {
			return false, ErrISSNotMatch
		}
		if int(claims["ID"].(float64)) == instance.ID {
			return true, nil
		}
		return false, ErrUserNotMatch
	}
	return false, err
}

/// 发送邮件

// SendEmail 向用户发送邮件
func (instance *User) SendEmail(tp *connects.EmailMessageTemplate, kwargs map[string]interface{}) {
	err := connects.Email.Send(tp, instance.Email, kwargs)
	if err != nil {
		log.Warn(map[string]interface{}{"error": err}, "send error")
	}
}

// SendVerifyToken 向用户发送邮件用于验证创建成功
func (instance *User) SendVerifyToken(baseURL string) (string, error) {
	token, err := instance.SignJWT(int64(2 * 60 * 60))
	verifyURL := baseURL + "/auth/verify/" + token
	if err != nil {
		return "", err
	}
	if script.Config["use_email"].(bool) {
		kwargs := map[string]interface{}{
			"name":      instance.Name,
			"verifyURL": verifyURL}
		tp := &connects.EmailMessageTemplate{
			Subject:      "创建用户验证",
			HTMLTemplate: `<h2> hello <b>{{ name }}</b>, verify url is {{ verifyURL }} </h2>`,
		}
		go instance.SendEmail(tp, kwargs)
	}
	return token, nil
}

//UserNew 创建新的一般用户
func UserNew(name string, email string, password string, baseURL string) (string, error) {
	if connects.DB.Ok == false {
		return "", connects.ErrDBProxyNotInited
	}
	user := User{
		Name:     name,
		Email:    email,
		Password: password,
		Right:    map[string]interface{}{"admin": false}}
	err := connects.DB.Cli.Insert(&user)
	if err != nil {
		return "", err
	}
	token, err := user.SendVerifyToken(baseURL)
	if err != nil {
		return "", err
	}
	return token, nil
}

//UserVerify 校验注册但未校验的用户
func UserVerify(tokenString string) (bool, error) {
	if connects.DB.Ok == false {
		return false, connects.ErrDBProxyNotInited
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(script.Config["secret"].(string)), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["iss"].(string) == script.Config["component_name"] {
			user := &User{ID: int(claims["ID"].(float64))}
			err := connects.DB.Cli.Select(user)
			if err != nil {
				return false, ErrUserNotMatch
			}
			user.Verified = true
			err = connects.DB.Cli.Update(user)
			if err != nil {
				return false, err
			}
			return true, nil
		}
		return false, ErrISSNotMatch

	}
	return false, err
}

// userRegistCallback 需要注册到db代理上的回调
func userRegistCallback(db *pg.DB) error {
	err := db.CreateTable(
		new(User),
		&orm.CreateTableOptions{IfNotExists: true})
	if err != nil {
		log.Warn(map[string]interface{}{"error": err, "place": "CreateTable"}, "userRegistCallback get error")
		return err
	}
	var users []User
	count, err := db.Model(&users).Count()
	if err != nil {
		log.Warn(map[string]interface{}{"error": err, "place": "Count"}, "userRegistCallback get error")
		return err
	}
	if count == 0 {
		email := ""
		if connects.Email.Ok {
			email = connects.Email.Options.Username
		}
		admin := User{
			Name:     "admin",
			Email:    email,
			Password: "admin",
			Verified: true,
			Right:    map[string]interface{}{"admin": true}}
		err := db.Insert(&admin)
		if err != nil {
			return err
		}
	}
	return nil
}
