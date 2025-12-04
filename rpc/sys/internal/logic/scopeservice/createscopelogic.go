package scopeservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"zero-admin/pkg/convert"
	"zero-admin/rpc/sys/db/mysql/model"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateScopeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateScopeLogic {
	return &CreateScopeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateScopeLogic) CreateScope(in *sysclient.CreateScopeRequest) (*sysclient.Scope, error) {
	exists, err := l.svcCtx.DB.ExistsScopeByCode(l.ctx, in.ScopeCode)
	if err != nil {
		logc.Errorf(l.ctx, "判断安全范围是否存在失败, scope code：%s, 错误：%s", in.ScopeCode, err.Error())
		return nil, status.Error(codes.Internal, "添加安全范围失败")
	}
	if exists {
		return nil, errors.New("安全范围已存在")
	}

	now := time.Now()
	operator := convert.ToString(in.OperatorId)
	scopeID, err := l.svcCtx.DB.CreateScopeTx(l.ctx, model.SysScope{
		ScopeName:   in.ScopeName,
		ScopeCode:   in.ScopeCode,
		Description: in.Description,
		Sort:        in.Sort,
		Creator:     operator,
		CreateTime:  now,
		Updater:     operator,
		UpdateTime:  now,
		DelFlag:     0,
	}, in.MenuIds)

	return NewGetScopeByIdLogic(l.ctx, l.svcCtx).GetScopeById(&sysclient.Int64Value{Value: scopeID})
}
