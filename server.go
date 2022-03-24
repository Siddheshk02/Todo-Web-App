package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

const (
	DB_USER     = "_USER_"
	DB_PASSWORD = "_PASSWORD_"
	DB_NAME     = "_NAME_"
)

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var res string //used for storing rows while scaning

	var todos []string //for storing all the rows

	rows, err := db.Query("SELECT * FROM todos") // db.Query is used when we expect a result from the database query
                                                     // db.Exec is used when no result id expected from  the database query

	defer rows.Close() //closing the rows to prevent further enumeration when the function completes

	if err != nil {
		log.Fatalln(err)

		c.JSON("An Error occured")
	}
	for rows.Next() {
		rows.Scan(&res)
		todos = append(todos, res)
	}

	return c.Render("index", fiber.Map{
		"Todos": todos,
	})
}

type todo struct {
	Item string
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	newTodo := todo{}

	if err := c.BodyParser(&newTodo); err != nil {
		log.Printf("An Error occured: %v", err)
		return c.SendString(err.Error())
	}

	fmt.Printf("%v", newTodo)

	if newTodo.Item != "" {
		_, err := db.Exec("Insert into todos VALUES ($1)", newTodo.Item)
		if err != nil {
			log.Fatalf("An Error occured while executing query: %v", err)
		}
	}
	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	olditem := c.Query("olditem")
	newitem := c.Query("newitem")
	db.Exec("UPDATE todos SET item=$1 WHERE item=$2", newitem, olditem)
	return c.Redirect("/")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	todoToDelete := c.Query("item")
	db.Exec("DELETE from todos WHERE item=$1", todoToDelete)
	return c.SendString("Successfully Deleted :)")
}

func main() {
	//app := fiber.New()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	//connStr := "postgresql://postgres:Sid@2002@192.168.1.206/todolist?sslmode=disable"

	//Connection to the DB
	//db, err := sql.Open("postgres", "postgres:Sid@2002@tcp(127.0.0.1:8080)/todolist")
	//db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Fatal(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})

	app.Static("/", "./public")

	log.Fatal(app.Listen(":8080"))
}
