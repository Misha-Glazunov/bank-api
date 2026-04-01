# Bank API Service

Микросервис для управления банковскими операциями с использованием Go, PostgreSQL и Docker.

### Предварительные требования
- Docker и Docker Compose
- Go 1.24+ (для локальной разработки)

### Установка
git clone https://github.com/Misha-Glazunov/bank-api.git
cd bank-api
cp .env.example .env 
Пример данных
# Настройки базы данных
DB_HOST=postgres
DB_PORT=5432
DB_USER=bank_user
DB_PASSWORD=secure_dev_password
DB_NAME=bank_dev
DB_SSLMODE=disable

# Настройки JWT
JWT_SECRET=9$5zLq#2rT!pX8vKsYmNfWbE@dC&H*Gj
JWT_LIFETIME=24h

# Настройки приложения
HTTP_PORT=8080.
docker-compose up --build

Основные команды

# Запуск всех сервисов
docker-compose up

# Пересборка с очисткой данных
docker-compose down -v && docker-compose up --build

# Просмотр логов приложения
docker-compose logs -f api

# Выполнение миграций вручную
docker-compose run migrate

Тестирование
Интеграционные тесты
bash
go test -v ./internal/integration_tests/...
Пример запроса через cURL
bash
# Регистрация пользователя
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com", "username":"testuser", "password":"securepassword"}'

# Получение токена
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com", "password":"securepassword"}'
Docker окружение
PostgreSQL 15 на порту 5432

Миграции через migrate

Go приложение на порту 8080

Лимиты и ограничения
Максимальная сумма перевода: 1,000,000 ₽

Минимальный баланс: -50,000 ₽ (овердрафт)

Лимит запросов: 100 RPM

Безопасность
Все транзакции записываются в audit log

JWT срок действия: 24 часа

Хеширование паролей с bcrypt
