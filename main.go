package main

import(
    "gorm.io/gorm" //для работы с бд (ORM)
	"gorm.io/driver/postgres" //драйвер для постгри
	"time"
	"log"
	"github.com/gofiber/fiber/v2"
)

func main(){
	type Transaction struct{ //модель
		ID uint
		Amount float64
		Type string //тип транзакции - трата или расход
		Category string //категория - еда, одежда и тд
		Description *string //может быть пустым
		Date time.Time
		CreatedAt time.Time //автоматически создается GORM
	}

	dsn := "host=localhost user=postgres password=postgres dbname=finance-tracker port=5432 sslmode=disable TimeZone=Europe/Moscow" //data source name
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
    	log.Fatal("Ошибка подключения к базе данных:", err) //выводит сообщение и завершает программу
	}

	db.AutoMigrate(&Transaction{}) //передаем указатель созданный пустой экземпляр структуры


	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}