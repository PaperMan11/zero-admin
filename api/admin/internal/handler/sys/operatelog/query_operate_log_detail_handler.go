// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package operatelog

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
	"zero-admin/api/admin/internal/logic/sys/operatelog"
	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"
)

func QueryOperateLogDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryOperateLogDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := operatelog.NewQueryOperateLogDetailLogic(r.Context(), svcCtx)
		resp, err := l.QueryOperateLogDetail(&req)
		if err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			// httpx.OkJsonCtx(r.Context(), w, resp)
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
