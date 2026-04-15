// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package home

import (
	"context"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHotSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHotSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotSearchLogic {
	return &GetHotSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHotSearchLogic) GetHotSearch(req *types.Empty) (resp *types.SearchSuggestResp, err error) {
	// todo: add your logic here and delete this line
	return &types.SearchSuggestResp{
		Keywords: []string{"计算机网络", "Go语言", "区块链", "机器学习", "算法", "数据结构"},
	}, nil
}
