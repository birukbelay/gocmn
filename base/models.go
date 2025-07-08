package base

import (
	"github.com/lib/pq"
	"time"
)

type Upload struct {
	Base      `mapstructure:",squash" `
	UploadDto `mapstructure:",squash" `
}
type UploadDto struct {
	Name   string `gorm:"index:,unique;not null" json:"name,omitempty"`
	Url    string `json:"url"`
	UserId string `json:"user_id,omitempty"`
	Path   string `json:"-"`
	//Todo the files hash
	Hash     string `json:"hash"`
	FileType string `json:"file_type"`
	Size     int64
}
type UploadQuery struct {
	ID     string `json:"id" form:"id"`
	Name   string `json:"name,omitempty" form:"name"`
	Hash   string `json:"hash" form:"hash"`
	UserId string `json:"user_id" form:"user_id"`
}

// func (m *Upload) BeforeCreate(tx *gorm.DB) (err error) {
// 	if m.ID == "" {
// 		m.ID = primitive.NewObjectID().Hex()
// 	}
// 	return nil
// }

type Notice struct {
	Base      `mapstructure:",squash" `
	NoticeDto `mapstructure:",squash" `
}
type NoticeDto struct {
	Title     string             `json:"title,omitempty"`
	Body      string             `json:"body"`
	UserId    *string            `json:"user_id,omitempty"`
	Type      NotificationType   `json:"notificationType,omitempty"`
	Target    NotificationTarget `json:"notificationTarget,omitempty"`
	NoticeFor pq.StringArray     `json:"noticeFor,omitempty" gorm:"type:text[]"`
	//Level NotificationType   `json:"status,omitempty"`

	//Inc        int        `json:"inc" gorm:"autoIncrement:true;unique"`
	Seen bool `json:"seen,omitempty"`
}
type NoticeFilter struct {
	ID     string             `query:"id"`
	Title  string             `query:"title"`
	UserId string             `query:"user_id"`
	Type   NotificationTarget `query:"notification_type"`
	Status NotificationType   `query:"status"`

	Seen bool `query:"seen"`
}
type NoticeQuery struct {
	CreatedAt      time.Time `query:"created_at,omitempty"`
	SelectedFields []string  `query:"selected_fields" enum:"title,body,user_id,type,status,module,seen,id,created_at,updated_at"`
	Sort           string    `query:"sort" enum:"title,body,user_id,type,status,module,seen,id,created_at,updated_at"`
}

// TagDto is used for updating the tags, of data
type TagDto struct {
	Tags pq.StringArray `msgpack:"-" gorm:"type:text[]" json:"tags,omitempty" swaggertype:"array,string"`
	Tag  string         `msgpack:"-" json:"tag,omitempty"`
}

type LogTrace struct {
	Base        `mapstructure:",squash" `
	LogTraceDto `mapstructure:",squash" `
}

type LogType string

const (
	ImplantLog = LogType("implantLog")
	GenericLog = LogType("genericLog")
)

type LogLevel string

const (
	LogDebug   = LogLevel("Debug")
	LogInfo    = LogLevel("Info")
	LogWarning = LogLevel("Warning")
	LogError   = LogLevel("Error")
	LogFatal   = LogLevel("Fatal")
)

type LogTraceDto struct {
	Title     string   `json:"title,omitempty"`
	Message   string   `json:"message,omitempty"`
	Type      LogType  `json:"type,omitempty"`
	ImplantID *string  `json:"implant_id,omitempty"`
	Function  string   `json:"function,omitempty"`
	Location  string   `json:"location,omitempty"`
	Level     LogLevel `json:"level,omitempty" enum:"Debug,Info,Warning,Error,Fatal"`
}
type LogTraceFilter struct {
	ID       string  `query:"id,omitempty"`
	Title    string  `query:"title,omitempty"`
	Message  string  `query:"message,omitempty"`
	Type     LogType `query:"type,omitempty" enum:"implantLog,genericLog"`
	Function string  `query:"function,omitempty"`
	Location string  `query:"location,omitempty"`
	Level    string  `query:"level,omitempty" enum:"Debug,Info,Warning,Error,Fatal"`
}
type LogTraceQuery struct {
	CreatedAt      time.Time `query:"created_at,omitempty"`
	SelectedFields []string  `query:"selected_fields" enum:"title,message,implant_id,function,location,status,id,created_at,updated_at"`
	Sort           string    `query:"sort" enum:"title,message,implant_id,function,location,status,id,created_at,updated_at"`
}
