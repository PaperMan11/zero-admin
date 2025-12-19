package scopeservicelogic

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

type ToggleScopeStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewToggleScopeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleScopeStatusLogic {
	return &ToggleScopeStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ToggleScopeStatusLogic) ToggleScopeStatus(in *sysclient.ToggleScopeStatusRequest) (*sysclient.Scope, error) {
	scope, err := l.svcCtx.DB.GetScopeByCode(l.ctx, in.ScopeCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("安全范围不存在")
		}
		logc.Errorf(l.ctx, "查询role_code失败, 参数：%+v, 异常: %s", in, err.Error())
		return nil, status.Error(codes.Internal, "禁用安全范围失败")
	}
	if scope.Status == 0 {
		in.Status = 1
	} else {
		in.Status = 0
	}
	err = l.svcCtx.DB.ToggleScopeStatus(l.ctx, scope.ID, in.Status, convert.ToString(in.OperatorId))
	if err != nil {
		logc.Errorf(l.ctx, "禁用安全范围失败, 安全范围ID：%d, 错误：%s", scope.ID, err.Error())
		return nil, status.Error(codes.Internal, "禁用/启用安全范围失败")
	}

	newRole, _ := l.svcCtx.DB.GetScopeByID(l.ctx, scope.ID)
	return logic.ConvertToRpcScope(newRole), nil
}
