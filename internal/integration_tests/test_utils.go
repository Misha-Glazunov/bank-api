package integration_tests

import (
    "database/sql"
    "fmt"
    "os"
    "testing"
    _ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    teardown()
    os.Exit(code)
}

func setup() {
    // 1. Подключение к тестовой БД
    connStr := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        "localhost", 5432, "postgres", "postgres", "bank_test",
    )
    
    var err error
    testDB, err = sql.Open("postgres", connStr)
    if err != nil {
        panic(fmt.Sprintf("DB connection failed: %v", err))
    }

    // 2. Запуск миграций
    runMigrations()

    // 3. Очистка тестовых данных
    cleanTestData()
}

func teardown() {
    // 1. Закрытие соединения с БД
    if testDB != nil {
        testDB.Close()
    }
    
    // 2. Дополнительная очистка (опционально)
}

func runMigrations() {
    // Реализация запуска миграций 
    // (можно использовать exec.Command для вызова migrate)
}

func cleanTestData() {
    // Очистка всех тестовых таблиц
    tables := []string{"users", "accounts", "transactions"}
    for _, table := range tables {
        _, err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
        if err != nil {
            panic(fmt.Sprintf("Failed to truncate table %s: %v", table, err))
        }
    }
}
