package scopeservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"zero-admin/pkg/response/xerr"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMenuLogic) UpdateMenu(in *sysclient.UpdateMenuRequest) (*sysclient.Menu, error) {
	menu, err := l.svcCtx.DB.GetMenuByID(l.ctx, in.Menu.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorMenuNotExist)
		}
		logc.Errorf(l.ctx, "更新菜单, 菜单ID：%d, 错误：%s", in.Menu.Id, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	exists, _ := l.svcCtx.DB.ExistsMenuByName(l.ctx, in.Menu.MenuName)
	if exists {
		return nil, xerr.NewErrCode(xerr.ErrorMenuExist)
	}
	exists, _ = l.svcCtx.DB.ExistsMenuByPath(l.ctx, in.Menu.Path)
	if exists {
		return nil, xerr.NewErrMsg("菜单路径已存在")
	}

	menu.MenuName = in.Menu.MenuName
	menu.Path = in.Menu.Path
	menu.Component = in.Menu.Component
	menu.Icon = in.Menu.Icon
	menu.Sort = in.Menu.Sort
	menu.ParentID = in.Menu.ParentId
	menu.MenuType = in.Menu.MenuType
	var hidden, affix, external, noCache int32
	if in.Menu.Hidden {
		hidden = 1
	}
	if in.Menu.Affix {
		affix = 1
	}
	if in.Menu.External {
		external = 1
	}
	if in.Menu.NoCache {
		noCache = 1
	}
	menu.Hidden = hidden
	menu.Affix = affix
	menu.External = external
	menu.NoCache = noCache
	menu.Status = in.Menu.Status
	err = l.svcCtx.DB.SaveMenu(l.ctx, *menu)
	if err != nil {
		logc.Errorf(l.ctx, "更新菜单, 菜单ID：%d, 错误：%s", in.Menu.Id, err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ErrorDb, "更新菜单失败")
	}

	return NewGetMenuByIdLogic(l.ctx, l.svcCtx).GetMenuById(&sysclient.Int64Value{Value: menu.ID})
}
