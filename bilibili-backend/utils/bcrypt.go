package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// 这是一个密码哈希加密函数，用来把用户输入的明文密码转换成不可逆的密文，存到数据库里。
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 检查用户密码是否一致的函数
func CheckPassword(password, hash string) bool {
	//bcrypt设计为验证的时候不需要对输入的密码进行加密才能验证，也就是内部处理后验证即可
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
