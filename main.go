package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres" //драйвер для постгри
	"gorm.io/gorm"            //для работы с бд (ORM)

	"github.com/joho/godotenv"
	"os"

	"github.com/gofiber/swagger"
	_ "finance-tracker/docs"

	"golang.org/x/crypto/bcrypt" //для хеширования пароля
	"github.com/golang-jwt/jwt/v5" //для генерации токена
	jwtware "github.com/gofiber/contrib/jwt" //для проверки токена
)

type Transaction struct { //модель
	ID          uint	`gorm:"primaryKey"`
	UserID      uint       `gorm:"not null"` // Привязка к пользователю
	Amount      float64  `gorm:"not null"`
	Type        string   `gorm:"not null; check:type_check,type IN ('income','expense')"` //тип транзакции - трата или расход
	Category    string  //категория - еда, одежда и тд
	Description *string //может быть пустым
	Date        time.Time 
	CreatedAt   time.Time //автоматически создается GORM
}

type User struct {
	ID uint `gorm:"primaryKey"`
	Email string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}

type authRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

//хэширование пароля
func GeneratePassword(p string) string { //возвращает хеш пароля
	hash, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(hash)
}
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

//генерация токена
func GenerateToken(id uint) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ //создаем токен
		"user_id": id,
	})

	secret_key := os.Getenv("JWT_SECRET")

	t, err := token.SignedString([]byte(secret_key)) //подписываем токен секретным ключом
		if err != nil {
			return "", err
	}

	return t, nil
	}

func VerifyToken(tokenString string) (bool, error) {
	secret_key := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
	})

	if err != nil {
		return false, err
	}

 	return token.Valid, nil
}

var db *gorm.DB

// @Summary Get all transactions
// @Description Retrieve all transactions for the authenticated user
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} Transaction "List of transactions"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/transactions [get]
func GetTransaction(c *fiber.Ctx) error { //обрабатываем HTTP-метод GET.
		// c *fiber.Ctx - указатель на контекст запроса
		// error - тип, возвращаемый функцией
		// определяем функцию, которая вызывается когда поступает запрос
		var transactions []Transaction //создаем срез для хранения списка транзакций из бд
		userID := c.Locals("user_id").(uint)

		db.Where("user_id = ?", userID).Find(&transactions)
		return c.JSON(transactions)
	}

// @Summary Create a new transaction
// @Summary Create a new transaction
// @Description Create a new transaction for the authenticated user
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param transaction body Transaction true "Transaction data"
// @Success 201 {object} Transaction "Created transaction"
// @Failure 400 {object} map[string]string "Invalid request body or parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/transactions [post]
func PostTransactions(c *fiber.Ctx) error {
		transaction := new(Transaction) //возвращаем указатель на пустую структуру
		err := c.BodyParser(transaction) //записывает данные из запроса в структуру
		if err != nil {       
			return c.Status(400).JSON(fiber.Map{"error": err.Error()}) //Map возвращает json
		}

		if transaction.Type == "" || transaction.Amount == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Type and Amount are required"})
		}

		if transaction.Type != "income" && transaction.Type != "expense" {
			return c.Status(400).JSON(fiber.Map{"error": "Type must be 'income' or 'expense'"})
		}

		// Если Date не передан в запросе, установим текущее время
		if transaction.Date.IsZero() {
			transaction.Date = time.Now()
		}

		transaction.UserID = c.Locals("user_id").(uint) // привязываем к пользователю

		db.Create(transaction)
		return c.Status(201).JSON(transaction)
	}

// @Summary Update a transaction
// @Description Fully update a transaction by ID for the authenticated user
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Transaction ID"
// @Param transaction body Transaction true "Full transaction data"
// @Success 200 {object} Transaction "Updated transaction"
// @Failure 400 {object} map[string]string "Invalid request body or parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]string "Transaction not found"
// @Router /api/transactions/{id} [put]
func PutTransaction(c *fiber.Ctx) error {
	id := c.Params("id") //получаем id из URL
	userID := c.Locals("user_id").(uint)

	var transaction Transaction
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}

	updated := new(Transaction)
	if err := c.BodyParser(updated); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	
	// Обновляем все поля (кроме ID и CreatedAt)
	transaction.Amount = updated.Amount
	transaction.Type = updated.Type
	transaction.Category = updated.Category
	transaction.Description = updated.Description
	transaction.Date = updated.Date

	db.Save(&transaction)

	return c.JSON(transaction)
}

