// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/pkg/convert"
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
			logc.Infof(r.Context(), "token已过期, uid=%d, uuid=%s, tokenID=%s", uid, uuid, tokenID)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
