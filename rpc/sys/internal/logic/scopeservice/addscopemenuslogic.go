package scopeservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddScopeMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddScopeMenusLogic {
	return &AddScopeMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddScopeMenusLogic) AddScopeMenus(in *sysclient.AddScopeMenusRequest) (*sysclient.ScopeInfo, error) {
	exists, err := l.svcCtx.DB.ExistsScope(l.ctx, in.ScopeId)
	if err != nil {
		logc.Errorf(l.ctx, "查询安全范围失败, 安全范围ID：%d, 错误：%s", in.ScopeId, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	if !exists {
		return nil, xerr.NewErrCode(xerr.ErrorScopeNotExist)
	}

	err = l.svcCtx.DB.AddScopeMenusTx(l.ctx, in.ScopeId, in.MenuIds)
	if err != nil {
		logc.Errorf(l.ctx, "添加安全范围菜单失败, 安全范围ID：%d, 错误：%s", in.ScopeId, err.Error())
		return nil, xerr.NewErrMsg("添加安全范围菜单失败")
	}
	return NewGetScopeMenusLogic(l.ctx, l.svcCtx).GetScopeMenus(&sysclient.Int64Value{Value: in.ScopeId})
}
