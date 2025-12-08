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

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuRequest) (resp *types.Menu, err error) {
	res, err := l.svcCtx.ScopeService.CreateMenu(l.ctx, &scopeservice.CreateMenuRequest{
		Menu: &scopeservice.Menu{
			Id:        req.ScopeId,
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
			ScopeId:   req.ScopeId,
			Remark:    req.Remark,
		},
		OperatorId: utils.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "创建菜单失败: %v", err)
		return nil, err
	}

	menu := utils.ConvertToTypesMenu(res)
	return &menu, nil
}
