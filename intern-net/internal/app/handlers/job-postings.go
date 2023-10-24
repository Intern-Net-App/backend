package handlers

import (
	"encoding/json"
	"intern-net/internal/app/repositories"
	"net/http"
	"strconv"
)

type JobPostingsHandler struct {
	JobRepository *repositories.JobRepository
}

func (h *JobPostingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetJobPostings(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func NewJobPostingsHandler(jobRepo *repositories.JobRepository) *JobPostingsHandler {
	return &JobPostingsHandler{
		JobRepository: jobRepo,
	}
}

func (h *JobPostingsHandler) GetJobPostings(w http.ResponseWriter, r *http.Request) {
	// Parse query params to determine how many job postings to skip and fetch
	skip, limit := parseSkipAndLimitQueryParams(r)

	// Retrieve job postings from the database
	jobPostings, err := h.JobRepository.GetJobs(r.Context(), skip, limit)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond with the job postings as JSON
	w.Header().Set("Content-Type", "applications/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(jobPostings)
}

func parseSkipAndLimitQueryParams(r *http.Request) (int, int) {
	skip := 0
	limit := 10 // Default limit to 10, we can adjust later

	// Parse param for skip and limit (example: ?skip=5&limit=10)
	query := r.URL.Query()
	if skipParam := query.Get("skip"); skipParam != "" {
		skip, _ = strconv.Atoi(skipParam)
	}
	if limitParam := query.Get("limit"); limitParam != "" {
		limit, _ = strconv.Atoi(limitParam)
	}

	return skip, limit
}
