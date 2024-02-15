package rest

import (
	"FriendsAdvice/internal/transport"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func Put(controller transport.IController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1.Take key
		vars := mux.Vars(r)
		key := vars["key"]
		if len(key) == 0 {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusBadRequest)
			return
		}

		// 2.Take body
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// 3.Work with (in according with expires)
		if expires, hasExpires := r.Header["Expires"]; hasExpires {
			if len(expires) != 1 {
				http.Error(w, http.StatusText(http.StatusRequestHeaderFieldsTooLarge), http.StatusRequestHeaderFieldsTooLarge)
				return
			}
			exp, err := strconv.Atoi(expires[0])
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			lifetime := time.Duration(exp) * time.Millisecond
			putStatus := controller.PutObjectWithExpires(key, string(reqData), lifetime)
			if putStatus {
				w.WriteHeader(http.StatusOK)
			} else {
				http.Error(w, http.StatusText(http.StatusFound), http.StatusFound)
			}
		} else {
			putStatus := controller.PutObject(key, string(reqData))
			if putStatus {
				w.WriteHeader(http.StatusOK)
			} else {
				http.Error(w, http.StatusText(http.StatusFound), http.StatusFound)
			}
		}
	}
}
