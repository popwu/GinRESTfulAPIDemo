package models

import (
	"apidemo/dbconfig"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Nickname string `json:"nickname" gorm:"uniqueIndex"`
	Password string `json:"-"` // 使用 json:"-" 确保密码不会被序列化到 JSON 中
}

func (User) TableName() string {
	return "users"
}

func CreateUser(nickname, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Nickname: nickname,
		Password: string(hashedPassword),
	}

	return dbconfig.DB.Create(&user).Error
}

func UserGetByNickname(nickname string) (*User, error) {
	var user User
	result := dbconfig.DB.Where("nickname = ?", nickname).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UserGetAll(page, limit, search, status string) (totalUsers int64, users []*User, err error) {
	// 将 users 初始化为空的指针切片
	users = []*User{}

	query := dbconfig.DB.Model(&User{})

	// 应用搜索条件
	if search != "" {
		query = query.Where("nickname LIKE ?", "%"+search+"%")
	}

	// 应用状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总用户数
	if err := query.Count(&totalUsers).Error; err != nil {
		return 0, nil, err
	}

	// 应用分页
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	if pageInt > 0 && limitInt > 0 {
		offset := (pageInt - 1) * limitInt
		query = query.Offset(offset).Limit(limitInt)
	}

	// 执行查询
	if err := query.Find(&users).Error; err != nil {
		return 0, nil, err
	}

	return totalUsers, users, nil
}
