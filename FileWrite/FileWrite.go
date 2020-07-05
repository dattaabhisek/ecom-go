package FileWrite

import (
	"encoding/json"
	"io/ioutil"
	"test-go/DbConnect"
	"test-go/awss3"
	"time"
	"os"
	"strconv"
)

type Product struct {
    Id				int64	`json:"id"`
	Name			string	`json:"name"`
	Price			float64	`json:"price"`
	Quantity		int64	`json:"quantity"`
	Description		string	`json:"description"`
	Photo			string	`json:"photo"`
	Discount 		float64	`json:"discount"`
	Coupon 			string	`json:"coupon"`
	Username 		string	`json:"username"`
	Cart_status 	string 	`json:"cart_status"`
	Subtotal		float64	`json:"subtotal"`
	Total			float64	`json:"total"`
	Time			string	`json:"time"`

}

func WriteToFile() {
	var Products []Product
	var Product Product
	
	var FilePrefix string
	t := time.Now().UTC().UnixNano() / 1000000
	FilePrefix="order_"+strconv.FormatInt(t,10)
	
	db := DbConnect.DbConnect()
	rows, err := db.Query("SELECT * FROM cart where Cart_status='OPEN' and Username=?",os.Getenv("GOECOM_USER"))
	CheckErr(err)
	//var products cartItem
	for rows.Next() {
		err:=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo,&Product.Discount,&Product.Coupon,&Product.Username,&Product.Cart_status,&Product.Subtotal,&Product.Total)		
		CheckErr(err)
		Product.Time=strconv.FormatInt(t,10)
		//pearsCount=pearsCount+products.Quantity
		Products=append(Products, Product)
	}
	
		
	file, _ := json.MarshalIndent(Products, "", " ")
	_ = ioutil.WriteFile("DataDump/"+FilePrefix+".json", file, 0644)
	
	awss3.WriteToS3("DataDump/"+FilePrefix+".json","goecom")

}

func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}