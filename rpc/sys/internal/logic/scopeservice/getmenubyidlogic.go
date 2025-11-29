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

type GetMenuByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuByIdLogic {
	return &GetMenuByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuByIdLogic) GetMenuById(in *sysclient.Int64Value) (*sysclient.Menu, error) {
	menu, err := l.svcCtx.DB.GetMenuByID(l.ctx, in.Value)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, xerr.NewErrCode(xerr.ErrorMenuNotExist)
	case err != nil:
		logc.Errorf(l.ctx, "查询菜单信息, 参数：%+v, 错误：%v", in, err)
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	return logic.ConvertToRpcMenu(menu), nil
}
