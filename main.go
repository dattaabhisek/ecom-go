package main

import (
    "fmt"
    "net/http"
    _"github.com/go-sql-driver/mysql"
    "database/sql"
    "os"
    "html/template"
	//"github.com/gorilla/sessions"
	//"tawesoft.co.uk/go/dialog"
	"strconv"
	"test-go/Payment"
	 "test-go/DbConnect"
	 "test-go/Coupon"
	 "test-go/authentication"
	 "test-go/mainlib"
	 "github.com/gorilla/sessions"
)

var (
    // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
    key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)
)

const (
    S3_REGION = "eu-central-1"
    S3_BUCKET = "goecom/order"
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

}

type cartItem struct {
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
}

type ProductModel struct {
}

//var store=sessions.NewCookieStore([]byte("mysession"))

func removecart(w http.ResponseWriter, r *http.Request) {

	//dialog.Alert(store.Get(r,"mysession"))
	//var tmpl = template.Must(template.ParseFiles("cart.html"))
	query :=r.URL.Query()
	
	//fmt(query.Get("id"))
	fmt.Println(query.Get("id"))
	id,_:=strconv.ParseInt(query.Get("id"),10,64)
	
	fmt.Println(query.Get("Name"))
	//var products []Product
	//var cartItems []cartItem
	
	
    //table := dbSearch(id)
	//fmt.Println(products[0])
	dbCartRemove(id)
	//table:=dbSearch(id,w) 
	table:=dbProductReload(w)
	//Abhisek
	//cartComboRule()
	//products=dbCartSelect()
	//query =r.URL.Query()
	
	//var tmplCart = template.Must(template.ParseFiles("product.html"))
	var tmpl = template.Must(template.ParseFiles("product.html"))
	//table := dbSelect()
	
    err := tmpl.ExecuteTemplate(w, "Index", table)
	checkErr(err)
	//http.Redirect(w,r,"/product",http.StatusSeeOther)
	//
	//viewProduct(id,w,r)
}


func addcart(w http.ResponseWriter, r *http.Request) {

	//dialog.Alert(store.Get(r,"mysession"))
	//var tmpl = template.Must(template.ParseFiles("cart.html"))
	query :=r.URL.Query()
	
	//fmt(query.Get("id"))
	fmt.Println(query.Get("id"))
	id,_:=strconv.ParseInt(query.Get("id"),10,64)
	
	fmt.Println(query.Get("Name"))
	//var products []Product
	//var cartItems []cartItem
	
	
    //table := dbSearch(id)
	//fmt.Println(products[0])
	dbCart(id)
	//table:=dbSearch(id,w) 
	table:=dbProductReload(w)
	//Abhisek
	//cartComboRule()
	//products=dbCartSelect()
	//query =r.URL.Query()
	
	//var tmplCart = template.Must(template.ParseFiles("product.html"))
	var tmpl = template.Must(template.ParseFiles("product.html"))
	//table := dbSelect()
	
    err := tmpl.ExecuteTemplate(w, "Index", table)
	checkErr(err)
	//http.Redirect(w,r,"/product",http.StatusSeeOther)
	//
	//viewProduct(id,w,r)
}
	

func viewcart(w http.ResponseWriter, r *http.Request) {
	var products []Product
	//Abhisek
	cartComboRule()
	products=dbCartSelect()
	//query =r.URL.Query()
	
	var tmplCart = template.Must(template.ParseFiles("cart.html"))
	
    err := tmplCart.ExecuteTemplate(w, "Cart", products)
	checkErr(err)
	//http.Redirect(w,r,"/cart",http.StatusSeeOther)
}
func viewProduct(w http.ResponseWriter, r *http.Request) {
	var products []Product
	//Abhisek
	//cartComboRule()
	products=dbSearch(1,w)
	//query =r.URL.Query()
	fmt.Println(products)
	var tmplCart = template.Must(template.ParseFiles("view_product.html"))
	
    err := tmplCart.ExecuteTemplate(w, "Index", products)
	checkErr(err)
	//http.ServeFile(w, r, "view_product.html")
	//http.Redirect(w,r,"/cart",http.StatusSeeOther)
}


