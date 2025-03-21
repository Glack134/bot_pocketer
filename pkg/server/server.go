package server

import (
	"net/http"
	"tg_bot/pkg/storage"

	"go.uber.org/zap"
)

type AuthServer struct {
	server *http.Server
	logger *zap.Logger

	storage storage.TokenStorage
	client  *pocket.Client

	redirectUrl string
}

func NewAuthServer(redirectUrl string, storage storage.TokenStorage, client *pocket.Client) *AuthServer {
	return &AuthServer{
		redirectUrl: redirectUrl,
		storage:     storage,
		client:      client,
	}
}

func (s *AuthServer) Start() error {
	s.server = &http.Server{
		Handler: s,
		Addr:    ":80",
	}

	logger, _ := zap.NewDevelopment(zap.Fields(
		zap.String("app", "authorization server")))
	defer logger.Sync()

	s.logger = logger

	return s.server.ListenAndServe()
}

func (s *AuthServer) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.logger.Debug("received unavailable HTTP method request",
			zap.String("method", r.Method))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//next
}
