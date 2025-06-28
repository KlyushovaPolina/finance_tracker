package main

import(
	//"github.com/gin-gonic/gin" //фреймворк для API
    "gorm.io/gorm" //для работы с бд (ORM)
	"gorm.io/driver/postgres" //драйвер для постгри
	"time"
	"log"
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
}