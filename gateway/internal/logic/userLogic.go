package logic

import (
	"context"
	"errors"

	"gateway/internal/code"
	"gateway/internal/database/crud"
	"gateway/internal/database/model"
	"gateway/pkg/encrypt"
	"gateway/pkg/jwt"
	"gateway/pkg/logger"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginLogic(ctx context.Context, userName, password string) (string, int64, error) {
	// 查询Id和password
	userId, hashedPassword, err := crud.GetPasswordAndIdByName(ctx, userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorf("userName %s not found", userName)
			return "", 0, code.ErrNotFound
		} else {
			logger.Errorf("database error: %v", err)
			return "", 0, code.ErrDatabase
		}
	}

	//对比密码
	err = encrypt.CheckPassword(hashedPassword, password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			logger.Errorf("password %s wrong", password)
			return "", 0, code.ErrPassword
		} else {
			logger.Errorf("unknown error: %s", err)
			return "", 0, code.ErrUnknown
		}
	}

	//初始化token
	token, expireAt, err := jwt.NewToken(userId, userName)
	if err != nil {
		logger.Errorf("NewToken: %v", err)
		return "", 0, code.ErrJwtCreate
	}

	return token, expireAt, nil
}

func RegisterLogic(ctx context.Context, userName, password string) error {
	//密码加盐哈希
	hashedPassword, err := encrypt.HashPassword(password)
	if err != nil {
		logger.Errorf("HashPassword: %v", err)
		return code.ErrEncrypt
	}

	u := &model.User{
		UserName: userName,
		Password: hashedPassword,
	}
	err = crud.CreateUser(ctx, u)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.Is(err, gorm.ErrDuplicatedKey) || (errors.As(err, &mysqlErr) && mysqlErr.Number == 1062) {
			logger.Errorf("register duplicate user_name: %s", userName)
			return code.ErrUserNameExists
		}
		logger.Errorf("database error: %v", err)
		return code.ErrDatabase
	}

	return nil
}
