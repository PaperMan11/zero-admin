package utils

import "fmt"

const (
	// token
	AccessTokenPrefix  = "access_token:uid:"
	RefreshTokenPrefix = "refresh_token:uid:"
)

func GetAccessTokenKey(userId int64) string {
	return fmt.Sprintf("%s%d", AccessTokenPrefix, userId)
}
func GetRefreshTokenKey(userId int64) string {
	return fmt.Sprintf("%s%d", RefreshTokenPrefix, userId)
}
