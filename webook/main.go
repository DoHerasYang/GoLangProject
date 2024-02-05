package main

import (
	"GoLangProject/webook/config"
	"GoLangProject/webook/internal/repository"
	"GoLangProject/webook/internal/repository/dao"
	"GoLangProject/webook/internal/service"
	"GoLangProject/webook/internal/web"
	"GoLangProject/webook/internal/web/middleware"
	"GoLangProject/webook/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func InitialDB() *gorm.DB {
	// Initialize DB ORM
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	_ = dao.InitUserTables(db)        // -- Initialized User Table in database
	_ = dao.InitUserProfileTables(db) // -- Initialized UserProfile Table in database
	return db
}

func InitialGin() *gin.Engine {
	// Gin Part
	ginEngine := gin.Default()
	ginEngine.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"}, // Allow front end to access backend response header
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	}))
	// JWT Middleware
	useJWT(ginEngine)
	// RateLimit
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	ginEngine.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	return ginEngine
}

func useJWT(server *gin.Engine) {
	loginJWT := &middleware.LoginJWTMiddlewareBuilder{}
	server.Use(loginJWT.CheckLogin())
}

func useSession(server *gin.Engine) {
	login := &middleware.LoginMiddlewareBuilder{}
	// Two Options to implement Cookie/memstore/Redis
	//store := memstore.NewStore(
	//	[]byte("LNTPPX5K7drGwzRk2yn7Fs3jbuxb71XLnRcIo9RnxAiSk7Hzrb6h5eTbl4kq7F7o"),
	//	[]byte("0BVTqBkxzN5BpWMzoTv0kvgBUiP879N8KO7vJieZXfKCRrvSUIA98OyCKnl3ZC1X"))
	store := cookie.NewStore([]byte("secret")) // 存储在Cookie中
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	[]byte("lENqlMABzIpbhm4bBzYoEBoXwmsfJ9ht"),
	//	[]byte("dn7Bx2dJ4SAaiSprPUXZmUeldjWksFxJ"),
	//)
	//if err != nil {
	//	panic(err)
	//}
	server.Use(
		sessions.Sessions("ssid", store),
		login.CheckLogin())
}

func initUserHandler(db *gorm.DB, ginEngine *gin.Engine) {
	userDaoStruct := dao.NewUserDAO(db)
	userRepositoryStruct := repository.NewUserRepository(userDaoStruct)
	userServiceStruct := service.NewUserService(userRepositoryStruct)
	userWebStruct := web.NewUserHandler(userServiceStruct)
	userWebStruct.RegisterRoutes(ginEngine)
}

func main() {
	// Initialized ORM MySQL DB
	db := InitialDB()
	// Initialized Gin
	ginEngine := InitialGin()
	// Register Handler
	initUserHandler(db, ginEngine)
	// Run Server
	err := ginEngine.Run(config.Config.GinHost.Addr)
	if err != nil {
		panic("无法启动服务器")
	}
}
