package xerr

// 成功返回
const OK uint32 = 200

/**(前3位代表业务,后三位代表具体功能)**/

// 全局错误码
const (
	ErrorServerCommon         uint32 = 100001
	ErrorRequestParam         uint32 = 100002
	ErrorTokenExpire          uint32 = 100003
	ErrorTokenGenerate        uint32 = 100004
	ErrorTokenInvalid         uint32 = 100005
	ErrorDb                   uint32 = 100006
	ErrorDbUpdateAffectedZero uint32 = 100007
	ErrorPermissionDenied     uint32 = 100008
)

// sys 模块
const (
	ErrorInvalidInput uint32 = 200001
	ErrorRegister     uint32 = 200004
	ErrorUserExist    uint32 = 200003
	ErrorUserNotExist uint32 = 200005
	ErrorOldPassword  uint32 = 200006
	ErrorCaptcha      uint32 = 200007
	ErrorUserPassword uint32 = 200008
	ErrorRepeatName   uint32 = 200009
)

// sys_base 模块
const (
	ErrorRepeatOssGetBucket = 300001
	ErrorRepeatOssPutBucket = 300002
	ErrorEmailCannotBeEmpty = 300003
	ErrorAlreadyExists      = 300004
)

// article
const (
	ErrorArticleNotExist = 400001
	ErrorCommentNotExist = 400002
)

// menu
const (
	ErrorCreateMenuFailed      = 500001
	ErrorUpdateMenuFailed      = 500002
	ErrorDeleteMenuFailed      = 500003
	ErrorGetMenuFailed         = 500004
	ErrorGetMenuListFailed     = 500005
	ErrorGetMenuTreeFailed     = 500006
	ErrorGetMenuRoleListFailed = 500007
	ErrorGetMenuRoleTreeFailed = 500008
)
