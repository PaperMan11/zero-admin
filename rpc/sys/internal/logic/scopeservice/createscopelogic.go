package scopeservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateScopeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateScopeLogic {
	return &CreateScopeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateScopeLogic) CreateScope(in *sysclient.CreateScopeRequest) (*sysclient.ScopeInfo, error) {
	operator := convert.ToString(in.GetOperatorId())
	scopeID, err := l.svcCtx.DB.CreateScope(l.ctx, model.SysScope{
		ScopeName:   in.ScopeName,
		ScopeCode:   in.ScopeCode,
		Description: in.Description,
		Sort:        in.Sort,
		Creator:     operator,
		Updater:     operator,
		DelFlag:     0,
	})
	if err != nil {
		logc.Errorf(l.ctx, "创建安全范围失败, 参数：%+v, 错误：%s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorCreateScope)
	}
	scope, _ := l.svcCtx.DB.GetScopeByID(l.ctx, scopeID)
	menus, _ := l.svcCtx.DB.GetMenusByScopeID(l.ctx, scopeID)
	return &sysclient.ScopeInfo{
		Id:          scope.ID,
		ScopeName:   scope.ScopeName,
		ScopeCode:   scope.ScopeCode,
		Description: scope.Description,
		Sort:        scope.Sort,
		Menus:       logic.ConvertToRpcMenus(menus),
	}, nil
}
