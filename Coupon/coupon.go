package Coupon

import (

	"fmt"
	//"github.com/captaincodeman/couponcode"
	"net/http"
	"time"
	"test-go/DbConnect"
	"html/template"
	//"tawesoft.co.uk/go/dialog"
	"os"
	
)

func GenerateCouponCode(w http.ResponseWriter, r *http.Request) {
	//var ccApplied string
	db := DbConnect.DbConnect()
	count:=0
	err := db.QueryRow("select count(Coupon) from cart where Name='Orange' and Coupon<>'-' and Username=?",os.Getenv("GOECOM_USER")).Scan(&count)
	checkErr(err)
	if count >0 {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w,"<p style='color:red;'>Coupon Code:"+"ORNG30"+" is already applied."+"</p>")
		fmt.Fprintf(w,"<br><br><a href='/product' style='font-size:20px'>Show Product</a><br> ")
		return
	}
			
	if os.Getenv("SES_VAL_AUTH") != "true" {
        http.Error(w, "Access denied. Please login and try again!", http.StatusForbidden)
        return
    }
	//code := couponcode.Generate()
	code:="ORNG30"
	//start := time.Now()
	//t := time.Now()
	//elapsed := t.Sub(start)
	fmt.Println(code)
	updateCouponToDB(code, w, r)
	//validated, err := couponcode.Validate(code)
	
	//checkErr(err)
	//fmt.Println(validated)
	//http.ServeFile(w, r, "cart.html")
	fmt.Println("couponcode====")
	//http.HandleFunc("/payment", showPaymentForm)
	
	//db = DbConnect.DbConnect()
	count=0
		
		err_1 := db.QueryRow("select count(*) from cart where Name='Orange' and Cart_status='OPEN' and Username=?",os.Getenv("GOECOM_USER")).Scan(&count)
		checkErr(err_1)
		
		if count==0 {
			//dialog.Alert("No oranges in the cart. This couponcode is only applied for oranges.")
			fmt.Fprintf(w,"No oranges in the cart. This couponcode is only valid for oranges.")
			//fmt.Sprintf("<%s>", string("No oranges in the cart. This couponcode is only applied for oranges.")) 
			return
		}
	
		rows, err1 := db.Query("select Total,if(locate('~',Coupon)=0,'',substring_index(Coupon, '~', 1)),if(locate('~',Coupon)=0,'',substring_index(Coupon, '~', -1)) from cart where Name='Orange' and Cart_status='OPEN' and Username=? limit 1",os.Getenv("GOECOM_USER"))
		//rows, err := db.Query("select Total,if(locate('~',Coupon)=0,'',substring_index(Coupon, '~', 1)),if(locate('~',Coupon)=0,'',substring_index(Coupon, '~', -1)) from cart where Name='Orange' and Username=? limit 1",os.Getenv("GOECOM_USER"))
		checkErr(err1)
		totalPayment:=0.0
		data := struct {
			Total float64
			Subtotal float64
			CouponCode string
			CouponCodeStatus string
			FinalAmt float64
		}{
			Total: totalPayment,
			
		}
		
		
		for rows.Next() {
			err=rows.Scan(&data.Total, &data.CouponCode, &data.CouponCodeStatus)
			if data.CouponCode=="" {
				data.CouponCode="Already Applied"
			}
			checkErr(err)
			
			
			//Products=append(Products, Product)
			
		}
		
		rows, err = db.Query("select subtotal from cart where Name='Orange' and Cart_status='OPEN' and Username=?",os.Getenv("GOECOM_USER"))
		checkErr(err)
		for rows.Next() {
			err=rows.Scan(&data.Subtotal)
			checkErr(err)
			
			
			//Products=append(Products, Product)
			
		}
		data.FinalAmt=data.Total-data.Subtotal*30/100
		//data.FinalAmt=float64(fmt.Printf("%0.2f", data.FinalAmt))
		
		fmt.Println("Abhisek="+data.CouponCode)
		//test := r.FormValue("cardnumber")
		//fmt.Println(test)
		tmpl := template.Must(template.ParseFiles("payment.html"))
		//var tmpl = template.Must(template.ParseFiles("product.html"))
		

        tmpl.ExecuteTemplate(w,"Payment", data)
		defer db.Close()
	
}
func updateCouponToDB(couponcode string,w http.ResponseWriter, r *http.Request) {
	db := DbConnect.DbConnect()
	count:=0
	
	fmt.Println("cc"+couponcode)
	//where
    err := db.QueryRow("SELECT count(Coupon) FROM cart where Coupon<>'-' and Name='Orange' and Cart_status='OPEN' and Username=?",os.Getenv("GOECOM_USER")).Scan(&count)
	//err := db.QueryRow("SELECT count(Coupon) FROM cart where Coupon<>'-' and Name='Orange' and Username=?",os.Getenv("GOECOM_USER")).Scan(&count)
	checkErr(err)
	fmt.Println(count)
	if count !=0 {
	
		//updateCart, errCartUpdate := db.Prepare("UPDATE cart SET Coupon=? where Name=? and Cart_status='OPEN' and Username=?",os.Getenv("GOECOM_USER"))
		//checkErr(errCartUpdate)
		//updateCart.Exec("INVALID","Orange")
		//fmt.Println("invalid")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w,"<p style='color:red;'>Coupon Code:"+couponcode+" is already applied."+"</p>")
		fmt.Fprintf(w,"<br><br><a href='/product' style='font-size:20px'>Show Product</a><br> ")
		
	} else {
		updateCart, errCartUpdate := db.Prepare("UPDATE cart SET Coupon=? where Name=? and Cart_status='OPEN' and Username=?")
		checkErr(errCartUpdate)
		updateCart.Exec(couponcode+"~"+time.Now().Format(time.RFC1123),"Orange",os.Getenv("GOECOM_USER"))
		//updateCart.Exec(couponcode)
		fmt.Println(time.Now().Format(time.RFC1123))
	}
	
	defer db.Close()
}


 func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}