package pgx_repository

import (
	"errors"

	"github.com/homework3/notification/internal/database"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

func (r *repository) GetMail(ctx context.Context, userId int32) (string, error) {
	const query = `
		SELECT mail
		FROM mail
		WHERE user_id = $1
	`

	var mail string
	err := r.pool.QueryRow(ctx, query, userId).Scan(&mail)

	if errors.Is(err, pgx.ErrNoRows) {
		return mail, database.ErrNotFound
	}

	return mail, nil

}
