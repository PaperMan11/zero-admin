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

type GetScopeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetScopeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScopeListLogic {
	return &GetScopeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScopeListLogic) GetScopeList(req *types.ScopeListRequest) (resp *types.ScopeListResponse, err error) {
	res, err := l.svcCtx.ScopeService.GetScopeList(l.ctx, &scopeservice.ScopeListRequest{
		PageRequest: &scopeservice.PageRequest{
			Page:     int32(req.Page),
			PageSize: int32(req.PageSize),
			Keyword:  req.Keyword,
		},
	})
	if err != nil {
		logc.Errorf(l.ctx, "查询安全范围列表失败: %v", err)
		return nil, err
	}

	scopeInfos := make([]types.ScopeInfo, 0, len(res.Scopes))
	for _, v := range res.Scopes {
		scopeInfos = append(scopeInfos, types.ScopeInfo{
			Scope: utils.ConvertToTypesScope(v.Scope),
			Menus: utils.ConvertToTypesMenus(v.Menus),
		})
	}
	return &types.ScopeListResponse{
		PageResponse: types.PageResponse{
			Total:     int64(res.PageResponse.Total),
			Page:      int64(res.PageResponse.Page),
			PageSize:  int64(res.PageResponse.PageSize),
			TotalPage: int64(res.PageResponse.TotalPage),
		},
		Scopes: scopeInfos,
	}, nil
}
