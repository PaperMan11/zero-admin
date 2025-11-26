package roleservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteRoleLogic) DeleteRole(in *sysclient.DeleteRoleRequest) (*sysclient.Empty, error) {
	err := l.svcCtx.DB.DeleteRoleTx(l.ctx, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "删除角色失败, 参数：%+v 错误：%s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDeleteRole)
	}
	return &sysclient.Empty{}, nil
}
