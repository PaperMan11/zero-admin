// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/scopeservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuListLogic {
	return &GetMenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuListLogic) GetMenuList(req *types.GetMenuListRequest) (resp *types.GetMenuListResponse, err error) {
	res, _ := l.svcCtx.ScopeService.GetMenuList(l.ctx, &scopeservice.GetMenuListRequest{
		PageRequest: &scopeservice.PageRequest{
			Page:     int32(req.Page),
			PageSize: int32(req.PageSize),
			Keyword:  req.Keyword,
		},
		Status: req.Status,
	})

	return &types.GetMenuListResponse{
		PageResponse: types.PageResponse{
			Total:     int64(res.PageResponse.Total),
			Page:      int64(res.PageResponse.Page),
			PageSize:  int64(res.PageResponse.PageSize),
			TotalPage: int64(res.PageResponse.TotalPage),
		},
		Menus: utils.ConvertToTypesMenus(res.Menus),
	}, nil
}
