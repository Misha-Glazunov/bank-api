package integration_tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"
    "github.com/stretchr/testify/assert"
)

type UserCredentials struct {
    Email    string `json:"email"`
    Username string `json:"username"`
    Password string `json:"password"`
}

func TestUserRegistrationAndLogin(t *testing.T) {
    baseURL := "http://localhost:8080"
    
    // Test registration
    user := UserCredentials{
        Email:    "testuser@example.com",
        Username: "testuser",
        Password: "securepassword123",
    }

    payload, _ := json.Marshal(user)
    resp, err := http.Post(baseURL+"/register", "application/json", bytes.NewBuffer(payload))
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // Test login
    loginData := map[string]string{
        "email":    user.Email,
        "password": user.Password,
    }
    payload, _ = json.Marshal(loginData)
    
    resp, err = http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(payload))
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var result struct {
        Token string `json:"token"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    
    assert.NotEmpty(t, result.Token)
}
