package domain

type MessageHandler interface {
    Handle(message Message) error
}