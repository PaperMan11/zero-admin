// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/logic"
	"zero-admin/rpc/sys/client/scopeservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.UpdateMenuRequest) (resp *types.Menu, err error) {
	res, err := l.svcCtx.ScopeService.UpdateMenu(l.ctx, &scopeservice.UpdateMenuRequest{
		Menu: &scopeservice.Menu{
			Id:        req.ID,
			ParentId:  req.ParentId,
			MenuName:  req.MenuName,
			MenuType:  req.MenuType,
			Path:      req.Path,
			Component: req.Component,
			Redirect:  req.Redirect,
			Icon:      req.Icon,
			Sort:      req.Sort,
			NoCache:   req.NoCache,
			Affix:     req.Affix,
			External:  req.External,
			Hidden:    req.Hidden,
			Status:    req.Status,
			Remark:    req.Remark,
		},
		OperatorId: logic.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "更新菜单失败: %v", err)
		return nil, err
	}

	menu := logic.ConvertToTypesMenu(res)
	return &menu, nil
}
