package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"zero-admin/pkg/response/xerr"
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

type AccessClaims struct {
	Uid  int64    `json:"uid"`
	Role []string `json:"role"`
	Uuid string   `json:"uuid"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(uuid string, issuer string, userID int64, role []string, secretKey string, accessExpire int64) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		Uid:  userID,
		Role: role,
		Uuid: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(accessExpire) * time.Second)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseToken(tokenString, secretKey string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		logx.Error("ParseToken: ", err)
		return nil, err
	}

	if !token.Valid {
		return nil, xerr.NewErrCode(xerr.ErrorTokenExpire)
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok {
		return nil, xerr.NewErrCode(xerr.ErrorTokenInvalid)
	}

	return claims, nil
}

type RefreshClaims struct {
	Uuid string `json:"uuid"`
	jwt.RegisteredClaims
}

// 生成 Refresh Token
func GenerateRefreshToken(uuid string, issuer string, refreshSecretKey string, refreshExpire int64) (string, error) {
	// Refresh Token
	now := time.Now()
	claims := &RefreshClaims{
		Uuid: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(refreshExpire) * time.Second)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshSecretKey))
}

func ParseRefreshToken(tokenString, refreshSecretKey string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshSecretKey), nil
	})
	if err != nil {
		logx.Error("ParseToken: ", err)
		return nil, err
	}

	if !token.Valid {
		return nil, xerr.NewErrCode(xerr.ErrorTokenExpire)
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok {
		return nil, xerr.NewErrCode(xerr.ErrorTokenInvalid)
	}

	return claims, nil
}
