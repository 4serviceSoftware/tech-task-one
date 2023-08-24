package handlers

import (
	"errors"
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
	w.Header().Set("Content-Type", "application/json")
	cache := nodes.NewCache(n.repo)
	if cache.Get(w) != nil {
		service := nodes.NewService(n.repo)
		err := service.WriteJsonNodesTree(w, 0)
		if err != nil {
			n.handleBadRequestError(w, err)
			return
		}
	}
}

func (n *Nodes) Post(w http.ResponseWriter, r *http.Request) {
	_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		n.handleBadRequestError(w, err)
		return
	}
	boundary, ok := params["boundary"]
	if !ok {
		// http.Error(w, "Not accepted Content-Type. Must be multipart/form-data with boundary.", http.StatusBadRequest)
		n.handleBadRequestError(w, errors.New("Not accepted Content-Type. Must be multipart/form-data with boundary."))
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
		n.handleBadRequestError(w, err)
		return
	}

	// Loop through each part of the request body
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			n.handleBadRequestError(w, err)
			return
		}

		// if part is a file than handle this part
		if len(part.FileName()) > 0 {
			err = service.SaveFromCarrier(part)
			if err != nil {
				n.handleBadRequestError(w, err)
				return
			}
		}
	}

	err = service.FinishSaving(n.logger)
	if err != nil {
		n.handleBadRequestError(w, err)
		return
	}
}

func (n *Nodes) handleBadRequestError(w http.ResponseWriter, err error) {
	n.logger.Println(err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
