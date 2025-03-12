package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kelseyhightower/envconfig"
	"knilson.org/accounts/account"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func sendError(w http.ResponseWriter, code int, message string) {
	type Error struct {
		Message string `json:"message"`
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Error{message})
}

func ifNotEmptyDate(s string) *account.Date {
	if s == "" {
		return nil
	}
	d, err := account.ParseDate(s)
	if err != nil {
		return nil
	}
	return &d
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

func (s *Server) GetApiTransactions(w http.ResponseWriter, r *http.Request) {
	acct, err := account.Open()
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
		return
	}
	query := account.QuerySpec{
		DateFrom:    ifNotEmptyDate(r.URL.Query().Get("DateFrom")),
		DateUntil:   ifNotEmptyDate(r.URL.Query().Get("DateUntil")),
		DescrLike:   ifNotEmpty(r.URL.Query().Get("DescrLike")),
		Category:    ifNotEmpty(r.URL.Query().Get("Category")),
		Subcategory: ifNotEmpty(r.URL.Query().Get("Subcategory")),
		State:       ifNotEmpty(r.URL.Query().Get("State")),
		Limit:       ifNotEmptyInt(r.URL.Query().Get("Limit")),
		Offset:      ifNotEmptyInt(r.URL.Query().Get("Offset")),
	}
	ch, err := acct.Query(query)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("error querying db: %v", err))
		return
	}
	results := make([]account.Record, 0)
	for r := range ch {
		results = append(results, *r)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (s *Server) GetApiCategories(w http.ResponseWriter, r *http.Request) {
	acct, err := account.Open()
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
		return
	}
	query := account.QuerySpec{
		DateFrom:  ifNotEmptyDate(r.URL.Query().Get("DateFrom")),
		DateUntil: ifNotEmptyDate(r.URL.Query().Get("DateUntil")),
		Limit:     ifNotEmptyInt(r.URL.Query().Get("Limit")),
		Offset:    ifNotEmptyInt(r.URL.Query().Get("Offset")),
	}
	ch, err := acct.AggregateCategories(query)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("error querying db: %v", err))
		return
	}
	results := make([]account.Aggregate, 0)
	for r := range ch {
		results = append(results, *r)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (s *Server) GetApiSummary(w http.ResponseWriter, r *http.Request) {
	acct, err := account.Open()
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
		return
	}
	query := account.QuerySpec{
		DateFrom:  ifNotEmptyDate(r.URL.Query().Get("DateFrom")),
		DateUntil: ifNotEmptyDate(r.URL.Query().Get("DateUntil")),
		Limit:     ifNotEmptyInt(r.URL.Query().Get("Limit")),
		Offset:    ifNotEmptyInt(r.URL.Query().Get("Offset")),
	}
	results, err := acct.Summary(query)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("error querying db: %v", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (s *Server) GetApiSummaryChart(w http.ResponseWriter, r *http.Request) {
	acct, err := account.Open()
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
		return
	}
	query := account.QuerySpec{
		DateFrom:  ifNotEmptyDate(r.URL.Query().Get("DateFrom")),
		DateUntil: ifNotEmptyDate(r.URL.Query().Get("DateUntil")),
		Limit:     ifNotEmptyInt(r.URL.Query().Get("Limit")),
		Offset:    ifNotEmptyInt(r.URL.Query().Get("Offset")),
	}
	results, err := acct.SummaryChart(query)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("error querying db: %v", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func (s *Server) PostApiImport(w http.ResponseWriter, r *http.Request) {
	if !contains(r.Header["Content-Type"], "text/csv") {
		sendError(w, http.StatusUnsupportedMediaType, "this route only accepts text/csv")
		return
	}
	acct, err := account.Open()
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
		return
	}
	_, err = acct.Import(r.Body)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error importing file: %v", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) insertOrUpdate(w http.ResponseWriter, _ *http.Request, rec account.Record) {
	acct, err := account.Open()
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
		return
	}
	err = acct.BeginUpdate()
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error beginning transaction: %v", err))
		return
	}
	defer acct.AbortUpdate()

	if rec.Id == "" {
		_, err = acct.Insert(&rec)
	} else {
		_, err = acct.Update(&rec)
	}
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("error on update/insert: %v", err))
		return
	}
	_, err = acct.UpdateLearning(&rec)
	if err != nil {
		log.Println("error updating learning:", err)
	}
	err = acct.CompleteUpdate()
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error completing update: %v", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) PostApiTransactions(w http.ResponseWriter, r *http.Request) {
	var rec account.Record
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid format for Transaction")
		return
	}
	s.insertOrUpdate(w, r, rec)
}

func (s *Server) PostApiTransactionsId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var rec account.Record
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid format for Transaction")
		return
	}
	if rec.Id == "" {
		rec.Id = id
	} else if rec.Id != id {
		sendError(w, http.StatusBadRequest, "ID must be empty or equal to path ID")
		return
	}
	s.insertOrUpdate(w, r, rec)
}

func (s *Server) DeleteApiTransactionsId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	acct, err := account.Open()
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error connecting to db: %v", err))
		return
	}
	err = acct.BeginUpdate()
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error beginning transaction: %v", err))
		return
	}
	defer acct.AbortUpdate()

	_, err = acct.Delete(id)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("error deleting transaction: %v", err))
		return
	}
	err = acct.CompleteUpdate()
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("error completing update: %v", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type specification struct {
	Port         int `default:"9000"`
	FrontendPath string
}

func main() {
	var spec specification
	err := envconfig.Process("accountserver", &spec)
	if err != nil {
		log.Fatal("error reading environment variables:", err)
	}

	s := NewServer()

	http.HandleFunc("GET /api/transactions", s.GetApiTransactions)
	http.HandleFunc("POST /api/transactions", s.PostApiTransactions)
	http.HandleFunc("POST /api/transactions/{id}", s.PostApiTransactionsId)
	http.HandleFunc("DELETE /api/transactions/{id}", s.DeleteApiTransactionsId)
	http.HandleFunc("POST /api/import", s.PostApiImport)
	http.HandleFunc("GET /api/categories", s.GetApiCategories)
	http.HandleFunc("GET /api/summary", s.GetApiSummary)
	http.HandleFunc("GET /api/summaryChart", s.GetApiSummaryChart)

	// bundled assets and static resources
	http.Handle("GET /assets/", http.FileServer(http.Dir(spec.FrontendPath)))
	http.Handle("GET /static/", http.FileServer(http.Dir(spec.FrontendPath)))
	// For other requests, serve up the frontend code
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", spec.FrontendPath))
	})

	log.Println("server listening on port", spec.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", spec.Port), nil))
}
