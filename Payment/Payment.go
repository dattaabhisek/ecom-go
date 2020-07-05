package Payment

 import (
         "fmt"
		 "net/http"
         "github.com/asaskevich/govalidator"
		// "html/template"
		// "github.com/go-sql-driver/mysql"
		"time"
		"test-go/DbConnect"
		"test-go/FileWrite"
		
		"html/template"
		"strings"
		"strconv"
		"os"
		"github.com/tkanos/gonfig"
		  
 )
type Configuration struct {
	coupon_timeout int64
}
 
 func Card_validation(w http.ResponseWriter, r *http.Request) {
 //where
 
	configuration := Configuration{}
	errConf := gonfig.GetConf("Config/goecom.json", &configuration)
	if errConf != nil {
		panic(errConf)
	}

	if os.Getenv("SES_VAL_AUTH") != "true" {
        http.Error(w, "Access denied. Please login and try again!", http.StatusForbidden)
        return
    }
	
		data := struct {
			
			CouponCodeStatus string
			CouponCode string
		}{
			CouponCodeStatus: "",
			
			
		}
		
		db := DbConnect.DbConnect()
		//count:=0
		
		rows, err := db.Query("select if(locate('~',Coupon)=0,'',substring_index(Coupon, '~', -1)) from cart where Name='Orange' and Cart_status='OPEN' and Username=? limit 1",os.Getenv("GOECOM_USER"))
		checkErr(err)
		for rows.Next() {
			err=rows.Scan(&data.CouponCodeStatus)
			data.CouponCode= "Not Applied"
			checkErr(err)
		}
		
		if data.CouponCodeStatus!="" {
			t := time.Now()
			fmt.Println(data.CouponCodeStatus)
			coupon_time, errtime := time.Parse(time.RFC1123, data.CouponCodeStatus)

			if errtime != nil {
				fmt.Println(errtime)
			}
			fmt.Println("test")
			fmt.Println(coupon_time)
			elapsed,_ := strconv.ParseInt(strings.Split(t.Sub(coupon_time).String(),".")[0], 10, 64)
			//checkErr(err_1)
			fmt.Println("elapsed")
			fmt.Println(elapsed)
			
			if elapsed > configuration.coupon_timeout {
				data.CouponCodeStatus="Coupon code expired(expiry time=10s)"
			} else {
				data.CouponCodeStatus="Coupon code validated and applied successfully"
			}
			
		} else {
			data.CouponCodeStatus="Not Applied"
		}
		
        //ccNumber := "5176865765334720"
		ccNumber := r.FormValue("cardnumber")
         validCreditCard := govalidator.IsCreditCard(ccNumber)
		if validCreditCard== true {
			tmpl := template.Must(template.ParseFiles("payment_suuccessful.html"))
			
		
		
		
			tmpl.ExecuteTemplate(w,"PaymentSuccessful", data)
			//fmt.Println("cc="+data.CouponCode)
			FileWrite.WriteToFile()
			
			closeCart()
			
			//http.ServeFile(w, r, "payment_suuccessful.html")
		} else {
			http.ServeFile(w, r, "payment_failed.html")
		}

         fmt.Printf("%s is a valid credit card : %v \n", ccNumber, validCreditCard)
		 
		defer db.Close()
 }
 
 func closeCart() {
	db := DbConnect.DbConnect()
		//count:=0
		
		rows, err := db.Prepare("update  cart set Cart_status='CLOSE' where Cart_status='OPEN' and Username=?")
		checkErr(err)
		rows.Exec(os.Getenv("GOECOM_USER"))
		defer db.Close()
 }
 
 
 func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
