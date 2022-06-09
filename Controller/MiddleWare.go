package Controller

import (
	"GitHub.com/mhthrh/GL/Model/Result"
	"GitHub.com/mhthrh/GL/Model/Transaction"
	"GitHub.com/mhthrh/GL/Model/User"
	"GitHub.com/mhthrh/GL/Util/ConfigUtil"
	"GitHub.com/mhthrh/GL/Util/DbUtil/DbPool"
	"GitHub.com/mhthrh/GL/Util/JsonUtil"
	"GitHub.com/mhthrh/GL/Util/ValidationUtil"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type Key struct{}

type Controller struct {
	l    *logrus.Entry
	v    *ValidationUtil.Validation
	db   *DbPool.DBs
	Conf *ConfigUtil.Config
}

var (
	InvalidPath = fmt.Errorf("invalid Path, path must be{s}")
	timeOut     = 5000 * time.Millisecond
)

type GenericError struct {
	Message string `json:"message"`
}
type GenericError1 struct {
	Message error `json:"message"`
}
type ValidationError struct {
	Messages []string `json:"messages"`
}

func New(l *logrus.Entry, v *ValidationUtil.Validation, db *DbPool.DBs, c *ConfigUtil.Config) *Controller {
	return &Controller{l, v, db, c}
}
func (b *Controller) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fNext := func(in interface{}) {
			//setting time out for methods
			cnt, _ := context.WithTimeout(context.WithValue(r.Context(), Key{}, in), timeOut)
			r = r.WithContext(cnt)
			next.ServeHTTP(w, r)
		}
		if r.Host != fmt.Sprintf("%s:%d", b.Conf.Server.IP, b.Conf.Server.Port) {
			err := errors.New("access denied")
			Result.New(1002, http.StatusForbidden, GenericError{Message: err.Error()}.Message).SendResponse(w)
			return
		}

		switch strings.ToLower(r.URL.Path) {
		case "/login":
			{
				b.l.Println("incoming login request.")
				obj := User.Request{}
				err := JsonUtil.New(nil, r.Body).FromJSON(&obj)
				if err != nil {
					b.l.Println(err)
					Result.New(1003, http.StatusBadRequest, GenericError{Message: err.Error()}.Message).SendResponse(w)
					return
				}
				errs := b.v.Validate(obj)
				if len(errs) != 0 {
					b.l.Println(errs)
					j := JsonUtil.New(nil, nil).Struct2Json(ValidationError{Messages: errs.Errors()}.Messages)
					Result.New(1004, http.StatusUnprocessableEntity, j).SendResponse(w)
					return
				}
				fNext(obj)
			}
		case "/transaction":
			{
				obj := Transaction.Request{}
				err := JsonUtil.New(nil, r.Body).FromJSON(&obj)
				if err != nil {
					b.l.Println(err)
					Result.New(1005, http.StatusBadRequest, GenericError{Message: err.Error()}.Message).SendResponse(w)
					return
				}
				sign, err := User.New(nil).CheckSignKey(obj.Sign)
				if err != nil {
					b.l.Println(err)
					Result.New(1006, http.StatusForbidden, GenericError{Message: err.Error()}.Message).SendResponse(w)
					return
				}
				signedTime, err := time.Parse(time.UnixDate, sign[1])
				if err != nil {
					b.l.Println(err)
					Result.New(1007, http.StatusForbidden, GenericError{Message: err.Error()}.Message).SendResponse(w)
					return
				}
				now, err := time.Parse(time.UnixDate, time.Now().Format(time.UnixDate))
				if err != nil {
					b.l.Println(err)
					Result.New(1008, http.StatusForbidden, GenericError{Message: err.Error()}.Message).SendResponse(w)
					return
				}
				if !signedTime.After(now) {
					b.l.Println(err)
					Result.New(1009, http.StatusForbidden, "check sign key.").SendResponse(w)
					return
				}

				errs := b.v.Validate(obj)
				if len(errs) != 0 {
					b.l.Println(err)
					j := JsonUtil.New(nil, nil).Struct2Json(ValidationError{Messages: errs.Errors()}.Messages)
					Result.New(1010, http.StatusUnprocessableEntity, j).SendResponse(w)
					return
				}
				fNext(obj)
			}
		case "/transactions":
			{
				obj := Transaction.Search{}
				err := JsonUtil.New(nil, r.Body).FromJSON(&obj)
				if err != nil {
					b.l.Println(err)
					Result.New(1011, http.StatusBadRequest, GenericError{Message: err.Error()}.Message).SendResponse(w)
					return
				}
				sign, err := User.New(nil).CheckSignKey(obj.Sign)
				if err != nil {
					b.l.Println(err)
					Result.New(1012, http.StatusForbidden, GenericError{Message: err.Error()}.Message).SendResponse(w)
					return
				}
				signedTime, err := time.Parse(time.UnixDate, sign[1])
				if err != nil {
					b.l.Println(err)
					Result.New(1013, http.StatusForbidden, GenericError{Message: err.Error()}.Message).SendResponse(w)
					return
				}
				now, err := time.Parse(time.UnixDate, time.Now().Format(time.UnixDate))
				if err != nil {
					b.l.Println(err)
					Result.New(1014, http.StatusForbidden, GenericError{Message: err.Error()}.Message).SendResponse(w)
					return
				}
				if !signedTime.After(now) {
					b.l.Println(err)
					Result.New(1015, http.StatusForbidden, "check sign key.").SendResponse(w)
					return
				}

				errs := b.v.Validate(obj)
				if len(errs) != 0 {
					b.l.Println(err)
					j := JsonUtil.New(nil, nil).Struct2Json(ValidationError{Messages: errs.Errors()}.Messages)
					Result.New(1016, http.StatusUnprocessableEntity, j).SendResponse(w)
					return
				}
				fNext(obj)

			}
		case "/page":
			{
				b.l.Println("Loading UI")
				http.ServeFile(w, r, "./Page/index.html")
			}
		default:
			{
				err := errors.New("pageNotFound")
				Result.New(1017, http.StatusNotFound, GenericError{Message: err.Error()}.Message).SendResponse(w)
			}
		}

	})
}

func Check(interf interface{}, w11 http.ResponseWriter) {

}
