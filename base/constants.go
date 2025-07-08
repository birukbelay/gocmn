package base

type NotificationType string

const (
	NotificationError   NotificationType = "error"
	NotificationSuccess NotificationType = "success"
	NotificationWarning NotificationType = "warning"
	NotificationInfo    NotificationType = "info"
)

type NotificationTarget string

const (
	GeneralNotification    NotificationTarget = "general"
	IndividualNotification NotificationTarget = "individual"
	GroupNotification      NotificationTarget = "group"
)