//
func helloWorld(w http.ResponseWriter, r *http.Request){
    name, err := os.Hostname()
    checkErr(err)
    fmt.Fprintf(w, "HOSTNAME : %s\n", name)
}

func dbConnect() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "bankura-2020"
    dbHost := "localhost"
    dbPort := "3306"
    dbName := "goecom"
    db, err := sql.Open(dbDriver, dbUser +":"+ dbPass +"@tcp("+ dbHost +":"+ dbPort +")/"+ dbName +"?charset=utf8")
    checkErr(err)
    return db
}

func dbSearch(id int64, w http.ResponseWriter) []Product{
    db := dbConnect()
    rows, err := db.Query("select * from product where id=?",id)
    checkErr(err)

	var Products []Product
	var Product Product
	
    for rows.Next() {

        //err = rows.Scan(&Id, &Name, &Price, &Quantity, &Description,&Photo)
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo)
		Products=append(Products, Product)

    }
	
	rows, err = db.Query("select Quantity from cart where id=? and Cart_status='OPEN'",id)
    checkErr(err)

	
    for rows.Next() {

        //err = rows.Scan(&data.Quantity)
		err=rows.Scan(&Product.Quantity)
		Products=append(Products, Product)
		//dialog.Alert(strconv.FormatInt(Product.Quantity, 10))
		//fmt.Fprintf(w,strconv.FormatInt(Product.Quantity, 10))
    }
	
    defer db.Close()
    return Products
}
func dbProductReload(w http.ResponseWriter) []Product{
    db := dbConnect()
    rows, err := db.Query("select Id,Name,Price,Description, Photo from product")
    checkErr(err)

	var Products []Product
	var Product Product
	
    for rows.Next() {
		
        //err = rows.Scan(&Id, &Name, &Price, &Quantity, &Description,&Photo)
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Description, &Product.Photo)
		checkErr(err)
		//fmt.Println(strconv.FormatInt(Product.Id, 10))
		rowsCart, errCart := db.Query("select Quantity from cart where Id=? and Cart_status='OPEN'",Product.Id)
		count:=0
		checkErr(errCart)
		for rowsCart.Next() {
			
			err=rowsCart.Scan(&Product.Quantity)
			//Products=append(Products, Product)
			//fmt.Println(strconv.FormatInt(Product.Quantity, 10))
			//fmt.Println(strconv.FormatInt(Product.Quantity, 10))
			count++
		}
		if count==0 {
			Product.Quantity=0
		}
		Products=append(Products, Product)

    }
	

    defer db.Close()
    return Products
}
func cartComboRule ()  {
//where
	var pearsCount int64
	var bananaCount int64
	pearsCount =0
	bananaCount =0
	
	db := dbConnect()
    //err := db.Query("SELECT Quantity FROM cart where Name=? and Cart_status='OPEN'","Pears").Scan(&pearsCount)
	
	rowsCartPears, errCartPears := db.Query("SELECT * FROM cart where Name=? and Cart_status='OPEN'","Pears")
	checkErr(errCartPears)
	var cartItemPears cartItem
	for rowsCartPears.Next() {
		err:=rowsCartPears.Scan(&cartItemPears.Id, &cartItemPears.Name, &cartItemPears.Price, &cartItemPears.Quantity, &cartItemPears.Description, &cartItemPears.Photo,&cartItemPears.Discount,&cartItemPears.Coupon,&cartItemPears.Username,&cartItemPears.Cart_status,&cartItemPears.Subtotal,&cartItemPears.Total)		
		checkErr(err)
		pearsCount=pearsCount+cartItemPears.Quantity
	}
	
	rowsCartBanana, errCartBanana := db.Query("SELECT * FROM cart where Name=? and Cart_status='OPEN'","Banana")
	checkErr(errCartBanana)
	var cartItemBanana cartItem
	for rowsCartBanana.Next() {
		err:=rowsCartBanana.Scan(&cartItemBanana.Id, &cartItemBanana.Name, &cartItemBanana.Price, &cartItemBanana.Quantity, &cartItemBanana.Description, &cartItemBanana.Photo,&cartItemBanana.Discount,&cartItemBanana.Coupon,&cartItemBanana.Username,&cartItemBanana.Cart_status,&cartItemBanana.Subtotal,&cartItemBanana.Total)		
		checkErr(err)
		bananaCount=bananaCount+cartItemBanana.Quantity
	}
	
	//checkErr(err)
	//err = db.QueryRow("SELECT COUNT(*) FROM cart where Name=? and Cart_status='OPEN'","Banana").Scan(&bananaCount)
	//checkErr(err)
	fmt.Println("Abhisek:pearsCount=")
	fmt.Println(pearsCount)
	fmt.Println("bananaCount=")
	fmt.Println(bananaCount)
	//Check if combo offer to apply or not
	if pearsCount <4 && bananaCount <2 {
		
		defer db.Close()
		return
		
	}
	updateCart, errCartUpdate := db.Prepare("UPDATE cart SET Description=? where Name in (?,?) and Cart_status='OPEN'")
	checkErr(errCartUpdate)
	updateCart.Exec("DELETE","Banana","Pears")
	
	rowsCartPears, errCartPears = db.Query("select * from cart where Description=? and Name=? and Cart_status='OPEN' limit 1","DELETE","Pears")
	checkErr(errCartPears)
	
	rowsCartBanana, errCartBanana = db.Query("select * from cart where Description=? and Name=? and Cart_status='OPEN' limit 1","DELETE","Banana")
	checkErr(errCartBanana)
	//var cartItemPears cartItem
	for rowsCartPears.Next() {
		err:=rowsCartPears.Scan(&cartItemPears.Id, &cartItemPears.Name, &cartItemPears.Price, &cartItemPears.Quantity, &cartItemPears.Description, &cartItemPears.Photo,&cartItemPears.Discount,&cartItemPears.Coupon,&cartItemPears.Username,&cartItemPears.Cart_status,&cartItemPears.Subtotal,&cartItemPears.Total)		
		checkErr(err)
	}
	//var cartItemBanana cartItem
	for rowsCartBanana.Next() {
		err:=rowsCartBanana.Scan(&cartItemBanana.Id, &cartItemBanana.Name, &cartItemBanana.Price, &cartItemBanana.Quantity, &cartItemBanana.Description, &cartItemBanana.Photo,&cartItemBanana.Discount,&cartItemBanana.Coupon,&cartItemBanana.Username,&cartItemBanana.Cart_status,&cartItemBanana.Subtotal,&cartItemBanana.Total)		
		checkErr(err)
	}
	fmt.Println("pearsCount=")
	fmt.Println(pearsCount)
	fmt.Println("bananaCount=")
	fmt.Println(bananaCount)
	subtotal:=0.0
	
	
	
	for pearsCount >=4 && bananaCount >=2 {
		pearsCount=pearsCount-4
		bananaCount=bananaCount-2
		
		//Insert to pears
		insCart, errCart := db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
		checkErr(errCart)
		subtotal=cartItemPears.Price * 4 * 70/100
		insCart.Exec(cartItemPears.Id, cartItemPears.Name,cartItemPears.Price,4,"Combo Pack offer applied",cartItemPears.Photo,30,"-","-","OPEN",subtotal,0)
		
		//Insert to banana
		insCart, errCart = db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?) ")
		checkErr(errCart)
		subtotal=cartItemBanana.Price * 2 * 70/100
		insCart.Exec(cartItemBanana.Id, cartItemBanana.Name,cartItemBanana.Price,2,"Combo Pack offer applied",cartItemBanana.Photo,30,"-","-","OPEN",subtotal,0)
		
		
		fmt.Println("pears30%,qty=4")
		fmt.Println("banana30%,qty=2")
		
		fmt.Println("--End---")
	
	}
	fmt.Println(pearsCount)
	fmt.Println(bananaCount)
	if pearsCount >0 {
		fmt.Println("0%pears, qty=")
		//Insert to pears
		insCart, errCart := db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
		checkErr(errCart)
		subtotal=cartItemPears.Price * float64(pearsCount)
		insCart.Exec(cartItemPears.Id, cartItemPears.Name,cartItemPears.Price,pearsCount,"",cartItemPears.Photo,0,"-","-","OPEN",subtotal,0)

		
	}
	if bananaCount >0 {
		fmt.Println("0%banana, qty=")
		//Insert to banana
		insCart, errCart:= db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?) ")
		checkErr(errCart)
		subtotal=cartItemBanana.Price * float64(bananaCount)
		insCart.Exec(cartItemBanana.Id, cartItemBanana.Name,cartItemBanana.Price,bananaCount,"",cartItemBanana.Photo,0,"-","-","OPEN",subtotal,0)

	}
	
	//Remove delete marked records
	delCart, delerrCart:= db.Prepare("DELETE FROM cart where Description=?")
	checkErr(delerrCart)
	delCart.Exec("DELETE")
	
	defer db.Close()
    //return pearsCount,bananaCount
}

