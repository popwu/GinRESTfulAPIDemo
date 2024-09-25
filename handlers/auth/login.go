package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"apidemo/models"
	"apidemo/utils"
)

/*
处理用户登录请求。

此函数验证用户凭据，如果验证成功，则生成并返回 JWT token。

参数 json格式:
nickname (str): 用户昵称
password (str): 密码

返回:
dict: 包含 JWT token 和用户信息的字典

返回示例:

	{
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"user": {
			"id": 1,
			"nickname": "example_user"
		}
	}

异常:
400 Bad Request: 如果请求数据无效或不完整
401 Unauthorized: 如果提供的凭据无效
403 Forbidden: 如果用户账户被禁用或没有登录权限
429 Too Many Requests: 如果用户在短时间内发送了过多请求
500 Internal Server Error: 如果服务器出现内部错误

注意:
- 所有的错误响应都会包含一个 'error' 键，其值为描述错误的字符串。
- 成功的登录会返回一个包含 JWT token 的 200 OK 响应。
- 客户端应当在后续请求的 Authorization 头中使用返回的 token。
*/
func Login(c *gin.Context) {
	// 解析请求中的 JSON 数据
	var loginData struct {
		Nickname string `json:"nickname" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证用户凭据
	user, err := models.UserGetByNickname(loginData.Nickname)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not found username"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username or password error"})
		return
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(user.ID, user.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 返回 JWT token 和用户信息
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"nickname": user.Nickname,
		},
	})
}
