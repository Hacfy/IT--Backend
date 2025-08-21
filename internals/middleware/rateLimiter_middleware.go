package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// ("5-M", 1*time.Minute)
// ("10-S", 2*time.Minute)

var blocked = struct {
	m map[string]time.Time
	sync.Mutex
}{m: make(map[string]time.Time)}

func CustomRateLimiter(rateString string, Time time.Duration) echo.MiddlewareFunc {
	rate, _ := limiter.NewRateFromFormatted(rateString)
	store := memory.NewStore()
	limiterInstance := limiter.New(store, rate)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			key := e.RealIP()

			blocked.Lock()
			until, isBlocked := blocked.m[key]
			blocked.Unlock()

			if isBlocked && time.Now().Before(until) {
				return e.JSON(http.StatusTooManyRequests, map[string]string{
					"message": fmt.Sprintf("Blocked. Try again in %v", time.Until(until).Round(time.Second)),
				})
			}

			ctx, err := limiterInstance.Get(e.Request().Context(), key)
			if err != nil {
				return e.JSON(http.StatusInternalServerError, echo.Map{
					"error": "Limiter error",
				})
			}

			e.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", ctx.Limit))
			e.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", ctx.Remaining))
			if ctx.Reached {
				blocked.Lock()
				blocked.m[key] = time.Now().Add(Time)
				blocked.Unlock()

				return e.JSON(http.StatusTooManyRequests, map[string]string{
					"message": "Too many requests. Temporarily blocked for " + Time.String(),
				})
			}

			return next(e)
		}
	}
}
