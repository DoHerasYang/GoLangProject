package main

import (
	"GoLangProject/webhook/internal/repository"
	"GoLangProject/webhook/internal/repository/dao"
	"GoLangProject/webhook/internal/service"
	"GoLangProject/webhook/internal/web"
	"GoLangProject/webhook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func InitialDB() *gorm.DB {
	// Initialize DB ORM
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
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
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	}))
	login := &middleware.LoginMiddlewareBuilder{}
	// Two Options to implement Cookie/memstore/Redis
	//store := memstore.NewStore(
	//	[]byte("LNTPPX5K7drGwzRk2yn7Fs3jbuxb71XLnRcIo9RnxAiSk7Hzrb6h5eTbl4kq7F7o"),
	//	[]byte("0BVTqBkxzN5BpWMzoTv0kvgBUiP879N8KO7vJieZXfKCRrvSUIA98OyCKnl3ZC1X"))
	//store := cookie.NewStore([]byte("secret")) // 存储在Cookie中
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("lENqlMABzIpbhm4bBzYoEBoXwmsfJ9ht"),
		[]byte("dn7Bx2dJ4SAaiSprPUXZmUeldjWksFxJ"),
	)
	if err != nil {
		panic(err)
	}
	ginEngine.Use(
		sessions.Sessions("ssid", store),
		login.CheckLogin())
	return ginEngine
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
	err := ginEngine.Run(":8080")
	if err != nil {
		panic("无法启动服务器")
	}
}
