package roleservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"zero-admin/pkg/response/xerr"
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
	role, err := l.svcCtx.DB.GetRoleByID(l.ctx, in.RoleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorRoleNotExist)
		}
		logc.Errorf(l.ctx, "查询角色失败, 角色ID：%d, 错误：%s", in.RoleId, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	role.RoleName = in.RoleName
	role.Description = in.Description
	role.Status = in.Status
	err = l.svcCtx.DB.SaveRole(l.ctx, *role)
	if err != nil {
		logc.Errorf(l.ctx, "更新角色失败, 角色ID：%d, 错误：%s", in.RoleId, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	return logic.ConvertToRpcRole(role), nil
}
