package middleware

import (
	"fmt"
	"gameAppProject/param"
	"gameAppProject/pkg/claim"
	"gameAppProject/pkg/errmsg"
	"gameAppProject/pkg/timestamp"
	"gameAppProject/service/presenceservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpsertPresence(service presenceservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			claims := claim.GetClaimsFromEchoContext(c)
			_, err = service.Upsert(c.Request().Context(), param.UpsertPresenceRequest{
				UserID:    claims.UserID,
				Timestamp: timestamp.Now(),
			})
			if err != nil {
				// TODO - log unexpected error
				fmt.Println("UpsertPresence err", err.Error())
				// we can just log the error and go to the next step(middleware, handler)
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.ErrorMsgSomethingWentWrong,
				})
			}

			return next(c)
		}
	}
}
