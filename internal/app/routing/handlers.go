package routing

import (
	"avito-tech/internal/app/core"
	"avito-tech/internal/apperrors"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (s *Server) writeError(w http.ResponseWriter, err error) {
	var statusCode int
	var errorCode string

	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		statusCode = http.StatusNotFound
		errorCode = "NOT_FOUND"
	case errors.Is(err, apperrors.ErrTeamExists):
		statusCode = http.StatusBadRequest
		errorCode = "TEAM_EXISTS"
	case errors.Is(err, apperrors.ErrPRExists):
		statusCode = http.StatusConflict
		errorCode = "PR_EXISTS"
	case errors.Is(err, apperrors.ErrPRMerged):
		statusCode = http.StatusConflict
		errorCode = "PR_MERGED"
	case errors.Is(err, apperrors.ErrNotAssigned):
		statusCode = http.StatusConflict
		errorCode = "NOT_ASSIGNED"
	case errors.Is(err, apperrors.ErrNoCandidate):
		statusCode = http.StatusConflict
		errorCode = "NO_CANDIDATE"
	default:
		statusCode = http.StatusInternalServerError
		errorCode = "INTERNAL_ERROR"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    errorCode,
			Message: err.Error(),
		},
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) AddTeamHandler(w http.ResponseWriter, r *http.Request) {
	var req core.AddTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	resp, err := s.impl.AddTeam(r.Context(), &req)
	if err != nil {
		s.writeError(w, err)
		return
	}

	s.writeJSON(w, http.StatusCreated, resp)
}

func (s *Server) GetTeamHandler(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		s.writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "BAD_REQUEST",
				Message: "team_name parameter is required",
			},
		})
		return
	}

	resp, err := s.impl.GetTeamByTeamName(r.Context(), teamName)
	if err != nil {
		s.writeError(w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, resp)
}

func (s *Server) SetIsActiveHandler(w http.ResponseWriter, r *http.Request) {
	var req core.SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	resp, err := s.impl.UserSetIsActive(r.Context(), req)
	if err != nil {
		s.writeError(w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, resp)
}

func (s *Server) CreatePullRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req core.CreatePullReqRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	resp, err := s.impl.CreatePullRequestFromCreateRequest(r.Context(), &req)
	if err != nil {
		s.writeError(w, err)
		return
	}

	s.writeJSON(w, http.StatusCreated, resp)
}

func (s *Server) MergePullRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req core.MergePullReqRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	resp, err := s.impl.MergePullRequest(r.Context(), req)
	if err != nil {
		s.writeError(w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, resp)
}

func (s *Server) ReassignPullRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req core.ReassignPullReqRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	resp, err := s.impl.ReassignPullRequest(r.Context(), &req)
	if err != nil {
		s.writeError(w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, resp)
}

func (s *Server) GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		s.writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "BAD_REQUEST",
				Message: "user_id parameter is required",
			},
		})
		return
	}

	resp, err := s.impl.GetReview(r.Context(), userID)
	if err != nil {
		s.writeError(w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, resp)
}
