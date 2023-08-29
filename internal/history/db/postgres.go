package history

import (
	"avito-segment/internal/history"
	"avito-segment/pkg"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) history.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetByUserAndDate(ctx context.Context, userID uint, date time.Time) (*[]history.HistoryRecord, error) {
	q := `
    SELECT user_id, slug, action, created_at FROM history WHERE user_id = $1
    AND EXTRACT(YEAR FROM created_at) = $2
    AND EXTRACT(MONTH FROM created_at) = $3;`

	rows, err := r.db.Query(ctx, q, userID, date.Year(), date.Month())
	if err != nil {
		return nil, handleError(err)
	}

	records := make([]history.HistoryRecord, 0)
	for rows.Next() {
		var record history.HistoryRecord
		err = rows.Scan(
			&record.UserID,
			&record.Slug,
			&record.Action,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return &records, nil
}

func (r *repository) CreateBulk(ctx context.Context, records *history.CreateHistoryRecordParams) error {
	q := `INSERT INTO history (user_id, slug, action, created_at) VALUES ($1, $2, $3, $4);`

	batch := &pgx.Batch{}

	for _, slug := range *records.Slugs {
		batch.Queue(q, records.UserID, slug, records.Action, records.CreatedAt)
	}

	results := r.db.SendBatch(ctx, batch)
	defer results.Close()

	for range *records.Slugs {
		_, err := results.Exec()
		if err != nil {
			return handleError(err)
		}
	}

	return nil
}

func handleError(err error) error {
	if newErr, ok := err.(*pgconn.PgError); ok {
		detailedErr := pkg.PgErrToErr(*newErr)
		log.Println(detailedErr)
		return fmt.Errorf("Something went wrong")
	}
	return err
}
