package permissionservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMenuLogic) CreateMenu(in *sysclient.CreateMenuRequest) (*sysclient.Menu, error) {
	menuID, err := l.svcCtx.DB.CreateMenus(l.ctx, []model.SysMenu{*logic.ConvertToModelMenu(in.OperatorId, in.Menu)})
	if err != nil {
		logc.Errorf(l.ctx, "创建菜单失败, 参数：%+v, 错误：%s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorCreateMenu)
	}
	menu, _ := l.svcCtx.DB.GetMenuByID(l.ctx, menuID)
	return logic.ConvertToRpcMenu(&menu), nil
}
