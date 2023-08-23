package handlers

import (
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/4serviceSoftware/tech-task/nodes"
)

type Nodes struct {
	repo   nodes.Repository
	logger *log.Logger
}

func NewNodes(r nodes.Repository, l *log.Logger) *Nodes {
	return &Nodes{
		repo:   r,
		logger: l,
	}
}

func (n *Nodes) Get(w http.ResponseWriter, r *http.Request) {
	n.logger.Println("Getting nodes")
}

func (n *Nodes) Post(w http.ResponseWriter, r *http.Request) {
	n.logger.Println("Posting nodes")

	_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	boundary, ok := params["boundary"]
	if !ok {
		http.Error(w, "Not accepted Content-Type. Must be multipart/form-data with boundary.", http.StatusBadRequest)
		return
	}

	// TODO: get max bytes limit from config
	r.Body = http.MaxBytesReader(w, r.Body, 128<<20+1024)
	// Create a new MultipartReader
	mr := multipart.NewReader(r.Body, boundary)

	service := nodes.NewService(n.repo)

	// initializing new saving session
	err = service.StartSaving()
	defer service.RollbackSaving()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Loop through each part of the request body
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// if part is a file than handle this part
		if len(part.FileName()) > 0 {
			err = service.SaveFromCarrier(part)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	service.FinishSaving()
}
