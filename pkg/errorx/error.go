package errorx

import "errors"

var (
	ErrorInternal     = errors.New("内部错误")
	ErrorInvalidParam = errors.New("参数无效")
	ErrorCaptcha      = errors.New("验证码错误")

	ErrorTokenGenerate = errors.New("token生成失败")
	ErrorTokenInvalid  = errors.New("token无效")
	ErrorTokenExpired  = errors.New("token过期")

	ErrorUserNotFound              = errors.New("用户不存在")
	ErrorInvalidUsernameOrPassword = errors.New("用户名或密码错误")
	ErrorUserDisabled              = errors.New("用户已禁用")
	ErrorUserDeregistration        = errors.New("用户已注销")
	ErrorInvalidOldPassword        = errors.New("旧密码错误")
	ErrorInvalidPasswordLen        = errors.New("密码长度无效")
	ErrorEmailRegistered           = errors.New("该邮箱已注册")
	ErrorMobileRegistered          = errors.New("该手机号已注册")
	ErrorUsernameRegistered        = errors.New("该用户名已注册")

	ErrorArticleNotFound   = errors.New("文章不存在")
	ErrorInvalidCategoryID = errors.New("分类ID无效")
	ErrorCategoryDisabled  = errors.New("分类已禁用")
	ErrorCategoryNotFound  = errors.New("分类不存在")
	ErrorCategoryNameExist = errors.New("分类名已存在")

	ErrorArticleCommentNotFound     = errors.New("文章评论不存在")
	ErrorArticleCommentDisabled     = errors.New("文章评论已禁用")
	ErrorArticleCommentDeleted      = errors.New("文章评论已删除")
	ErrorArticleCommentNoPermission = errors.New("您没有权限删除该文章评论")
)
