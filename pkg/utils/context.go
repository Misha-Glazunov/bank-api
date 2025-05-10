package utils

import (
	"context"
	"errors"
)

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}
