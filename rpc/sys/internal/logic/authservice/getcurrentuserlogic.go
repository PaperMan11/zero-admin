package authservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/db"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/internal/logic"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserLogic {
	return &GetCurrentUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取当前用户信息
func (l *GetCurrentUserLogic) GetCurrentUser(in *sysclient.GetCurrentUserRequest) (*sysclient.UserInfo, error) {
	user, err := l.svcCtx.DB.GetUserByID(l.ctx, in.UserId)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, xerr.NewErrCode(xerr.ErrorUserNotExist)
	case err != nil:
		logc.Errorf(l.ctx, "查询用户信息, 参数：%+v, 异常: %s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	// 用户角色信息
	userRoles, roleCodes := GetUserRoles(l.ctx, l.svcCtx.DB, user.ID)
	// menus
	menus, _ := l.svcCtx.DB.GetMenusByRoles(l.ctx, roleCodes)

	resp := &sysclient.UserInfo{
		Id:       user.ID,
		Username: user.Username,
		Status:   user.Status,
		Roles:    logic.ConvertToRpcRoles(userRoles),
		MenuTree: logic.BuildMenuTree(menus, 0),
		Email:    user.Email,
		Mobile:   user.Mobile,
		RealName: user.RealName,
		Gender:   user.Gender,
		Avatar:   user.Avatar,
	}
	return resp, nil
}

func GetUserRoles(ctx context.Context, db db.DB, userID int64) (roles []model.SysRole, roleCodes []string) {
	roles, _ = db.GetRolesByUserID(ctx, userID)
	roleCodes = make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.RoleCode)
	}
	return
}
