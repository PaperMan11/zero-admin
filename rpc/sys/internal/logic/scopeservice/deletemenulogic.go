package scopeservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMenuLogic) DeleteMenu(in *sysclient.DeleteMenuRequest) (*sysclient.Empty, error) {
	menu, err := l.svcCtx.DB.GetMenuByID(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("菜单不存在")
		}
		logc.Errorf(l.ctx, "判断菜单是否存在失败, 菜单ID：%d, 错误：%s", in.Id, err.Error())
		return nil, status.Error(codes.Internal, "删除菜单失败")
	}
	if menu.ScopeID > 0 {
		return nil, errors.New("该菜单已关联权限，请先解除关联关系")
	}
	err = l.svcCtx.DB.DeleteMenu(l.ctx, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "删除菜单失败, 菜单ID：%d, 错误：%s", in.Id, err.Error())
		return nil, status.Error(codes.Internal, "删除菜单失败")
	}
	return &sysclient.Empty{}, nil
}
