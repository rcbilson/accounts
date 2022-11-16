package account

//go:generate oapi-codegen -package account -o api.gen.go api.yaml

import (
	"database/sql"
	"errors"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
)

type Context struct {
	db             *sql.DB
	tx             *sql.Tx
	deleteStmt     *sql.Stmt
	updateStmt     *sql.Stmt
	learnStmt      *sql.Stmt
	insertStmt     *sql.Stmt
	findAmtStmt    *sql.Stmt
	findApproxStmt *sql.Stmt
}

type specification struct {
	DbFile string
}

func Open() (*Context, error) {
	var ctx Context

	var s specification
	err := envconfig.Process("accounts", &s)
	if err != nil {
		return nil, err
	}

	ctx.db, err = sql.Open("sqlite3", s.DbFile)
	if err != nil {
		return nil, err
	}

	return &ctx, nil
}

func (ctx *Context) Close() {
	ctx.db.Close()
}

func (ctx *Context) BeginUpdate() error {
	if ctx.tx != nil {
		return errors.New("Attempt to begin an update while an update is already in progress.")
	}
	var err error
	ctx.tx, err = ctx.db.Begin()
	return err
}

func (ctx *Context) AbortUpdate() {
	if ctx.tx != nil {
		ctx.tx.Rollback()
		ctx.tx = nil
	}
}

func (ctx *Context) CompleteUpdate() error {
	if ctx.tx == nil {
		return errors.New("Attempt to complete an update without an update in progress.")
	}
	err := ctx.tx.Commit()
	ctx.tx = nil
	return err
}

type Record struct {
	Id          string
	Date        time.Time
	Descr       string
	Amount      string
	Category    string
	Subcategory string
}