func dbCart(id int64) {
    db := dbConnect()
    rows, err := db.Query("select * from product where id=?",id)
    checkErr(err)
	//var Products []Product
	var Product Product
    for rows.Next() {

        //err = rows.Scan(&Id, &Name, &Price, &Quantity, &Description,&Photo)
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo)
		//Products=append(Products, Product)

		
    }
	var count int64
	//count=0
	var subtotal float64
	var discount float64
	//var total float64
	//total=0
	subtotal=Product.Price * float64(Product.Quantity)
	discount=0

	err_cart := db.QueryRow("SELECT COUNT(*) FROM cart where id=? and Cart_status='OPEN'",id).Scan(&count)
	
    checkErr(err_cart)
	rowsCart, errCart1 := db.Query("select * from cart where id=? and Cart_status='OPEN'",id)
    checkErr(errCart1)
	var cartItem cartItem
    for rowsCart.Next() {

        //err = rows.Scan(&Id, &Name, &Price, &Quantity, &Description,&Photo)
		err=rowsCart.Scan(&cartItem.Id, &cartItem.Name, &cartItem.Price, &cartItem.Quantity, &cartItem.Description, &cartItem.Photo,&cartItem.Discount,&cartItem.Coupon,&cartItem.Username,&cartItem.Cart_status,&cartItem.Subtotal,&cartItem.Total)
		//Products=append(Products, Product)
		//total=total + cartItem.Subtotal
		
    }
	
	//fmt.Println(total)
	if count==0 {
		fmt.Println("Insert")
		insCart, errCart := db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
		checkErr(errCart)
		insCart.Exec(Product.Id, Product.Name,Product.Price,Product.Quantity,Product.Description,Product.Photo,0,"-","-","OPEN",subtotal,0)
		
	} else {
		
		fmt.Println("Update")
		var qty=cartItem.Quantity

		qty++
		//Cart Rule 1#: For Apple >=7, 10% discount
		if Product.Name=="Apple" {
			fmt.Println(qty)
			if qty >=7 {
			fmt.Println(Product.Name)
			subtotal=Product.Price * float64(qty) * 90/100
			discount=10
			fmt.Println("Dis=")
			fmt.Println(discount)
			}
			
  		} else{
			//pearsCount,bananaCount := cartComboRule()
			//fmt.Println("pearsCount=")
			//fmt.Println(pearsCount)
			//fmt.Println(bananaCount)
			
			
			
			subtotal=Product.Price * float64(qty)
			
		}
		updateCart, errCartUpdate := db.Prepare("UPDATE cart SET Quantity=?,Subtotal=?, Discount=?, Total=0 where id=? and Cart_status='OPEN'")
		checkErr(errCartUpdate)
		updateCart.Exec(qty,subtotal,discount,Product.Id)
		
		
	}

    defer db.Close()
    
}
func dbCartRemove(id int64) {
    db := dbConnect()
    rows, err := db.Query("select * from product where id=?",id)
    checkErr(err)
	//var Products []Product
	var Product Product
    for rows.Next() {

        //err = rows.Scan(&Id, &Name, &Price, &Quantity, &Description,&Photo)
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo)
		//Products=append(Products, Product)

		
    }
	var count int64
	//count=0
	var subtotal float64
	var discount float64
	//var total float64
	//total=0
	subtotal=Product.Price * float64(Product.Quantity)
	discount=0
	var cartQty int64=0

	err_cart := db.QueryRow("SELECT COUNT(*) FROM cart where id=? and Cart_status='OPEN'",id).Scan(&count)
	
    checkErr(err_cart)
	rowsCart, errCart1 := db.Query("select * from cart where id=? and Cart_status='OPEN'",id)
    checkErr(errCart1)
	var cartItem cartItem
    for rowsCart.Next() {

        //err = rows.Scan(&Id, &Name, &Price, &Quantity, &Description,&Photo)
		err=rowsCart.Scan(&cartItem.Id, &cartItem.Name, &cartItem.Price, &cartItem.Quantity, &cartItem.Description, &cartItem.Photo,&cartItem.Discount,&cartItem.Coupon,&cartItem.Username,&cartItem.Cart_status,&cartItem.Subtotal,&cartItem.Total)
		cartQty=cartQty+cartItem.Quantity
		//Products=append(Products, Product)
		//total=total + cartItem.Subtotal
		
    }
	
	//fmt.Println(total)
	if count==0 {
		fmt.Println("nothing to delete")
		//insCart, errCart := db.Prepare("delete from cart where Id=? and Cart_status='OPEN'")
		//checkErr(errCart)
		//insCart.Exec(Product.Id, Product.Name,Product.Price,Product.Quantity,Product.Description,Product.Photo,0,"-","-","OPEN",subtotal,0)
		
	} else {
		
		fmt.Println("Update")
		var qty=cartQty
		fmt.Println("qty--")
		fmt.Println(strconv.FormatInt(qty, 10))
		qty--
		
		//Cart Rule 1#: For Apple >=7, 10% discount
		if Product.Name=="Apple" {
			fmt.Println(qty)
			if qty >=7 {
			fmt.Println(Product.Name)
			subtotal=Product.Price * float64(qty) * 90/100
			discount=10
			fmt.Println("Dis=")
			fmt.Println(discount)
			}
			
  		} else{
			//pearsCount,bananaCount := cartComboRule()
			//fmt.Println("pearsCount=")
			//fmt.Println(pearsCount)
			//fmt.Println(bananaCount)
			
			
			
			subtotal=Product.Price * float64(qty)
			
		}
		updateCart, errCartUpdate := db.Prepare("delete from cart where id=? and Cart_status='OPEN'")
		
		checkErr(errCartUpdate)
		updateCart.Exec(id)
		fmt.Println("baana==="+Product.Name)
		insCart, errCart := db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
		checkErr(errCart)
		insCart.Exec(Product.Id, Product.Name,Product.Price,qty,Product.Description,Product.Photo,0,"-","-","OPEN",subtotal,0)

		
		updateCart, errCartUpdate = db.Prepare("UPDATE cart SET Quantity=?,Subtotal=?, Discount=?, Total=0 where id=? and Cart_status='OPEN'")
		checkErr(errCartUpdate)
		updateCart.Exec(qty,subtotal,discount,Product.Id)
				
		
		
	}

    defer db.Close()
    
}


