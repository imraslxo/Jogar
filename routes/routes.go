package routes

import (
	"futbikSecond/api"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://jogar3.web.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/auth", api.AuthHandler)
	r.GET("/auth/users", api.GetAuthUser)

	r.GET("/users/team/:team_id", api.GetUsersByTeamID)
	r.GET("/teams/players/count/:team_name", api.GetNumberOfPlayersInTeam)
	r.GET("/users", api.GetUsers)
	r.POST("/users", api.CreateUser)
	r.POST("/users/:user_id/team/:team_id", api.JoinTeam)
	r.DELETE("/users/:user_id/team", api.LeaveTeam)

	r.GET("/profiles", api.GetProfiles)
	r.POST("/profiles/by-tg/:tg_userid/profile", api.PostProfileFirstPg)
	r.POST("/profiles/second/:tg_userid/profile", api.PostProfileSecondPg)

	r.GET("/teams", api.GetTeams)
	r.POST("/teams", api.PostTeam)

	return r
}
