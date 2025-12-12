// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
	xhttp "github.com/zeromicro/x/http"
	"net/http"
	"strings"
	"zero-admin/pkg/convert"
)

type VerifyPermissionMiddleware struct {
	Casbin      *casbin.SyncedCachedEnforcer
	ExcludeUrls []string
}

func NewVerifyPermissionMiddleware(enforcer *casbin.SyncedCachedEnforcer, ExcludeUrls ...string) *VerifyPermissionMiddleware {
	return &VerifyPermissionMiddleware{
		Casbin:      enforcer,
		ExcludeUrls: ExcludeUrls,
	}
}

func (m *VerifyPermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 跳过无需校验的路径（如登录、健康检查）
		if m.skipPermissionCheck(r.RequestURI) {
			next(w, r)
			return
		}

		// 用户身份
		userId, userRoles, err := parseUserInfo(r)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, fmt.Errorf("用户身份解析失败: %v", err))
			return
		}

		// 解析请求资源和操作（适配Casbin的sub, obj, act）
		resource := parseResource(r) // 资源（如/api/v1/user）
		action := parseAction(r)     // 操作（如GET/POST/PUT/DELETE）

		// Casbin权限校验（sub=角色, obj=资源, act=操作）
		for _, role := range userRoles {
			if ok, _ := m.Casbin.Enforce(role, resource, action); !ok {
				logx.Infof("用户无权限: userId=%s, roles=%v, resource=%s, action=%s", userId, userRoles, resource, action)
				xhttp.JsonBaseResponseCtx(r.Context(), w, errors.New("用户无此操作权限"))
				return
			}
		}

		next(w, r)
	}
}

func (m *VerifyPermissionMiddleware) skipPermissionCheck(url string) bool {
	for _, excludeUrl := range m.ExcludeUrls {
		if excludeUrl == url {
			return true
		}
	}
	return false
}

func parseUserInfo(r *http.Request) (userId string, userRoles []string, err error) {
	userId = convert.ToString(r.Context().Value("uid"))
	userRoles = strings.Split(convert.ToString(r.Context().Value("roles")), ",")
	return userId, userRoles, nil
}

func parseResource(r *http.Request) string {
	return r.URL.Path
}

func parseAction(r *http.Request) string {
	switch r.Method {
	case http.MethodGet:
		return "READ"
	case http.MethodPost:
		return "CREATE"
	case http.MethodPut:
		return "UPDATE"
	case http.MethodDelete:
		return "DELETE"
	default:
		return ""
	}
}
