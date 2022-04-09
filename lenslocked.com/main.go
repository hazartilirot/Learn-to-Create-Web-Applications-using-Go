package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/username/project-name/controllers"
	"github.com/username/project-name/email"
	"github.com/username/project-name/middleware"
	"github.com/username/project-name/models"
	"github.com/username/project-name/rand"
	"net/http"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Provide this flag in production. This ensures "+
		"a .config file is provided before the application starts")

	cfg := LoadConfig(*boolPtr)
	dbCfg := cfg.Database
	dbCfgInfo := dbCfg.ConnectionInfo()
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(dbCfgInfo)),
		models.WithUser(cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
	)
	must(err)
	//services.ResetDB()
	services.AutoMigrate()

	mgCfg := cfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("LensLocked Support Team", "support@sandbox4aa3d393b9fd46b0a05c53f15a863611.mailgun.org"),
		email.WithMailgun(mgCfg.Domain, mgCfg.APIKey),
	)

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User, emailer)

	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	/*middleware*/
	n, err := rand.Bytes(32)
	must(err)
	csrfMw := csrf.Protect(n, csrf.Secure(cfg.isProd()))

	UserMw := middleware.User{
		UserService: services.User,
	}

	requireUserMw := middleware.RequireUser{
		User: UserMw,
	}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.Handle("/signup", usersC.NewView).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/signin", usersC.LoginView).Methods("GET")
	r.HandleFunc("/signin", usersC.Login).Methods("POST")
	r.HandleFunc("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")
	r.Handle("/recovery", usersC.ForgotPwView).Methods("GET")
	r.HandleFunc("/recovery", usersC.InitiateReset).Methods("POST")
	r.HandleFunc("/reset", usersC.ResetPw).Methods("GET")
	r.HandleFunc("/reset", usersC.CompleteReset).Methods("POST")

	/*Assets*/
	assetsHandler := http.FileServer(http.Dir("./assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetsHandler))
	/*Image routes*/
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	/*Gallery routes*/
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name("show_gallery")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(
		"edit_gallery")
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete",
		requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")

	fmt.Printf("The server is running on :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(UserMw.Apply(r)))

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
