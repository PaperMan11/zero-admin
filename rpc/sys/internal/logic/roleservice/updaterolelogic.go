package roleservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"time"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/response/xerr"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(in *sysclient.UpdateRoleRequest) (*sysclient.RoleInfo, error) {
	role, err := l.svcCtx.DB.GetRoleByID(l.ctx, in.Id)
	// 1.判断角色是否存在
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, xerr.NewErrCode(xerr.ErrorRoleNotExist)
	case err != nil:
		logc.Errorf(l.ctx, "查询角色异常, 请求参数：%+v, 异常信息: %s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	// 2.查询角色名称是否存在
	exists, err := l.svcCtx.DB.ExistsRoleByName(l.ctx, in.RoleName)
	if err != nil {
		logc.Errorf(l.ctx, "根据角色名称：%s,查询角色失败,异常:%s", in.RoleName, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	if exists {
		return nil, xerr.NewErrCode(xerr.ErrorRoleExist)
	}

	// 3.更新角色
	operator := convert.ToString(in.OperatorId)
	role.RoleName = in.RoleName
	role.Description = in.Description
	role.Status = in.Status
	role.Updater = operator
	role.UpdateTime = time.Now()
	l.svcCtx.DB.SaveRole(l.ctx, role)

	return NewGetRoleByIdLogic(l.ctx, l.svcCtx).GetRoleById(&sysclient.Int64Value{Value: in.Id})
}
