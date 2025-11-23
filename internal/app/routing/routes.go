package routing

import (
	"github.com/gorilla/mux"
)

// NewRouter создает роутер со всеми маршрутами
func NewRouter(server *Server) *mux.Router {
	router := mux.NewRouter()

	// Teams
	router.HandleFunc("/team/add", server.AddTeamHandler).Methods("POST")
	router.HandleFunc("/team/get", server.GetTeamHandler).Methods("GET")

	// Users
	router.HandleFunc("/users/setIsActive", server.SetIsActiveHandler).Methods("POST")
	router.HandleFunc("/users/getReview", server.GetReviewHandler).Methods("GET")

	// PullRequests
	router.HandleFunc("/pullRequest/create", server.CreatePullRequestHandler).Methods("POST")
	router.HandleFunc("/pullRequest/merge", server.MergePullRequestHandler).Methods("POST")
	router.HandleFunc("/pullRequest/reassign", server.ReassignPullRequestHandler).Methods("POST")

	return router
}
