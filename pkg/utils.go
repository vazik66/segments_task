package pkg

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type EmptyArgs struct{}

type EmptyResponse struct{}

type JsonRPCSuccessResponse struct {
	JsonRPC string      `json:"jsonrpc" example:"2.0"`
	Result  interface{} `json:"result"`
	ID      int         `json:"id" example:"1"`
}

type JsonRPCErrorResponse struct {
	JsonRPC string `json:"jsonrpc" example:"2.0"`
	Error   Error  `json:"error"`
	ID      int    `json:"id" example:"1"`
}

type JsonRPCRequest struct {
	JsonRPC string      `json:"jsonrpc" example:"2.0"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id" example:"1"`
}

func PgErrToErr(err pgconn.PgError) error {
	return fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", err.Message, err.Detail, err.Where, err.Code, err.SQLState())
}
