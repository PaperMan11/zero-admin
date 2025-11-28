package scopeservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateScopeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateScopeLogic {
	return &UpdateScopeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateScopeLogic) UpdateScope(in *sysclient.UpdateScopeRequest) (*sysclient.Scope, error) {
	scope, err := l.svcCtx.DB.GetScopeByID(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorScopeNotExist)
		}
		logc.Errorf(l.ctx, "判断安全范围是否存在失败, scope code：%s, 错误：%s", in.ScopeCode, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	scope.ScopeName = in.ScopeName
	scope.Description = in.Description
	scope.Sort = in.Sort
	err = l.svcCtx.DB.SaveScope(l.ctx, scope)
	if err != nil {
		logc.Errorf(l.ctx, "更新安全范围失败, scope code：%s, 错误：%s", in.ScopeCode, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	return logic.ConvertToRpcScope(&scope), nil
}
