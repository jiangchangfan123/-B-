package utils

import (
	"errors"
	"time"

	"bilibili-backend/config"

	"github.com/golang-jwt/jwt/v5"
)

// 包装用户id，用户名和对应的jwt令牌
type Claims struct {
	UserID               uint64 `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        //嵌入标准说明
}

func GenerateToken(userID uint64, username string) (string, error) {
	//先设置过期时间
	exp := time.Duration(config.C.JWT.ExpiresIn) * time.Second
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{ //JWT标准字段
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)), //过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),          //签发时间
		},
	}
	/*
		根据jwt令牌生成token，jwt.SigningMethodHS256，这个是一个签名算法，
		根据claims这个组装的数据载荷进行计算
	*/
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//签名并生成字符串
	return token.SignedString([]byte(config.C.JWT.Secret))
}

// 解析令牌的函数
func ParseToken(tokenStr string) (*Claims, error) {
	//检查解析的令牌是否和浏览器保存的令牌一致
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.C.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无效 token")
}
