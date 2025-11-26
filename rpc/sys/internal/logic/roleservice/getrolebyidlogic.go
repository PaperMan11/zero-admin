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

type GetRoleByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleByIdLogic {
	return &GetRoleByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRoleByIdLogic) GetRoleById(in *sysclient.Int64Value) (*sysclient.RoleInfo, error) {
	role, err := l.svcCtx.DB.GetRoleByID(l.ctx, in.GetValue())
	if err != nil {
		logc.Errorf(l.ctx, "查询角色失败, 参数：%+v, 错误：%s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	scopes, _ := l.svcCtx.DB.GetScopesByRoleCode(l.ctx, role.RoleCode)
	menus, _ := l.svcCtx.DB.GetMenusByRoleCode(l.ctx, role.RoleCode)
	return &sysclient.RoleInfo{
		Id:          role.ID,
		RoleName:    role.RoleName,
		RoleCode:    role.RoleCode,
		Description: role.Description,
		Status:      role.Status,
		Menus:       logic.ConvertToRpcMenus(menus),
		Scopes:      logic.ConvertToRpcScopes(scopes),
	}, nil
}
