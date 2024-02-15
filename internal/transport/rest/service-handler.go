package rest

import (
	"FriendsAdvice/internal/transport"
	"fmt"
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
		keyStr := vars["key"]
		if len(keyStr) == 0 {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusBadRequest)
			return
		}

		key := -1
		var err error = nil
		if key, err = strconv.Atoi(keyStr); err != nil {
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

func Get(controller transport.IController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1.Take key
		vars := mux.Vars(r)
		keyStr := vars["key"]
		if len(keyStr) == 0 {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusBadRequest)
			return
		}

		key := -1
		var err error = nil
		if key, err = strconv.Atoi(keyStr); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusBadRequest)
			return
		}

		// 2.Get object due to controller
		obj, haveObject := controller.GetObject(key)
		if haveObject {
			w.Write([]byte(obj))
			fmt.Printf("Object is %s, key is %d", obj, key)
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}
