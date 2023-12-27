package utls

import (
	"crypto/md5"
	"fmt"
)

// Md5Encryption 进行MD5加密
func Md5Encryption(password string) string {
	data := []byte(password)
	hash := md5.Sum(data)
	res := fmt.Sprintf("%x", hash)
	return res
}

// CheckPassword 检查密码是否通过
func CheckPassword(password, truePassword string) bool {
	return Md5Encryption(password) == truePassword
}
