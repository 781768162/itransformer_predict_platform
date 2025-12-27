package crud

import (
	"context"

	"gateway/internal/database"
	"gateway/internal/database/model"
)

func GetUserByName(ctx context.Context, userName string) (*model.User, error) {
	// 全量数据
	var u model.User
	err := database.DB.WithContext(ctx).
		Where("user_name = ?", userName).
		First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetPasswordAndIdByName(ctx context.Context, userName string) (int64, string, error) {
	//仅userId和password
	var u model.User
	err := database.DB.WithContext(ctx).
		Select("user_id", "password").
		Where("user_name = ?", userName).
		First(&u).Error
	if err != nil {
		return 0, "", err
	}
	return u.UserID, u.Password, nil
}

func CreateUser(ctx context.Context, u *model.User) error {
	//userId应为0
	return database.DB.WithContext(ctx).Create(u).Error
}
