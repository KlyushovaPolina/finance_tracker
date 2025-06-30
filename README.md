# Finance Tracker API

REST API на Go для отслеживания личных финансов: доходов, расходов и баланса с поддержкой аутентификации пользователей.

## Технологии

- **Go**: Основной язык программирования.
- **Fiber**: Лёгкий и быстрый веб-фреймворк.
- **GORM**: ORM для работы с PostgreSQL.
- **Swagger**: Автоматическая документация API.
- **JWT**: Аутентификация на основе JSON Web Tokens.
- **bcrypt**: Хэширование паролей для безопасного хранения.


## Функционал

- **Аутентификация**:
  - Регистрация пользователей (`/auth/register`).
  - Вход пользователей с выдачей JWT-токена (`/auth/login`).
- **Транзакции**:
  - CRUD-операции для транзакций (создание, получение, обновление, удаление).
  - Привязка транзакций к аутентифицированному пользователю.
- **Баланс**:
  - Получение текущего баланса (доходы минус расходы).
- **Документация**:
  - Интерактивная Swagger-документация API.


##  Установка

1. Клонируй репозиторий:

```bash
git clone https://github.com/KlyushovaPolina/finance_tracker.git
cd finance-tracker
```

2. Создай `.env` файл:

```
DB_USER=your_postgres_username
DB_PASSWORD=your_postgres_password
DB_NAME=your_database_name
DB_PORT=your_db_port
JWT_SECRET=your-secret-key
```

3. Установи зависимости:

```bash
go mod tidy
```

4. Запусти проект:

```bash
go run main.go
```


5. Откройте Swagger UI для тестирования API:

```
http://localhost:3000/swagger/
```

## Эндпоинты

Все маршруты, начинающиеся с `/api/`, требуют JWT-токена в заголовке `Authorization: Bearer <token>`,



