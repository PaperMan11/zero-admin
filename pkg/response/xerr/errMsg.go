package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "SUCCESS"
	message[ErrorServerCommon] = "服务器开小差啦,稍后再来试一试"
	message[ErrorRequestParam] = "参数错误"
	message[ErrorTokenExpire] = "token失效，请重新登陆"
	message[ErrorTokenGenerate] = "生成token失败"
	message[ErrorTokenInvalid] = "无效token"
	message[ErrorDb] = "数据库繁忙,请稍后再试"
	message[ErrorDbUpdateAffectedZero] = "更新数据影响行数为0"
	message[ErrorCaptcha] = "验证码错误"
	message[ErrorInvalidInput] = "无效输入"
	message[ErrorRegister] = "注册失败"

	message[ErrorPermissionDenied] = "权限不足"

	// sys_base 模块
	message[ErrorRepeatOssGetBucket] = "获取oss_bucket实例失败"
	message[ErrorRepeatOssPutBucket] = "上传oss_bucket失败"
	message[ErrorEmailCannotBeEmpty] = "邮箱不能为空错误"
	message[ErrorAlreadyExists] = "邮箱已发送"

	// article
	message[ErrorArticleNotExist] = "文章不存在"
	message[ErrorCommentNotExist] = "评论不存在"

	// user
	message[ErrorUserNotExist] = "用户不存在"
	message[ErrorUserExist] = "用户已存在"
	message[ErrorOldPassword] = "旧密码错误"

	message[ErrorUserPassword] = "用户或者密码错误"
	message[ErrorRepeatName] = "注册失败,不能添加同样的name"
	message[ErrorCreateUser] = "创建用户失败"
	message[ErrorUpdateUser] = "更新用户失败"
	message[ErrorDeleteUser] = "删除用户失败"

	// menu
	message[ErrorCreateMenu] = "创建菜单失败"
	message[ErrorGetMenuRoleTree] = "获取菜单角色树失败"
	message[ErrorGetMenuRoleList] = "获取菜单角色列表失败"
	message[ErrorGetMenuTree] = "获取菜单树失败"
	message[ErrorGetMenuList] = "获取菜单列表失败"
	message[ErrorGetMenu] = "获取菜单失败"
	message[ErrorDeleteMenu] = "删除菜单失败"
	message[ErrorUpdateMenu] = "更新菜单失败"
	message[ErrorMenuExist] = "菜单已存在"
	message[ErrorMenuNotExist] = "菜单不存在"

	// role
	message[ErrorCreateRole] = "创建角色失败"
	message[ErrorUpdateRole] = "修改角色信息失败"
	message[ErrorDeleteRole] = "删除角色信息失败"
	message[ErrorGetRole] = "获取角色信息失败"
	message[ErrorGetRoleList] = "获取角色列表失败"
	message[ErrorRoleExist] = "角色已存在"
	message[ErrorRoleNotExist] = "角色不存在"
	message[ErrorAddRoleScope] = "添加角色权限失败"
	message[ErrorGetRolePerms] = "获取角色权限失败"
	message[ErrorGetRoleAssociated] = "查询角色关联用户失败"

	// scope
	message[ErrorCreateScope] = "创建安全范围失败"
	message[ErrorUpdateScope] = "修改安全范围失败"
	message[ErrorDeleteScope] = "删除安全范围失败"
	message[ErrorGetScope] = "获取安全范围失败"
	message[ErrorGetScopeList] = "获取安全范围表失败"
	message[ErrorScopeExist] = "安全范围已存在"
	message[ErrorScopeNotExist] = "安全范围不存在"
}

func MapErrMsg(errcode uint32) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

func IsCodeErr(errcode uint32) bool {
	if _, ok := message[errcode]; ok {
		return true
	} else {
		return false
	}
}
