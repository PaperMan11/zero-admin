// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"github.com/ua-parser/uap-go/uaparser"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
	"zero-admin/api/admin/internal/logic/sys/auth"
	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterRequest
		if err := httpx.Parse(r, &req); err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}
		/*
			获取客户端真实的ip
			nginx的配置
			location /api {
					proxy_set_header Host $host;
					proxy_set_header X-Real-IP $remote_addr;
					proxy_set_header REMOTE-HOST $remote_addr;
					proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
					proxy_pass http://127.0.0.1:8888;
			}
		*/
		userAgent := r.Header.Get("User-Agent")
		parser := uaparser.NewFromSaved()
		ua := parser.Parse(userAgent)
		browser := ua.UserAgent.Family + " " + ua.UserAgent.Major
		os := ua.Os.Family + " " + ua.Os.Major
		l := auth.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req, httpx.GetRemoteAddr(r), os, browser)
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
