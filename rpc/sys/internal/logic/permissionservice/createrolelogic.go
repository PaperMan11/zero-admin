package permissionservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/internal/logic"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRoleLogic) CreateRole(in *sysclient.CreateRoleRequest) (*sysclient.RoleInfo, error) {
	operator := convert.ToString(in.GetOperatorId())
	roleID, err := l.svcCtx.DB.CreateRole(l.ctx, model.SysRole{
		RoleName:    in.RoleName,
		RoleCode:    in.RoleCode,
		Description: in.Description,
		Status:      in.Status,
		Creator:     operator,
		Updater:     operator,
		DelFlag:     0,
	})
	if err != nil {
		logc.Errorf(l.ctx, "创建角色失败, 参数：%+v, 错误：%s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorCreateRole)
	}
	role, _ := l.svcCtx.DB.GetRoleByID(l.ctx, roleID)
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
