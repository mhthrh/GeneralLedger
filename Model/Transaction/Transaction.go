package Transaction

import (
	"GitHub.com/mhthrh/GL/Util/DbUtil/PgSql"
	"database/sql"
	"fmt"
	"github.com/pborman/uuid"
	"time"
)

var Transactions []*Transaction

type Transaction struct {
	Account string  `json:"account" validate:"required"`
	Amount  float64 `json:"amount" validate:"required"`
	Note    string  `json:"note"`
	Action  int     `json:"action" validate:"required,gte=1,lte=2"`
}
type Request struct {
	Transaction Transaction `json:"transaction"`
	Sign        string      `json:"sign"`
}

type Search struct {
	Account  string `json:"account" validate:"required"`
	FromDate string `json:"fromDate" validate:"required"`
	ToDate   string `json:"toDate" validate:"required"`
	Sign     string `json:"sign"`
}
type tool struct {
	db *sql.DB
}

func New(db *sql.DB) *tool {
	return &tool{db: db}
}

func (to *tool) Create(t *Transaction, ch *chan bool, response *chan error) {
	var billan string
	transaction, err := to.db.Begin()
	if err != nil {
		*response <- err
		return
	}
	commit := false
	defer func(tx *sql.Tx, commit *bool) {
		select {
		case commit := <-*ch:
			if !commit {
				tx.Rollback()
				*response <- err
				return
			}
			tx.Commit()
			return

		}

	}(transaction, &commit)
	id := uuid.NewUUID()

	transaction.QueryRow(fmt.Sprintf("SELECT  \"AccNo\" FROM public.\"Accounts\" where type=%d", t.Action)).Scan(&billan)
	switch t.Action {
	case 1: //deposit
		{
			if err := insertRow(transaction, id, billan, t.Note, -t.Amount); err != nil {
				*ch <- false
				*response <- err
				return
			}
			if err := insertRow(transaction, id, t.Account, t.Note, t.Amount); err != nil {
				*ch <- false
				*response <- err
				return
			}

		}
	case 2: //withdrew
		{
			if err := insertRow(transaction, id, t.Account, t.Note, -t.Amount); err != nil {
				*ch <- false
				*response <- err
				return
			}
			if err := insertRow(transaction, id, billan, t.Note, t.Amount); err != nil {
				*ch <- false
				*response <- err
				return
			}
		}

	}
	*ch <- true
	*response <- nil
}

func (to *tool) Load(s *Search) (*[]Transaction, error) {

	var tran Transaction
	var trans []Transaction
	rows, err := PgSql.RunQuery(to.db, fmt.Sprintf("SELECT  \"Account\", \"Amount\" FROM public.\"Transactions\" where \"Account\"='%s' and date between '%s' and '%s'", s.Account, s.FromDate, s.ToDate))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rows.Scan(&tran.Account, &tran.Amount)
		trans = append(trans, tran)
	}

	return &trans, nil
}

func insertRow(t *sql.Tx, id uuid.UUID, acc, note string, amount float64) error {
	date := time.Now().Format("2006-01-02")
	time := time.Now().Format(time.RFC1123Z)
	var balance float64
	var typ int
	rows := t.QueryRow(fmt.Sprintf("SELECT \"Balance\",type FROM public.\"Accounts\" where \"AccNo\"='%s'", acc))

	if err := rows.Scan(&balance, &typ); err != nil {
		return fmt.Errorf("account not found")
	}

	if balance+amount < 0 && typ == 0 {
		return fmt.Errorf("no money")
	}
	result, err := t.Exec(fmt.Sprintf("INSERT INTO public.\"Transactions\" (\"ID\", \"TransactionId\", \"Account\", \"Amount\", \"Note\", date,time) VALUES ('%s', '%s', '%s', %2f, '%s', '%s','%s')", uuid.NewUUID(), id, acc, amount, note, date, time))
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count != 1 {
		return err
	}

	result, err = t.Exec(fmt.Sprintf("UPDATE public.\"Accounts\" SET  \"Balance\"=\"Balance\"+'%2f' WHERE \"AccNo\"='%s' ", amount, acc))
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count != 1 {
		return err
	}
	return nil
}
