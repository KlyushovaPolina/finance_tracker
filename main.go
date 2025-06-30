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
)

type Transaction struct { //модель
	ID          uint	`gorm:"primaryKey"`
	Amount      float64  `gorm:"not null"`
	Type        string   `gorm:"not null; check:type_check,type IN ('income','expense')"` //тип транзакции - трата или расход
	Category    string  //категория - еда, одежда и тд
	Description *string //может быть пустым
	Date        time.Time 
	CreatedAt   time.Time //автоматически создается GORM
}

var db *gorm.DB

// @Summary Get all transactions
// @Accept json
// @Produce json
// @Success 200 {array} Transaction
// @Failure 500 {object} map[string]string "Error response"
// @Router /transactions [get]
func GetTransaction(c *fiber.Ctx) error { //обрабатываем HTTP-метод GET.
		// c *fiber.Ctx - указатель на контекст запроса
		// error - тип, возвращаемый функцией
		// определяем функцию, которая вызывается когда поступает запрос
		var transactions []Transaction //создаем срез для хранения списка транзакций из бд
		db.Find(&transactions)
		return c.JSON(transactions)
	}

// @Summary Create a new transaction
// @Accept json
// @Produce json
// @Param transaction body Transaction true "Transaction data"
// @Success 201 {object} Transaction
// @Failure 400 {object} map[string]string "Error response"
// @Failure 500 {object} map[string]string "Error response"
// @Router /transactions [post]
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

		db.Create(transaction)
		return c.Status(201).JSON(transaction)
	}

// PATCH /transactions/:id
// @Summary Update a transaction
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param transaction body Transaction true "Transaction data"
// @Success 200 {object} Transaction
// @Failure 400 {object} map[string]string "Error response"
// @Failure 404 {object} map[string]string "Error response"
// @Router /transactions/{id} [patch]
func UpdateTransaction(c *fiber.Ctx) error {
	id := c.Params("id") //получаем id из URL

	var transaction Transaction
	if err := db.First(&transaction, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}

	updateData := new(Transaction)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Проверим что Type корректный, если он передан
	if updateData.Type != "" && updateData.Type != "income" && updateData.Type != "expense" {
		return c.Status(400).JSON(fiber.Map{"error": "Type must be 'income' or 'expense'"})
	}

	// Обновляем только переданные поля
	db.Model(&transaction).Updates(updateData)

	return c.JSON(transaction)
}

// DELETE /transactions/:id
// @Summary Delete a transaction
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} map[string]string "Success response"
// @Failure 404 {object} map[string]string "Error response"
// @Router /transactions/{id} [delete]
func DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction Transaction
	if err := db.First(&transaction, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}

	db.Delete(&transaction)

	return c.JSON(fiber.Map{"message": "Transaction deleted successfully"})
}

// @Summary Get balance
// @Description Calculate and return the balance
// @Accept json
// @Produce json
// @Success 200 {object} map[string]float64 "Balance response"
// @Failure 500 {object} map[string]string "Error response"
// @Router /balance [get]
func GetBalance(c *fiber.Ctx) error {
	var totalIncome float64
	var totalExpense float64

	db.Model(&Transaction{}).
		Where("type = ?", "income").
		Select("COALESCE(SUM(amount),0)"). //COALESCE заменит NULL на 0 если подходящие записи не найдены
		Scan(&totalIncome) //запишет результат запроса в переменную

	db.Model(&Transaction{}).
		Where("type = ?", "expense").
		Select("COALESCE(SUM(amount),0)").
		Scan(&totalExpense)

	balance := totalIncome - totalExpense

	return c.Status(200).JSON(fiber.Map{"balance": balance})
}

// @title Finance tracker API
// @description API for tracking personal finance transactions
// @host localhost:3000
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	db_user := os.Getenv("DB_USER")
    db_name := os.Getenv("DB_NAME")
	db_password := os.Getenv("DB_PASSWORD")
	dsn := "host=localhost user=" + db_user + " password=" + db_password + " dbname=" + db_name + " port=5432 sslmode=disable TimeZone=Europe/Moscow" //data source name

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err) //выводит сообщение и завершает программу
	}

	db.AutoMigrate(&Transaction{}) //передаем указатель на созданный пустой экземпляр структуры

	app := fiber.New() //экземпляр fiber

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/transactions", GetTransaction)
	app.Post("/transactions", PostTransactions)
	app.Patch("/transactions/:id", UpdateTransaction)
	app.Delete("/transactions/:id", DeleteTransaction)
	app.Get("/balance", GetBalance)

	app.Listen(":3000")
}
