package roleservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/db/common"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除角色
func (l *DeleteRoleLogic) DeleteRole(in *sysclient.DeleteRoleRequest) (*sysclient.Empty, error) {
	if common.IsSuperUser(in.RoleCode) {
		logc.Errorf(l.ctx, "超级管理员角色不允许修改, 角色：%s", in.RoleCode)
		return nil, status.Error(codes.PermissionDenied, common.ErrSuperUserDoNotEdit.Error())
	}
	//operator := convert.ToString(in.OperatorId)
	role, err := l.svcCtx.DB.GetRoleByCode(l.ctx, in.RoleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorRoleNotExist)
		}
		logc.Errorf(l.ctx, "查询角色失败, 角色ID：%s, 错误：%s", in.RoleCode, err.Error())
		return nil, status.Error(codes.Internal, "删除角色失败")
	}

	// 检查该角色是否已被用户关联，如果有则不能删除
	count, err := l.svcCtx.DB.CountUserRoles(l.ctx, role.RoleCode)
	if err != nil {
		logc.Errorf(l.ctx, "查询角色关联用户失败, 角色ID：%s, 错误：%s", in.RoleCode, err.Error())
		return nil, status.Error(codes.Internal, "查询角色关联用户失败")
	}
	if count > 0 {
		return nil, errors.New("该角色已被用户关联，请先解除关联关系")
	}
	// 同时删除角色与菜单的关联关系
	err = l.svcCtx.DB.DeleteRoleTx(l.ctx, role.RoleCode)
	if err != nil {
		logc.Errorf(l.ctx, "删除角色菜单关联失败, roleId: %s, err: %v", in.RoleCode, err)
		return nil, status.Error(codes.Internal, "删除角色菜单关联失败")
	}

	return &sysclient.Empty{}, nil
}
