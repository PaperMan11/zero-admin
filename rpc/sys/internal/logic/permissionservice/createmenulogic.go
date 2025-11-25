package permissionservicelogic

import (
	"context"
	"time"
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
	menuID, err := l.svcCtx.DB.CreateMenus(l.ctx, []model.SysMenu{*ConvertToModelMenu(in.OperatorId, in.Menu)})
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorCreateMenuFailed)
	}
	menu, _ := l.svcCtx.DB.GetMenuByID(l.ctx, menuID)
	return ConvertToRpcMenu(menu), nil
}

func ConvertToRpcMenu(menu *model.SysMenu) *sysclient.Menu {
	m := &sysclient.Menu{
		Id:        menu.ID,
		ParentId:  menu.ParentID,
		MenuName:  menu.MenuName,
		MenuType:  menu.MenuType,
		Status:    menu.Status,
		ScopeId:   menu.ScopeID,
		Path:      menu.Path,
		Component: menu.Component,
		Redirect:  menu.Redirect,
		Icon:      menu.Icon,
		Sort:      menu.Sort,
	}
	if menu.NoCache == 1 {
		m.NoCache = true
	}
	if menu.Affix == 1 {
		m.Affix = true
	}
	if menu.External == 1 {
		m.External = true
	}
	if menu.Hidden == 1 {
		m.Hidden = true
	}
	return m
}

func ConvertToModelMenu(operatorID int64, menu *sysclient.Menu) *model.SysMenu {
	now := time.Now()
	operator := convert.ToString(operatorID)
	m := &model.SysMenu{
		ScopeID:    menu.ScopeId,
		ParentID:   menu.ParentId,
		MenuName:   menu.MenuName,
		MenuType:   menu.MenuType,
		Path:       menu.Path,
		Component:  menu.Component,
		Redirect:   menu.Redirect,
		Icon:       menu.Icon,
		Sort:       menu.Sort,
		Creator:    operator,
		CreateTime: now,
		Updater:    operator,
		UpdateTime: now,
		DelFlag:    0,
		Remark:     menu.Remark,
	}
	if menu.NoCache {
		m.NoCache = 1
	}
	if menu.Affix {
		m.Affix = 1
	}
	if menu.External {
		m.External = 1
	}
	if menu.Hidden {
		m.Hidden = 1
	}
	if menu.Status == 1 {
		m.Status = 1
	}
	return m
}
