package postgre

import (
	"context"
	"database/sql"

	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/11Petrov/gopherloyal/internal/models"
	storageErrors "github.com/11Petrov/gopherloyal/internal/storage/errors"
	"github.com/11Petrov/gopherloyal/internal/utils"
)

func (d *Database) UserRegister(ctx context.Context, user *models.UserAuth) error {
	log := logger.LoggerFromContext(ctx)

	// Проверяем, занят ли логин пользователя
	var count int
	err := d.db.QueryRow(ctx, "SELECT COUNT(*) FROM Users WHERE login = $1", user.Login).Scan(&count)
	if err != nil {
		log.Errorf("error checking login: %s", err)
		return err
	}
	if count > 0 {
		log.Warn("login already taken")
		return storageErrors.ErrLoginTaken
	}

	// Добавляем нового пользователя
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Errorf("error hashing password: %s", err)
		return err
	}
	_, err = d.db.Exec(ctx, "INSERT INTO Users (login, password_hash) VALUES ($1, $2)", user.Login, hashedPassword)
	if err != nil {
		log.Errorf("error inserting user: %s", err)
		return err
	}

	log.Info("user successfully registered")
	return nil
}

func (d *Database) UserLogin(ctx context.Context, user *models.UserAuth) error {
	log := logger.LoggerFromContext(ctx)

	// Проверяем учетные данные пользователя
	var storedPasswordHash string
	err := d.db.QueryRow(ctx, "SELECT password_hash FROM Users WHERE login = $1", user.Login).Scan(&storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn("user not found")
			return storageErrors.ErrUserNotFound
		}
		log.Errorf("error checking user credentials: %s", err)
		return err
	}

	// Сравниваем пароль с hash-паролем
	if !utils.CheckPasswordHash(user.Password, storedPasswordHash) {
		log.Warn("invalid password")
		return storageErrors.ErrInvalidPassword
	}

	log.Info("user successfully logged in")
	return nil
}
