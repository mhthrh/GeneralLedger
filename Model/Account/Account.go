package Account

import (
	"GitHub.com/mhthrh/GL/Util/DbUtil/PgSql"
	"database/sql"
	"fmt"
	pb "github.com/pborman/uuid"
)

var (
	Accounts []*Account
)

type Account struct {
	ID      pb.UUID `json:"id" validate:"required,gt=0,uuid4"`
	UserId  pb.UUID `json:"userId" validate:"required,gt=0,uuid4"`
	AccNo   string  `json:"accNo" validate:"required,gt=0"`
	Balance float64 `json:"balance" validate:"required,gt=0"`
}

type tool struct {
	db *sql.DB
}

func New(db *sql.DB) *tool {
	return &tool{db: db}
}

func (t *tool) Load(userId pb.UUID) ([]Account, error) {

	var account Account
	var accounts []Account
	rows, err := PgSql.RunQuery(t.db, fmt.Sprintf("SELECT \"ID\", \"UserId\", \"AccNo\", \"Balance\" FROM public.\"Accounts\" where \"UserId\"='%s'", userId))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rows.Scan(&account.ID, &account.UserId, &account.AccNo, &account.Balance)
		accounts = append(accounts, account)
	}

	return accounts, nil
}
