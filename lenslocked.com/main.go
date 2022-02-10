package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/username/project-name/controllers"
	"github.com/username/project-name/middleware"
	"github.com/username/project-name/models"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "lenslocked_dev"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	services, err := models.NewServices(dsn)
	must(err)
	//services.ResetDB()
	services.AutoMigrate()

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)

	galleriesC := controllers.NewGalleries(services.Gallery, r)

	/*middleware*/
	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.Handle("/signup", usersC.NewView).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/signin", usersC.LoginView).Methods("GET")
	r.HandleFunc("/signin", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	/*Gallery routes*/
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name("show_gallery")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	http.ListenAndServe(":3000", r)

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
