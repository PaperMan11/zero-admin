package userservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewToggleUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleUserStatusLogic {
	return &ToggleUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ToggleUserStatusLogic) ToggleUserStatus(in *sysclient.ToggleUserStatusRequest) (*sysclient.User, error) {
	user, err := l.svcCtx.DB.GetUserByID(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorUserNotExist)
		}
		logc.Errorf(l.ctx, "查询用户信息, 参数：%+v, 异常: %s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}
	err = l.svcCtx.DB.UpdateUserByID(l.ctx, in.UserId, map[string]interface{}{"status": in.Status, "updater": convert.ToString(in.OperatorId)})
	if err != nil {
		logc.Errorf(l.ctx, "更新用户状态, 参数：%+v, 错误：%s", in, err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ErrorDb, "更新用户状态失败")
	}
	user.Status = in.Status
	return logic.ConvertToRpcUser(user), nil
}
