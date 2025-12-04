package roleservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleListLogic {
	return &GetRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 角色列表
func (l *GetRoleListLogic) GetRoleList(in *sysclient.RoleListRequest) (*sysclient.RoleListResponse, error) {
	roles, err := l.svcCtx.DB.GetRolesPagination(l.ctx, in.Status, int(in.PageRequest.Page), int(in.PageRequest.PageSize))
	if err != nil {
		logc.Errorf(l.ctx, "获取角色列表失败, 参数：%+v, 错误：%s", in, err.Error())
		return nil, status.Error(codes.Internal, "获取角色列表失败")
	}
	total, _ := l.svcCtx.DB.CountRoles(l.ctx, in.Status)

	return &sysclient.RoleListResponse{
		PageResponse: &sysclient.PageResponse{
			Total:     int32(total),
			Page:      in.PageRequest.Page,
			PageSize:  in.PageRequest.PageSize,
			TotalPage: int32(total)/in.PageRequest.PageSize + 1,
		},
		Roles: logic.ConvertToRpcRoles(roles),
	}, nil
}
