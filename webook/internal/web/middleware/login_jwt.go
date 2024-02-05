package middleware

import (
	"GoLangProject/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/signup" || ctx.Request.URL.Path == "/users/login" {
			return
		}
		authCode := ctx.GetHeader("Authorization") // Bearer
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return []byte(web.JWTKey), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if uc.UserAgent != ctx.Request.UserAgent() {
			// TODO: Add Prometheus Monitor embedding point
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Flush Session
		expireTime := uc.ExpiresAt
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			newToken, err := token.SignedString([]byte(web.JWTKey))
			if err != nil {
				log.Println(err)
			} else {
				ctx.Header("x-jwt-token", newToken)
			}
		}
		// Set context user to Obtain From
		ctx.Set("user", uc)
	}
}
