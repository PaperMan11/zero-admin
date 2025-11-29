package scopeservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/db/mysql/model"

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
	exists, err := l.svcCtx.DB.ExistsMenuByName(l.ctx, in.Menu.MenuName)
	if err != nil {
		logc.Errorf(l.ctx, "判断菜单是否存在失败, 菜单名：%s, 错误：%s", in.Menu.MenuName, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	if exists {
		return nil, xerr.NewErrCode(xerr.ErrorMenuExist)
	}
	exists, _ = l.svcCtx.DB.ExistsMenuByPath(l.ctx, in.Menu.Path)
	if exists {
		return nil, xerr.NewErrMsg("菜单路径已存在")
	}

	newMenu := in.Menu
	operator := convert.ToString(in.OperatorId)
	var noCache, affix, external, hidden int32
	if newMenu.NoCache {
		noCache = 1
	}
	if newMenu.Affix {
		affix = 1
	}
	if newMenu.External {
		external = 1
	}
	if newMenu.Hidden {
		hidden = 1
	}
	menuID, err := l.svcCtx.DB.CreateMenu(l.ctx, model.SysMenu{
		ScopeID:   newMenu.ScopeId,
		ParentID:  newMenu.ParentId,
		MenuName:  newMenu.MenuName,
		MenuType:  newMenu.MenuType,
		Path:      newMenu.Path,
		Component: newMenu.Component,
		Redirect:  newMenu.Redirect,
		Icon:      newMenu.Icon,
		Sort:      newMenu.Sort,
		NoCache:   noCache,
		Affix:     affix,
		External:  external,
		Hidden:    hidden,
		Status:    newMenu.Status,
		Creator:   operator,
		Updater:   operator,
		DelFlag:   0,
		Remark:    newMenu.Remark,
	})
	if err != nil {
		logc.Errorf(l.ctx, "创建菜单失败, 菜单名：%s, 错误：%s", in.Menu.MenuName, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	return NewGetMenuByIdLogic(l.ctx, l.svcCtx).GetMenuById(&sysclient.Int64Value{Value: menuID})
}
