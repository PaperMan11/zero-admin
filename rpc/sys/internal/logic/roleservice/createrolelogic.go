package roleservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-admin/pkg/convert"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建角色
func (l *CreateRoleLogic) CreateRole(in *sysclient.CreateRoleRequest) (*sysclient.Role, error) {
	operator := convert.ToString(in.OperatorId)
	// 1. 是否有重复的
	exists, err := l.svcCtx.DB.ExistsRoleByCode(l.ctx, in.RoleCode)
	if err != nil {
		logc.Errorf(l.ctx, "查询role_code失败, 参数：%+v, 异常: %s", in, err.Error())
		return nil, status.Error(codes.Internal, "创建角色失败")
	}
	if exists {
		logc.Errorf(l.ctx, "角色已存在, 参数：%+v", in)
		return nil, errors.New("角色已存在")
	}

	// 2. 创建
	roleID, err := l.svcCtx.DB.CreateRole(l.ctx, model.SysRole{
		RoleName:    in.RoleName,
		RoleCode:    in.RoleCode,
		Description: in.Description,
		Status:      in.Status,
		Creator:     operator,
		Updater:     operator,
	})
	if err != nil {
		logc.Errorf(l.ctx, "创建角色失败, 参数：%+v, 异常: %s", in, err.Error())
		return nil, status.Error(codes.Internal, "创建角色失败")
	}

	role, _ := l.svcCtx.DB.GetRoleByID(l.ctx, roleID)
	return logic.ConvertToRpcRole(role), nil
}
