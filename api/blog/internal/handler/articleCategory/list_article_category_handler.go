// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package articleCategory

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
	"zero-admin/api/blog/internal/logic/articleCategory"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
)

func ListArticleCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListArticleCategoryReq
		if err := httpx.Parse(r, &req); err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := articleCategory.NewListArticleCategoryLogic(r.Context(), svcCtx)
		resp, err := l.ListArticleCategory(&req)
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
