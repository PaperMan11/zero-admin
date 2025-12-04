package scopeservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-admin/pkg/convert"
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
		return nil, status.Error(codes.Internal, "创建菜单失败")
	}
	if exists {
		return nil, errors.New("菜单名已存在")
	}
	exists, _ = l.svcCtx.DB.ExistsMenuByPath(l.ctx, in.Menu.Path)
	if exists {
		return nil, errors.New("菜单路径已存在")
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
		return nil, status.Error(codes.Internal, "创建菜单失败")
	}

	return NewGetMenuByIdLogic(l.ctx, l.svcCtx).GetMenuById(&sysclient.Int64Value{Value: menuID})
}
