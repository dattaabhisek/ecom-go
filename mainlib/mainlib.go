package mainlib

import (
    "fmt"
    "net/http"
    _"github.com/go-sql-driver/mysql"
    "os"
    "html/template"
	"strconv"
	 "test-go/DbConnect"
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
	TotalString		string	`json:"TotalString"`

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
	TotalString		string	`json:"TotalString"`
}

type ProductModel struct {
}

//To remove items from cart, one by one
func Removecart(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SES_VAL_AUTH") != "true" {
		http.Error(w, "Access denied. Please login and try again!", http.StatusForbidden)
        return
    }
	query :=r.URL.Query()
	
	fmt.Println(query.Get("id"))
	id,_:=strconv.ParseInt(query.Get("id"),10,64)
	
	fmt.Println(query.Get("Name"))

	DbCartRemove(id)
	
	table:=DbProductReload(w)
	var tmpl = template.Must(template.ParseFiles("product.html"))
	
    err := tmpl.ExecuteTemplate(w, "Index", table)
	CheckErr(err)
}

//Add items to cart, one by one
func Addcart(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SES_VAL_AUTH") != "true" {
		http.Error(w, "Access denied. Please login and try again!", http.StatusForbidden)
        return
    }
	
	query :=r.URL.Query()
	
	fmt.Println(query.Get("id"))
	id,_:=strconv.ParseInt(query.Get("id"),10,64)
	
	fmt.Println(query.Get("Name"))
	
	DbCart(id)
	table:=DbProductReload(w)
	var tmpl = template.Must(template.ParseFiles("product.html"))
	
    err := tmpl.ExecuteTemplate(w, "Index", table)
	CheckErr(err)
}
	

//View cart detail
func Viewcart(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SES_VAL_AUTH") != "true" {
        http.Error(w, "Access denied. Please login and try again!", http.StatusForbidden)
        return
    }
	var products []Product

	CartComboRule()
	products=DbCartSelect()
	
	var tmplCart = template.Must(template.ParseFiles("cart.html"))
	
    err := tmplCart.ExecuteTemplate(w, "Cart", products)
	CheckErr(err)

}

//View product detail
func ViewProduct(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SES_VAL_AUTH") != "true" {
        http.Error(w, "Access denied. Please login and try again!", http.StatusForbidden)
        return
    }
	var products []Product
	products=DbSearch(1,w)
	fmt.Println(products)
	var tmplCart = template.Must(template.ParseFiles("view_product.html"))
	
    err := tmplCart.ExecuteTemplate(w, "Index", products)
	CheckErr(err)
}


//Search items in Product and Cart table
func DbSearch(id int64, w http.ResponseWriter) []Product{
    db :=  DbConnect.DbConnect()
    rows, err := db.Query("select * from product where id=?",id)
    CheckErr(err)

	var Products []Product
	var Product Product
	
    for rows.Next() {
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo)
		Products=append(Products, Product)

    }
	
	rows, err = db.Query("select Quantity from cart where id=? and Cart_status='OPEN' and Username=?",id,os.Getenv("GOECOM_USER"))
    CheckErr(err)

	
    for rows.Next() {
		err=rows.Scan(&Product.Quantity)
		Products=append(Products, Product)
    }
	
    defer db.Close()
    return Products
}

