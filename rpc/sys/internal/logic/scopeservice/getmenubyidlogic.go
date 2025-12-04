package scopeservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
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
		return nil, errors.New("该菜单不存在")
	case err != nil:
		logc.Errorf(l.ctx, "查询菜单信息, 参数：%+v, 错误：%v", in, err)
		return nil, status.Error(codes.Internal, "查询菜单信息失败")
	}
	return logic.ConvertToRpcMenu(menu), nil
}
