# blog_aggregator

## Prerequisites

You'll need [Postgres](https://www.postgresql.org/download/) and [Go](https://go.dev/doc/install) installed. Use the links to download and install them. Once installed, you can confirm:

- Postgres is running by typing `psql --version` in your terminal.
- Go is installed by typing `go version`.

## Installation

Install the `gator` CLI globally with this command:

```bash
go install github.com/realquiller/blog_aggregator@latest

This places the CLI in your $GOPATH/bin directory (or $HOME/go/bin by default).
Ensure this directory is in your terminal's PATH to use gator from anywhere.


## Commands
You can run the program with the following commands:

gator register <name>
Register a new user by specifying their name and immediately logging in as them.

gator login <name>
Log into the program as the specified user.

gator agg <time-format>
Start the aggregator (e.g., 5s, 1m, 1h).

gator feeds
Show all feeds currently in the database.

gator addfeed <name> <url>
Add a new feed and follow it (e.g., Boot.dev Blog and https://blog.boot.dev/index.xml).

gator follow <name>
Follow an existing feed by its name.

gator following
List all feeds you are currently following.

gator unfollow <url>
Unfollow a feed by its URL.

gator browse <optional_number_limit>
Browse posts from feeds you're following. You can specify an optional number limit.
(Defaults to 2 if not provided)

## Examples

gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
gator browse
gator browse 5