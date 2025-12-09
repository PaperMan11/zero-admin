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
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新角色
func (l *UpdateRoleLogic) UpdateRole(in *sysclient.UpdateRoleRequest) (*sysclient.Role, error) {
	if common.IsSuperUser(in.RoleCode) {
		logc.Errorf(l.ctx, "超级管理员角色不允许修改, 角色：%s", in.RoleCode)
		return nil, status.Error(codes.PermissionDenied, common.ErrSuperUserDoNotEdit.Error())
	}

	role, err := l.svcCtx.DB.GetRoleByCode(l.ctx, in.RoleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorRoleNotExist)
		}
		logc.Errorf(l.ctx, "查询角色失败, 角色ID：%s, 错误：%s", in.RoleCode, err.Error())
		return nil, status.Error(codes.Internal, "更新角色失败")
	}

	role.RoleName = in.RoleName
	role.Description = in.Description
	role.Status = in.Status
	err = l.svcCtx.DB.SaveRole(l.ctx, *role)
	if err != nil {
		logc.Errorf(l.ctx, "更新角色失败, 角色ID：%s, 错误：%s", in.RoleCode, err.Error())
		return nil, status.Error(codes.Internal, "更新角色失败")
	}

	return logic.ConvertToRpcRole(role), nil
}
