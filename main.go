package main

import (
  "fmt"
  "html/template"
  "net/http"
  

  "log"

  "github.com/gorilla/pat"
  "github.com/markbates/goth"
  "github.com/markbates/goth/gothic"
  "github.com/markbates/goth/providers/google"
  "github.com/gorilla/sessions"
)


func main() {
  
  key := "Secret-session-key" // Replace this with a secure, randomly generated key for production
  maxAge := 86400 * 30        // 30 days
  isProd := false             // Set to true if you’re using HTTPS in production
  
  store := sessions.NewCookieStore([]byte(key))
  store.MaxAge(maxAge)
  store.Options.Path = "/"
  store.Options.HttpOnly = true   // HttpOnly should be enabled for security
  store.Options.Secure = isProd   // Ensure this is true in production
  gothic.Store = store

  goth.UseProviders(
    google.New("", "", "http://127.0.0.1:5500/auth/google/callback", "email", "profile"),
  )

  p := pat.New()
  p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
    user, err := gothic.CompleteUserAuth(res, req)
    if err != nil {
        log.Printf("Error completing user auth: %v", err)
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Printf("User authenticated: %+v", user)
    t, _ := template.ParseFiles("templates/success.html")
    t.Execute(res, user)
})

  p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
    gothic.BeginAuthHandler(res, req)
  })

  p.Get("/", func(res http.ResponseWriter, req *http.Request) {
    t, _ := template.ParseFiles("templates/index.html")
    t.Execute(res, false)
  })
  log.Println("listening on 127.0.0.1:5500")
  log.Fatal(http.ListenAndServe(":5500", p))
}