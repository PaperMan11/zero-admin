package scopeservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteScopeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScopeLogic {
	return &DeleteScopeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteScopeLogic) DeleteScope(in *sysclient.DeleteScopeRequest) (*sysclient.Empty, error) {
	scope, err := l.svcCtx.DB.GetScopeByID(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("安全范围不存在")
		}
		logc.Errorf(l.ctx, "判断安全范围是否存在失败, scope id：%s, 错误：%s", in.Id, err.Error())
		return nil, status.Error(codes.Internal, "删除安全范围失败")
	}

	roles, _ := l.svcCtx.DB.GetRolePermsByScopeCode(l.ctx, scope.ScopeCode)
	if len(roles) > 0 {
		return nil, errors.New("该安全范围已被角色关联，请先解除关联关系")
	}

	err = l.svcCtx.DB.DeleteScopeTx(l.ctx, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "删除安全范围失败, scope id：%s, 错误：%s", in.Id, err.Error())
		return nil, status.Error(codes.Internal, "删除安全范围失败")
	}
	return &sysclient.Empty{}, nil
}
