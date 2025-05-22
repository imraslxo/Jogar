package api

import (
	"context"
	"futbikSecond/config"
	"futbikSecond/models"
	"github.com/gin-gonic/gin"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"log"
	"net/http"
)

// NewAuthHandler godoc
// @Summary Авторизация по initData
// @Description Принимает строку initData в теле запроса и создает или обновляет пользователя в базе
// @Tags Авторизация
// @Accept json
// @Produce json
// @Param input body models.InitDataRequest true "Строка initData"
// @Success 201 {object} map[string]interface{} "Пользователь успешно авторизован"
// @Failure 400 {object} map[string]string "Неверные данные запроса"
// @Failure 500 {object} map[string]string "Не удалось авторизовать пользователя"
// @Router /auth [post]
func NewAuthHandler(c *gin.Context) {
	var req models.InitDataRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}
	input, err := initdata.Parse(req.InitData)
	if err != nil {
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
		input.User.Username,
		input.User.FirstName,
		input.User.LastName,
		input.User.PhotoURL,
		input.User.IsPremium,
		input.User.LanguageCode,
		input.User.AllowsWriteToPm,
		input.AuthDateRaw,
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
