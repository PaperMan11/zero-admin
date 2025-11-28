package roleservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"
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
	exists, err := l.svcCtx.DB.ExistsRoleByID(l.ctx, in.RoleId)
	if err != nil {
		logc.Errorf(l.ctx, "查询角色失败, 角色ID：%d, 错误：%s", in.RoleId, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	if !exists {
		return nil, xerr.NewErrCode(xerr.ErrorRoleNotExist)
	}

	err = l.svcCtx.DB.DeleteRoleScopes(l.ctx, in.RoleId, in.ScopeCodes)
	if err != nil {
		return nil, xerr.NewErrMsg("删除角色权限失败")
	}

	perms, err := NewGetRolePermsLogic(l.ctx, l.svcCtx).GetRolePerms(&sysclient.Int64Value{Value: in.RoleId})
	if err != nil {
		logc.Errorf(l.ctx, "获取角色权限失败, 角色ID：%d, 异常: %s", in.RoleId, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorGetRolePerms)
	}

	return perms, nil
}
