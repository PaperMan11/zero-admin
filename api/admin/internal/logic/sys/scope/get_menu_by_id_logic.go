// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/scopeservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuByIdLogic {
	return &GetMenuByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuByIdLogic) GetMenuById(req *types.IdValue) (resp *types.Menu, err error) {
	res, err := l.svcCtx.ScopeService.GetMenuById(l.ctx, &scopeservice.Int64Value{Value: req.Id})
	if err != nil {
		logc.Errorf(l.ctx, "查询菜单失败: %v", err)
		return nil, err
	}

	menu := utils.ConvertToTypesMenu(res)
	return &menu, nil
}
