package scopeservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScopeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScopeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScopeListLogic {
	return &GetScopeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 安全范围管理
func (l *GetScopeListLogic) GetScopeList(in *sysclient.ScopeListRequest) (*sysclient.ScopeListResponse, error) {
	scopes, err := l.svcCtx.DB.GetScopesPagination(l.ctx, int(in.PageRequest.Page), int(in.PageRequest.PageSize))
	if err != nil {
		logc.Errorf(l.ctx, "获取安全范围列表失败: %v", err)
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	total, _ := l.svcCtx.DB.CountScopes(l.ctx)
	scopeInfos := make([]*sysclient.ScopeInfo, 0, len(scopes))
	getScopeInfosLogic := NewGetScopeMenusLogic(l.ctx, l.svcCtx)
	for _, scope := range scopes {
		scopeInfo, err := getScopeInfosLogic.GetScopeMenus(&sysclient.Int64Value{Value: scope.ID})
		if err != nil {
			logc.Errorf(l.ctx, "获取安全范围菜单列表失败: %v", err)
			continue
		}
		scopeInfos = append(scopeInfos, scopeInfo)
	}

	return &sysclient.ScopeListResponse{
		PageResponse: &sysclient.PageResponse{
			Total:     int32(total),
			Page:      in.PageRequest.Page,
			PageSize:  in.PageRequest.PageSize,
			TotalPage: int32(total)/in.PageRequest.GetPageSize() + 1,
		},
		Scopes: scopeInfos,
	}, nil
}
