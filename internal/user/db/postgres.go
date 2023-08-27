package user

import (
	"avito-segment/internal/user"
	"avito-segment/pkg"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) user.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context) (*user.User, error) {
	q := `INSERT INTO users DEFAULT VALUES RETURNING id;`
	usr := new(user.User)

	if err := r.db.QueryRow(ctx, q).Scan(&usr.ID); err != nil {
		return &user.User{}, handleError(err)
	}

	return usr, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	q := `DELETE FROM users WHERE users.id = $1;`

	tag, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return handleError(err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("Not found")
	}

	return nil

}

func (r *repository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	q := `SELECT id FROM users WHERE users.id = $1;`

	usr := new(user.User)
	err := r.db.QueryRow(ctx, q, id).Scan(&usr.ID)
	if err != nil {
		return &user.User{}, handleError(err)
	}

	return usr, nil
}

func (r *repository) List(ctx context.Context) (*[]user.User, error) {
	q := `SELECT id FROM users;`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, handleError(err)
	}

	users := make([]user.User, 0)
	for rows.Next() {
		var usr user.User

		err = rows.Scan(&usr.ID)
		if err != nil {
			return nil, err
		}

		users = append(users, usr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func handleError(err error) error {
	if newErr, ok := err.(*pgconn.PgError); ok {
		detailedErr := pkg.PgErrToErr(*newErr)
		log.Println(detailedErr)
		return fmt.Errorf("Something went wrong")
	}

	if err.Error() == "no rows in result set" {
		return fmt.Errorf("Not found")
	}

	return err
}
