package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title     string
	IsLogin   bool
	UserName  string
	FlashData string
}

var Data = MetaData{
	Title: "Personal Web",
}

type Blog struct {
	Id          int
	Title       string
	Post_date   time.Time
	Format_date string
	Author      string
	Content     string
	Image       string
	IsLogin     bool
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

var Blogs = []Blog{}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnection()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	route.HandleFunc("/", helloWorld).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET").Name("home")
	route.HandleFunc("/blog", blogs).Methods("GET")
	route.HandleFunc("/blog/{id}", blogDetail).Methods("GET")
	route.HandleFunc("/add-blog", formBlog).Methods("GET")
	route.HandleFunc("/blog", middleware.UploadFile(addBlog)).Methods("POST")
	route.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")
	route.HandleFunc("/contact", contactMe).Methods("GET")

	route.HandleFunc("/register", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/login", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

	// port := 5000
	fmt.Println("Server is running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world!"))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func blogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	rows, _ := connection.Conn.Query(context.Background(), "SELECT tb_blog.id, title, images, content, post_at, tb_user.name as author FROM tb_blog LEFT JOIN tb_user ON tb_user.id = tb_blog.author_id ORDER BY id DESC")

	var result []Blog
	for rows.Next() {
		var each = Blog{}

		var err = rows.Scan(&each.Id, &each.Title, &each.Image, &each.Content, &each.Post_date, &each.Author)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// each.Author = "Abel Dustin"
		each.Format_date = each.Post_date.Format("2 January 2006")

		if session.Values["IsLogin"] != true {
			each.IsLogin = false
		} else {
			each.IsLogin = session.Values["IsLogin"].(bool)
		}

		result = append(result, each)
	}

	respData := map[string]interface{}{
		"Data":  Data,
		"Blogs": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/blog-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// code here
	BlogDetail := Blog{}

	//kongteks nya apa cuy?
	err = connection.Conn.QueryRow(context.Background(), "SELECT tb_blog.id, title, images, content, post_at, tb_user.name as author FROM tb_blog LEFT JOIN tb_user ON tb_user.id = tb_blog.author_id WHERE tb_blog.id=$1", id).Scan(
		&BlogDetail.Id, &BlogDetail.Title, &BlogDetail.Image, &BlogDetail.Content, &BlogDetail.Post_date, &BlogDetail.Author,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// BlogDetail.Author = "Abel Dustin"
	BlogDetail.Format_date = BlogDetail.Post_date.Format("2 January 2006")

	resp := map[string]interface{}{
		"Data": Data,
		"Blog": BlogDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func formBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func addBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	author := session.Values["Id"].(int)
	fmt.Println(author)

	// code here
	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_blog(title, content, images, author_id) VALUES ($1, $2, $3, $4)", title, content, image, author)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	//code here
	// var newBlog = Blog{
	// 	Title:   title,
	// 	Author:  "Abel Dustin",
	// 	Content: content,
	// }

	// Blogs = append(Blogs, newBlog)

	// fmt.Println(Blogs)

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// code here
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_blog WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// fmt.Println(id)

	// Blogs = append(Blogs[:id], Blogs[id+1:]...)

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=ansi")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.AddFlash("Successfully register!", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// cookie = storing data
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func login(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["Id"] = user.Id
	session.Options.MaxAge = 10800 // 1 jam = 3600 detik | 3 jam = 10800

	session.AddFlash("successfully login!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logout.")
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	session.Options.MaxAge = -1 // gak boleh kurang dari 0
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
