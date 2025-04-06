# Gator

Gator is a command-line tool that helps developers fetch and manage RSS feeds from websites. It stores the feeds in a PostgreSQL database, allowing you to monitor status, synchronize, while also providing login features to enable different users register and perform common operations on Feeds.

## Features

- Integrated with PostgreSQL for persistent storage
- Simple CLI interface for efficient workflow
- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post


## Prerequisites

Before running this program, make sure you have the following installed:

Go (version 1.23.6 or later) – [Install Go](https://go.dev/doc/install)
PostgreSQL (version 15.4 or later) – [Install PostgreSQL](https://www.postgresql.org/download/)

## Installation

To install the gator CLI, ensure you have Go installed on your system, then run:

go install github.com/hamstimusprime/gator@latest

This will compile the program and install the binary to your $GOPATH/bin directory. Make sure this directory is in your PATH to access the gator command from anywhere.


## Configuration

To use Gator, you'll need to set up a configuration file named `config.json` in your home directory:

1. Create a file at `~/.gator/config.json` (Linux/Mac) or `C:\Users\YourUsername\.gator\config.json` (Windows)
2. Add the following configuration, updating with your database details:
```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "your_password",
    "dbname": "gator"
  }
}
```
Before running Gator, make sure your PostgreSQL database server is up and running. You can typically start or check the status of the database with the following commands (depending on your operating system): 

Linux:
- sudo systemctl start postgresql
- sudo systemctl status postgresql

MacOs(using Homebrew):
- brew services start postgresql
- brew services list

Windows:
Ensure PostgreSQL is running using the Services app or check the Task Manager for the postgres.exe processes


## Commands
- gator login username - Logs in a user 
- gator register username - Registers a name as a user
- gator reset - Resets configuration, removing all users and their entries
- gator users - Lists all registered users
- gator agg - Aggregates and saves all feeds added by the logged in user
- gator addfeed name [feed-url] - Stores feed with a given name to records of logged in user
