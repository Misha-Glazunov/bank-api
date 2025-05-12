package integration_tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestMoneyTransfer(t *testing.T) {
    baseURL := "http://localhost:8080"
    
    // 1. Регистрация и логин
    user := map[string]string{
        "email":    "transfer_test@example.com",
        "username": "transfer_user",
        "password": "transfer_pass123",
    }
    token := registerAndLogin(t, user)

    // 2. Создание счетов
    fromAccount := createAccount(t, token)
    toAccount := createAccount(t, token)

    // 3. Тест перевода
    transferData := map[string]interface{}{
        "from_account": fromAccount,
        "to_account":   toAccount,
        "amount":       100.50,
    }
    payload, _ := json.Marshal(transferData)
    
    req, _ := http.NewRequest("POST", baseURL+"/transfer", bytes.NewBuffer(payload))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func createAccount(t *testing.T, token string) string {
    req, _ := http.NewRequest("POST", "http://localhost:8080/accounts", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    client := &http.Client{}
    resp, err := client.Do(req)
    assert.NoError(t, err)
    
    var account struct {
        ID string `json:"id"`
    }
    json.NewDecoder(resp.Body).Decode(&account)
    return account.ID
}

// Добавьте эту функцию если её нет в других файлах
func registerAndLogin(t *testing.T, user map[string]string) string {
    // Регистрация
    payload, _ := json.Marshal(user)
    resp, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(payload))
    assert.NoError(t, err)
    
    // Логин
    loginData := map[string]string{
        "email":    user["email"],
        "password": user["password"],
    }
    payload, _ = json.Marshal(loginData)
    
    resp, err = http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(payload))
    assert.NoError(t, err)
    
    var result struct {
        Token string `json:"token"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    return result.Token
}
