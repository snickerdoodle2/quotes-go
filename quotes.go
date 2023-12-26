package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Router chi.Router
	DB     *pgxpool.Pool
}

type QuoteCreate struct {
	Book  string
	Quote string
}

type Quote struct {
	Id         uuid.UUID
	Book       string
	Quote      string
	InsertedAt time.Time
	UpdatedAt  time.Time
}

func newQuote(book string, quote string) Quote {
	now := time.Now().UTC()
	return Quote{
		Id:         uuid.New(),
		Book:       book,
		Quote:      quote,
		InsertedAt: now,
		UpdatedAt:  now,
	}

}

func (app *App) createQuote(w http.ResponseWriter, r *http.Request) {
	var q QuoteCreate
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quote := newQuote(q.Book, q.Quote)

    _, db_err := app.DB.Exec(context.Background(),
    `INSERT INTO quotes (id, book, quote, inserted_at, updated_at) 
    VALUES ($1, $2, $3, $4, $5)
    `, quote.Id, quote.Book, quote.Quote, quote.InsertedAt, quote.UpdatedAt)
    if db_err != nil {
        http.Error(w, db_err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(quote)
	w.WriteHeader(http.StatusCreated)
}

func (app *App) readQuotes(w http.ResponseWriter, r *http.Request) {
    var quotes []*Quote

    pgxscan.Select(context.Background(), app.DB, &quotes, `SELECT id, book, quote, inserted_at, updated_at FROM quotes`)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(quotes)
}

func (app *App) updateQuote(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var q QuoteCreate
	err = json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    now := time.Now().UTC()

    res, err := app.DB.Exec(context.Background(), `
    UPDATE quotes
    SET book = $1, quote = $2, updated_at = $3
    WHERE id = $4
    `, q.Book, q.Quote, now, id)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if res.RowsAffected() == 0 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (app *App) deleteQuote(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    res, err := app.DB.Exec(context.Background(), `DELETE FROM quotes WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }
    if res.RowsAffected() == 0 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)

}

func (app *App) closeConnection() {
    app.DB.Close()
}

func getApp(dbUrl string) (*App, error) {
	conn, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		return nil, err
	}

    r := chi.NewRouter()

	return &App{
		Router: r,
		DB:     conn,
	}, nil
}

func (app *App) mountHandlers() {
    app.Router.Post("/", app.createQuote)
    app.Router.Get("/", app.readQuotes)
    app.Router.Put("/{id}", app.updateQuote)
    app.Router.Delete("/{id}", app.deleteQuote)
}
