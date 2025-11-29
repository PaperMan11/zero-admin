package scopeservicelogic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"zero-admin/pkg/response/xerr"
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
		return nil, xerr.NewErrCode(xerr.ErrorScopeNotExist)
	case err != nil:
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	return logic.ConvertToRpcScope(scope), nil
}
