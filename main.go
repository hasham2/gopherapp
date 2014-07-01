package main

import (
  "github.com/codegangsta/martini"
  "github.com/martini-contrib/render"
  "net/http"
  "database/sql"
  _ "github.com/lib/pq"
)

type User struct {
  Name string
  Email string
}

func SetupDB() *sql.DB{
  db, err := sql.Open("postgres", "dbname=gh_development sslmode=disable")
  PanicIf(err)
  return db
}

func PanicIf(err error){
  if err != nil {
    panic(err)
  }
}

func main() {
  m := martini.Classic()
  m.Map(SetupDB())
  m.Use(render.Renderer(render.Options{Layout: "layout",}))
  m.Get("/", func(ren render.Render, r *http.Request, db *sql.DB){
    searchTerm := "%" + r.URL.Query().Get("q") + "%"
    rows, err := db.Query("SELECT name, email FROM users WHERE name ILIKE $1 OR email ILIKE $1 AND email != 'admin@example.com'", searchTerm)
    PanicIf(err)
    defer rows.Close()
    users := []User{}
    for rows.Next() {
      u := User{}
      PanicIf(rows.Err())
      err := rows.Scan(&u.Name, &u.Email)
      PanicIf(err)
      users = append(users, u)
    }
    ren.HTML(200, "name", users)
  })
  m.Run()
}
