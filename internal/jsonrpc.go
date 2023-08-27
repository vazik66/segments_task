package internal

import (
	"avito-segment/internal/segment"
	segments "avito-segment/internal/segment/db"
	segmentRpc "avito-segment/internal/segment/transport"

	usersegment "avito-segment/internal/user_segments"
	usersegments "avito-segment/internal/user_segments/db"
	usersegmentsRpc "avito-segment/internal/user_segments/transport"

	"avito-segment/internal/user"
	users "avito-segment/internal/user/db"
	userRpc "avito-segment/internal/user/transport"

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

type RPCServer struct {
	cfg    *Config
	rpc    *mux.Router
	pgConn *pgx.Conn
}

// @title Test Task
// @version 1.0
// @BasePath /
func NewRPCServer(cfg *Config) (*RPCServer, error) {
	pgConn, err := pkgDb.NewPostgresClient(context.Background(), cfg.Db.BuildDsn("postgresql"))
	if err != nil {
		log.Fatal(err)
	}

	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")

	if err = initServices(s, pgConn); err != nil {
		log.Fatalf("Could not initialize services. %v", err)
	}

	mux := mux.NewRouter()
	mux.Handle("/jsonrpc", s)
	mux.Handle("/jsonrpc/{.+$}", s) // Swagger workaround
	mux.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", cfg.Port)),
	))

	return &RPCServer{
		cfg:    cfg,
		rpc:    mux,
		pgConn: pgConn,
	}, nil
}

func (s *RPCServer) Run() error {
	defer s.pgConn.Close(context.Background())
	log.Println("Starting server...")
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", s.cfg.Port),
		s.rpc,
	)
	return err
}

func initServices(server *rpc.Server, dbConn *pgx.Conn) error {
	userRepo := users.NewRepository(dbConn)
	userService := user.NewService(userRepo)
	userRpcHandler := userRpc.NewHandler(userService)

	segmentRepo := segments.NewRepository(dbConn)
	segmentService := segment.NewService(segmentRepo)
	segmentRpcHandler := segmentRpc.NewHandler(segmentService)

	userSegmentsRepo := usersegments.NewRepository(dbConn)
	userSegmentsService := usersegment.NewService(userSegmentsRepo, userRepo)
	userSegmentsHandler := usersegmentsRpc.NewHandler(userSegmentsService)

	if err := server.RegisterService(userRpcHandler, "users"); err != nil {
		return err
	}

	if err := server.RegisterService(segmentRpcHandler, "segments"); err != nil {
		return err
	}

	if err := server.RegisterService(userSegmentsHandler, "usersegments"); err != nil {
		return err
	}

	return nil
}
