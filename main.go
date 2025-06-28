package main

import(
	"github.com/gin-gonic/gin" //фреймворк для API
    "gorm.io/gorm"
	"gorm.io/driver/postgres"
	"time"
)

func main(){
	type Transaction struct{
		ID uint
		Amount float64
		Type string
		Category string
		Description string
		Date time.Time
		CreatedAt time.Time
	}
}