//Reload product page after adding/removing items
func DbProductReload(w http.ResponseWriter) []Product{
    db :=  DbConnect.DbConnect()
    rows, err := db.Query("select Id,Name,Price,Description, Photo from product")
    CheckErr(err)

	var Products []Product
	var Product Product
	
    for rows.Next() {
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Description, &Product.Photo)
		CheckErr(err)
		rowsCart, errCart := db.Query("select Quantity from cart where Id=? and Cart_status='OPEN' and Username=?",Product.Id,os.Getenv("GOECOM_USER"))
		count:=0
		CheckErr(errCart)
		for rowsCart.Next() {
			
			err=rowsCart.Scan(&Product.Quantity)
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

//Apply Combo Rule for Banana and Pears
func CartComboRule ()  {
//where
	var pearsCount int64
	var bananaCount int64
	pearsCount =0
	bananaCount =0
	
	db :=  DbConnect.DbConnect()
	rowsCartPears, errCartPears := db.Query("SELECT * FROM cart where Name=? and Cart_status='OPEN' and Username=?","Pears",os.Getenv("GOECOM_USER"))
	CheckErr(errCartPears)
	var cartItemPears cartItem
	for rowsCartPears.Next() {
		err:=rowsCartPears.Scan(&cartItemPears.Id, &cartItemPears.Name, &cartItemPears.Price, &cartItemPears.Quantity, &cartItemPears.Description, &cartItemPears.Photo,&cartItemPears.Discount,&cartItemPears.Coupon,&cartItemPears.Username,&cartItemPears.Cart_status,&cartItemPears.Subtotal,&cartItemPears.Total)		
		CheckErr(err)
		pearsCount=pearsCount+cartItemPears.Quantity
	}
	
	rowsCartBanana, errCartBanana := db.Query("SELECT * FROM cart where Name=? and Cart_status='OPEN' and Username=?","Banana",os.Getenv("GOECOM_USER"))
	CheckErr(errCartBanana)
	var cartItemBanana cartItem
	for rowsCartBanana.Next() {
		err:=rowsCartBanana.Scan(&cartItemBanana.Id, &cartItemBanana.Name, &cartItemBanana.Price, &cartItemBanana.Quantity, &cartItemBanana.Description, &cartItemBanana.Photo,&cartItemBanana.Discount,&cartItemBanana.Coupon,&cartItemBanana.Username,&cartItemBanana.Cart_status,&cartItemBanana.Subtotal,&cartItemBanana.Total)		
		CheckErr(err)
		bananaCount=bananaCount+cartItemBanana.Quantity
	}

	//Check if combo offer to apply or not
	if pearsCount <4 && bananaCount <2 {
		
		defer db.Close()
		return
		
	}
	updateCart, errCartUpdate := db.Prepare("UPDATE cart SET Description=? where Name in (?,?) and Cart_status='OPEN' and Username=?")
	CheckErr(errCartUpdate)
	updateCart.Exec("DELETE","Banana","Pears",os.Getenv("GOECOM_USER"))
	
	rowsCartPears, errCartPears = db.Query("select * from cart where Description=? and Name=? and Cart_status='OPEN' and Username=? limit 1","DELETE","Pears",os.Getenv("GOECOM_USER"))
	CheckErr(errCartPears)
	
	rowsCartBanana, errCartBanana = db.Query("select * from cart where Description=? and Name=? and Cart_status='OPEN' and Username=? limit 1","DELETE","Banana",os.Getenv("GOECOM_USER"))
	CheckErr(errCartBanana)

	for rowsCartPears.Next() {
		err:=rowsCartPears.Scan(&cartItemPears.Id, &cartItemPears.Name, &cartItemPears.Price, &cartItemPears.Quantity, &cartItemPears.Description, &cartItemPears.Photo,&cartItemPears.Discount,&cartItemPears.Coupon,&cartItemPears.Username,&cartItemPears.Cart_status,&cartItemPears.Subtotal,&cartItemPears.Total)		
		CheckErr(err)
	}

	for rowsCartBanana.Next() {
		err:=rowsCartBanana.Scan(&cartItemBanana.Id, &cartItemBanana.Name, &cartItemBanana.Price, &cartItemBanana.Quantity, &cartItemBanana.Description, &cartItemBanana.Photo,&cartItemBanana.Discount,&cartItemBanana.Coupon,&cartItemBanana.Username,&cartItemBanana.Cart_status,&cartItemBanana.Subtotal,&cartItemBanana.Total)		
		CheckErr(err)
	}

	subtotal:=0.0
	
	
	
	for pearsCount >=4 && bananaCount >=2 {
		pearsCount=pearsCount-4
		bananaCount=bananaCount-2
		
		//Insert to pears
		insCart, errCart := db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
		CheckErr(errCart)
		subtotal=cartItemPears.Price * 4 * 70/100
		insCart.Exec(cartItemPears.Id, cartItemPears.Name,cartItemPears.Price,4,"Combo Pack offer applied",cartItemPears.Photo,30,"-",os.Getenv("GOECOM_USER"),"OPEN",subtotal,0)
		
		//Insert to banana
		insCart, errCart = db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?) ")
		CheckErr(errCart)
		subtotal=cartItemBanana.Price * 2 * 70/100
		insCart.Exec(cartItemBanana.Id, cartItemBanana.Name,cartItemBanana.Price,2,"Combo Pack offer applied",cartItemBanana.Photo,30,"-",os.Getenv("GOECOM_USER"),"OPEN",subtotal,0)
	
	}
	
	if pearsCount >0 {
		fmt.Println("0%pears, qty=")
		//Insert to pears
		insCart, errCart := db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
		CheckErr(errCart)
		subtotal=cartItemPears.Price * float64(pearsCount)
		insCart.Exec(cartItemPears.Id, cartItemPears.Name,cartItemPears.Price,pearsCount,"",cartItemPears.Photo,0,"-",os.Getenv("GOECOM_USER"),"OPEN",subtotal,0)

		
	}
	if bananaCount >0 {
		//Insert to banana
		insCart, errCart:= db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?) ")
		CheckErr(errCart)
		subtotal=cartItemBanana.Price * float64(bananaCount)
		insCart.Exec(cartItemBanana.Id, cartItemBanana.Name,cartItemBanana.Price,bananaCount,"",cartItemBanana.Photo,0,"-",os.Getenv("GOECOM_USER"),"OPEN",subtotal,0)

	}
	
	//Remove delete marked records
	delCart, delerrCart:= db.Prepare("DELETE FROM cart where Description=? and Username=?")
	CheckErr(delerrCart)
	delCart.Exec("DELETE",os.Getenv("GOECOM_USER"))
	
	defer db.Close()

}

