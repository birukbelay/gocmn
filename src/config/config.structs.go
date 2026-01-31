package config

type S3Config struct {
	Endpoint         string `koanf:"S3_ENDPOINT"`
	S3WebEndpoint    string `koanf:"S3_WEB_ENDPOINT"`
	Region           string `koanf:"S3_REGION"`
	S3ForcePathStyle bool   `koanf:"S3_FORCE_PATH_STYLE"` // TODO: Remove?
	AccessKey        string `koanf:"S3_ACCESS_KEY_ID"`
	SecretKey        string `koanf:"S3_SECRET_ACCESS_KEY"`
	BucketName       string `koanf:"S3_BUCKET_NAME"`
}
type CloudinaryConfig struct {
	CloudinaryCloudName string `koanf:"CLOUDINARY_CLOUD_NAME"`
	CloudinaryAPIKey    string `koanf:"CLOUDINARY_API_KEY"`
	CloudinaryAPISecret string `koanf:"CLOUDINARY_API_SECRET"`
	CloudinaryFolder    string `koanf:"CLOUDINARY_FOLDER"`
}

type SqlDbConfig struct {
	DbName   string `koanf:"SQL_DB_NAME"`
	Username string `koanf:"SQL_USERNAME"`
	Password string `koanf:"SQL_PASSWORD"`
	SqlHost  string `koanf:"SQL_HOST"`
	SqlPort  string `koanf:"SQL_PORT"`
	Driver   string `koanf:"SQL_DRIVER"`
	SSLMode  string `koanf:"SSL_MODE"`
}
type KeyValConfig struct {
	KVDbName   int    `koanf:"KV_DB"`
	KVUsername string `koanf:"KV_USER"`
	KVPassword string `koanf:"KV_PASSWORD"`
	KVHost     string `koanf:"KV_HOST"`
	KVPort     string `koanf:"KV_PORT"`
}

type JwtVar struct {
	AccessSecret  string `koanf:"ACCESS_SECRET"`
	RefreshSecret string `koanf:"REFRESH_SECRET"`

	AccessExpireMin  int `koanf:"ACCESS_SECRET_EXPIRE_MIN"`
	RefreshExpireMin int `koanf:"REFRESH_SECRET_EXPIRES_MIN"`
}
type SmtpCred struct {
	SmtpHost     string `koanf:"SMTP_HOST"`
	SmtpPort     string `koanf:"SMTP_PORT"`
	SmtpPwd      string `koanf:"SMTP_PWD"`
	SmtpUsername string `koanf:"SMTP_USERNAME"`
}

type LocalFileConfig struct {
	FileServingUrl string `koanf:"FILE_SERVE_URL"`
	PanicFile      string `koanf:"PANIC_FILE"`
	FileUploadPath string `koanf:"FILE_SERVE_URL"`
	//IMage related errors
	ImageErrLogFile string `koanf:"UPLOAD_ERR_LOG"`
	ImageUploadPath string `koanf:"IMAGE_UPLOAD_PATH"`
	MaxUploadSize   string `koanf:"MAX_UPLOAD_SIZE"`
}

type FirebaseConfig struct {
	Type                    string `koanf:"FIREBASE_TYPE"`
	ProjectID               string `koanf:"FIREBASE_PROJECT_ID"`
	PrivateKey              string `koanf:"FIREBASE_PRIVATE_KEY"`
	PrivateKeyID            string `koanf:"FIREBASE_PRIVATE_KEY_ID"`
	ClientEmail             string `koanf:"FIREBASE_CLIENT_EMAIL"`
	ClientID                string `koanf:"FIREBASE_CLIENT_ID"`
	AuthURI                 string `koanf:"FIREBASE_AUTH_URI"`
	TokenUrl                string `koanf:"FIREBASE_TOKEN_URI"`
	AuthProviderX509CertUrl string `koanf:"FIREBASE_AUTH_PROVIDER_X509_CERT_URL"`
	ClientX509CertUrl       string `koanf:"FIREBASE_CLIENT_X509_CERT_URL"`
}
