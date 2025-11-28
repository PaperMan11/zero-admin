package userservicelogic

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

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserLogic) DeleteUser(in *sysclient.DeleteUserRequest) (*sysclient.Empty, error) {
	_, err := l.svcCtx.DB.GetUserByID(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorUserNotExist)
		}
		logc.Errorf(l.ctx, "查询用户失败,参数：%+v, 错误：%s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	err = l.svcCtx.DB.DeleteUserTx(l.ctx, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "删除用户失败,参数：%+v, 错误：%s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDeleteUser)
	}
	return &sysclient.Empty{}, nil
}
