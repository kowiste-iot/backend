package domain

type SubjectGenerator interface {
    Generate(params ...string) string
}