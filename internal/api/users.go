package api

import (
	"context"
	"encoding/json"
	"go-project-template/internal/domain"
	"go-project-template/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Upsert godoc
// @Summary Create or Update User
// @Description Accepts a JSON model and on conflict of a present UUID will instead update the user fields
// @Tags  Users
// @Accept json
// @Produce json
// @Param payload body domain.User true "User"
// @Success 200
// @Router /users [post]
func (a *api) upsertUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	u := &domain.User{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	u, err := a.userRepo.CreateOrUpdate(ctx, u)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(u)
}

// Delete godoc
// @Summary Delete User
// @Description Accepts an ID and processes the delete. Returns a 200 header if no errors occured during the delete.
// @Tags  Users
// @Param userid path string true "userid"
// @Success 200
// @Router /users/{userid} [delete]
func (a *api) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	parsedUserId, err := uuid.Parse(chi.URLParam(r, "userid"))
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	err = a.userRepo.Delete(ctx, parsedUserId)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Get godoc
// @Summary Get User
// @Description Accepts an ID and returns a JSON model
// @Tags  Users
// @Produce json
// @Param userid path string true "userid"
// @Success 200 {object} domain.User
// @Router /users/{userid} [get]
func (a *api) getByIdUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	parsedUserId, err := uuid.Parse(chi.URLParam(r, "userid"))
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	user, err := a.userRepo.GetByID(ctx, parsedUserId)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

// Get List godoc
// @Summary Get List of Users
// @Description Accepts pagination based query parameters and returns a paginated response.
// @Tags  Users
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query string false "filter name" Format(orderBy)
// @Param orderDir query string false "filter name" Format(orderDir)
// @Success 200 {object} []domain.User
// @Failure 400
// @Router /users [get]
func (a *api) getUserListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	pagQuery, err := utils.GetPaginationFromRequest(r)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	user, err := a.userRepo.GetList(ctx, pagQuery)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}
