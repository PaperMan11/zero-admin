// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleRoleStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleRoleStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleRoleStatusLogic {
	return &ToggleRoleStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleRoleStatusLogic) ToggleRoleStatus(req *types.ToggleRoleStatusRequest) (resp *types.Role, err error) {
	// todo: add your logic here and delete this line

	return
}
