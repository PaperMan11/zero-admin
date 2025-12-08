// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/scopeservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteScopeMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScopeMenusLogic {
	return &DeleteScopeMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteScopeMenusLogic) DeleteScopeMenus(req *types.DeleteScopeMenusRequest) (resp *types.ScopeInfo, err error) {
	res, err := l.svcCtx.ScopeService.DeleteScopeMenus(l.ctx, &scopeservice.DeleteScopeMenusRequest{
		ScopeId:    req.ScopeId,
		MenuIds:    req.MenuIds,
		OperatorId: utils.GetOperateID(l.ctx),
		DeleteAll:  req.DeleteAll,
	})
	if err != nil {
		logc.Errorf(l.ctx, "删除安全范围菜单失败: %v", err)
		return nil, err
	}

	return &types.ScopeInfo{
		Scope: types.Scope{
			Id:          res.Scope.Id,
			ScopeName:   res.Scope.ScopeName,
			ScopeCode:   res.Scope.ScopeCode,
			Description: res.Scope.Description,
			Sort:        res.Scope.Sort,
		},
		Menus: utils.ConvertToTypesMenus(res.Menus),
	}, nil
}