func dbSelect() []Product{
    db := dbConnect()
    rows, err := db.Query("select * from product")
    checkErr(err)

    //Product :=Product{}
    //Products := []Product{}
			var Products []Product
			var Product Product
	
    for rows.Next() {

        //err = rows.Scan(&Id, &Name, &Price, &Quantity, &Description,&Photo)
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo)
		//Product.Quantity=0
		rowsCart, errCart := db.Query("select Quantity from cart where Id=? and Cart_status='OPEN'",Product.Id)
		count:=0
		var qty int64=0
		Product.Quantity=0
		checkErr(errCart)
		for rowsCart.Next() {
			
			//err=rowsCart.Scan(&Product.Quantity)
			err=rowsCart.Scan(&qty)
			Product.Quantity=Product.Quantity+qty
			//Products=append(Products, Product)
			//fmt.Println(strconv.FormatInt(Product.Quantity, 10))
			//fmt.Println(strconv.FormatInt(Product.Quantity, 10))
			count++
		}
		if count==0 {
			Product.Quantity=0
		}
		
		Products=append(Products, Product)

    }
	
	//dialog.Alert(strconv.FormatInt(Product.Quantity, 10))
	
    defer db.Close()
    return Products
}



func dbCartSelect() []Product{
	var total float64
	total=0
    db := dbConnect()
    rows, err := db.Query("select * from cart where Cart_status='OPEN'")
    checkErr(err)

			var Products []Product
			var Product Product
	
    for rows.Next() {

        //err = rows.Scan(&Id, &Name, &Price, &Quantity, &Description,&Photo)
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo,&Product.Discount,&Product.Coupon,&Product.Username,&Product.Cart_status,&Product.Subtotal,&Product.Total)
		
	
		total=total + Product.Subtotal
			
		
	
		//Products=append(Products, Product)
		
    }
	
	s := fmt.Sprintf("%0.2f", total)
	updateCart, errCartUpdate := db.Prepare("UPDATE cart SET Total=? where Cart_status='OPEN'")
	checkErr(errCartUpdate)
	updateCart.Exec(total)
	rows, err = db.Query("select * from cart where Cart_status='OPEN'")
	for rows.Next() {

        //err = rows.Scan(&Id, &Name, &Price, &Quantity, &Description,&Photo)
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo,&Product.Discount,&Product.Coupon,&Product.Username,&Product.Cart_status,&Product.Subtotal,&Product.Total)
		
		
		Products=append(Products, Product)
		
    }
	
	fmt.Println("Total="+s)
	fmt.Println(total)
    defer db.Close()
    return Products
}

