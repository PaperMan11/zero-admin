package scopeservicelogic

import (
	"context"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuListLogic {
	return &GetMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuListLogic) GetMenuList(in *sysclient.GetMenuListRequest) (*sysclient.GetMenuListResponse, error) {
	menus, _ := l.svcCtx.DB.GetMenus(l.ctx, in.Status, int(in.PageRequest.Page), int(in.PageRequest.PageSize))
	total, _ := l.svcCtx.DB.CountMenus(l.ctx, in.Status)
	return &sysclient.GetMenuListResponse{
		PageResponse: &sysclient.PageResponse{
			Total:     int32(total),
			Page:      in.PageRequest.Page,
			PageSize:  in.PageRequest.PageSize,
			TotalPage: int32(total)/in.PageRequest.PageSize + 1,
		},
		Menus: logic.ConvertToRpcMenus(menus),
	}, nil
}
