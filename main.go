package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres" //драйвер для постгри
	"gorm.io/gorm"            //для работы с бд (ORM)
)

type Transaction struct { //модель
	ID          uint	`gorm:"primaryKey"`
	Amount      float64  `gorm:"not null"`
	Type        string   `gorm:"not null"` //тип транзакции - трата или расход
	Category    string  //категория - еда, одежда и тд
	Description *string //может быть пустым
	Date        time.Time 
	CreatedAt   time.Time //автоматически создается GORM
}

var db *gorm.DB


func GetTransaction(c *fiber.Ctx) error { //обрабатываем HTTP-метод GET.
		// c *fiber.Ctx - указатель на контекст запроса
		// error - тип, возвращаемый функцией
		// определяем функцию, которая вызывается когда поступает запрос
		var transactions []Transaction //создаем срез для хранения списка транзакций из бд
		db.Find(&transactions)
		return c.JSON(transactions)
	}


func PostTransactions(c *fiber.Ctx) error {
		transaction := new(Transaction) //возвращаем указатель на пустую структуру
		c.BodyParser(transaction)       //записывает данные из запроса в структуру

		// Если Date не передан в запросе, установим текущее время
		if transaction.Date.IsZero() {
			transaction.Date = time.Now()
		}

		db.Create(transaction)
		return c.Status(201).JSON(transaction)
	}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=finance-tracker port=5432 sslmode=disable TimeZone=Europe/Moscow" //data source name
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err) //выводит сообщение и завершает программу
	}

	db.AutoMigrate(&Transaction{}) //передаем указатель на созданный пустой экземпляр структуры

	app := fiber.New() //экземпляр fiber

	app.Get("/transactions", GetTransaction)
	app.Post("/transactions", PostTransactions)

	app.Listen(":3000")
}
