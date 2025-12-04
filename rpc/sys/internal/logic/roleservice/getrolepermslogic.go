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

type GetRolePermsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRolePermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePermsLogic {
	return &GetRolePermsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取角色权限
func (l *GetRolePermsLogic) GetRolePerms(in *sysclient.Int64Value) (*sysclient.RoleInfo, error) {
	role, err := l.svcCtx.DB.GetRoleByID(l.ctx, in.Value)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorRoleNotExist)
		}
		logc.Errorf(l.ctx, "查询角色失败, 角色ID：%d, 错误：%s", in.Value, err.Error())
		return nil, status.Error(codes.Internal, "获取角色权限失败")
	}

	roleScopeInfos, err := l.svcCtx.DB.GetRoleScopesPerm(l.ctx, role.RoleCode)
	if err != nil {
		logc.Errorf(l.ctx, "查询角色权限失败, 角色ID：%d, 错误：%s", in.Value, err.Error())
		return nil, status.Error(codes.Internal, "查询角色权限失败")
	}

	scopes := make([]*sysclient.RoleScopeInfo, 0, len(roleScopeInfos))
	for _, roleScopeInfo := range roleScopeInfos {
		scopes = append(scopes, &sysclient.RoleScopeInfo{
			Scope: &sysclient.Scope{
				Id:          roleScopeInfo.ID,
				ScopeName:   roleScopeInfo.ScopeName,
				ScopeCode:   roleScopeInfo.ScopeCode,
				Description: roleScopeInfo.Description,
				Sort:        roleScopeInfo.Sort,
			},
			Perms: common.PermissionMap[roleScopeInfo.Perm],
		})
	}

	return &sysclient.RoleInfo{
		Role:   logic.ConvertToRpcRole(role),
		Scopes: scopes,
	}, nil
}
