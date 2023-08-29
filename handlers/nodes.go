package handlers

import (
	"errors"
	"log"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/4serviceSoftware/tech-task/internal/config"
	"github.com/4serviceSoftware/tech-task/internal/nodes"
)

type Nodes struct {
	service *nodes.Service
	logger  *log.Logger
	config  *config.Config
}

func NewNodes(service *nodes.Service, logger *log.Logger, config *config.Config) *Nodes {
	return &Nodes{
		service: service,
		logger:  logger,
		config:  config,
	}
}

func (n *Nodes) Get(w http.ResponseWriter, r *http.Request) {
	err := n.service.WriteCachedJsonNodesTree(w)
	if err != nil {
		n.handleBadRequestError(w, err)
		return
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
		n.handleBadRequestError(w, errors.New("Not accepted Content-Type. Must be multipart/form-data with boundary."))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(n.config.ServerMaxUploadSize)<<20+1024)
	// Create a new MultipartReader
	multipartReader := multipart.NewReader(r.Body, boundary)

	err = n.service.SaveFromMultipartReader(multipartReader)
	if err != nil {
		n.handleBadRequestError(w, err)
		return
	}
}

func (n *Nodes) handleBadRequestError(w http.ResponseWriter, err error) {
	n.logger.Println(err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
