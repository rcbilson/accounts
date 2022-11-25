package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"knilson.org/accounts/account"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendError(ctx echo.Context, code int, message string) error {
	type Error struct {
		message string
	}
	err := ctx.JSON(code, Error{message})
	return err
}

func ifNotEmptyTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}

func ifNotEmptyInt(s string) *int {
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &i
}

func ifNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// (GET /transactions)
func (s *Server) GetApiTransactions(ctx echo.Context) error {
	acct, err := account.Open()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
	}
	query := account.QuerySpec{
		DateFrom:    ifNotEmptyTime(ctx.QueryParam("DateFrom")),
		DateUntil:   ifNotEmptyTime(ctx.QueryParam("DateUntil")),
		DescrLike:   ifNotEmpty(ctx.QueryParam("DescrLike")),
		Category:    ifNotEmpty(ctx.QueryParam("Category")),
		Subcategory: ifNotEmpty(ctx.QueryParam("Subcategory")),
		State:       ifNotEmpty(ctx.QueryParam("State")),
		Limit:       ifNotEmptyInt(ctx.QueryParam("Limit")),
		Offset:      ifNotEmptyInt(ctx.QueryParam("Offset")),
	}
	ch, err := acct.Query(query)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("error querying db: %v", err))
	}
	results := make([]account.Record, 0)
	for r := range ch {
		results = append(results, *r)
	}
	return ctx.JSON(http.StatusOK, results)
}

func emptyIfNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// (POST /transactions)
func (s *Server) insertOrUpdate(ctx echo.Context, r account.Record) error {
	acct, err := account.Open()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
	}
	err = acct.BeginUpdate()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error beginning transaction: %v", err))
	}
	defer acct.AbortUpdate()

	if r.Id == "" {
		_, err = acct.Insert(&r)
	} else {
		_, err = acct.Update(&r)
	}
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("error on update/insert: %v", err))
	}
	err = acct.CompleteUpdate()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error completing update: %v", err))
	}

	return ctx.NoContent(http.StatusNoContent)
}

// (POST /transactions)
func (s *Server) PostApiTransactions(ctx echo.Context) error {
	var r account.Record
	err := ctx.Bind(&r)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for Transaction")
	}
	return s.insertOrUpdate(ctx, r)
}

// (POST /transactions/id)
func (s *Server) PostApiTransactionsId(ctx echo.Context) error {
	id := ctx.Param("id")
	var r account.Record
	err := ctx.Bind(&r)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for Transaction")
	}
	if r.Id == "" {
		r.Id = id
	} else if r.Id != id {
		return sendError(ctx, http.StatusBadRequest, "ID must be empty or equal to path ID")
	}
	return s.insertOrUpdate(ctx, r)
}

// (DELETE /transactions/id)
func (s *Server) DeleteApiTransactionsId(ctx echo.Context) error {
	id := ctx.Param("id")
	acct, err := account.Open()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
	}
	err = acct.BeginUpdate()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error beginning transaction: %v", err))
	}
	defer acct.AbortUpdate()

	_, err = acct.Delete(id)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("error deleting transaction: %v", err))
	}
	err = acct.CompleteUpdate()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error completing update: %v", err))
	}
	return ctx.NoContent(http.StatusNoContent)
}

func main() {
	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	s := NewServer()

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/api/transactions", s.GetApiTransactions)
	e.POST("/api/transactions", s.PostApiTransactions)
	e.POST("/api/transactions/:id", s.PostApiTransactionsId)
	e.DELETE("/api/transactions/:id", s.DeleteApiTransactionsId)

	e.Static("/", "../../frontend/build")
	e.File("/", "../../frontend/build/index.html")

	// And we serve HTTP until the world ends.
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *port)))
}
