package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	oapifilter "github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"knilson.org/accounts/account"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

// (GET /transactions)
func (s *Server) GetTransactions(ctx echo.Context, params account.GetTransactionsParams) error {
	log.Println("GetTransactions")
	return ctx.String(http.StatusOK, "GetTransactions")
}

// (POST /transactions)
func (s *Server) PostTransactions(ctx echo.Context) error {
	log.Println("PostTransactions")
	return ctx.String(http.StatusOK, "PostTransactions")
}

func main() {
	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	swagger, err := account.GetSwagger()
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
	account.RegisterHandlers(e, s)

	// And we serve HTTP until the world ends.
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *port)))
}
