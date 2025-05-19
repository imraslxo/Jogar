package api

import (
	"context"
	"futbikSecond/config"
	"futbikSecond/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// AuthHandler godoc
//
//	@Summary		Авторизация пользователя
//	@Description	Авторизирует пользователя
//	@Tags			Авторизация
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.AuthRequestDTO	true	"Данные пользователя"
//	@Success		201		{object}	map[string]interface{}	"Пользователь успешно авторизован"
//	@Failure		400		{object}	map[string]string	"Неверные данные запроса"
//	@Failure		500		{object}	map[string]string	"Не удалось авторизовать пользователя"
//	@Router			/auth [post]
//	@Security		BearerAuth
func AuthHandler(c *gin.Context) {
	var input models.AuthRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при парсинге: " + err.Error()})
		return
	}

	conn, err := config.DB.Acquire(c.Request.Context())
	if err != nil {
		log.Println("Ошибка подключения к базе данных: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer conn.Release()

	tx, err := conn.Begin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании транзакции: " + err.Error()})
	}

	query := "INSERT INTO \"user\" (tg_username, tg_first_name, tg_last_name, photo_url, is_premium, ui_language_code, allows_write_to_pm, auth_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	log.Println("Выполняется запрос: ", query)

	var userID uint64
	err = tx.QueryRow(
		c.Request.Context(),
		query,
		input.TgUsername,
		input.TgFirstName,
		input.TgLastName,
		input.PhotoURL,
		input.IsPremium,
		input.UILanguageCode,
		input.AllowsWriteToPM,
		input.AuthDate,
	).Scan(&userID)

	if err != nil {
		tx.Rollback(c.Request.Context())
		log.Println("Ошибка при создании пользователя: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		return
	}

	if err := tx.Commit(c.Request.Context()); err != nil {
		log.Println("Ошибка при коммите транзакции: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      userID,
		"message": "Пользователь успешно авторизован",
	})
}

// GetAuthUser godoc
//
// @Summary Получить список пользователей
// @Description Возвращает всех пользователей из таблицы "user".
// @Tags Авторизация
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /auth/users [get]
func GetAuthUser(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), "SELECT id, tg_username, tg_first_name, tg_last_name, photo_url, is_premium, ui_language_code, allows_write_to_pm, auth_date FROM \"user\"")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка в запросе: " + err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.TgUsername, &user.TgFirstName, &user.TgLastName, &user.PhotoURL, &user.IsPremium, &user.UILanguageCode, &user.AllowsWriteToPM, &user.AuthDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сканировании: " + err.Error()})
			return
		}
		users = append(users, user)
	}

	c.IndentedJSON(http.StatusOK, users)
}
