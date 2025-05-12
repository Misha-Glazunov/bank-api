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
    token := authenticateUser(t)

    // Create two accounts
    fromAccount := createAccount(t, token)
    toAccount := createAccount(t, token)

    // Test transfer
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
