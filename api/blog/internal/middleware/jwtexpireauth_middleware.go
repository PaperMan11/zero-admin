// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"context"
	"net/http"
	"time"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/errorx"
	"zero-admin/pkg/utils"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type JwtExpireAuthMiddleware struct {
	localCache *collection.Cache
	redis      *redis.Redis
}

func NewJwtExpireAuthMiddleware(redisCli *redis.Redis, localCache *collection.Cache) *JwtExpireAuthMiddleware {
	return &JwtExpireAuthMiddleware{
		localCache: localCache,
		redis:      redisCli,
	}
}

func (m *JwtExpireAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := convert.ToInt64(convert.ToString(r.Context().Value("uid")))
		tokenID := convert.ToString(r.Context().Value("uuid"))
		key := utils.GetAccessTokenKey(uid, tokenID)
		expiredAt, _ := m.localCache.Take(key, func() (any, error) {
			return m.redis.GetCtx(r.Context(), key)
		})
		if convert.ToInt64(expiredAt) < time.Now().Unix() {
			httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
				"code": -1,
				"msg":  errorx.ErrorTokenExpired,
			})
			return
		}
		next(w, r)
	}
}

func GetUidFromContext(ctx context.Context) int64 {
	return convert.ToInt64(ctx.Value("uid"))
}

func GetTokenIDFromContext(ctx context.Context) string {
	return convert.ToString(ctx.Value("uuid"))
}
