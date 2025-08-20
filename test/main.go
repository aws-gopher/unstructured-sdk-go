// Package test provides testing utilities and examples for the unstructured SDK.
package test

import (
	"net/http"
	"strings"
)

// FakeAPIKey is a fake API key for testing.
const FakeAPIKey = "91pmLBeETAbXCpNylRsLq11FdiZPTk"

// Mux is a HTTP server that mocks the Unstructured API.
type Mux struct {
	mux *http.ServeMux

	// Destination handlers
	CreateDestination                 func(w http.ResponseWriter, r *http.Request)
	ListDestinations                  func(w http.ResponseWriter, r *http.Request)
	GetDestination                    func(w http.ResponseWriter, r *http.Request)
	UpdateDestination                 func(w http.ResponseWriter, r *http.Request)
	DeleteDestination                 func(w http.ResponseWriter, r *http.Request)
	CreateConnectionCheckDestinations func(w http.ResponseWriter, r *http.Request)
	GetConnectionCheckDestinations    func(w http.ResponseWriter, r *http.Request)

	// Source handlers
	ListSources                  func(w http.ResponseWriter, r *http.Request)
	CreateSource                 func(w http.ResponseWriter, r *http.Request)
	GetSource                    func(w http.ResponseWriter, r *http.Request)
	UpdateSource                 func(w http.ResponseWriter, r *http.Request)
	DeleteSource                 func(w http.ResponseWriter, r *http.Request)
	CreateConnectionCheckSources func(w http.ResponseWriter, r *http.Request)
	GetConnectionCheckSources    func(w http.ResponseWriter, r *http.Request)

	// Job handlers
	ListJobs          func(w http.ResponseWriter, r *http.Request)
	GetJob            func(w http.ResponseWriter, r *http.Request)
	CancelJob         func(w http.ResponseWriter, r *http.Request)
	DownloadJobOutput func(w http.ResponseWriter, r *http.Request)
	GetJobDetails     func(w http.ResponseWriter, r *http.Request)
	GetJobFailedFiles func(w http.ResponseWriter, r *http.Request)

	// Workflow handlers
	CreateWorkflow func(w http.ResponseWriter, r *http.Request)
	ListWorkflows  func(w http.ResponseWriter, r *http.Request)
	GetWorkflow    func(w http.ResponseWriter, r *http.Request)
	UpdateWorkflow func(w http.ResponseWriter, r *http.Request)
	DeleteWorkflow func(w http.ResponseWriter, r *http.Request)
	RunWorkflow    func(w http.ResponseWriter, r *http.Request)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1")
	r.URL.RawPath = r.URL.Path
	m.mux.ServeHTTP(w, r)
}

// NewMux creates a new Mux with all the routes for the API.
func NewMux() *Mux {
	m := &Mux{mux: http.NewServeMux()}

	try := func(f *func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			if f != nil && *f != nil {
				(*f)(w, r)
			} else {
				var msg string
				if r.Pattern != "" {
					msg = "handler for " + r.Pattern + " is nil"
				}

				http.Error(w, msg, http.StatusMethodNotAllowed)
			}
		}
	}

	// Destination routes
	m.mux.HandleFunc("POST /destinations/", try(&m.CreateDestination))
	m.mux.HandleFunc("GET /destinations/", try(&m.ListDestinations))
	m.mux.HandleFunc("GET /destinations/{id}", try(&m.GetDestination))
	m.mux.HandleFunc("PUT /destinations/{id}", try(&m.UpdateDestination))
	m.mux.HandleFunc("DELETE /destinations/{id}", try(&m.DeleteDestination))
	m.mux.HandleFunc("POST /destinations/{id}/connection-check", try(&m.CreateConnectionCheckDestinations))
	m.mux.HandleFunc("GET /destinations/{id}/connection-check", try(&m.GetConnectionCheckDestinations))

	// Source routes
	m.mux.HandleFunc("GET /sources/", try(&m.ListSources))
	m.mux.HandleFunc("POST /sources/", try(&m.CreateSource))
	m.mux.HandleFunc("GET /sources/{id}", try(&m.GetSource))
	m.mux.HandleFunc("PUT /sources/{id}", try(&m.UpdateSource))
	m.mux.HandleFunc("DELETE /sources/{id}", try(&m.DeleteSource))
	m.mux.HandleFunc("POST /sources/{id}/connection-check", try(&m.CreateConnectionCheckSources))
	m.mux.HandleFunc("GET /sources/{id}/connection-check", try(&m.GetConnectionCheckSources))

	// Job routes
	m.mux.HandleFunc("GET /jobs/", try(&m.ListJobs))
	m.mux.HandleFunc("GET /jobs/{id}", try(&m.GetJob))
	m.mux.HandleFunc("POST /jobs/{id}/cancel", try(&m.CancelJob))
	m.mux.HandleFunc("GET /jobs/{id}/download", try(&m.DownloadJobOutput))
	m.mux.HandleFunc("GET /jobs/{id}/details", try(&m.GetJobDetails))
	m.mux.HandleFunc("GET /jobs/{id}/failed-files", try(&m.GetJobFailedFiles))

	// Workflow routes
	m.mux.HandleFunc("POST /workflows/", try(&m.CreateWorkflow))
	m.mux.HandleFunc("GET /workflows/", try(&m.ListWorkflows))
	m.mux.HandleFunc("GET /workflows/{id}", try(&m.GetWorkflow))
	m.mux.HandleFunc("PUT /workflows/{id}", try(&m.UpdateWorkflow))
	m.mux.HandleFunc("DELETE /workflows/{id}", try(&m.DeleteWorkflow))
	m.mux.HandleFunc("POST /workflows/{id}/run", try(&m.RunWorkflow))

	return m
}
