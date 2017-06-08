package ranger_http

import (
	"log"
	"net/http"
	"time"

	"github.com/fesposito/go-ranger/ranger_logger"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
)

// Server ...
type Server struct {
	*httprouter.Router
	ResponseWriter
	middlewares []alice.Constructor
}

const (
	defaultResponseCacheTimeout = time.Duration(5) * time.Minute
)

var (
	logger ranger_logger.LoggerInterface
)

// NewHTTPServer ...
func NewHTTPServer(l ranger_logger.LoggerInterface) *Server {
	logger = l
	responseWriter := ResponseWriter{}
	router := httprouter.New()

	return &Server{
		Router:         router,
		ResponseWriter: responseWriter,
	}
}

func (s *Server) WithDefaultErrorRoute() {
	s.PanicHandler = PanicHandler(s.ResponseWriter)
	s.NotFound = NotFoundHandler(s.ResponseWriter)
}

func (s *Server) WithHealthCheckFor(services ...interface{}) {
	s.GET("/health/check", HealthCheckHandler(services))
	s.GET("/health/check/lb", HealthCheckHandlerLB())
}

func (s *Server) WithMiddleware(middlewares ...func(http.Handler) http.Handler) {
	for _, v := range middlewares {
		s.middlewares = append(s.middlewares, v)
	}
}

func (s *Server) WithThrottle(handler *http.HandlerFunc) http.Handler {
	// @todo learn more about this memstore
	store, err := memstore.New(65536)
	if err != nil {
		logger.Error(err.Error(), nil)
	}

	quota := throttled.RateQuota{throttled.PerMin(20), 5}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		logger.Error(err.Error(), nil)
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true},
	}

	return httpRateLimiter.RateLimit(handler)
}

func (s *Server) Start(addr string) {
	chain := alice.New(s.middlewares...)
	logger.Info("Listening to address:", map[string]interface{}{"addr": addr})
	log.Fatal(http.ListenAndServe(addr, chain.Then(s.Router)))
}

// @todo add cache headers to response
