package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type APIService struct {
	listenAddr string
	database   Database
}

func NewAPIService(listenAddr string, storage Database) *APIService {
	return &APIService{listenAddr: listenAddr, database: storage}
}

func (s *APIService) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/accounts/transfer", makeHTTPHandlerFunc(s.handlerTransfer))
	router.HandleFunc("/accounts", makeHTTPHandlerFunc(s.handlerAccount))
	router.HandleFunc("/accounts/{id}", makeHTTPHandlerFunc(s.handlerAccountById))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIService) handlerAccountById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	if r.Method == http.MethodGet && id != "" {
		return s.handlerFindAccountById(w, r)
	} else if r.Method == http.MethodPatch && id != "" {
		return s.handlerUpdateAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIService) handlerAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handlerFindAccounts(w, r)
	} else if r.Method == http.MethodPost {
		return s.handlerCreateAccount(w, r)
	} else if r.Method == http.MethodDelete {
		return s.handlerDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIService) handlerFindAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := toId(r)
	if err != nil {
		return err
	}
	account, errQuery := s.database.FindAccountById(id)
	if errQuery != nil {
		return errQuery
	}
	return WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "success", "data": account})
}

func (s *APIService) handlerFindAccounts(w http.ResponseWriter, _ *http.Request) error {
	accounts, err := s.database.FindAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "success", "data": accounts})
}

func (s *APIService) handlerCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	account := NewAccount(req.FirstName, req.LastName)
	if err := s.database.CreateAccount(account); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusCreated, map[string]interface{}{"message": "account created"})
}

func (s *APIService) handlerUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := toId(r)
	if err != nil {
		return err
	}

	accountExist, errFind := s.database.FindAccountById(id)
	if errFind != nil || accountExist == nil {
		return WriteJSON(w, http.StatusCreated, map[string]interface{}{"error": fmt.Sprintf("can't find account: %v", id)})
	}

	req := new(UpdateAccountRequest)
	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account := UpdateAccount(int64(id), req.FirstName, req.LastName, req.Balance)
	updatedAccount, errUpdate := s.database.UpdateAccount(account)
	if errUpdate != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "account updated", "data": updatedAccount})
}

func (s *APIService) handlerDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := toId(r)
	if err != nil {
		return err
	}
	err = s.database.DeleteAccount(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "account deleted"})
}

func (s *APIService) handlerTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		req := new(TransferRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			return err
		}
		defer r.Body.Close()

		accountExist, err := s.database.FindAccountById(req.ToAccount)
		if err != nil {
			return err
		}

		accountExist.Balance += int64(req.Amount)
		account, errUpdate := s.database.UpdateAccount(accountExist)
		if errUpdate != nil {
			return errUpdate
		}

		return WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "transfer succeed", "data": account})
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func WriteJSON(w http.ResponseWriter, status int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func toId(r *http.Request) (int64, error) {
	id := mux.Vars(r)["id"]
	intId, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return int64(intId), nil
}
