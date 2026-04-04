package idgen

import "github.com/google/uuid"

type Generator interface {
	Generate() string
}

func UUIDV4() Generator {
	return &uuidV4Impl{}
}

type uuidV4Impl struct{}

func (g *uuidV4Impl) Generate() string {
	return uuid.New().String()
}
