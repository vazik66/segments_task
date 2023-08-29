package history

import (
	usersegments "avito-segment/internal/user_segments"
	"avito-segment/pkg/events"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

const (
	host          = "localhost"
	port          = "8080"
	filesDir      = "files"
	filesEndpoint = "files"
)

type HistoryService struct {
	repo   Repository
	events events.EventManager
}

func NewService(historyRepo Repository, events events.EventManager) *HistoryService {
	s := &HistoryService{
		repo:   historyRepo,
		events: events,
	}
	s.postInit()
	return s
}

func (s *HistoryService) postInit() {
	s.events.Subscribe("usersegments.deleted", func(i interface{}) error {
		log.Printf("Event: usersegments.deleted, Started, data: %+v", i)

		var params CreateHistoryRecordParams
		err := mapstructure.Decode(i, &params)
		if err != nil {
			return err
		}
		params.Action = Removed

		err = s.CreateBulk(context.Background(), &params)
		log.Printf("Event: usersegments.deleted, Finished")
		return err
	})

	s.events.Subscribe("usersegments.ttlDeleted", func(i interface{}) error {
		log.Printf("Event: usersegments.ttlDeleted, Started, data: %+v", i)

		var deleted *[]usersegments.DeadUserSegment
		err := mapstructure.Decode(i, &deleted)
		if err != nil {
			return err
		}

		userSlugsMap := make(map[uint][]string)
		for _, del := range *deleted {
			userSlugsMap[del.UserID] = append(userSlugsMap[del.UserID], del.Slug)
		}

		for userID, slugs := range userSlugsMap {
			params := CreateHistoryRecordParams{
				UserID:    userID,
				Slugs:     &slugs,
				CreatedAt: (*deleted)[0].CreatedAt,
				Action:    Removed,
			}
			err = s.CreateBulk(context.Background(), &params)
			if err != nil {
				log.Printf("%+v, %v", params, err)
			}
		}

		log.Printf("Event: usersegments.ttlDeleted, Finished")
		return nil
	})

	s.events.Subscribe("usersegments.added", func(i interface{}) error {
		log.Printf("Event: usersegments.added, Started, data: %+v", i)

		var params CreateHistoryRecordParams
		err := mapstructure.Decode(i, &params)
		if err != nil {
			return err
		}
		params.Action = Added

		err = s.CreateBulk(context.Background(), &params)
		log.Printf("Event: usersegments.added, Finished")
		return err
	})

	s.events.Subscribe("task.create_report", func(i interface{}) error {
		log.Printf("Task: create_report, Started, data: %+v", i)
		params, ok := i.(*CreateReportEventParams)
		if !ok {
			err := fmt.Errorf("could not type cast to CreateReportEventParams")
			log.Println(err)
			return err
		}

		err := s.writeToCsv(context.Background(), params)
		log.Printf("Task: create_report, Finished")
		return err
	})
}

func (s *HistoryService) CreateReport(ctx context.Context, userID uint, date time.Time) (string, error) {
	filename := uuid.New().String() + ".csv"

	// using events as bg task)
	_ = s.events.Publish(
		"task.create_report",
		&CreateReportEventParams{
			Filename: filename,
			UserID:   userID,
			Date:     date,
		},
	)
	url := fmt.Sprintf("http://%s:%s/%s/%s", host, port, filesEndpoint, filename)
	return url, nil
}

func (s *HistoryService) CreateBulk(ctx context.Context, params *CreateHistoryRecordParams) error {
	err := s.repo.CreateBulk(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (s *HistoryService) writeToCsv(ctx context.Context, params *CreateReportEventParams) error {
	records, err := s.repo.GetByUserAndDate(ctx, params.UserID, params.Date)
	if err != nil {
		return err
	}

	err = createDir(filesDir)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", filesDir, params.Filename)
	vals := historyToSlice(records)
	err = saveToCsvFile(vals, path)
	if err != nil {
		log.Printf("Error writing to file: %v. %v", path, err)
		return err
	}

	return nil
}

func saveToCsvFile(data *[][]string, path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	err = w.WriteAll(*data)
	if err != nil {
		return err
	}
	return nil
}

func createDir(name string) error {
	_, err := os.Stat(name)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(filesDir, 0777)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func historyToSlice(h *[]HistoryRecord) *[][]string {
	vals := [][]string{{"UserID", "Slug", "Action", "CreatedAt"}}

	for _, val := range *h {
		v := []string{fmt.Sprintf("%d", val.UserID), string(val.Slug), string(val.Action), val.CreatedAt.Format(time.RFC3339)}
		vals = append(vals, v)
	}

	return &vals
}
