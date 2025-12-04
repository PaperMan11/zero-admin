package scopeservicelogic

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScopeByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScopeByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScopeByIdLogic {
	return &GetScopeByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetScopeByIdLogic) GetScopeById(in *sysclient.Int64Value) (*sysclient.Scope, error) {
	scope, err := l.svcCtx.DB.GetScopeByID(l.ctx, in.Value)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errors.New("该安全范围不存在")
	case err != nil:
		return nil, status.Error(codes.Internal, "查询安全范围失败")
	}
	return logic.ConvertToRpcScope(scope), nil
}
