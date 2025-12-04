package userservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/pkg/response/xerr"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignUserRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignUserRoleLogic {
	return &AssignUserRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignUserRoleLogic) AssignUserRole(in *sysclient.AssignUserRoleRequest) (*sysclient.Empty, error) {
	user, err := l.svcCtx.DB.GetUserByID(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorUserNotExist)
		}
		logc.Errorf(l.ctx, "查询用户信息, 参数：%+v, 错误：%v", in, err)
		return nil, status.Error(codes.Internal, "查询用户信息异常")
	}

	err = l.svcCtx.DB.AddUserRolesTx(l.ctx, user.ID, in.RoleCodes)
	if err != nil {
		logc.Errorf(l.ctx, "添加用户角色, 参数：%+v, 错误：%v", in, err)
		return nil, status.Error(codes.Internal, "添加用户角色异常")
	}

	return &sysclient.Empty{}, nil
}
