# Finance Tracker API

REST API на Go для отслеживания личных финансов: доходов и расходов.

## Технологии

- **Go** 
- **Fiber** – web-фреймворк
- **GORM** – ORM для работы с PostgreSQL
- **Swagger** – авто-документация API



## Функционал

- CRUD операции для транзакций  
- Получение текущего баланса  
- Swagger-документация



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
```

3. Установи зависимости:

```bash
go mod tidy
```

4. Запусти проект:

```bash
go run main.go
```



##  Эндпоинты

| Метод | Route                 | Описание                  |
|-------|-----------------------|---------------------------|
| GET   | /transactions         | Получить все транзакции   |
| POST  | /transactions         | Создать новую транзакцию  |
| PUT   | /transactions/:id     | Изменить транзакцию по ID |
| DELETE | /transactions/:id    | Удалить транзакцию по ID  |
| GET   | /balance              | Получить текущий баланс   |
| GET   | /swagger/*            | Swagger UI документация   |




##  Документация API

Swagger доступен по адресу:

```
http://localhost:3000/swagger/index.html
```


