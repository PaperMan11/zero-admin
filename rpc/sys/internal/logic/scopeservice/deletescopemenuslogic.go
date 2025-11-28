package scopeservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteScopeMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScopeMenusLogic {
	return &DeleteScopeMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteScopeMenusLogic) DeleteScopeMenus(in *sysclient.DeleteScopeMenusRequest) (*sysclient.ScopeInfo, error) {
	exists, err := l.svcCtx.DB.ExistsScope(l.ctx, in.ScopeId)
	if err != nil {
		logc.Errorf(l.ctx, "判断安全范围是否存在失败, scope id：%s, 错误：%s", in.ScopeId, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	if !exists {
		return nil, xerr.NewErrCode(xerr.ErrorScopeNotExist)
	}

	err = l.svcCtx.DB.DeleteScopeMenus(l.ctx, in.ScopeId)
	if err != nil {
		logc.Errorf(l.ctx, "删除安全范围菜单失败, scope id：%s, 错误：%s", in.ScopeId, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	return NewGetScopeMenusLogic(l.ctx, l.svcCtx).GetScopeMenus(&sysclient.Int64Value{Value: in.ScopeId})
}
