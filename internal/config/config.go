package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/realquiller/blog_aggregator/internal/database"
	"github.com/realquiller/blog_aggregator/internal/feeds"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	Db     *database.Queries
	Config *Config
}

type Command struct {
	Name string
	Args []string
}

func HandlerLogin(s *State, cmd Command) error {
	// Check if args are empty
	if len(cmd.Args) < 1 {
		return fmt.Errorf("missing username")
	}

	// Check if the user already exists
	user, err := s.Db.GetUser(context.Background(), cmd.Args[0])

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Any error other than "no rows" is unexpected
			return fmt.Errorf("user %s doesn't exist", cmd.Args[0])
		}
		// Handle any other error
		return fmt.Errorf("error checking user existence: %w", err)

	}

	// Set the current user to the provided username
	if err := s.Config.SetUser(user.Name); err != nil {
		return fmt.Errorf("failed to set current user: %w", err)
	}

	// Print the changed username
	fmt.Printf("Username set to: %s\n", user.Name)

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	// Check if username is provided in command arguments
	if len(cmd.Args) < 1 {
		return fmt.Errorf("missing username")
	}

	// Check if the user already exists
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])

	if err == nil {
		// if no error, the user exists
		return fmt.Errorf("user already exists")
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Any error other than "no rows" is unexpected
		return fmt.Errorf("error checking user existence: %w", err)
	}

	// Create a new user in the database
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	// Set the current user to the provided username
	if err := s.Config.SetUser(user.Name); err != nil {
		return fmt.Errorf("error setting current user: %w", err)
	}

	// Print the changed username
	fmt.Printf("User %s was created.\n", user.Name)

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.Reset(context.Background())

	if err != nil {
		return fmt.Errorf("error resetting database: %w", err)
	}

	fmt.Println("Database reset successful!")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	items, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving users: %w", err)
	}
	current := ""

	for _, item := range items {
		if item.Name == s.Config.CurrentUserName {
			current = " (current)"
		}
		fmt.Printf("* %s%s\n", item.Name, current)
		current = ""
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("expected one argument: time_between_reqs (e.g. '5s', '1m')")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	if err := scrapeFeeds(s); err != nil {
		fmt.Printf("initial scrape error: %v\n", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("scrape error: %v\n", err)
		}
	}

	// Needed for compilation
	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	// check the length of the command arguments
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	now := time.Now()
	user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user ID: %w", err)
	}

	feed, err := s.Db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("error adding feed: %w", err)
	}

	// Print the added feed
	fmt.Println("Feed added:\n", feed)

	// Create a new feed follow in the database
	follow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	fmt.Printf("%s now follows %s\n", user.Name, follow.FeedName)

	// Yaaaay, success!
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	// Fetch the feeds from the database
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving feeds: %w", err)
	}

	// Print the feeds
	for _, feed := range feeds {
		fmt.Printf("Feed name: %s\n", feed.Name)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("Added by: %s\n\n", feed.UserName)
	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	// Check if the command arguments are provided
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: follow <url>")
	}
	// Get the feed ID from the database
	feed, err := s.Db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}

	// Create a new feed follow in the database
	follow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	fmt.Printf("%s now follows %s\n", user.Name, follow.FeedName)

	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	// Fetch the feeds that the user is following
	follows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error retrieving feed follows: %w", err)
	}

	// Print the feeds that the user is following
	for _, follow := range follows {
		fmt.Printf("%s follows: %s\n", user.Name, follow.FeedName)
	}

	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	// Check if the command arguments are provided
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}

	// Get the feed ID from the database
	feed, err := s.Db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}

	// Delete the feed follow from the database
	err = s.Db.DeleteFollow(context.Background(), database.DeleteFollowParams{
		Name: user.Name,
		Url:  feed.Url,
	})
	if err != nil {
		return fmt.Errorf("error deleting feed follow: %w", err)
	}

	fmt.Printf("%s no longer follows %s\n", user.Name, feed.Name)

	return nil
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	// 1. Gets the user
	// 2. Calls the original handler with the user
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error getting user: %w", err)
		}
		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *State) error {
	// Get the feed from database
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %w", err)
	}

	// Mark the feed as fetched
	err = s.Db.MarkFeedFetched(context.Background(), feed.ID)

	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %w", err)
	}

	// Fetch the feed
	fetchedFeed, err := feeds.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	// Print the fetched feeds
	for _, item := range fetchedFeed.Channel.Item {
		if item.Title == "" {
			fmt.Println("Skipping item with empty title")
			continue
		}

		fmt.Printf("Feed title: %s\n", item.Title)
	}

	return nil
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(*State, Command) error)
	}
	c.Handlers[name] = f
}
func (c *Commands) Run(s *State, cmd Command) error {
	if handler, ok := c.Handlers[cmd.Name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("command not found: %s", cmd.Name)
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(homePath, configFileName)
	return filePath, nil
}

func Read() (Config, error) {
	// Read the file
	filePath, err := getConfigFilePath()

	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(filePath)

	if err != nil {
		return Config{}, err
	}

	// Parse JSON into struct
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func write(cfg Config) error {
	// get the file path
	filePath, err := getConfigFilePath()

	if err != nil {
		return err
	}

	// convert struct to JSON

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)
}
