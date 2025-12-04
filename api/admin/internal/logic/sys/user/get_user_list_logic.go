// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/logic"
	"zero-admin/rpc/sys/client/userservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserListLogic) GetUserList(req *types.UserListRequest) (resp *types.UserListResponse, err error) {
	res, err := l.svcCtx.UserService.GetUserList(l.ctx, &userservice.UserListRequest{
		PageRequest: &userservice.PageRequest{
			Page:     int32(req.Page),
			PageSize: int32(req.PageSize),
			Keyword:  req.Keyword,
		},
		Status: req.Status,
	})
	if err != nil {
		logc.Errorf(l.ctx, "获取用户列表失败: %v", err)
		return nil, err
	}

	return &types.UserListResponse{
		PageResponse: types.PageResponse{
			Total:     int64(res.PageResponse.Total),
			Page:      int64(res.PageResponse.Page),
			PageSize:  int64(res.PageResponse.PageSize),
			TotalPage: int64(res.PageResponse.TotalPage),
		},
		Users: logic.ConvertToTypesUsers(res.Users),
	}, nil
}
