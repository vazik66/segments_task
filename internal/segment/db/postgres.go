package segment

import (
	"avito-segment/internal/segment"
	"avito-segment/pkg"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) segment.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, segmentParams *segment.CreateSegmentParams) (*segment.Segment, error) {
	q := `INSERT INTO segments (slug) VALUES ($1) RETURNING slug;`

	seg := new(segment.Segment)
	if err := r.db.QueryRow(ctx, q, segmentParams.Slug).Scan(&seg.Slug); err != nil {
		return &segment.Segment{}, handleError(err)
	}

	return seg, nil
}

func (r *repository) Delete(ctx context.Context, slug string) error {
	q := `DELETE FROM segments WHERE segments.slug = $1;`

	tag, err := r.db.Exec(ctx, q, slug)
	if err != nil {
		return handleError(err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("Not found")
	}

	return err
}

func (r *repository) GetBySlug(ctx context.Context, slug string) (*segment.Segment, error) {
	q := `SELECT slug FROM segments WHERE segments.slug = $1;`

	seg := new(segment.Segment)
	err := r.db.QueryRow(ctx, q, slug).Scan(&seg.Slug)
	if err != nil {
		return &segment.Segment{}, handleError(err)
	}

	return seg, nil
}

func (r *repository) List(ctx context.Context) (*[]segment.Segment, error) {
	q := `SELECT slug FROM segments;`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, handleError(err)
	}

	segments := make([]segment.Segment, 0)
	for rows.Next() {
		var seg segment.Segment

		err = rows.Scan(&seg.Slug)
		if err != nil {
			return nil, err
		}

		segments = append(segments, seg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &segments, nil
}

func (r *repository) Update(ctx context.Context, segment *segment.Segment) error {
	return nil
}

func handleError(err error) error {
	if newErr, ok := err.(*pgconn.PgError); ok {
		detailedErr := pkg.PgErrToErr(*newErr)
		log.Println(detailedErr)

		switch newErr.Code {
		case pgerrcode.UniqueViolation:
			return fmt.Errorf("Segment with such slug already exists")
		default:
			return fmt.Errorf("Something went wrong")
		}
	}
	if err.Error() == "no rows in result set" {
		return fmt.Errorf("Not found")
	}
	return err
}
