package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/realquiller/blog_aggregator/internal/config"
	"github.com/realquiller/blog_aggregator/internal/database"
)

func main() {
	// connect to the DB
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		os.Exit(1)
	}

	// Check if the connection is successful
	if err := db.Ping(); err != nil {
		fmt.Println("Error pinging the database:", err)
		os.Exit(1)
	}

	// If the connection is successful, print a message
	fmt.Println("Database connection successful!")

	// Create a new database query object
	dbQueries := database.New(db)

	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	state := config.State{Config: &cfg, Db: dbQueries}
	cmds := config.Commands{}
	cmds.Register("reset", config.HandlerReset)
	cmds.Register("register", config.HandlerRegister)
	cmds.Register("login", config.HandlerLogin)
	cmds.Register("agg", config.HandlerAgg)
	cmds.Register("feeds", config.HandlerFeeds)
	cmds.Register("addfeed", config.MiddlewareLoggedIn(config.HandlerAddFeed))
	cmds.Register("follow", config.MiddlewareLoggedIn(config.HandlerFollow))
	cmds.Register("following", config.MiddlewareLoggedIn(config.HandlerFollowing))
	cmds.Register("unfollow", config.MiddlewareLoggedIn(config.HandlerUnfollow))

	if len(os.Args) < 2 {
		fmt.Println("Error: not enough arguments were provided")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := []string{}
	if len(os.Args) > 2 {
		cmdArgs = os.Args[2:]
	}

	cmd := config.Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.Run(&state, cmd)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}

// DB LOGIN: psql "postgres://postgres:postgres@localhost:5432/gator"
// DB LOGIN: psql "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

// goose postgres postgres://postgres:postgres@localhost:5432/gator up

// go run . reset

// go run . register "username"
// go run . login "username"

// go run . addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
// go run . addfeed "TechCrunch" "https://techcrunch.com/feed/"
// go run . addfeed "Hacker News" "https://news.ycombinator.com/rss"

// go run . feeds
// go run . follow feed
// go run . following

// go run . agg 30s
