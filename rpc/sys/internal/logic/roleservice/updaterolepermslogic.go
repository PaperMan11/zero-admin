package roleservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
	"zero-admin/rpc/sys/db/common"
	"zero-admin/rpc/sys/db/mysql/model"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRolePermsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRolePermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRolePermsLogic {
	return &UpdateRolePermsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新角色权限
func (l *UpdateRolePermsLogic) UpdateRolePerms(in *sysclient.UpdateRolePermsRequest) (*sysclient.RoleInfo, error) {
	role, err := l.svcCtx.DB.GetRoleByCode(l.ctx, in.RoleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		logc.Errorf(l.ctx, "查询角色失败, 角色ID：%s, 错误：%s", in.RoleCode, err.Error())
		return nil, status.Error(codes.Internal, "更新角色权限失败")
	}

	now := time.Now()
	updates := make([]model.SysRoleScope, 0, len(in.RoleScopes))
	for _, roleScope := range in.RoleScopes {
		updates = append(updates, model.SysRoleScope{
			RoleCode:   in.RoleCode,
			ScopeCode:  roleScope.ScopeCode,
			Perm:       common.ParsePermission(roleScope.Perms),
			CreateTime: now,
		})
	}
	err = l.svcCtx.DB.UpdateRoleScopesTx(l.ctx, in.RoleCode, updates)
	if err != nil {
		logc.Errorf(l.ctx, "更新角色权限失败, 角色ID：%d, 错误：%s", in.RoleCode, err.Error())
		return nil, status.Error(codes.Internal, "更新角色权限失败")
	}

	return NewGetRolePermsLogic(l.ctx, l.svcCtx).GetRolePerms(&sysclient.GetRolePermsRequest{RoleCode: role.RoleCode})
}
