// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/roleservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleListLogic {
	return &GetRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleListLogic) GetRoleList(req *types.RoleListRequest) (resp *types.RoleListResponse, err error) {
	res, err := l.svcCtx.RoleService.GetRoleList(l.ctx, &roleservice.RoleListRequest{
		PageRequest: &roleservice.PageRequest{
			Page:     int32(req.PageRequest.Page),
			PageSize: int32(req.PageRequest.PageSize),
			Keyword:  req.PageRequest.Keyword,
		},
		Status: req.Status,
	})
	if err != nil {
		logc.Errorf(l.ctx, "获取角色列表失败：%v", err)
		return nil, err
	}

	roles := make([]types.Role, 0, len(res.Roles))
	for _, role := range res.Roles {
		roles = append(roles, utils.ConvertToTypesRole(role))
	}
	return &types.RoleListResponse{
		PageResponse: types.PageResponse{
			Total:     int64(res.PageResponse.Total),
			Page:      int64(res.PageResponse.Page),
			PageSize:  int64(res.PageResponse.PageSize),
			TotalPage: int64(res.PageResponse.TotalPage),
		},
		Roles: roles,
	}, nil
}
