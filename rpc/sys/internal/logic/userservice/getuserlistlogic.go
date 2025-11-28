package userservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户管理
func (l *GetUserListLogic) GetUserList(in *sysclient.UserListRequest) (*sysclient.UserListResponse, error) {
	users, err := l.svcCtx.DB.GetUsersPagination(l.ctx, in.Status, int(in.PageRequest.Page), int(in.PageRequest.PageSize))
	if err != nil {
		logc.Errorf(l.ctx, "获取用户列表失败: %v", err)
		return nil, xerr.NewErrCodeMsg(xerr.ErrorDb, "获取用户列表失败")
	}

	total, _ := l.svcCtx.DB.CountUsers(l.ctx, in.Status)

	return &sysclient.UserListResponse{
		PageResponse: &sysclient.PageResponse{
			Total:     int32(total),
			Page:      in.PageRequest.Page,
			PageSize:  in.PageRequest.PageSize,
			TotalPage: int32(total)/in.PageRequest.PageSize + 1,
		},
		Users: nil,
	}, nil
}
