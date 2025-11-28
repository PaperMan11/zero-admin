package scopeservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteScopeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScopeLogic {
	return &DeleteScopeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteScopeLogic) DeleteScope(in *sysclient.DeleteScopeRequest) (*sysclient.Empty, error) {
	exists, err := l.svcCtx.DB.ExistsScope(l.ctx, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "判断安全范围是否存在失败, scope id：%s, 错误：%s", in.Id, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	if !exists {
		return nil, xerr.NewErrCode(xerr.ErrorScopeNotExist)
	}

	menus, _ := l.svcCtx.DB.GetMenusByScopeID(l.ctx, in.Id)
	if len(menus) > 0 {
		return nil, xerr.NewErrMsg("该安全范围已关联菜单，请先解除关联关系")
	}

	err = l.svcCtx.DB.DeleteScope(l.ctx, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "删除安全范围失败, scope id：%s, 错误：%s", in.Id, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	return &sysclient.Empty{}, nil
}
