package scopeservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-admin/rpc/sys/internal/logic"
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
	scopes, err := l.svcCtx.DB.GetScopesPagination(l.ctx, in.Status, int(in.PageRequest.Page), int(in.PageRequest.PageSize))
	if err != nil {
		logc.Errorf(l.ctx, "获取安全范围列表失败: %v", err)
		return nil, status.Error(codes.Internal, "获取安全范围列表失败")
	}

	total, _ := l.svcCtx.DB.CountScopes(l.ctx)
	return &sysclient.ScopeListResponse{
		PageResponse: &sysclient.PageResponse{
			Total:     int32(total),
			Page:      in.PageRequest.Page,
			PageSize:  in.PageRequest.PageSize,
			TotalPage: int32(total)/in.PageRequest.GetPageSize() + 1,
		},
		Scopes: logic.ConvertToRpcScopes(scopes),
	}, nil
}
