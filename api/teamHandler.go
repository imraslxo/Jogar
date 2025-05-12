package api

import (
	"context"
	"futbikSecond/config"
	"futbikSecond/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// GetTeam godoc
// @Summary Получить список команд
// @Description Возвращает массив команд
// @Tags Команды
// @Accept json
// @Produce json
// @Success 200 {array} models.Team
// @Router /teams [get]
func GetTeams(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), "SELECT t.id, t.team_name, t.photo, t.playing_in, t.stadium, t.description FROM teams t")
	if err != nil {
		log.Println("Ошибка при выполнении запроса: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		err = rows.Scan(&team.ID, &team.TeamName, &team.Photo, &team.PlayingIn, &team.Stadium, &team.Description)
		if err != nil {
			log.Println("Ошибка при сканировании: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		teams = append(teams, team)
	}
	c.JSON(http.StatusOK, gin.H{"data": teams})
}

// CreateTeam godoc
//
//	@Summary		Создать команду
//	@Description	Добавляет новую команду в систему
//	@Tags			Команды
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.TeamCreateRequest	true	"Данные команды"
//	@Success		201		{object}	models.Team
//
// @Failure      404  {object}  map[string]string
//
//	@Router			/teams [post]
//	@Security		BearerAuth
func PostTeam(c *gin.Context) {
	var input models.TeamCreateRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := config.DB.Acquire(c.Request.Context())
	if err != nil {
		log.Println("Ошибка подключения к базе данных: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer conn.Release()

	query := "INSERT INTO teams(team_name, photo, playing_in, stadium, description) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	log.Println("Выполняется запрос: ", query)

	var teamID uint64
	err = conn.QueryRow(c.Request.Context(), query, input.TeamName, input.Photo, input.PlayingIn, input.Stadium, input.Description).Scan(&teamID)
	if err != nil {
		log.Println("Ошибка при выполнении запроса: ", err)
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      teamID,
		"message": "Команда успешно добавлена",
	})
}
