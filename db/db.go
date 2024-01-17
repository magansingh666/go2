package db

import (
	"encoding/json"
	"fmt"

	"github.com/magansingh666/go2/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBRef *gorm.DB

func InitDB() {
	var e error
	DBRef, e = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if e != nil {
		fmt.Print("Error in opening DB connection : ", e)
	}
	DBRef.AutoMigrate(&models.Product{})

	// Create dummy Products
	//createDummyProducts()
	printAllProducts()

}

func createDummyProducts() {
	DBRef.Create(&models.Product{Code: "P01", Name: "Product01"})
	DBRef.Create(&models.Product{Code: "P02", Name: "Product02"})
	DBRef.Create(&models.Product{Code: "P03", Name: "Product03"})
	DBRef.Create(&models.Product{Code: "P04", Name: "Product04"})
	DBRef.Create(&models.Product{Code: "P05", Name: "Product05"})

}

func printAllProducts() {
	p := []models.Product{}
	e := DBRef.Find(&p).Error
	if e != nil {
		fmt.Println(e)
	}
	if len(p) < 1 {
		createDummyProducts()
	}
	fmt.Println(p)
}

func CreateProuduct(b []byte) error {
	fmt.Println("crating product....")
	p := models.Product{}
	e := json.Unmarshal(b, &p)
	if e != nil {
		fmt.Print(e)
		return e
	}
	e = DBRef.Create(&p).Error
	if e != nil {
		fmt.Print(e)
		return e
	}
	return nil

}
