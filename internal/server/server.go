package server

import (
	"net"
	"net/http"
	"strings"

	"github.com/YEgorLu/kafka-test/internal/controllers"
	"github.com/YEgorLu/kafka-test/internal/logger"
	"github.com/YEgorLu/kafka-test/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	http.Server
	router *http.ServeMux
	log    logger.Logger
	addr   string
}

type ServerConfig struct {
	Port string
	Log  logger.Logger
	Addr string
}

var defaultConfig = ServerConfig{
	Port: ":8080",
	Addr: "localhost",
}

func NewServer(conf *ServerConfig) *Server {
	conf = parseConfig(conf)
	server := &Server{
		Server: http.Server{
			Addr: conf.Port,
		},
		log:  conf.Log,
		addr: conf.Addr,
	}
	return server
}

func (s *Server) Configure() *Server {
	router, err := controllers.GetRoutes()
	if err != nil {
		s.log.Error(err)
		panic(err)
	}
	s.router = router
	return s
}

func (s *Server) WithSwagger() *Server {
	if s.router == nil {
		s.log.Error("WithSwagger must be placed after Configure()")
		panic("WithSwagger must be placed after Configure()")
	}
	s.router.HandleFunc("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://"+s.addr+s.Addr+"/swagger/doc.json")))
	return s
}

func (s *Server) Run() error {
	s.Handler = middleware.Combine(
		middleware.CORS(middleware.CORSConfig{AllowedOrigins: []string{"*"}, AllowCredentials: true}),
		middleware.Recover(s.log),
		middleware.Logger(s.log),
	)(s.router)
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		s.log.Error(err)
		return err
	}

	s.log.Info("server is listening on ", s.Addr)

	return s.Server.Serve(l)
}

func parseConfig(c *ServerConfig) *ServerConfig {
	conf := *c
	if conf.Port == "" {
		conf.Port = defaultConfig.Port
	} else if !strings.HasPrefix(conf.Port, ":") {
		conf.Port = ":" + conf.Port
	}

	if conf.Addr == "" {
		conf.Addr = defaultConfig.Addr
	}
	return &conf
}
