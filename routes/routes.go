package routes

import (
	"futbikSecond/api"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/auth", api.AuthHandler)

	r.GET("/users/team/:team_id", api.GetUsersByTeamID)
	r.GET("/teams/players/count/:team_name", api.GetNumberOfPlayersInTeam)
	r.GET("/users", api.GetUsers)
	r.POST("/users", api.CreateUser)
	r.POST("/users/:user_id/team/:team_id", api.JoinTeam)
	r.DELETE("/users/:user_id/team", api.LeaveTeam)

	r.GET("/profiles", api.GetProfiles)
	r.POST("/profiles", api.PostProfile)

	r.GET("/teams", api.GetTeams)
	r.POST("/teams", api.PostTeam)

	return r
}
