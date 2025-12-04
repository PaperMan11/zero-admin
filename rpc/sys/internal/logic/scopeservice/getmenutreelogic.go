package scopeservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuTreeLogic {
	return &GetMenuTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 菜单管理
func (l *GetMenuTreeLogic) GetMenuTree(in *sysclient.MenuListRequest) (*sysclient.MenuTreeResponse, error) {
	menus, err := l.svcCtx.DB.GetMenus(l.ctx, in.Status, 0, -1)
	if err != nil {
		logc.Errorf(l.ctx, "查询菜单信息, 参数：%+v, 错误：%v", in, err)
		return nil, status.Error(codes.Internal, "查询菜单信息异常")
	}

	tree := logic.BuildMenuTree(menus, 0)
	return &sysclient.MenuTreeResponse{
		Menus: tree,
	}, nil
}
