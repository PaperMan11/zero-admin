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

type GetScopeMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScopeMenusLogic {
	return &GetScopeMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetScopeMenusLogic) GetScopeMenus(in *sysclient.Int64Value) (*sysclient.ScopeInfo, error) {
	scope, err := l.svcCtx.DB.GetScopeByID(l.ctx, in.Value)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorScopeNotExist)
		}
		logc.Errorf(l.ctx, "判断安全范围是否存在失败, scope id：%s, 错误：%s", in.Value, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	menus, err := l.svcCtx.DB.GetMenusByScopeID(l.ctx, in.Value)
	if err != nil {
		logc.Errorf(l.ctx, "获取安全范围菜单列表, 参数：%+v, 错误：%v", in, err)
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	return &sysclient.ScopeInfo{
		Scope: logic.ConvertToRpcScope(scope),
		Menus: logic.ConvertToRpcMenus(menus),
	}, nil
}
