# blog_aggregator

## Prerequisites
You'll need [Postgres](https://www.postgresql.org/download/) and [Go](https://go.dev/doc/install) installed. Use the links to download and install them. Once installed, you can confirm:
- Postgres is running by typing `psql --version` in your terminal.
- Go is installed by typing `go version`.


## Installation
Install the `gator` CLI globally with this command:

go install github.com/realquiller/blog_aggregator@latest

This places the CLI in your $GOPATH/bin directory (or $HOME/go/bin by default).
Ensure this directory is in your terminal's PATH to use gator from anywhere.


## Commands
You can run the program with the following commands:

gator register <name> // Register a new user by specifying their name and immediately login as them.
gator login <name> // Log into the program as the specified user.
gator agg <time-format> // Start the aggregator with a time format (e.g., 5s, 1m, 1h).
gator feeds // Show every feed from the database.
gator addfeed <name> <url> // Add a new feed (e.g., "Boot.dev Blog" "https://blog.boot.dev/index.xml") and immediately follow it.
gator follow <name> // Follow an existing feed by its name.
gator following //  List all feeds you are currently following.
gator unfollow <url> // Unfollow a feed by its url.
gator browse <optional_number_limit> // Browse posts from feeds you're following. You can specify an optional number limit (default is `2`).

## Examples

gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
gator browse
gator browse 5





