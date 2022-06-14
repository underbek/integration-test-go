package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AndreyAndreevich/articles/user_service/api"
	"github.com/AndreyAndreevich/articles/user_service/domain"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

type UseCase interface {
	CreateUser(ctx context.Context, name string) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	UpdateBalance(ctx context.Context, id int, amount decimal.Decimal) (domain.User, error)
}

type Handler struct {
	useCase UseCase
}

func New(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func userToApi(user domain.User) api.User {
	return api.User{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	request := api.CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.useCase.CreateUser(r.Context(), request.Name)
	if err != nil {
		fmt.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := api.CreateUserResponse(userToApi(user))

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		fmt.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := h.useCase.GetUser(r.Context(), id)
	if err != nil {
		fmt.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := api.CreateUserResponse(userToApi(user))

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		fmt.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DepositBalance(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	request := api.DepositBalanceRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.useCase.UpdateBalance(r.Context(), request.ID, request.Amount)
	if err != nil {
		fmt.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := api.DepositBalanceResponse(userToApi(user))

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		fmt.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
