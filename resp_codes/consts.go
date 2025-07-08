package resp_const

type RespCode string

func (r RespCode) Msg() string {

	data := errorText[r]
	if data != "" {
		return data
	}
	return r.ToStr()
}

func (r RespCode) ToStr() string {
	return string(r)
}

// Generic Http Status codes
const (
	// Success Messages

	Success       = RespCode("SUCCESS")        //200
	CreateSuccess = RespCode("CREATE_SUCCESS") //201
	UpdateSuccess = RespCode("UPDATE_SUCCESS") //201
	DeleteSuccess = RespCode("DELETE_SUCCESS") //204
	//HTTP CODES

	DataNotFound  = RespCode("DATA_NOT_FOUND") //404
	BadRequest    = RespCode("BAD_REQUEST")    //400
	InternalError = RespCode("INTERNAL_ERROR") //500

)

// Operation Messages
const (
	NoRecordsUpdated           = RespCode("NO_RECORDS_UPDATED")
	NoRecordsDeleted           = RespCode("NO_RECORDS_DELETED")
	UpdatingAssociationsFailed = RespCode("UPDATING_ASSOCIATIONS_FAILED")
	NotModified                = RespCode("NOT_MODIFIED")
	//model related

	UserNotCreated = RespCode("USER_NOT_CREATED")
	UserCreated    = RespCode("USER_CREATED")
	//Generic

	Generic = RespCode("SOME_THING_WRONG")
	FAIL    = RespCode("FAILURE")
)

// Validation errors
const (
	Validation   = RespCode("VALIDATION_ERROR")
	EmptyFilter  = RespCode("EMPTY_FILTER")
	EmptyBody    = RespCode("EMPTY_BODY")
	UnknownValue = RespCode("UNKNOWN_VALUE")
	InvalidData  = RespCode("INVALID_DATA")
)

// Auth Status Codes
const (
	EmailOrPassword     = RespCode("EMAIL_OR_PASSWORD_WRONG")
	InfoOrCode          = RespCode("INFO_OR_CODE_WRONG")
	InvalidEmailOrPhone = RespCode("INVALID_EMAIL_OR_PHONE")
	PhoneExists         = RespCode("PHONE_EXISTS")
	EmailExists         = RespCode("EMAIL_EXISTS")
	Password            = RespCode("PASSWORD_NOT_STORED")
	CouldNotAssignRole  = RespCode("ROLE_NOT_ASSIGNED")
	TokenDontMatch      = RespCode("TOKEN_DONT_MATCH")
	InvalidToken        = RespCode("INVALID_TOKEN")
	UserNotFound        = RespCode("USER_NOT_FOUND")
	UserExists          = RespCode("USER_EXISTS")
)

// File related api responses
const (
	MaxUploadSizeExceeded = RespCode("MAX_UPLOAD_SIZE_EXCEEDED")
	ParseFile             = RespCode("CANT_PARSE_FORM")
	InvalidFile           = RespCode("INVALID_FILE")
	InvalidFileType       = RespCode("INVALID_FILE_TYPE")
	FileTooBig            = RespCode("FILE_TOO_BIG")
	CantReadFileType      = RespCode("CANT_READ_FILE_TYPE")
	CantReadFile          = RespCode("CANT_READ_FILE")
	CantWriteFile         = RespCode("CANT_WRITE_FILE")

	ImageUploadError = RespCode("IMAGE_UPLOAD_ERROR")
	ImageUpdateError = RespCode("IMAGE_UPDATE_ERROR")
	ImageDeleteError = RespCode("IMAGE_DELETE_ERROR")
	JsonUnmarshal    = RespCode("JSON_UNMARSHAL_ERROR")
)

var (
	errorText = map[RespCode]string{
		//Success responses
		Success:       "Success",
		CreateSuccess: "Successfully Created Record",
		UpdateSuccess: "Successfully Updated Record",
		DeleteSuccess: "Successfully Deleted Record",
		//Generic error Responses
		Generic: "Some thing may be wrong please try again",
		FAIL:    "Operation Failed please try again",
		//validation error
		Validation:   "validation failed",
		UnknownValue: "this value is not expected",
		EmptyFilter:  "the filter is empty",
		EmptyBody:    "the request body is empty",
		//operation errors
		NoRecordsUpdated:           "NO records were Updated",
		NoRecordsDeleted:           "NO records were Deleted",
		UpdatingAssociationsFailed: "Updating Associated Records Failed",
		//status code messages
		DataNotFound:  "data not found",
		BadRequest:    "bad request",
		InternalError: "internal server error",
		//

		//Auth Messages

		EmailOrPassword:     "Your email address or password is wrong",
		InvalidEmailOrPhone: " please enter a valid email address or phone number",
		InfoOrCode:          "please enter valid code or info",
		InvalidData:         "This Data Is Invalid",
		PhoneExists:         "Phone Already Exists",
		EmailExists:         "UserId Already Exists",
		Password:            "Password Could not be stored",
		CouldNotAssignRole:  "could not assign role to the user",
		TokenDontMatch:      "the Token Doesnt match",
		UserNotFound:        "user Not Found",
		InvalidToken:        "the token is not valid",
		//FILE
		ParseFile:             "Could not parse multipart form",
		InvalidFile:           "Invalid File",
		InvalidFileType:       "File Type Not Allowed",
		FileTooBig:            "File is bigger than max allowed size",
		MaxUploadSizeExceeded: `Total Image UploadFromPath Exceeded, more than 30MB`,
		CantReadFileType:      "Cant Read The File SupplierId",
		CantReadFile:          "Cant Read The File",
		CantWriteFile:         "Cant Write The File",

		ImageUploadError: "Image UploadFromPath Resp",
		ImageDeleteError: "IMAGE DELETE ERROR",
		ImageUpdateError: "Image UpdateForm Resp",

		//------------------------------ Entities -----------------
		UserCreated:    " User Created",
		UserNotCreated: "User Not Created",
		UserExists:     "This user already Exists",
		JsonUnmarshal:  "Json Unmarshal Resp",
	}
)

func ErrorText(ErrorCode RespCode) string {
	return errorText[ErrorCode]
}
