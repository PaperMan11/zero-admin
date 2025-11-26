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

type GetRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleListLogic {
	return &GetRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 角色管理
func (l *GetRoleListLogic) GetRoleList(in *sysclient.RoleListRequest) (*sysclient.RoleListResponse, error) {
	roles, err := l.svcCtx.DB.GetRolesPagination(l.ctx, in.GetStatus(), int(in.GetPageRequest().Page), int(in.GetPageRequest().PageSize))
	if err != nil {
		logc.Errorf(l.ctx, "查询角色失败, 参数：%+v, 错误：%s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	roleInfos := make([]*sysclient.RoleInfo, 0, len(roles))
	for _, role := range roles {
		scopes, _ := l.svcCtx.DB.GetScopesByRoleCode(l.ctx, role.RoleCode)
		menus, _ := l.svcCtx.DB.GetMenusByRoleCode(l.ctx, role.RoleCode)
		roleInfos = append(roleInfos, &sysclient.RoleInfo{
			Id:          role.ID,
			RoleName:    role.RoleName,
			RoleCode:    role.RoleCode,
			Description: role.Description,
			Status:      role.Status,
			Menus:       logic.ConvertToRpcMenus(menus),
			Scopes:      logic.ConvertToRpcScopes(scopes),
		})
	}
	total, _ := l.svcCtx.DB.CountRoles(l.ctx)
	return &sysclient.RoleListResponse{
		PageResponse: &sysclient.PageResponse{
			Total:     int32(total),
			Page:      in.GetPageRequest().GetPage(),
			PageSize:  in.GetPageRequest().GetPageSize(),
			TotalPage: int32(total) / (in.GetPageRequest().GetPageSize()),
		},
		Roles: roleInfos,
	}, nil
}
