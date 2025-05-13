package api

import (
	"context"
	"futbikSecond/config"
	"futbikSecond/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/jackc/pgx/v5/pgconn"
)

// GetUsers godoc
// @Summary Получить список пользователей
// @Description Возвращает массив пользователей
// @Tags Пользователи
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), "SELECT u.id, u.tg_username, u.tg_first_name, u.tg_last_name, u.photo_url, u.tg_last_login, u.registered_at, u.ui_language_code, u.allows_write_to_pm, p.id, p.pref_position, p.height, p.foot, p.age, p.playing_frequency, p.games_played, p.city, p.country,t.id, t.team_name, t.photo, t.playing_in, t.stadium, t.discription FROM users u LEFT JOIN profiles p ON u.profile_id = p.id LEFT JOIN teams t ON u.team_id = t.id")
	if err != nil {
		log.Println("Ошибка при выполнении запроса: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		user := models.User{}
		profile := models.Profile{}
		team := models.Team{}
		var profileID *uint64 // NULL-safe

		err := rows.Scan(&user.ID, &user.TgUsername, &user.TgFirstName, &user.TgLastName,
			&user.PhotoURL, &user.IsPremium, &user.AuthDate,
			&user.UILanguageCode, &user.AllowsWriteToPM,
			&profileID, &profile.PrefPosition, &profile.Height,
			&profile.Foot, &profile.Age, &profile.PlayingFrequency,
			&profile.GamesPlayed, &profile.City, &profile.Country,
			&team.ID, &team.TeamName, &team.Photo, &team.PlayingIn,
			&team.Stadium, &team.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if profile.ID != 0 {
			user.Profile = &profile
		}
		if team.ID != 0 {
			user.Team = &team
		}

		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser godoc
//
//	@Summary		Создать пользователя с профилем
//	@Description	Создает нового пользователя и профиль для него в одном запросе
//	@Tags			Пользователи
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.UserCreateRequest	true	"Данные пользователя"
//	@Success		201		{object}	map[string]interface{}	"Пользователь успешно создан"
//	@Failure		400		{object}	map[string]string	"Неверные данные запроса"
//	@Failure		500		{object}	map[string]string	"Не удалось создать пользователя"
//	@Router			/users [post]
//	@Security		BearerAuth
func CreateUser(c *gin.Context) {
	var input models.UserCreateRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
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
		log.Println("Ошибка при создании транзакции: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction error"})
	}

	query := "INSERT INTO users(tg_username, tg_first_name, tg_last_name, photo_url, ui_language_code, allows_write_to_pm, registered_at) VALUES ($1, $2, $3, $4, $5, $6, NOW()) RETURNING id"
	log.Println("Выполняется запрос: ", query)

	var userID uint64
	err = tx.QueryRow(
		c.Request.Context(),
		query,
		input.TgUsername,
		input.TgFirstName,
		input.TgLastName,
		input.PhotoURL,
		input.UILanguageCode,
		input.AllowsWriteToPM,
	).Scan(&userID)

	if err != nil {
		tx.Rollback(c.Request.Context())
		log.Println("Ошибка при создании пользователя: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		return
	}

	queryProfile := "INSERT INTO profiles(pref_position, height, foot, age, city, country, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	var profileID uint64
	err = tx.QueryRow(c.Request.Context(), queryProfile, input.ProfileCreateRequest.PrefPosition, input.ProfileCreateRequest.Height, input.ProfileCreateRequest.Foot, input.ProfileCreateRequest.Age, input.ProfileCreateRequest.City, input.ProfileCreateRequest.Country, userID).Scan(&profileID)
	if err != nil {
		tx.Rollback(c.Request.Context())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании профиля: " + err.Error()})
		return
	}

	updateUserQuery := "UPDATE users SET profile_id = $1 WHERE id = $2"
	_, err = tx.Exec(c.Request.Context(), updateUserQuery, profileID, userID)
	if err != nil {
		tx.Rollback(c.Request.Context())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении пользователя с профилем"})
		return
	}

	if err := tx.Commit(c.Request.Context()); err != nil {
		log.Println("Ошибка при коммите транзакции: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      userID,
		"message": "Пользователь и профиль успешно созданы",
	})
}

// JoinTeam добавляет пользователя в команду
// @Summary Присоединение пользователя к команде
// @Description Добавляет пользователя в команду по ID команды и ID пользователя
// @Tags Пользователь и команды
// @Accept  json
// @Produce  json
// @Param team_id path int true "ID команды"
// @Param user_id path int true "ID пользователя"
// @Success 200 {object} map[string]string "Пользователь присоединен к команде"
// @Failure 400 {object} map[string]string "Недопустимые параметры запроса"
// @Failure 500 {object} map[string]string "Не удалось присоединиться к команде"
// @Router /users/{user_id}/team/{team_id} [post]
func JoinTeam(c *gin.Context) {
	userID := c.Param("user_id")
	teamID := c.Param("team_id")

	query := "UPDATE \"user\" SET team_id = $1 WHERE id = $2"

	_, err := config.DB.Exec(context.Background(), query, teamID, userID)
	if err != nil {
		log.Println("Ошибка при добавлении пользователя в команду: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при присоединении к команде"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь присоединился к команде"})
}

// LeaveTeam удаляет пользователя из команды
// @Summary Выход пользователя из команды
// @Description Убирает пользователя из команды, устанавливая team_id в NULL
// @Tags Пользователь и команды
// @Accept  json
// @Produce  json
// @Param user_id path int true "ID пользователя"
// @Success 200 {object} map[string]string "Пользователь покинул команду"
// @Failure 400 {object} map[string]string "Недопустимые параметры запроса"
// @Failure 500 {object} map[string]string "Не удалось покинуть команду"
// @Router /users/{user_id}/team [delete]
func LeaveTeam(c *gin.Context) {
	userID := c.Param("user_id")

	query := "UPDATE \"user\" SET team_id = NULL WHERE id = $1"

	_, err := config.DB.Exec(context.Background(), query, userID)
	if err != nil {
		log.Println("Ошибка при выходе пользователя из команды: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выходе из команды"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь вышел из команды"})
}

// GetUsersByTeamID godoc
// @Summary Получение пользователей по ID команды
// @Description Возвращает список пользователей, принадлежащих к команде с указанным ID
// @Tags Пользователь и команды
// @Param team_id path string true "ID команды"
// @Produce json
// @Success 200 {array} models.AuthRequestDTO
// @Failure 500 {object} map[string]string "Ошибка при выполнении запроса или сканировании"
// @Router /users/team/{team_id} [get]
func GetUsersByTeamID(c *gin.Context) {
	teamID := c.Param("team_id")

	query := "SELECT tg_username, tg_first_name, tg_last_name, photo_url, is_premium, ui_language_code, allows_write_to_pm, auth_date FROM \"user\" JOIN teams ON \"user\".team_id = teams.id WHERE teams.id = $1"
	log.Println("Выполняется запрос... ", query)

	rows, err := config.DB.Query(context.Background(), query, teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выводе пользователей: " + err.Error()})
		return
	}
	defer rows.Close()

	var users []models.AuthRequestDTO
	for rows.Next() {
		var user models.AuthRequestDTO
		err := rows.Scan(
			&user.TgUsername,
			&user.TgFirstName,
			&user.TgLastName,
			&user.PhotoURL,
			&user.IsPremium,
			&user.UILanguageCode,
			&user.AllowsWriteToPM,
			&user.AuthDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сканировании: " + err.Error()})
			return
		}

		users = append(users, user)
	}

	c.IndentedJSON(http.StatusOK, users)
}

// GetNumberOfPlayersInTeam godoc
// @Summary Получить количество игроков в команде
// @Description Возвращает количество игроков в команде по имени. Если команд с таким именем несколько, вернёт массив.
// @Tags Пользователь и команды
// @Param team_name path string true "Название команды"
// @Produce json
// @Success 200 {array} models.TeamWithCount
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /teams/players/count/{team_name} [get]
func GetNumberOfPlayersInTeam(c *gin.Context) {
	teamName := c.Param("team_name")

	query := "SELECT teams.team_name, COUNT(\"user\".team_id) FROM teams LEFT JOIN \"user\" ON \"user\".team_id = teams.id WHERE teams.team_name = $1 GROUP BY teams.team_name"
	log.Println("Выполняется запрос... ", query)

	rows, err := config.DB.Query(context.Background(), query, teamName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка в запросе: " + err.Error()})
		return
	}

	var teams []models.TeamWithCount
	for rows.Next() {
		var team models.TeamWithCount
		err := rows.Scan(&team.Name, &team.Players_count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сканировании: " + err.Error()})
			return
		}
		teams = append(teams, team)
	}

	c.IndentedJSON(http.StatusOK, teams)
}
