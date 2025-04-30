package middleware

import (
	"net/http"
	"strings"
	"sync"

	"spmb/back-end/models/response"
	"spmb/back-end/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Success: false,
				Message: "Tidak memiliki Akses",
				Error:   map[string]string{"authorization": "No Header Found"},
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Success: false,
				Message: "Tidak memiliki Akses",
				Error:   map[string]string{"authorization": "Invalid Token"},
			})
			return
		}

		var claims map[string]interface{}
		var err error
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			claims, err = utils.ParseToken(tokenStr)
		}()

		wg.Wait()

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Success: false,
				Message: "Tidak memiliki Akses",
				Error:   map[string]string{"authorization": err.Error()},
			})
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Success: false,
				Message: "Server Error",
				Error:   map[string]string{"user_id": "User  ID format is invalid"},
			})
			return
		}

		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Success: false,
				Message: "Server Error",
				Error:   map[string]string{"user_id": "User  ID format is invalid"},
			})
			return
		}

		// Store it as uuid.UUID in context
		ctx.Set("user_id", userUUID)
		ctx.Next()
	}
}
