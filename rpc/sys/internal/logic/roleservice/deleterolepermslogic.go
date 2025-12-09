package roleservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/rpc/sys/db/common"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRolePermsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRolePermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRolePermsLogic {
	return &DeleteRolePermsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除角色权限
func (l *DeleteRolePermsLogic) DeleteRolePerms(in *sysclient.DeleteRolePermsRequest) (*sysclient.RoleInfo, error) {
	if common.IsSuperUser(in.RoleCode) {
		logc.Errorf(l.ctx, "超级管理员角色不允许修改, 角色：%s", in.RoleCode)
		return nil, status.Error(codes.PermissionDenied, common.ErrSuperUserDoNotEdit.Error())
	}
	
	role, err := l.svcCtx.DB.GetRoleByCode(l.ctx, in.RoleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		logc.Errorf(l.ctx, "查询角色失败, 角色ID：%s, 错误：%s", in.RoleCode, err.Error())
		return nil, status.Error(codes.Internal, "删除角色权限失败")
	}

	err = l.svcCtx.DB.DeleteRoleScopes(l.ctx, in.RoleCode, in.ScopeCodes)
	if err != nil {
		return nil, status.Error(codes.Internal, "删除角色权限失败")
	}

	perms, err := NewGetRolePermsLogic(l.ctx, l.svcCtx).GetRolePerms(&sysclient.GetRolePermsRequest{RoleCode: role.RoleCode})
	if err != nil {
		logc.Errorf(l.ctx, "获取角色权限失败, 角色ID：%s, 异常: %s", in.RoleCode, err.Error())
		return nil, status.Error(codes.Internal, "获取角色权限失败")
	}

	return perms, nil
}
