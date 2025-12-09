package roleservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/pkg/convert"
	"zero-admin/rpc/sys/db/common"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleRoleStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewToggleRoleStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleRoleStatusLogic {
	return &ToggleRoleStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 禁用角色
func (l *ToggleRoleStatusLogic) ToggleRoleStatus(in *sysclient.ToggleRoleStatusRequest) (*sysclient.Role, error) {
	if common.IsSuperUser(in.RoleCode) {
		logc.Errorf(l.ctx, "超级管理员角色不允许修改, 角色：%s", in.RoleCode)
		return nil, status.Error(codes.PermissionDenied, common.ErrSuperUserDoNotEdit.Error())
	}

	role, err := l.svcCtx.DB.GetRoleByCode(l.ctx, in.RoleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		logc.Errorf(l.ctx, "查询role_code失败, 参数：%+v, 异常: %s", in, err.Error())
		return nil, status.Error(codes.Internal, "禁用角色失败")
	}

	err = l.svcCtx.DB.ToggleRoleStatus(l.ctx, role.ID, in.Status, convert.ToString(in.OperatorId))
	if err != nil {
		logc.Errorf(l.ctx, "禁用角色失败, 角色ID：%d, 错误：%s", role.ID, err.Error())
		return nil, status.Error(codes.Internal, "禁用/启用角色失败")
	}

	newRole, _ := l.svcCtx.DB.GetRoleByID(l.ctx, role.ID)
	return logic.ConvertToRpcRole(newRole), nil
}
