package middlewares

import (
	"bytes"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// CacheMiddleware caches GET requests and serves cached responses
func CacheMiddleware(redisClient *redis.Client, expiration time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != http.MethodGet {
				return next(c)
			}

			// Start a new span
			tracer := otel.Tracer("middlewares/CacheMiddleware")
			ctx, span := tracer.Start(c.Request().Context(), "CacheMiddleware", trace.WithAttributes(
				semconv.HTTPMethodKey.String(c.Request().Method),
				semconv.HTTPURLKey.String(c.Request().URL.Path),
			))
			defer span.End()

			cacheKey := "cache:" + c.Request().URL.Path + c.Request().URL.RawQuery
			cachedResponse, err := redisClient.Get(ctx, cacheKey).Result()

			if err == nil && cachedResponse != "" {
				c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c.Response().WriteHeader(http.StatusOK)
				_, err := c.Response().Write([]byte(cachedResponse))
				if err != nil {
					return err
				}
				return nil
			}

			recorder := newResponseRecorder(c.Response().Writer)
			c.Response().Writer = recorder

			err = next(c)
			if err != nil {
				return err
			}

			if recorder.status == http.StatusOK {
				_, err := redisClient.Set(c.Request().Context(), cacheKey, recorder.body.Bytes(), expiration).Result()
				if err != nil {
					return err
				}
			}

			return nil
		}
	}
}

// InvalidateCacheMiddleware invalidates cache after successful POST requests
func InvalidateCacheMiddleware(redisClient *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == http.MethodPost {
				defer func() {
					if c.Response().Status == http.StatusOK || c.Response().Status == http.StatusCreated {
						cacheKey := "cache:" + c.Request().URL.Path
						keys, _ := redisClient.Keys(c.Request().Context(), cacheKey+"*").Result()
						for _, key := range keys {
							redisClient.Del(c.Request().Context(), key)
						}
					}
				}()
			}
			return next(c)
		}
	}
}

// responseRecorder is used to capture the response for caching
type responseRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{ResponseWriter: w, body: new(bytes.Buffer)}
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
