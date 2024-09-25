package api

import (
	"net/http"
	"strconv"

	"apidemo/models"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	// 获取查询参数并转换为 int64
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		page = 1 // 使用默认值
	}

	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		limit = 10 // 使用默认值
	}

	search := c.Query("search")
	status := c.Query("status")

	// 将 page 和 limit 转换回字符串，因为 UserGetAll 期望字符串参数
	pageStr := strconv.FormatInt(page, 10)
	limitStr := strconv.FormatInt(limit, 10)

	totalUsers, users, err := models.UserGetAll(pageStr, limitStr, search, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (totalUsers + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
			"total_users": totalUsers,
		},
	})
}
