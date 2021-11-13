package mysql

import (
	"context"
	"database/sql"

	"github.com/Lambels/usho/encoding/base64"
	"github.com/Lambels/usho/repo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type storage struct {
	db *sqlx.DB

	stmtCreate *sql.Stmt
	stmtGet    *sql.Stmt
}

func New(dsn string) (repo.Repo, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS url (
		Id INT,
		Initial INT,
		PRIMARY KEY(id))
	`); err != nil {
		return nil, err
	}

	qCreate, err := db.Prepare("INSERT INTO url(Initial) VALUES (?)")
	if err != nil {
		return nil, err
	}

	qGet, err := db.Prepare("SELECT * WHERE Id = ?")
	if err != nil {
		return nil, err
	}

	return &storage{
		db:         db,
		stmtCreate: qCreate,
		stmtGet:    qGet,
	}, nil
}

func (s *storage) New(ctx context.Context, in repo.URLRequest) (repo.URLResponse, error) {
	var out repo.URLResponse

	rsp, err := s.stmtCreate.ExecContext(ctx, in.Intial)
	if err != nil {
		return out, err
	}

	id, err := rsp.LastInsertId()
	if err != nil {
		return out, err
	}

	out.ID = uint64(id)
	out.Initial = in.Intial
	out.Short = base64.EncodeID(out.ID)

	return out, nil
}

func (s *storage) Get(ctx context.Context, in string) (repo.URLResponse, error) {
	var out repo.URLResponse

	inID, err := base64.DecodeID(in)
	if err != nil {
		return out, err
	}

	rws, err := s.stmtGet.QueryContext(ctx, inID)
	if err != nil {
		return out, err
	}

	rws.Next()
	if err := rws.Scan(&out); err != nil {
		return out, err
	}

	return out, rws.Close()
}

func (s *storage) Close() error {
	if err := s.stmtCreate.Close(); err != nil {
		return err
	}

	if err := s.stmtGet.Close(); err != nil {
		return err
	}

	return s.db.Close()
}
