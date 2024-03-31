package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/utsavgupta/go-todo/uc"
)

type RouterDependencies struct {
	GetUC    uc.GetTaskUC
	ListUC   uc.ListTasksUC
	CreateUC uc.CreateTaskUC
	UpdateUC uc.UpdateTaskUC
	DeleteUC uc.DeleteTaskUC
}

func NewRouter(dependencies *RouterDependencies) http.Handler {

	router := mux.NewRouter()

	fmt.Printf("%#v\n", dependencies)

	router.NewRoute().Path("/").HandlerFunc(newListHandler(dependencies.ListUC)).Methods(http.MethodGet)
	router.NewRoute().Path("/{id}").HandlerFunc(newGetHandler(dependencies.GetUC)).Methods(http.MethodGet)
	router.NewRoute().Path("/{id}").HandlerFunc(newDeleteHandler(dependencies.DeleteUC)).Methods(http.MethodDelete)
	router.NewRoute().Path("/").HandlerFunc(newCreateHandler(dependencies.CreateUC)).Methods(http.MethodPost)

	return router
}

func newGetHandler(getUC uc.GetTaskUC) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		idParam, _ := vars["id"]
		id, _ := strconv.Atoi(idParam)

		details, err := getUC(r.Context(), id)

		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if details == nil {

			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(details)
	}
}

func newDeleteHandler(deleteUC uc.DeleteTaskUC) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		idParam, _ := vars["id"]
		id, _ := strconv.Atoi(idParam)

		err := deleteUC(r.Context(), id)

		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func newListHandler(listUC uc.ListTasksUC) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tasks, err := listUC(r.Context())

		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
	}
}

func newCreateHandler(listUC uc.CreateTaskUC) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var b []byte

		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {

			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}