//DB call to update cart
func DbCart(id int64) {
    db :=  DbConnect.DbConnect()
    rows, err := db.Query("select * from product where id=?",id)
    CheckErr(err)
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

	err_cart := db.QueryRow("SELECT COUNT(*) FROM cart where id=? and Cart_status='OPEN' and Username=?",id,os.Getenv("GOECOM_USER")).Scan(&count)
	
    CheckErr(err_cart)
	rowsCart, errCart1 := db.Query("select * from cart where id=? and Cart_status='OPEN' and Username=?",id,os.Getenv("GOECOM_USER"))
    CheckErr(errCart1)
	var cartItem cartItem
    for rowsCart.Next() {
		err=rowsCart.Scan(&cartItem.Id, &cartItem.Name, &cartItem.Price, &cartItem.Quantity, &cartItem.Description, &cartItem.Photo,&cartItem.Discount,&cartItem.Coupon,&cartItem.Username,&cartItem.Cart_status,&cartItem.Subtotal,&cartItem.Total)
    }
	
	if count==0 {
		insCart, errCart := db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
		CheckErr(errCart)
		insCart.Exec(Product.Id, Product.Name,Product.Price,Product.Quantity,Product.Description,Product.Photo,0,"-",os.Getenv("GOECOM_USER"),"OPEN",subtotal,0)
		
	} else {
		
		var qty=cartItem.Quantity

		qty++
		//Cart Rule 1#: For Apple >=7, 10% discount
		if Product.Name=="Apple" {
			fmt.Println(qty)
			if qty >=7 {
			fmt.Println(Product.Name)
			subtotal=Product.Price * float64(qty) * 90/100
			discount=10
			}
			
  		} else{
			subtotal=Product.Price * float64(qty)
		}
		updateCart, errCartUpdate := db.Prepare("UPDATE cart SET Quantity=?,Subtotal=?, Discount=?, Total=0 where id=? and Cart_status='OPEN' and Username=?")
		CheckErr(errCartUpdate)
		updateCart.Exec(qty,subtotal,discount,Product.Id,os.Getenv("GOECOM_USER"))
		
		
	}

    defer db.Close()
    
}

//Remove items from cart
func DbCartRemove(id int64) {
    db :=  DbConnect.DbConnect()
    rows, err := db.Query("select * from product where id=?",id)
    CheckErr(err)
	var Product Product
    for rows.Next() {
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo)
    }
	var count int64
	var subtotal float64
	var discount float64
	subtotal=Product.Price * float64(Product.Quantity)
	discount=0
	var cartQty int64=0

	err_cart := db.QueryRow("SELECT COUNT(*) FROM cart where id=? and Cart_status='OPEN' and Username=?",id,os.Getenv("GOECOM_USER")).Scan(&count)
	
    CheckErr(err_cart)
	rowsCart, errCart1 := db.Query("select * from cart where id=? and Cart_status='OPEN' and Username=?",id,os.Getenv("GOECOM_USER"))
    CheckErr(errCart1)
	var cartItem cartItem
    for rowsCart.Next() {
		err=rowsCart.Scan(&cartItem.Id, &cartItem.Name, &cartItem.Price, &cartItem.Quantity, &cartItem.Description, &cartItem.Photo,&cartItem.Discount,&cartItem.Coupon,&cartItem.Username,&cartItem.Cart_status,&cartItem.Subtotal,&cartItem.Total)
		cartQty=cartQty+cartItem.Quantity
		
    }

	if count!=0 {

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

			subtotal=Product.Price * float64(qty)
			
		}
		updateCart, errCartUpdate := db.Prepare("delete from cart where id=? and Cart_status='OPEN' and Username=?")
		
		CheckErr(errCartUpdate)
		updateCart.Exec(id,os.Getenv("GOECOM_USER"))
		if qty>0  {
			insCart, errCart := db.Prepare("INSERT INTO cart(Id, Name,Price,Quantity,Description,Photo,Discount,Coupon,Username,Cart_status,Subtotal,Total) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
			CheckErr(errCart)
		
			insCart.Exec(Product.Id, Product.Name,Product.Price,qty,Product.Description,Product.Photo,0,"-",os.Getenv("GOECOM_USER"),"OPEN",subtotal,0)

		
			updateCart, errCartUpdate = db.Prepare("UPDATE cart SET Quantity=?,Subtotal=?, Discount=?, Total=0 where id=? and Cart_status='OPEN' and Username=?")
			CheckErr(errCartUpdate)
			updateCart.Exec(qty,subtotal,discount,Product.Id,os.Getenv("GOECOM_USER"))
		}	
		
		
	}

    defer db.Close()
    
}

