package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

const updateTimeKey = "update_time"

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// gob to
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/signup" || ctx.Request.URL.Path == "/users/login" {
			return
		}
		sess := sessions.Default(ctx)
		userID := sess.Get("userId")
		if userID == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Flush Session
		// 看到了 17 分钟
		timeSessionVal := sess.Get(updateTimeKey)
		timeNow := time.Now()
		lastUpdateTime, ok := timeSessionVal.(time.Time)
		if timeSessionVal == nil || !ok || timeNow.Sub(lastUpdateTime) > time.Minute {
			resetSession(sess, userID)
		}
	}
}

func resetSession(sess sessions.Session, userID interface{}) {
	sess.Set(updateTimeKey, time.Now())
	sess.Set("userId", userID)
	err := sess.Save()
	if err != nil {
		fmt.Printf("Error Session: %v\n", err)
	}
}
