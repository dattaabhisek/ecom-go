package authentication

import ("database/sql"
		"fmt"
		_ "github.com/go-sql-driver/mysql"

		"golang.org/x/crypto/bcrypt"

		"net/http"

		"html/template"
		"test-go/DbConnect"
		//"test-go/mainlib"
		//"github.com/dgrijalva/jwt-go"
		//"time"
		"github.com/gorilla/sessions"
		"os"
		"time"
		"strconv"
		"test-go/awss3"
		"encoding/json"
		"io/ioutil"
		"log"
		)

var db *sql.DB
var err error


type Newuser struct {
	Newusername			string	`json:"username"`

}


var (
    // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
    key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)
)

func SignupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	fmt.Println("register")
	var user string
	db := DbConnect.DbConnect()
	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error-1, unable to create your account.", 500)
			panic(err)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error-2, unable to create your account.", 500)
			panic(err)
			return
		}
		//***
		var Newusers []Newuser
		var Newuser Newuser
	
		var FilePrefix string
		t := time.Now().UTC().UnixNano() / 1000000
		FilePrefix="newuser_"+strconv.FormatInt(t,10)
	
		
		
		Newuser.Newusername=username
		
		//pearsCount=pearsCount+products.Quantity
		Newusers=append(Newusers, Newuser)
		
	
		
		file, err1 := json.MarshalIndent(Newusers, "", " ")
		
		if err1 != nil {
			panic(err1)
		}	
		err1 = ioutil.WriteFile("DataDump/"+FilePrefix+".json", file, 0644)
		if err1 != nil {
			panic(err1)
		}	
		awss3.WriteToS3("DataDump/"+FilePrefix+".json","goecom-userverification")
		//***
		res.Header().Set("Content-Type", "text/html")
		log.Print("New User created!")
		//fmt.Fprintf(res,"<a href='/login' style='font-size:20px'>Login</a><br> ")
		res.Write([]byte("User created!"))
		res.Write([]byte("Please visit your email and confirm the email verification."))
		fmt.Fprintf(res,"<br><br><a href='/login' style='font-size:20px'>Login</a><br> ")
		
		return
	case err != nil:
		http.Error(res, "Server error-3, unable to create your account.", 500)
		panic(err)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}

func LoginPage(res http.ResponseWriter, req *http.Request) {

	
	 session, _ := store.Get(req, "cookie-name")
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string
	db := DbConnect.DbConnect()
	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		res.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(res,"<p style='color:red;'>"+"Username/Password is wrong. Please retry"+"</p>")
		http.Redirect(res, req, "/login", 301)
		//http.ServeFile(res, req, "login.html")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}
	session.Values["authenticated"] = true
	
    session.Save(req, res)
	err = os.Setenv("SES_VAL_AUTH", "true")
    if err != nil {
        fmt.Println(err)
    }
	err = os.Setenv("GOECOM_USER", username)
    if err != nil {
        fmt.Println(err)
    }
	log.Print(os.Getenv("GOECOM_USER")+" logged in")
	
	http.Redirect(res, req, "/product", 301)

}
func Logout(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Revoke users authentication
    session.Values["authenticated"] = false
    session.Save(r, w)
	err = os.Setenv("SES_VAL_AUTH", "false")
    if err != nil {
        fmt.Println(err)
    }
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w,"<p style='color:red;'>"+"USer "+os.Getenv("GOECOM_USER")+" is successfully logged out"+"</p>")
	log.Print(os.Getenv("GOECOM_USER")+" logged out")
	err = os.Setenv("GOECOM_USER", "")
    if err != nil {
        fmt.Println(err)
    }
	
	http.Redirect(w, r, "/login", 301)
		
}

func productPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "product.html")
		return
	}


	var t = template.Must(template.ParseFiles("product.html"))
	req.ParseForm()
	t.Execute(res, "qq")
	//tmp.Execute(response,data)
}

func HomePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "home.html")
}



