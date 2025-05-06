package middlewares

import (
	"gotickitz/pkg"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) AccsessGate(allowedRole ...string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		payloads, exits := ctx.Get("payloads")
		if !exits {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "please login first",
			})
			return
		}
		userPayload, ok := payloads.(*pkg.Payload)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "your login identity is malformed, please login again",
			})
			return
		}
		if !slices.Contains(allowedRole, userPayload.Role) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"msg": "you do not have permission to access",
			})
			return
		}
		ctx.Next()
	}
}
