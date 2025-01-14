package domain

type Connection interface {
    Connect() error
    Close() error
    IsConnected() bool
}