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

	message[ErrorInvalidInput] = "无效输入"
	message[ErrorRegister] = "注册失败"
	message[ErrorUserNotExist] = "用户不存在"
	message[ErrorUserExist] = "用户已存在"
	message[ErrorOldPassword] = "旧密码错误"
	message[ErrorCaptcha] = "验证码错误"
	message[ErrorUserPassword] = "用户或者密码错误"
	message[ErrorRepeatName] = "注册失败,不能添加同样的name"
	message[ErrorPermissionDenied] = "权限不足"

	// sys_base 模块
	message[ErrorRepeatOssGetBucket] = "获取oss_bucket实例失败"
	message[ErrorRepeatOssPutBucket] = "上传oss_bucket失败"
	message[ErrorEmailCannotBeEmpty] = "邮箱不能为空错误"
	message[ErrorAlreadyExists] = "邮箱已发送"

	// article
	message[ErrorArticleNotExist] = "文章不存在"
	message[ErrorCommentNotExist] = "评论不存在"

	// menu
	message[ErrorCreateMenuFailed] = "创建菜单失败"
	message[ErrorGetMenuRoleTreeFailed] = "获取菜单角色树失败"
	message[ErrorGetMenuRoleListFailed] = "获取菜单角色列表失败"
	message[ErrorGetMenuTreeFailed] = "获取菜单树失败"
	message[ErrorGetMenuListFailed] = "获取菜单列表失败"
	message[ErrorGetMenuFailed] = "获取菜单失败"
	message[ErrorDeleteMenuFailed] = "删除菜单失败"
	message[ErrorUpdateMenuFailed] = "更新菜单失败"
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
