package domain

type Connection interface {
	Send(message []byte) error
	Close() error
	GetTenantID() string
	GetUserID() string
}
