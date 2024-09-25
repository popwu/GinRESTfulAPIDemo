package store

import (
	"fmt"
	"reflect"
	"strings"

	"apidemo/dbconfig"
	"apidemo/models"
)

// ModelsList 包含所有需要迁移和清理的模型
var ModelsList = []interface{}{
	&models.User{},
	// 在这里添加其他模型...
}

func AutoMigrate() error {
	for _, model := range ModelsList {
		if err := dbconfig.DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("自动迁移失败 %T: %w", model, err)
		}
	}
	return nil
}

func ClearTestDB() error {
	// 禁用外键检查
	if err := dbconfig.DB.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("禁用外键检查失败: %w", err)
	}

	// 删除所有模型对应的表
	for _, model := range ModelsList {
		tableName := getTableName(model)
		if err := dbconfig.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tableName)).Error; err != nil {
			return fmt.Errorf("删除表 %s 失败: %w", tableName, err)
		}
	}

	// 重新启用外键检查
	if err := dbconfig.DB.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return fmt.Errorf("启用外键检查失败: %w", err)
	}

	return nil
}

// getTableName 获取模型对应的表名
func getTableName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.ToLower(t.Name()) + "s"
}
