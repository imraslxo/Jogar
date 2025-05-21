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

// GetProfile godoc
// @Summary Получить список профилей
// @Description Возвращает массив профилей
// @Tags Профили
// @Accept json
// @Produce json
// @Success 200 {array} models.Profile
// @Router /profiles [get]
func GetProfiles(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), "SELECT p.id, p.pref_position, p.height, p.foot, p.age, p.playing_frequency, p.games_played, p.city, p.country FROM profiles p")
	if err != nil {
		log.Println("Ошибка при выполнении запроса: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var profile models.Profile

		err := rows.Scan(&profile.ID, &profile.PrefPosition, &profile.Height, &profile.Foot, &profile.Age, &profile.PlayingFrequency, &profile.GamesPlayed, &profile.City, &profile.Country)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		profiles = append(profiles, profile)
	}
	c.JSON(http.StatusOK, profiles)
}

// CreateProfile godoc
//
//	@Summary		Создать профиль
//	@Description	Добавляет новый профиль в систему
//	@Tags			Профили
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.ProfileCreateRequest	true	"Данные профиля"
//	@Success		201		{object}	models.Profile
//
// @Failure      404  {object}  map[string]string
//
//	@Router			/profiles [post]
//	@Security		BearerAuth
func PostProfile(c *gin.Context) {
	var input models.ProfileCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := config.DB.Acquire(c.Request.Context())
	if err != nil {
		log.Println("Ошибка подключения к базе данных: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer conn.Release()

	query := "INSERT INTO profiles(pref_position, height, foot, age, city, country, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	log.Println("Выполняется запрос: ", query)

	var profileID uint64
	err = conn.QueryRow(c.Request.Context(), query, input.PrefPosition, input.Height, input.Foot, input.Age, input.City, input.Country).Scan(&profileID)
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
		"id":      profileID,
		"message": "Профиль успешно добавлен",
	})
}
