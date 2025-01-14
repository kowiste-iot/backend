package domain

import "context"

type Publisher interface {
    Publish(ctx context.Context, subjectGen SubjectGenerator, msg Message, params ...string) error
}

type Subscriber interface {
    Publisher
    Subscribe(ctx context.Context, subjectGen SubjectGenerator, handler MessageHandler, params ...string) error
    Unsubscribe(subjectGen SubjectGenerator, params ...string) error
    Close() error
    HandleReconnect(ctx context.Context)
}