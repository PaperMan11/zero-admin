package permissionservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"

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
	err := l.svcCtx.DB.DeleteMenu(l.ctx, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "删除菜单失败, 参数：%+v 错误：%s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDeleteMenu)
	}
	return &sysclient.Empty{}, nil
}
