// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
	"net/http"
	"zero-admin/api/admin/internal/logic/sys/auth"
	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"
	"zero-admin/pkg/response/xerr"
)

func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshTokenRequest
		if err := httpx.Parse(r, &req); err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := auth.NewRefreshTokenLogic(r.Context(), svcCtx)
		resp, err := l.RefreshToken(&req)
		if err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			// code-data 响应格式
			if errors.Is(err, auth.ErrRefreshTokenExpired) {
				//xhttp.JsonBaseResponseCtx(r.Context(), w, err)
				httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
					"code": xerr.ErrorRefreshTokenExpired,
					"msg":  xerr.MapErrMsg(xerr.ErrorRefreshTokenExpired),
				})
			} else {
				httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
					"code": xerr.ErrorTokenInvalid,
					"msg":  xerr.MapErrMsg(xerr.ErrorTokenInvalid),
				})
			}
		} else {
			// httpx.OkJsonCtx(r.Context(), w, resp)
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
