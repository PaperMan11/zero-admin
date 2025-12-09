package userservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/pkg/bcrypt"
	"zero-admin/pkg/convert"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(in *sysclient.UpdatePasswordRequest) (*sysclient.Empty, error) {
	user, err := l.svcCtx.DB.GetUserByID(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		logc.Errorf(l.ctx, "查询用户信息, 参数：%+v, 异常: %s", in, err.Error())
		return nil, status.Error(codes.Internal, "管理员更新用户密码失败")
	}

	hashPass := bcrypt.HashPassword(in.Password + user.Salt)
	err = l.svcCtx.DB.UpdateUserByID(l.ctx, in.UserId, map[string]interface{}{"password": hashPass, "updater": convert.ToString(in.OperatorId)})
	if err != nil {
		logc.Errorf(l.ctx, "更新用户密码, 错误：%s", err.Error())
		return nil, status.Error(codes.Internal, "管理员更新用户密码失败")
	}

	return &sysclient.Empty{}, nil
}
