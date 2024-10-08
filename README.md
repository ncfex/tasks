# Tasks CLI

A simple and powerful command-line todo application written in Go. Manage your tasks efficiently with support for multiple storage backends (JSON, CSV, and SQL).

## Features

- Multiple storage backends (JSON, CSV, SQL)
- Flexible column display configuration
- Human-friendly date parsing
- Persistent configuration
- Task completion tracking
- Customizable task views

## Installation

```bash
go install github.com/ncfex/tasks@latest
```

## Configuration

The application stores its data in the `~/.tasks` directory. For SQL storage mode, you need to set the `DB_URL` environment variable.

Default configuration will be created automatically on first run.

## Usage

### Basic Commands

```bash
tasks [command] [flags]
```

### Global Flags

- `-m, --format string`: Storage format (json, csv, or sql) (json is default configured mode)
- `-h, --help`: Help for any command

### Available Commands

### Add a Task

```bash
tasks add [description] [flags]

Flags:
  -d, --due string   Due date for the task (defaults to "tomorrow")
```

The `--due` flag supports human-readable time formats:

1. Special keywords:

   - `tomorrow` - Sets due date to tomorrow

2. Relative past time (for recurring tasks):

   ```bash
   tasks add "Weekly review" --due "1 week ago"
   ```

   Supported formats:

   - `X seconds ago`
   - `X minutes ago`
   - `X hours ago`
   - `X days ago`
   - `X weeks ago`
   - `X months ago`
   - `X years ago`

3. Future time:
   ```bash
   tasks add "Team meeting" --due "in 2 hours"
   ```
   Supported formats:
   - `in X seconds`
   - `in X minutes`
   - `in X hours`
   - `in X days`
   - `in X weeks`
   - `in X months`
   - `in X years`

Examples:

```bash
# Due tomorrow
tasks add "Review pull requests" --due tomorrow

# Due in the future
tasks add "Team meeting" --due "in 2 hours"
tasks add "Quarterly review" --due "in 3 months"
tasks add "Weekly sync" --due "in 1 week"

# Recurring tasks from past
tasks add "Weekly report" --due "1 week ago"
tasks add "Monthly maintenance" --due "1 month ago"
```

Note: The time expressions are case-insensitive and support both singular and plural units (e.g., both "1 hour" and "2 hours" work).

When listing tasks, the due dates are automatically formatted into human-readable relative time using the `timediff` package.

#### List Tasks

```bash
tasks list [flags]

Flags:
  -a, --all              Show all tasks (including completed)
  -c, --columns strings  Columns to display
  -s, --save            Save selected columns to config
```

Available columns:

- `id`: Task identifier
- `description`: Task description
- `iscompleted`: Completion status
- `createdat`: Creation timestamp
- `duedate`: Due date

Example:

```bash
tasks list --columns id,description,duedate --save
```

#### Complete a Task

```bash
tasks complete [task_id]
```

Example:

```bash
tasks complete abc123
```

#### Delete a Task

```bash
tasks delete [task_id]
```

Example:

```bash
tasks delete abc123
```

#### Set Storage Mode

```bash
tasks set-mode [mode]
```

Supported modes:

- `json`: Store tasks in a JSON file
- `csv`: Store tasks in a CSV file
- `sql`: Store tasks in a SQL database

Example:

```bash
tasks set-mode json
```

## Storage Backends

### JSON Storage

Tasks are stored in `~/.tasks/tasks.json`

### CSV Storage

Tasks are stored in `~/.tasks/tasks.csv`

### SQL Storage

Requires a database connection string in the `DB_URL` environment variable.

## Example Usage Workflow

1. Add a new task:

```bash
tasks add "Write blog post" --due "next monday"
```

2. List all pending tasks:

```bash
tasks list
```

3. Complete a task:

```bash
tasks complete abc123
```

4. View all tasks including completed:

```bash
tasks list --all
```

5. Customize column display:

```bash
tasks list --columns id,description,duedate --save
```

## Development

### Project Structure

```
.
├── internal/
│   ├── cli/         # CLI implementation
│   ├── config/      # Configuration management
│   ├── storage/     # Storage backends
│   │   ├── csv/
│   │   ├── json/
│   │   └── sql/
│   ├── task/       # Core task domain
│   └── utils/      # Utility functions
└── main.go
```

### Adding New Features

1. Implement new functionality in appropriate package
2. Add new command in `internal/cli/commands.go`
3. Register command in `setupCommands()` in `internal/cli/app.go`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
