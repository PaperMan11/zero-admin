// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/response/xerr"
)

type JwtExpireAuthMiddleware struct {
	redis      *redis.Redis
	localCache *collection.Cache
}

func NewJwtExpireAuthMiddleware(redisCli *redis.Redis, localCache *collection.Cache) *JwtExpireAuthMiddleware {
	return &JwtExpireAuthMiddleware{
		redis:      redisCli,
		localCache: localCache,
	}
}

func (m *JwtExpireAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := convert.ToInt64(convert.ToString(r.Context().Value("uid")))
		tokenID := convert.ToString(r.Context().Value("uuid"))
		key := utils.GetAccessTokenKey(uid)
		uuid, _ := m.localCache.Take(key, func() (any, error) {
			return m.redis.GetCtx(r.Context(), key)
		})
		if uuid != tokenID {
			httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
				"code": xerr.ErrorTokenExpired,
				"msg":  xerr.MapErrMsg(xerr.ErrorTokenExpired),
			})
			return
		}
		next(w, r)
	}
}
