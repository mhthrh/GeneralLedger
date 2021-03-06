package User

import (
	"GitHub.com/mhthrh/GL/Model/Account"
	"GitHub.com/mhthrh/GL/Util/CryptoUtil"
	"GitHub.com/mhthrh/GL/Util/DbUtil/PgSql"
	"database/sql"
	"errors"
	"fmt"
	"github.com/pborman/uuid"
	"strings"
	"time"
)

type User struct {
	Id       uuid.UUID
	UserName string
	Password string
	Email    string
}
type Request struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}
type Response struct {
	Username  string            `json:"Username,omitempty"`
	IsActive  bool              `json:"IsActive,omitempty"`
	ValidTill string            `json:"ValidTill"`
	Sign      string            `json:"Sign,omitempty"`
	Accounts  []Account.Account `json:"Accounts,omitempty"`
	Err       string            `json:"Err,omitempty"`
}
type tool struct {
	db *sql.DB
}

var (
	validationDuration time.Duration = 180 * time.Second
	dateFormat                       = time.UnixDate
)

func New(db *sql.DB) *tool {
	return &tool{db: db}
}

func (t *tool) Login(r *Request) *Response {
	var user User
	res := Response{
		Username:  r.Username,
		IsActive:  false,
		ValidTill: time.Time{}.Format(time.UnixDate),
		Sign:      "",
		Accounts: []Account.Account{
			{
				ID:      nil,
				UserId:  nil,
				AccNo:   "",
				Balance: 0,
			},
		},
		Err: "",
	}
	SignedPassword := CryptoUtil.NewKey()
	SignedPassword.Text = r.Password
	SignedPassword.Sha256()
	rows, err := PgSql.RunQuery(t.db, fmt.Sprintf("SELECT \"ID\", \"UserName\", \"Email\" FROM public.\"Users\" where \"UserName\"='%s' and \"Password\"='%s'", r.Username, SignedPassword.Result))
	if err != nil {
		res.Err = err.Error()
		return &res
	}

	if rows.Next() {
		rows.Scan(&user.Id, &user.UserName, &user.Email)
	} else {
		res.Err = errors.New("user or pass invalid").Error()
		return &res
	}
	accounts, _ := Account.New(t.db).Load(user.Id)

	res.IsActive = true
	res.ValidTill = time.Now().Add(validationDuration).Format(dateFormat)
	res.Accounts = accounts
	res.Sign = t.GenerateSignKey(r.Username, validationDuration)

	return &res
}

func (t *tool) GenerateSignKey(userName string, validationDuration time.Duration) string {
	j := CryptoUtil.NewKey()
	j.Text = fmt.Sprintf("%s#%s", userName, time.Now().Add(validationDuration).Format(dateFormat))
	j.Encrypt()
	return j.Result
}

func (t *tool) CheckSignKey(signedKey string) ([]string, error) {

	k := CryptoUtil.NewKey()
	k.Text = signedKey
	err := k.Decrypt()
	if err != nil {
		return nil, err
	}
	spl := strings.Split(k.Result, "#")
	if len(spl) > 1 {
		return spl, nil
	}
	return nil, fmt.Errorf("error in sign key")

}
