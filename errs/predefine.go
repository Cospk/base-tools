package errs

// 预定义的错误码常量
const (
	ServerInternalError      = 500  // 服务器内部错误
	ArgsError                = 1001 // 输入参数错误
	NoPermissionError        = 1002 // 权限不足错误
	DuplicateKeyError        = 1003 // 重复键错误
	RecordNotFoundError      = 1004 // 记录不存在错误
	TokenExpiredError        = 1501 // Token已过期
	TokenInvalidError        = 1502 // Token无效
	TokenMalformedError      = 1503 // Token格式错误
	TokenNotValidYetError    = 1504 // Token尚未生效
	TokenUnknownError        = 1505 // Token未知错误
	TokenKickedError         = 1506 // Token被踢出
	TokenNotExistError       = 1507 // Token不存在
	OrgUserNoPermissionError = 1520 // 组织用户无权限
)

// 预定义的CodeError实例，可直接使用或通过WrapMsg添加上下文
var (
	ErrArgs                     = NewCodeError(ArgsError, "ArgsError")
	ErrNoPermission             = NewCodeError(NoPermissionError, "NoPermissionError")
	ErrInternalServer           = NewCodeError(ServerInternalError, "ServerInternalError")
	ErrRecordNotFound           = NewCodeError(RecordNotFoundError, "RecordNotFoundError")
	ErrDuplicateKey             = NewCodeError(DuplicateKeyError, "DuplicateKeyError")
	ErrTokenExpired             = NewCodeError(TokenExpiredError, "TokenExpiredError")
	ErrTokenInvalid             = NewCodeError(TokenInvalidError, "TokenInvalidError")
	ErrTokenMalformed           = NewCodeError(TokenMalformedError, "TokenMalformedError")
	ErrTokenNotValidYet         = NewCodeError(TokenNotValidYetError, "TokenNotValidYetError")
	ErrTokenUnknown             = NewCodeError(TokenUnknownError, "TokenUnknownError")
	ErrTokenKicked              = NewCodeError(TokenKickedError, "TokenKickedError")
	ErrTokenNotExist            = NewCodeError(TokenNotExistError, "TokenNotExistError")
	ErrOrgUserNoPermissionError = NewCodeError(OrgUserNoPermissionError, "OrgUserNoPermissionError")
)
