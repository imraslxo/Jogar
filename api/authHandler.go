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

// AuthHandler godoc
// @Summary Авторизация по initData
// @Description Принимает строку initData и создает нового пользователя, если он еще не зарегистрирован
// @Tags Авторизация
// @Accept json
// @Produce json
// @Param input body models.InitDataRequest true "InitData от Telegram"
// @Success 201 {object} map[string]interface{} "Пользователь успешно авторизован (создан)"
// @Success 200 {object} map[string]interface{} "Пользователь уже зарегистрирован"
// @Failure 400 {object} map[string]interface{} "Неверные данные запроса или ошибка парсинга"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера (БД, транзакция и т.п.)"
// @Router /auth [post]
func AuthHandler(c *gin.Context) {
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

	commited := false
	defer func() {
		if !commited {
			tx.Rollback(c.Request.Context())
		}
	}()

	var exist bool
	err = tx.QueryRow(c.Request.Context(), "SELECT EXISTS (SELECT 1 FROM \"user\" WHERE tg_userid = $1)", input.User.ID).Scan(&exist)
	if err != nil {
		log.Println("Ошибка в запросе проверки ID: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось проверить ID пользователя: "})
		return
	}

	if exist {
		c.JSON(http.StatusOK, gin.H{
			"message": "Пользователь уже зарегистрирован",
			"exist":   exist,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exist": exist,
	})

	query := "INSERT INTO \"user\" (tg_username, tg_userid, tg_first_name, tg_last_name, photo_url, is_premium, ui_language_code, allows_write_to_pm, auth_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9) RETURNING id"
	log.Println("Выполняется запрос: ", query)

	var userID uint64
	err = tx.QueryRow(
		c.Request.Context(),
		query,
		input.User.Username,
		input.User.ID,
		input.User.FirstName,
		input.User.LastName,
		input.User.PhotoURL,
		input.User.IsPremium,
		input.User.LanguageCode,
		input.User.AllowsWriteToPm,
		input.AuthDateRaw,
	).Scan(&userID)

	if err := tx.Commit(c.Request.Context()); err != nil {
		log.Println("Ошибка при коммите транзакции: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit error"})
		return
	}
	commited = true
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
	rows, err := config.DB.Query(context.Background(), "SELECT id, tg_username, tg_userid, tg_first_name, tg_last_name, photo_url, is_premium, ui_language_code, allows_write_to_pm, auth_date FROM \"user\"")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка в запросе: " + err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.TgUsername, &user.TgUserID, &user.TgFirstName, &user.TgLastName, &user.PhotoURL, &user.IsPremium, &user.UILanguageCode, &user.AllowsWriteToPM, &user.AuthDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сканировании: " + err.Error()})
			return
		}
		users = append(users, user)
	}

	c.IndentedJSON(http.StatusOK, users)
}
