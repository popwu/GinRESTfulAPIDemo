package auth

import (
	"apidemo/dbconfig"
	"apidemo/models"
	"apidemo/store"
	"apidemo/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("Failed to load environment variables: %v", err)
	}

	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 初始化测试数据库
	err = dbconfig.InitTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	// 清理测试数据库
	err = store.ClearTestDB()
	if err != nil {
		t.Fatalf("Failed to clear test database: %v", err)
	}

	err = store.AutoMigrate()
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// 创建测试用户
	err = models.CreateUser("testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// 验证用户是否正确创建
	user, err := models.UserGetByNickname("testuser")
	if err != nil {
		t.Fatalf("Failed to retrieve test user: %v", err)
	}
	t.Logf("Test user created: %+v", user)

	// 测试用例
	tests := []struct {
		name           string
		inputJSON      string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Valid Login",
			inputJSON:      `{"nickname": "testuser", "password": "password123"}`,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"user": map[string]interface{}{
					"id":       float64(1), // JSON 将整数解码为 float64
					"nickname": "testuser",
				},
			},
		},
		{
			name:           "Invalid Password",
			inputJSON:      `{"nickname": "testuser", "password": "wrongpassword"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid credentials",
			},
		},
		{
			name:           "User Not Found",
			inputJSON:      `{"nickname": "nonexistentuser", "password": "password123"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid credentials",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个新的 Gin 上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 模拟请求
			c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(tt.inputJSON))
			c.Request.Header.Set("Content-Type", "application/json")

			// 调用被测试的函数
			Login(c)

			// 断言
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			t.Logf("Response: %+v", response)

			if tt.expectedStatus == http.StatusOK {
				// 检查 token 是否存在，但不检查具体值
				assert.Contains(t, response, "token")
				token := response["token"]
				delete(response, "token")
				t.Logf("Generated token: %v", token)

				assert.Equal(t, tt.expectedBody, response)
			} else {
				// 检查错误信息
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}

}
