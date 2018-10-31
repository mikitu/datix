package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// Header used to get/set request id
const HeaderKey = "X-Request-Id"

// RequestID middleware
func RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Request().Header.Get(HeaderKey)

		if id == "" {
			uuid, err := uuid.NewUUID()
			if err != nil {
				c.Error(err)
			}
			id = uuid.String()
		}
		c.Set("RequestID", id)
		log.Debugf("Set RequestId %+v ", id)
		c.Response().Header().Set(HeaderKey, id)
		err := next(c)
		return err
	}
}
