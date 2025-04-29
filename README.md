#  rssy

**rssy** is a simple command-line RSS Feed Aggregator written in Go that helps you stay on top of your favorite content feeds.

### Features

- **Feed Management**: Add, list, and manage RSS feeds
- **User System**: Register, login, and view other users
- **Social Following**: Follow/unfollow feeds discovered from other users
- **Automatic Aggregation**: Schedule regular updates of your feeds
- **Terminal Viewing**: Browse post summaries directly in your terminal
- **PostgreSQL Storage**: Reliable storage for all your feed data

### Tools Used
- Go 1.18 or higher
- PostgreSQL database
- `goose` for database migrations
- `sqlc` for type-safe SQL
- Linux

## Installation
clone the repository and  then run:
```bash
go mod download # Install dependencies
go build -o rssy # Build the binary
sudo mv rssy /usr/local/bin/ # Add to your PATH
```

#### Database setup
```bash
# Create database
psql -U postgres -c "CREATE DATABASE rssy;"

# Apply migrations
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/rssy" up

# Generate Go code from SQL
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
sqlc generate
```

#### Configuration

You have to create a configuration file at `~/.rssyconfig`:
```json
{
  "conn_str": "postgres://username:password@localhost:5432/rssy?sslmode=disable",
  "username": ""
}
```
put the connection string in the `conn_str`
and leave the `username` field empty 

## Usage

```bash
# Get help
rssy help

# Register a new user
rssy register username

# Login (set the username in the config file)
rssy login username

# List all users
rssy users

# Add a new feed
rssy addfeed "Feed Name" https://example.com/feed.xml

# List all feeds
rssy feeds

# Follow a feed from another user
rssy follow https://example.com/feed.xml

# Unfollow a feed
rssy unfollow https://example.com/feed.xml

# List feeds you're following
rssy following

rssy aggregate 1m0s  # fetch new content every 1 min and 0 sec

rssy browse # Browse recent posts (default: 2 posts)
rssy browse 10  # Show 10 recent posts
```

### License
This project is licensed under the MIT License - see the LICENSE file for details.