//DB call to get Product and Cart data
func DbSelect() []Product{
    db :=  DbConnect.DbConnect()
    rows, err := db.Query("select * from product")
    CheckErr(err)
			var Products []Product
			var Product Product
	
    for rows.Next() {

       
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo)
		
		rowsCart, errCart := db.Query("select Quantity from cart where Id=? and Cart_status='OPEN' and Username=?",Product.Id,os.Getenv("GOECOM_USER"))
		count:=0
		var qty int64=0
		Product.Quantity=0
		CheckErr(errCart)
		for rowsCart.Next() {

			err=rowsCart.Scan(&qty)
			Product.Quantity=Product.Quantity+qty
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


//DB call to get Cart data
func DbCartSelect() []Product{
	var total float64
	count:=0
	total=0
    db :=  DbConnect.DbConnect()
    rows, err := db.Query("select * from cart where Cart_status='OPEN' and Username=?",os.Getenv("GOECOM_USER"))
    CheckErr(err)

			var Products []Product
			var Product Product
	
    for rows.Next() {
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo,&Product.Discount,&Product.Coupon,&Product.Username,&Product.Cart_status,&Product.Subtotal,&Product.Total)
		Product.TotalString="Total:"+fmt.Sprintf("%0.2f", total)
		//if count==0 {
			total=total + Product.Subtotal
		//}
		count++
		
    }
	
	
	//fmt.Println(Product)
	
	updateCart, errCartUpdate := db.Prepare("UPDATE cart SET Total=? where Cart_status='OPEN' and Username=?")
	CheckErr(errCartUpdate)
	updateCart.Exec(total,os.Getenv("GOECOM_USER"))
	rows, err = db.Query("select * from cart where Cart_status='OPEN' and Username=?",os.Getenv("GOECOM_USER"))
	for rows.Next() {    
		err=rows.Scan(&Product.Id, &Product.Name, &Product.Price, &Product.Quantity, &Product.Description, &Product.Photo,&Product.Discount,&Product.Coupon,&Product.Username,&Product.Cart_status,&Product.Subtotal,&Product.Total)
		Products=append(Products, Product)
    }
    defer db.Close()
    return Products
}

var tmpl = template.Must(template.ParseFiles("product.html"))

//Load data to product page
func DbTableHtml(w http.ResponseWriter, r *http.Request){

    // Check if user is authenticated
   if os.Getenv("SES_VAL_AUTH") != "true" {
        http.Error(w, "Access denied. Please login and try again!", http.StatusForbidden)
        return
    }
    table := DbSelect()
    err := tmpl.ExecuteTemplate(w, "Index", table)

	
	CheckErr(err)
}

//Show payment form
func ShowPaymentForm(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SES_VAL_AUTH") != "true" {
        http.Error(w, "Access denied. Please login and try again!", http.StatusForbidden)
        return
    }
	db := DbConnect.DbConnect()
		//rows, err := db.Query("select Total,if(locate('~',Coupon)=0,'',substring_index(Coupon, '~', 1)),if(locate('~',Coupon)=0,'',substring_index(Coupon, '~', -1)),Subtotal from cart where Name='Orange' and Cart_status='OPEN' limit 1")
		rows, err := db.Query("select Total,if(locate('~',Coupon)=0,'',substring_index(Coupon, '~', -1)),Subtotal from cart where Cart_status='OPEN' and Username=? limit 1",os.Getenv("GOECOM_USER"))

		CheckErr(err)
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
		
		rows1, err1:= db.Query("select if(locate('~',Coupon)=0,'',substring_index(Coupon, '~', 1)) from cart where Name='Orange' and Cart_status='OPEN' and Username=? limit 1",os.Getenv("GOECOM_USER"))
		
		CheckErr(err1)
		count:=0
		for rows.Next() {
			for rows1.Next() {
				err=rows1.Scan(&data.CouponCode)
				count++
				CheckErr(err)
			}
			err2:=rows.Scan(&data.Total,&data.CouponCodeStatus,&data.Subtotal)
			CheckErr(err2)
			if data.CouponCode=="" {
				data.CouponCode= "Not Applied"
			}
			
			//Products=append(Products, Product)

		}
		if data.CouponCode=="Not Applied" {
			data.FinalAmt=data.Total
			
		} else {
			data.FinalAmt=data.Total-data.Subtotal*30/100
		}

		tmpl := template.Must(template.ParseFiles("payment.html"))
		//var tmpl = template.Must(template.ParseFiles("product.html"))
		

        tmpl.ExecuteTemplate(w,"Payment", data)
		defer db.Close()
}


//Check for error, generate panic

func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}