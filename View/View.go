package View

import (
	"GitHub.com/mhthrh/GL/Controller"
	"GitHub.com/mhthrh/GL/Util/ConfigUtil"
	"GitHub.com/mhthrh/GL/Util/DbUtil/DbPool"
	"GitHub.com/mhthrh/GL/Util/ValidationUtil"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RunApiOnRouter(sm *mux.Router, log *logrus.Entry, db *DbPool.DBs, config *ConfigUtil.Config) {
	ph := Controller.New(log, ValidationUtil.NewValidation(), db, config)
	sm.Use(ph.Middleware)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/login", ph.SignIn)
	postR.HandleFunc("/transaction", ph.Transaction)
	postR.HandleFunc("/transactions", ph.Transactions)

	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/page", nil)

}
