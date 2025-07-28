package resp_const

import (
	"errors"
)

var (
	// AUTH Errors

	EmailOrPasswordErr  = errors.New(EmailOrPassword.Msg())
	PwdDontMatch        = errors.New(PasswordDontMatch.Msg())
	UserExistError      = errors.New(UserExists.Msg())
	UserNotFoundError   = errors.New(UserNotFound.Msg())
	DataNotFoundError   = errors.New(DataNotFound.Msg())
	TokenDontMatchError = errors.New(TokenDontMatch.Msg())
	InvalidTokenError   = errors.New(InvalidToken.Msg())
	UserNotActiveErr    = errors.New(UserNotActive.Msg())
	IBadRequestError    = errors.New(BadRequest.Msg())
	InfoOrCodeErr       = errors.New(InfoOrCode.Msg())
	EmailExistsErr      = errors.New(EmailExists.Msg())
)
