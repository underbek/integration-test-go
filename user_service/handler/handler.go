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
)

type UseCase interface {
	CreateUser(ctx context.Context, name string) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
}

type Handler struct {
	useCase UseCase
}

func New(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
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

	response := api.CreateUserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}

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

	response := api.CreateUserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		fmt.Println("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
