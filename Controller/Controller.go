package Controller

import (
	"GitHub.com/mhthrh/GL/Model/Result"
	"GitHub.com/mhthrh/GL/Model/Transaction"
	"GitHub.com/mhthrh/GL/Model/User"
	"GitHub.com/mhthrh/GL/Util/JsonUtil"
	"fmt"
	"net/http"
)

func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	var i User.Request
	i = r.Context().Value(Key{}).(User.Request)
	d := c.db.Pull()
	response := User.New(d.Db).Login(&User.Request{
		Username: i.Username,
		Password: i.Password,
	})
	c.db.Push(d)
	Result.New(1, http.StatusOK, JsonUtil.New(nil, nil).Struct2Json(response)).SendResponse(w)
}
func (c *Controller) Transaction(w http.ResponseWriter, r *http.Request) {
	var i Transaction.Request
	cc := make(chan bool, 1)
	req := make(chan error, 1)
	i = r.Context().Value(Key{}).(Transaction.Request)
	go func(ch *chan bool) {
		d := c.db.Pull()
		Transaction.New(d.Db).Create(&Transaction.Transaction{
			Account: i.Transaction.Account,
			Amount:  i.Transaction.Amount,
			Note:    i.Transaction.Note,
			Action:  i.Transaction.Action,
		}, &cc, &req)
		c.db.Push(d)

	}(&cc)
	for {
		select {
		case <-r.Context().Done():
			cc <- false
			fmt.Println("time out")
			Result.New(1, http.StatusOK, fmt.Sprintf("%s%s", "kir shodi, time out", r.Context().Err())).SendResponse(w)
			return
		case response := <-req:
			if response != nil {
				Result.New(1, http.StatusOK, JsonUtil.New(nil, nil).Struct2Json(response.Error())).SendResponse(w)
				return
			}
			Result.New(1, http.StatusOK, "transaction commit").SendResponse(w)
			return
		}
	}
}

func (c *Controller) Transactions(w http.ResponseWriter, r *http.Request) {
	var i Transaction.Search
	i = r.Context().Value(Key{}).(Transaction.Search)
	d := c.db.Pull()
	response, err := Transaction.New(d.Db).Load(&Transaction.Search{
		Account:  i.Account,
		FromDate: i.FromDate,
		ToDate:   i.ToDate,
	})
	c.db.Push(d)
	if err != nil {
		Result.New(1, http.StatusOK, JsonUtil.New(nil, nil).Struct2Json(err.Error())).SendResponse(w)
	}
	Result.New(1, http.StatusOK, JsonUtil.New(nil, nil).Struct2Json(response)).SendResponse(w)
}
