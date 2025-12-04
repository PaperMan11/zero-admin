package userservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/pkg/convert"
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
			return nil, errors.New("用户不存在")
		}
		logc.Errorf(l.ctx, "查询用户信息, 参数：%+v, 异常: %s", in, err.Error())
		return nil, status.Error(codes.Internal, "查询用户信息异常")
	}
	err = l.svcCtx.DB.UpdateUserByID(l.ctx, in.UserId, map[string]interface{}{"status": in.Status, "updater": convert.ToString(in.OperatorId)})
	if err != nil {
		logc.Errorf(l.ctx, "更新用户状态, 参数：%+v, 错误：%s", in, err.Error())
		return nil, status.Error(codes.Internal, "更新用户状态失败")
	}
	user.Status = in.Status
	return logic.ConvertToRpcUser(user), nil
}
