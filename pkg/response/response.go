package response

import (
	"net/http"
	"zero-admin/pkg/response/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

type Body struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body

	if err == nil {
		body.Code = xerr.OK
		body.Msg = xerr.MapErrMsg(xerr.OK)
		body.Data = resp
	} else {
		//错误返回
		errcode := xerr.ErrorServerCommon
		errmsg := "服务器开小差啦，稍后再来试一试"

		causeErr := errors.Cause(err) // err类型

		if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型
			//自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
				grpcCode := uint32(gstatus.Code())
				if xerr.IsCodeErr(grpcCode) { //区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
					errcode = grpcCode
					errmsg = gstatus.Message()
				} else { // 临时加上
					errcode = grpcCode
					errmsg = "grpc err:" + gstatus.Message()
				}
			}
		}

		logx.Errorf("【API-ERR】 : %+v ", err)
		body.Code = errcode
		body.Msg = errmsg
		body.Data = resp
	}
	httpx.OkJson(w, body)
}
