package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/deepmap/oapi-codegen/examples/petstore-expanded/echo/api/models"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	oapifilter "github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"knilson.org/accounts/account"
	"knilson.org/accounts/account/api"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendError(ctx echo.Context, code int, message string) error {
	sendErr := models.Error{
		Message: message,
	}
	err := ctx.JSON(code, sendErr)
	return err
}

// (GET /transactions)
func (s *Server) GetTransactions(ctx echo.Context, params api.GetTransactionsParams) error {
	acct, err := account.Open()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
	}
	query := account.QuerySpec{
		DescrLike:   params.DescrLike,
		Category:    params.Category,
		Subcategory: params.Subcategory,
		State:       params.State,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}
	if params.DateFrom != nil {
		query.DateFrom = &params.DateFrom.Time
	}
	if params.DateUntil != nil {
		query.DateUntil = &params.DateUntil.Time
	}
	ch, err := acct.Query(query)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("error querying db: %v", err))
	}
	results := make([]api.Transaction, 0)
	for r := range ch {
		var t api.Transaction
		t.Id = &r.Id
		t.Date = &openapi_types.Date{r.Date}
		t.Description = &r.Descr
		t.Amount = &r.Amount
		t.Category = &r.Category
		t.Subcategory = &r.Subcategory
		results = append(results, t)
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
func (s *Server) PostTransactions(ctx echo.Context) error {
	log.Println("PostTransactions", ctx.Request().Body)
	var t api.Transaction
	err := ctx.Bind(&t)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for Transaction")
	}

	acct, err := account.Open()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
	}
        err = acct.BeginUpdate()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error beginning transaction: %v", err))
	}
        defer acct.AbortUpdate()

	var r account.Record
	if t.Date != nil {
                r.Date = t.Date.Time
        }
	r.Id = emptyIfNil(t.Id)
	r.Descr = emptyIfNil(t.Description)
	r.Amount = emptyIfNil(t.Amount)
	r.Category = emptyIfNil(t.Category)
	r.Subcategory = emptyIfNil(t.Subcategory)

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

func main() {
	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create an instance of our handler which satisfies the generated interface
	s := NewServer()

	// This is how you set up a basic Echo router
	e := echo.New()
	// Log all requests
	e.Use(echomiddleware.Logger())

	validatorOptions := &middleware.Options{}

	validatorOptions.Options.AuthenticationFunc = func(c context.Context, input *oapifilter.AuthenticationInput) error {
		return nil
	}

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	e.Use(middleware.OapiRequestValidatorWithOptions(swagger, validatorOptions))

	// We now register our petStore above as the handler for the interface
	api.RegisterHandlers(e, s)

	// And we serve HTTP until the world ends.
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *port)))
}
