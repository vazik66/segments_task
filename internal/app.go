package internal

import (
	"avito-segment/internal/segment"
	segmentDb "avito-segment/internal/segment/db"
	segmentRpc "avito-segment/internal/segment/transport"
	"avito-segment/pkg/events"
	"time"

	usersegment "avito-segment/internal/user_segments"
	usersegmentDb "avito-segment/internal/user_segments/db"
	usersegmentRpc "avito-segment/internal/user_segments/transport"

	"avito-segment/internal/user"
	userDb "avito-segment/internal/user/db"
	userRpc "avito-segment/internal/user/transport"

	"avito-segment/internal/history"
	historyDb "avito-segment/internal/history/db"
	historyTransport "avito-segment/internal/history/transport"

	_ "avito-segment/docs"
	pkgDb "avito-segment/pkg/db"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/jackc/pgx/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type App struct {
	cfg    *Config
	mux    *mux.Router
	pgConn *pgx.Conn
	em     events.EventManager
}

// @title		Test Task
// @version	1.0
// @BasePath	/
func NewApp(ctx context.Context, cfg *Config) (*App, error) {
	pgConn, err := pkgDb.NewPostgresClient(ctx, cfg.Db.BuildDsn("postgresql"))
	if err != nil {
		log.Fatal(err)
	}

	em := events.NewLocalEventManager()
	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")

	if err = initServices(s, pgConn, em); err != nil {
		log.Fatalf("Could not initialize services. %v", err)
	}

	mux := mux.NewRouter()
	mux.Handle("/jsonrpc", s)
	mux.Handle("/jsonrpc/{.+$}", s) // Swagger workaround
	mux.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", cfg.Port)),
	))
	mux.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))

	return &App{
		cfg:    cfg,
		mux:    mux,
		pgConn: pgConn,
		em:     em,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	ctxLocal, cancel := context.WithCancel(ctx)
	defer cancel()
	defer app.pgConn.Close(ctxLocal)

	log.Println("Starting server...")
	app.startScheduler(ctxLocal)
	srv := &http.Server{
		Handler:      app.mux,
		Addr:         fmt.Sprintf(":%s", app.cfg.Port),
		WriteTimeout: 6 * time.Second,
		ReadTimeout:  6 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func initServices(server *rpc.Server, dbConn *pgx.Conn, em events.EventManager) error {
	userRepo := userDb.NewRepository(dbConn)
	userService := user.NewService(userRepo)
	userRpcHandler := userRpc.NewHandler(userService)

	segmentRepo := segmentDb.NewRepository(dbConn)
	segmentService := segment.NewService(segmentRepo)
	segmentRpcHandler := segmentRpc.NewHandler(segmentService)

	userSegmentsRepo := usersegmentDb.NewRepository(dbConn)
	userSegmentsService := usersegment.NewService(userSegmentsRepo, userRepo, em)
	userSegmentsRpcHandler := usersegmentRpc.NewHandler(userSegmentsService)

	historyRepo := historyDb.NewRepository(dbConn)
	historyService := history.NewService(historyRepo, em)
	historyRpcHandler := historyTransport.NewHandler(historyService)

	if err := server.RegisterService(userRpcHandler, "users"); err != nil {
		return err
	}

	if err := server.RegisterService(segmentRpcHandler, "segments"); err != nil {
		return err
	}

	if err := server.RegisterService(userSegmentsRpcHandler, "usersegments"); err != nil {
		return err
	}

	if err := server.RegisterService(historyRpcHandler, "history"); err != nil {
		return err
	}

	return nil
}

func (app *App) startScheduler(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				_ = app.em.Publish("task.deleteUserSegmentsTTL", nil)
			case <-ctx.Done():
				return
			}

		}
	}()
}
