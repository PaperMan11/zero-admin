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

// 删除角色
func (l *DeleteRoleLogic) DeleteRole(in *sysclient.DeleteRoleRequest) (*sysclient.Empty, error) {
	//operator := convert.ToString(in.OperatorId)
	exists, err := l.svcCtx.DB.ExistsRoleByID(l.ctx, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "查询角色失败, 角色ID：%d, 错误：%s", in.Id, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	if !exists {
		return nil, xerr.NewErrCode(xerr.ErrorRoleNotExist)
	}

	// 检查该角色是否已被用户关联，如果有则不能删除
	count, err := l.svcCtx.DB.CountRoleAssociated(l.ctx, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "查询角色关联用户失败, 角色ID：%d, 错误：%s", in.Id, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorGetRoleAssociated)
	}
	if count > 0 {
		return nil, xerr.NewErrMsg("该角色已被用户关联，无法删除")
	}
	// 同时删除角色与菜单的关联关系
	err = l.svcCtx.DB.DeleteRoleTx(l.ctx, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "删除角色菜单关联失败, roleId: %d, err: %v", in.Id, err)
		return nil, xerr.NewErrMsg("删除角色菜单关联失败")
	}

	return &sysclient.Empty{}, nil
}