// @Summary Delete a transaction
// @Description Delete a transaction by ID for the authenticated user
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} map[string]string "Success response"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]string "Transaction not found"
// @Router /api/transactions/{id} [delete]
func DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uint)

	var transaction Transaction
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}

	db.Delete(&transaction)

	return c.JSON(fiber.Map{"message": "Transaction deleted successfully"})
}

// @Summary Get user balance
// @Description Calculate and return the balance for the authenticated user
// @Tags transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]float64 "Balance response"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/balance [get]
func GetBalance(c *fiber.Ctx) error {
	var totalIncome float64
	var totalExpense float64

	userID := c.Locals("user_id").(uint)

	db.Model(&Transaction{}).
		Where("type = ? AND user_id = ?", "income", userID).
		Select("COALESCE(SUM(amount),0)"). //COALESCE заменит NULL на 0 если подходящие записи не найдены
		Scan(&totalIncome) //запишет результат запроса в переменную

	db.Model(&Transaction{}).
		Where("type = ? AND user_id = ?", "expense", userID).
		Select("COALESCE(SUM(amount),0)").
		Scan(&totalExpense)

	balance := totalIncome - totalExpense

	return c.Status(200).JSON(fiber.Map{"balance": balance})
}

// @Summary Register a new user
// @Description Create a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body authRequest true "User credentials"
// @Success 201 {object} map[string]string "Success response"
// @Failure 400 {object} map[string]string "Invalid request body or email already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
    var req authRequest //структура для десереализации JSON
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": err.Error(),
        })
    }
    user := User{
        Email:        req.Email,
        PasswordHash: GeneratePassword(req.Password),
    }
    res := db.Create(&user)
    if res.Error != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": res.Error.Error(),
        })
    }
    return c.Status(201).JSON(fiber.Map{
        "message": "user created",
    })
}

// @Summary Login a user
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body authRequest true "User credentials"
// @Success 200 {object} map[string]interface{} "Token response"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Invalid email or password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
    var req authRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": err.Error(),
        })
    }
    var user User
    res := db.Where("email = ?", req.Email).First(&user)
    if res.Error != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "user not found",
        })
    }
	if !ComparePassword(user.PasswordHash, req.Password) {
   		return c.Status(400).JSON(fiber.Map{
        	"message": "incorrect password",
    	})
	}
    token, err := GenerateToken(user.ID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": err.Error(),
        })
    }
    return c.JSON(fiber.Map{
        "token": token,
    })
}

func JWTProtected(c *fiber.Ctx) error { //middleware проверяющее JWT-токен
    return jwtware.New(jwtware.Config{ //возвращает функцию
        SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
        ContextKey: "jwt",
        ErrorHandler: func(c *fiber.Ctx, err error) error { //обработчик ошибок если токен отсутствует
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": true,
                "msg":   err.Error(),
            })
        },
    })(c) //созданный мидлвар вызывается текущим контекстом
}

func ExtractUserIDMiddleware(c *fiber.Ctx) error {
    token := c.Locals("jwt").(*jwt.Token) // Изменено с "user" на "jwt"
    claims := token.Claims.(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))
    c.Locals("user_id", userID)
    return c.Next()
}

// @title Finance Tracker API
// @description API for tracking personal finance transactions
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	db_user := os.Getenv("DB_USER")
    db_name := os.Getenv("DB_NAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_port := os.Getenv("DB_PORT")
	dsn := "host=localhost user=" + db_user + " password=" + db_password + " dbname=" + db_name + " port=" + db_port + " sslmode=disable TimeZone=Europe/Moscow" //data source name

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err) //выводит сообщение и завершает программу
	}

	db.AutoMigrate(&Transaction{}, &User{}) //передаем указатель на созданный пустой экземпляр структуры

	app := fiber.New() //экземпляр fiber

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api") // защищённые маршруты

	api.Use(JWTProtected) 
	api.Use(ExtractUserIDMiddleware)

	api.Get("/transactions", GetTransaction)
	api.Post("/transactions", PostTransactions)
	api.Put("/transactions/:id", PutTransaction)
	api.Delete("/transactions/:id", DeleteTransaction)
	api.Get("/balance", GetBalance)

	// Auth
	auth := app.Group("/auth")
	auth.Post("/login", Login)
	auth.Post("/register", Register)

	app.Listen(":3000")
}
