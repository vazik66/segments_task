package usersegments

import (
	usersegments "avito-segment/internal/user_segments"
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

func NewRepository(db *pgx.Conn) usersegments.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetByUserID(ctx context.Context, userID uint) (*[]usersegments.UserSegment, error) {
	q := `SELECT 
    user_segments.user_id, segments.slug, user_segments.ttl, user_segments.created_at
    FROM user_segments
    INNER JOIN segments ON user_segments.segment_id = segments.id 
    WHERE user_id = $1;`

	rows, err := r.db.Query(ctx, q, userID)

	if err != nil {
		return nil, handleError(err)
	}

	segments := make([]usersegments.UserSegment, 0)
	for rows.Next() {
		var seg usersegments.UserSegment

		err = rows.Scan(&seg.UserID, &seg.SegmentSlug, &seg.Ttl, &seg.CreatedAt)
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

func (r *repository) AddSegmentsToUser(ctx context.Context, userID uint, slugs *[]string, ttl uint) error {
	q := `INSERT INTO user_segments (user_id, segment_id, ttl)
    VALUES ($1, (SELECT segments.id FROM segments WHERE segments.slug = $2), $3);`

	batch := &pgx.Batch{}
	for _, slug := range *slugs {
		batch.Queue(q, userID, slug, ttl)
	}

	results := r.db.SendBatch(ctx, batch)
	defer results.Close()

	for range *slugs {
		_, err := results.Exec()
		if err != nil {
			return handleError(err)
		}

	}

	return nil
}

func (r *repository) RemoveSegmentsFromUser(ctx context.Context, userID uint, slugs *[]string) error {
	q := `DELETE FROM user_segments
    WHERE user_segments.user_id = $1
    AND user_segments.segment_id = (SELECT id FROM segments WHERE slug = $2);`

	batch := &pgx.Batch{}
	for _, slug := range *slugs {
		batch.Queue(q, userID, slug)
	}

	results := r.db.SendBatch(ctx, batch)
	defer results.Close()

	for range *slugs {
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
		switch newErr.Code {
		case pgerrcode.UniqueViolation:
			return fmt.Errorf("Segment with such slug already added to this user")
		case pgerrcode.ForeignKeyViolation:
			return fmt.Errorf("User does not exists")
		case pgerrcode.NotNullViolation:
			return fmt.Errorf("Segment with slug does not exists")
		default:
			return fmt.Errorf("Something went wrong")
		}
	}
	return err
}
