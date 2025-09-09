package consts

type ContextKey string

func (o ContextKey) Str() string {
	return string(o)
}

var CtxClaims = ContextKey("USER_CLAIMS")

var CTXCompany_ID = ContextKey("CTX_COMPANY_ID") //used to set the companies id on
var CTXUser_ID = ContextKey("CTX_USER_ID")

// ===============. authKeys

// AUTH_FIELD is the key used on the database
type AUTH_FIELD string

func (o AUTH_FIELD) S() string {
	return string(o)
}

var COMPANY_ID = AUTH_FIELD("company_id")
var USER_ID = AUTH_FIELD("user_id")
