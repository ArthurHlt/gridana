package gridana

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type JsonError struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

func panicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errMsg := JsonError{
					Status:  http.StatusInternalServerError,
					Title:   http.StatusText(http.StatusInternalServerError),
					Details: fmt.Sprint(err),
				}
				if httpError, ok := err.(HttpError); ok {
					log.WithField("verb", httpError.Verb).Error(httpError.Error())
					errMsg.Title = http.StatusText(httpError.Status)
					errMsg.Status = httpError.Status
					errMsg.Details = httpError.Details
				} else {
					log.Error(fmt.Sprint(err))
				}
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(errMsg.Status)
				body, err := json.Marshal(errMsg)
				if err != nil {
					panic(err)
				}
				if _, err := w.Write(body); err != nil {
					panic(err)
				}
			}
		}()
		next.ServeHTTP(w, req)
	})
}

type HttpError struct {
	Status  int
	Details string
	Verb    string
}

func (e HttpError) Error() string {
	return fmt.Sprintf("%s (%d): %s", http.StatusText(e.Status), e.Status, e.Details)
}
