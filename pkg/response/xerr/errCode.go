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
	ErrorInvalidInput         uint32 = 100009
	ErrorCaptcha              uint32 = 100010
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

// -----------------------------------后台管理系统----------------------------------

// user
const (
	ErrorRegister     uint32 = 200004
	ErrorUserExist    uint32 = 200003
	ErrorUserNotExist uint32 = 200005
	ErrorOldPassword  uint32 = 200006
	ErrorUserPassword uint32 = 200008
	ErrorRepeatName   uint32 = 200009
	ErrorCreateUser   uint32 = 200010
	ErrorUpdateUser   uint32 = 200011
	ErrorDeleteUser   uint32 = 200012
)

// menu
const (
	ErrorCreateMenu      = 500001
	ErrorUpdateMenu      = 500002
	ErrorDeleteMenu      = 500003
	ErrorGetMenu         = 500004
	ErrorGetMenuList     = 500005
	ErrorGetMenuTree     = 500006
	ErrorGetMenuRoleList = 500007
	ErrorGetMenuRoleTree = 500008
	ErrorMenuExist       = 500009
	ErrorMenuNotExist    = 500010
)

// role
const (
	ErrorCreateRole   = 600001
	ErrorUpdateRole   = 600002
	ErrorDeleteRole   = 600003
	ErrorGetRole      = 600004
	ErrorGetRoleList  = 600005
	ErrorRoleExist    = 600006
	ErrorRoleNotExist = 600007
	ErrorAddRoleScope = 600008
)

// scope
const (
	ErrorCreateScope   = 700001
	ErrorUpdateScope   = 700002
	ErrorDeleteScope   = 700003
	ErrorGetScope      = 700004
	ErrorGetScopeList  = 700005
	ErrorScopeExist    = 700006
	ErrorScopeNotExist = 700007
)
