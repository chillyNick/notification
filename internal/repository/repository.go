package repository

import "golang.org/x/net/context"

type Repository interface {
	GetMail(ctx context.Context, userId int32) (string, error)
}
