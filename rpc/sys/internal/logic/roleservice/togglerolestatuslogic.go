package roleservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"
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
	exists, err := l.svcCtx.DB.ExistsRoleByID(l.ctx, in.RoleId)
	if err != nil {
		logc.Errorf(l.ctx, "查询role_code失败, 参数：%+v, 异常: %s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	if !exists {
		return nil, xerr.NewErrCode(xerr.ErrorRoleNotExist)
	}

	err = l.svcCtx.DB.ToggleRoleStatus(l.ctx, in.RoleId, in.Status)
	if err != nil {
		logc.Errorf(l.ctx, "禁用角色失败, 角色ID：%d, 错误：%s", in.RoleId, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	role, _ := l.svcCtx.DB.GetRoleByID(l.ctx, in.RoleId)
	return logic.ConvertToRpcRole(&role), nil
}
