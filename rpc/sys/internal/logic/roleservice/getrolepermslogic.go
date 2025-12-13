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
	"zero-admin/rpc/sys/db/mysql/model"
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
func (l *GetRolePermsLogic) GetRolePerms(in *sysclient.GetRolePermsRequest) (*sysclient.RoleInfo, error) {
	role, err := l.svcCtx.DB.GetRoleByCode(l.ctx, in.RoleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorRoleNotExist)
		}
		logc.Errorf(l.ctx, "查询角色失败, 角色ID：%s, 错误：%s", in.RoleCode, err.Error())
		return nil, status.Error(codes.Internal, "获取角色权限失败")
	}

	var roleScopeInfos []model.RoleScopeInfo
	if common.IsSuperUser(role.RoleCode) {
		scopes, _ := l.svcCtx.DB.GetAllScopes(l.ctx)
		for _, scope := range scopes {
			roleScopeInfos = append(roleScopeInfos, model.RoleScopeInfo{
				RoleID:   role.ID,
				RoleName: role.RoleName,
				RoleCode: role.RoleCode,
				Perm:     common.PERM_ALL,
				SysScope: model.SysScope{
					ID:          scope.ID,
					ScopeName:   scope.ScopeName,
					ScopeCode:   scope.ScopeCode,
					Description: scope.Description,
					Sort:        scope.Sort,
					Creator:     scope.Creator,
					CreateTime:  scope.CreateTime,
					Updater:     scope.Updater,
					UpdateTime:  scope.UpdateTime,
					DelFlag:     scope.DelFlag,
				},
			})
		}
	} else {
		roleScopeInfos, err = l.svcCtx.DB.GetRoleScopesPerm(l.ctx, role.RoleCode)
	}
	if err != nil {
		logc.Errorf(l.ctx, "查询角色权限失败, 角色ID：%s, 错误：%s", in.RoleCode, err.Error())
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
