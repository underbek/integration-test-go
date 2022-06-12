package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AndreyAndreevich/articles/user_service/api"
	"github.com/AndreyAndreevich/articles/user_service/domain"
)

type UseCase interface {
	CreateUser(ctx context.Context, name string) (domain.User, error)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
