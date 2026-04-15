package utils

import "fmt"

const (
	// token
	AccessTokenPrefix  = "access_token:uid:"
	RefreshTokenPrefix = "refresh_token:uid:"
)

func GetAccessTokenKey(userId int64, uuid string) string {
	return fmt.Sprintf("%s%d:%s", AccessTokenPrefix, userId, uuid)
}

func GetRefreshTokenKey(userId int64, uuid string) string {
	return fmt.Sprintf("%s%d:%s", RefreshTokenPrefix, userId, uuid)
}
