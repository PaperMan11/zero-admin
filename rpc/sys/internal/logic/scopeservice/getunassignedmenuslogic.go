package scopeservicelogic

import (
	"context"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnassignedMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUnassignedMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnassignedMenusLogic {
	return &GetUnassignedMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUnassignedMenusLogic) GetUnassignedMenus(in *sysclient.Empty) (*sysclient.UnassignedMenusResponse, error) {
	menus, _ := l.svcCtx.DB.GetUnassignedMenus(l.ctx)
	return &sysclient.UnassignedMenusResponse{
		Menus: logic.ConvertToRpcMenus(menus),
	}, nil
}
