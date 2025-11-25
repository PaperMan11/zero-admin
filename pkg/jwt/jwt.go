package jwt

import (
	"errors"
	"time"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/response/xerr"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

/*

┌─────────┐                ┌─────────┐                ┌─────────┐
│  前端   │                 │  后端    │                │  存储   │
└────┬────┘                └────┬────┘                └────┬────┘
     │                          │                         │
     │  1. 登录（账号密码）        │                         │
     │────────────────────────> │                         │
     │                          │                         │
     │  2. 返回 Access/Refresh   │                         │
     │<──────────────────────── │                         │
     │                          │                         │
     │  3. 调用接口（带 Access）  │                         │
     │────────────────────────>│                         │
     │                         │  4. 验证 Access         │
     │                         │<────────────────────────│
     │                         │                         │
     │  5. 返回数据（成功）      │                         │
     │<────────────────────────│                         │
     │                         │                         │
     │  6. Access 过期，调用接口失败 │                    │
     │────────────────────────>│                         │
     │                         │  7. 返回 401（过期）     │
     │<────────────────────────│                         │
     │                         │                         │
     │  8. 用 Refresh 换 Access  │                         │
     │────────────────────────>│                         │
     │                         │  9. 验证 Refresh        │
     │                         │<────────────────────────│
     │                         │                         │
     │  10. 返回新 Access       │                         │
     │<────────────────────────│                         │
     │                         │                         │
     │  11. 重试接口（带新 Access） │                    │
     │────────────────────────>│                         │
     │                         │                         │

*/

type CustomClaims struct {
	Uid  int64    `json:"uid"`
	Role []string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(Issuer string, userID int64, role []string, secretKey string, accessExpire int64) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		Uid:  userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(accessExpire) * time.Second)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseToken(tokenString, secretKey string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		logx.Error("ParseToken: ", err)
		return nil, err
	}

	if !token.Valid {
		return nil, xerr.NewErrCode(xerr.ErrorTokenExpire)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, xerr.NewErrCode(xerr.ErrorTokenInvalid)
	}

	return claims, nil
}

// 生成 Refresh Token
func GenerateRefreshToken(Issuer string, userID int64, refreshSecretKey string, refreshExpire int64) (string, error) {
	expireTime := time.Now().Add(time.Duration(refreshExpire) * time.Second)
	// Refresh Token 载荷可简化（仅需用户 ID 和过期时间）
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expireTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    Issuer,
		Subject:   convert.ToString(userID), // 存储用户 ID（字符串格式）
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshSecretKey))
}

func ParseRefreshToken(Issuer, tokenString, refreshSecretKey string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(refreshSecretKey), nil
		},
	)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token")
	}
	// 验证 Issuer（可选，增强安全性）
	if claims.Issuer != Issuer {
		return "", errors.New("invalid issuer")
	}
	return claims.Subject, nil // 返回用户 ID
}
