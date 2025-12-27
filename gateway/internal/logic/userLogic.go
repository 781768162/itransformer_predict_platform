package logic

import (
	"context"
	"errors"

	"gateway/internal/code"
	"gateway/internal/database/crud"
	"gateway/internal/database/model"
	"gateway/pkg/encrypt"
	"gateway/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginLogic(ctx context.Context, userName, password string) (string, int64, error) {
	// 查询Id和password
	 userId, hashedPassword, err := crud.GetPasswordAndIdByName(ctx, userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", 0, code.ErrNotFound
		}else {
			return "", 0, code.ErrDatabase
		}
	}

	//对比密码
	err = encrypt.CheckPassword(hashedPassword, password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", 0, code.ErrPassword
		}else {
			return "", 0, code.ErrUnknown
		}
	}

	//初始化token
	token, expireAt, err := jwt.NewToken(userId, userName)
	if err != nil {
		return "", 0, code.ErrJwtCreate
	}

	return token, expireAt, nil
}

func RegisterLogic(ctx context.Context, userName, password string) error {
	//密码加盐哈希
	hashedPassword, err := encrypt.HashPassword(password)
	if err != nil {
		return code.ErrEncrypt
	}

	u := &model.User{
		UserName: userName,
		Password: hashedPassword,
	}
	err = crud.CreateUser(ctx, u)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return code.ErrUserNameExists
		}else {
			return code.ErrDatabase
		}
	}

	return nil
}