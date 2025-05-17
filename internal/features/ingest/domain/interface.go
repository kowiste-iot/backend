package domain

type Consumer interface {
	Start() error
	Stop() error
	Subscribe(handler func(msg *Message) error) error
}

type Publisher interface {
	Publish(topic string, msg *Message) error
	Close() error
}
