package postgre

import (
	"context"
	"database/sql"

	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/models"
	storageErrors "github.com/11Petrov/gopherloyal/internal/storage/errors"
	"github.com/11Petrov/gopherloyal/internal/utils"
)

func (d *Database) UserRegister(ctx context.Context, user *models.Users) (int, error) {
	log := logger.FromContext(ctx)

	// Проверяем, занят ли логин пользователя
	var count int
	err := d.db.QueryRow(ctx, "SELECT COUNT(*) FROM Users WHERE login = $1", user.Login).Scan(&count)
	if err != nil {
		log.Errorf("error checking login: %s", err)
		return 0, err
	}
	if count > 0 {
		log.Warn("login already taken")
		return 0, storageErrors.ErrLoginTaken
	}

	// Добавляем нового пользователя
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Errorf("error hashing password: %s", err)
		return 0, err
	}

	var userID int
	err = d.db.QueryRow(ctx, "INSERT INTO Users (login, password_hash) VALUES ($1, $2) RETURNING user_id", user.Login, hashedPassword).Scan(&userID)
	if err != nil {
		log.Errorf("error inserting user: %s", err)
		return 0, err
	}

	log.Info("user successfully registered")
	return userID, nil
}

func (d *Database) UserLogin(ctx context.Context, user *models.Users) (*models.Users, error) {
	log := logger.FromContext(ctx)

	// Проверяем учетные данные пользователя
	var storedPasswordHash string
	err := d.db.QueryRow(ctx, "SELECT user_id, login, password_hash FROM Users WHERE login = $1", user.Login).Scan(&user.ID, &user.Login, &storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn("user not found")
			return nil, storageErrors.ErrUserNotFound
		}
		log.Errorf("error checking user credentials: %s", err)
		return nil, err
	}

	// Сравниваем пароль с hash-паролем
	if !utils.CheckPasswordHash(user.Password, storedPasswordHash) {
		log.Warn("invalid password")
		return nil, storageErrors.ErrInvalidPassword
	}

	log.Info("user successfully logged in")
	return user, nil
}
