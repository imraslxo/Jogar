package api

import (
	"context"
	"futbikSecond/config"
	"futbikSecond/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

// PostProfileFirstPg godoc
//
// @Summary      Создание профиля и привязка к пользователю
// @Description  Добавляет запись в таблицу profiles и обновляет поле profile_id у пользователя
// @Tags         Профили
// @Accept       json
// @Produce      json
// @Param        tg_userid   path      int                          true  "Telegram ID пользователя"
// @Param        input     body      models.ProfileCreateFirstDTO true  "Данные профиля"
// @Success      200       {object}  map[string]interface{}        "Профиль успешно создан и привязан"
// @Failure      400       {object}  map[string]string             "Неверный ввод"
// @Failure      500       {object}  map[string]string             "Ошибка сервера"
// @Router       /users/by-tg/{tg_userid}/profile [post]
func PostProfileFirstPg(c *gin.Context) {
	var input models.ProfileCreateFirstDTO
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

	tx, err := conn.Begin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании транзакции: " + err.Error()})
	}

	commited := false
	defer func() {
		if !commited {
			tx.Rollback(c.Request.Context())
		}
	}()

	tguserID := c.Param("tg_userid")
	query := "INSERT INTO profiles (app_language, country, city, tg_user_id) VALUES ($1, $2, $3, $4) RETURNING id"
	log.Println("Выполняется запрос: ", query)

	var profileID int64
	err = tx.QueryRow(c.Request.Context(), query, input.AppLanguage, input.Country, input.City, tguserID).Scan(&profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выполнении запроса: " + err.Error()})
		return
	}

	secondQuery := "UPDATE \"user\" SET profile_id = $1 WHERE id = $2"
	log.Println("Выполняется запрос: ", query)

	_, err = tx.Exec(c.Request.Context(), secondQuery, profileID, tguserID)

	if err := tx.Commit(c.Request.Context()); err != nil {
		log.Println("Ошибка при коммите транзакции: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit error"})
		return
	}
	commited = true
}
