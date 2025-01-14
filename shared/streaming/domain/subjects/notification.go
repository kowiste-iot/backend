package subjects

import "fmt"

type NotificationSubject struct{}

func NewNotificationSubject() *NotificationSubject {
    return &NotificationSubject{}
}

func (s *NotificationSubject) Generate(params ...string) string {
    if len(params) != 2 {
        return "notifications.*.*"
    }
    return fmt.Sprintf("notifications.%s.%s", params[0], params[1]) // tenant.userID
}