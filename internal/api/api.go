package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/madalinpopa/memoflash/internal/client"
)

type ApiServer struct {
	Config *Config
	Client *client.MemosClient
}

func NewApiServer(c *Config) *ApiServer {
	return &ApiServer{
		Config: c,
		Client: client.NewMemosClient(c.UseMemosHost, c.UseMemosToken),
	}
}

func (s *ApiServer) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.RootHandler)
	mux.HandleFunc("GET /memos", s.GetMemoHandler)

	handler := Chain(mux,
		LoggerMiddleWare,
		JSONContentTypeMiddleware,
	)

	server := &http.Server{
		Addr:    s.Config.ListenAddr,
		Handler: handler,
	}

	log.Printf("Listening on %s", s.Config.ListenAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server error: %v", err)
	}
}

func (s *ApiServer) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	response := map[string]string{"message": "Welcome to MemoFlash API"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *ApiServer) GetMemoHandler(w http.ResponseWriter, r *http.Request) {
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	pageToken := r.URL.Query().Get("pageToken")
	filter := r.URL.Query().Get("filter")
	tag := r.URL.Query().Get("tag")

	log.Printf("Requesting memos with pageSize: %d pageToken: %s filter: %s tag: %s", pageSize, pageToken, filter, tag)

	response, err := s.Client.GetMemos(pageSize, pageToken, filter, tag)
	if err != nil {
		log.Printf("Error getting memos: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