var tmpl = template.Must(template.ParseFiles("product.html"))
//var tmpl = template.Must(template.ParseGlob("layout.html"))
func dbTableHtml(w http.ResponseWriter, r *http.Request){
	 session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Fprintf(w,"Access forbidden. Please login first.")
        return
    }
	
    table := dbSelect()
    err := tmpl.ExecuteTemplate(w, "Index", table)

	
	checkErr(err)
}




func dbTable(w http.ResponseWriter, r *http.Request){
    table := dbSelect()
    for i := range(table) {
        prd := table[i]
        fmt.Fprintf(w,"%12s|%12s|%12s|%12s|%30s|%30s|\n" ,prd.Id ,prd.Name ,prd.Price ,prd.Quantity,prd.Description,prd.Photo)
    }
}

func showPaymentForm(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Changing totalPayment")
	db := DbConnect.DbConnect()
		rows, err := db.Query("select Total,if(locate('~',Coupon)=0,'',substring_index(Coupon, '~', 1)),if(locate('~',Coupon)=0,'',substring_index(Coupon, '~', -1)),Subtotal from cart where Name='Orange' and Cart_status='OPEN' limit 1")
		checkErr(err)
		totalPayment:=0.0
		data := struct {
			Total float64
			CouponCode string
			CouponCodeStatus string
			Subtotal float64
			FinalAmt float64
		}{
			Total: totalPayment,
			
		}
		
		
		for rows.Next() {
			err=rows.Scan(&data.Total, &data.CouponCode, &data.CouponCodeStatus, &data.Subtotal)
			checkErr(err)
			if data.CouponCode=="" {
				data.CouponCode= "Not Applied"
			}
			
			//Products=append(Products, Product)

		}
		data.FinalAmt=data.Total-data.Subtotal*30/100
		
		
		
		
		//fmt.Println(r.FormValue["cardnumber"])
		//test := r.FormValue("cardnumber")
		//fmt.Println(test)
		tmpl := template.Must(template.ParseFiles("payment.html"))
		//var tmpl = template.Must(template.ParseFiles("product.html"))
		

        tmpl.ExecuteTemplate(w,"Payment", data)
		defer db.Close()
}


func main() {

	fs:=http.FileServer(http.Dir("images"))
	http.Handle("/images/",http.StripPrefix("/images/",fs))
	http.HandleFunc("/signup", authentication.SignupPage)
	http.HandleFunc("/login", authentication.LoginPage)
	http.HandleFunc("/logout", authentication.Logout)
	http.HandleFunc("/", authentication.HomePage)
    //http.HandleFunc("/", mainlib.DbTableHtml)
    http.HandleFunc("/product", mainlib.DbTableHtml) 
    //http.HandleFunc("/raw", mainlib.DbTable)
	http.HandleFunc("/cart", mainlib.Addcart)
	http.HandleFunc("/removecart", mainlib.Removecart)
	http.HandleFunc("/viewcart", mainlib.Viewcart)
	http.HandleFunc("/payment_status", Payment.Card_validation)
	http.HandleFunc("/payment", mainlib.ShowPaymentForm)
	http.HandleFunc("/coupon", Coupon.GenerateCouponCode)
	http.HandleFunc("/viewProduct", mainlib.ViewProduct)
	//payment_finalyze
	
	
	
    http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}