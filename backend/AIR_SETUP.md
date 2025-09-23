# Air Live Reloading Setup for StaffFind Backend

Air provides automatic live reloading for Go development, making development faster and more efficient.

## 🚀 Quick Start

### Using Makefile (Recommended)

```bash
# Start development environment and backend with Air
make quick-start

# Or start backend with Air (after starting dev environment)
make backend-air

# Simple mode without Air
make backend-dev-simple
```

### Using Startup Script

```bash
# Start development environment with Air
./start-docker.sh dev-air

# Start development environment only
./start-docker.sh dev
```

### Manual Usage

```bash
# Start development database first
make dev

# Then start backend with Air
cd backend
air
```

## ⚙️ Configuration

### Air Configuration File

The Air configuration is in `backend/.air.toml` with the following key settings:

```toml
[build]
  cmd = "go build -o ./tmp/main ./cmd/server"
  bin = "tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html", "sql", "yaml", "yml"]
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "node_modules", ".git"]
  exclude_regex = ["_test.go"]
  log = "errors.log"

[misc]
  clean_on_exit = true

[screen]
  clear_on_rebuild = true
  keep_scroll = true
```

### Watched Files

Air watches for changes in:
- `.go` files (main application code)
- `.sql` files (migration files)
- `.yaml`/`.yml` files (configuration files)
- `.html`/`.tpl`/`.tmpl` files (templates)

### Ignored Files

Air ignores:
- Test files (`*_test.go`)
- Vendor directory
- Node modules
- Git directory
- Temporary files
- IDE files

## 🔧 Development Workflow

### 1. Start Development Environment

```bash
# Start PostgreSQL + pgAdmin
make dev
```

### 2. Start Backend with Air

```bash
# With live reloading (recommended)
make backend-air

# Or simple mode
make backend-dev
```

### 3. Make Changes

Any changes to Go files, SQL files, or configuration files will automatically:
1. Stop the current server
2. Rebuild the application
3. Restart the server
4. Show build logs

### 4. View Logs

Air provides colored output:
- **Magenta**: Main application logs
- **Cyan**: File watcher notifications
- **Yellow**: Build process logs
- **Green**: Server runner logs

## 📁 File Structure

```
backend/
├── .air.toml              # Air configuration
├── .gitignore             # Includes Air temp files
├── tmp/                   # Air temporary files (auto-created)
├── errors.log             # Air error logs (auto-created)
└── cmd/server/main.go     # Main application entry point
```

## 🛠️ Available Commands

### Makefile Commands

| Command | Description |
|---------|-------------|
| `make backend-air` | Run backend with Air (live reloading) |
| `make backend-dev` | Run backend in simple development mode |
| `make backend-dev-simple` | Run backend in simple development mode |
| `make quick-start` | Quick start with Air (database + backend) |
| `make quick-start-simple` | Quick start without Air |

### Startup Script Commands

| Command | Description |
|---------|-------------|
| `./start-docker.sh dev-air` | Start development with Air |
| `./start-docker.sh dev` | Start development environment only |

## 🔍 Troubleshooting

### Common Issues

#### 1. Air Not Found
```bash
Error: command not found: air
```
**Solution**: Install Air:
```bash
go install github.com/air-verse/air@latest
```

#### 2. Build Errors
```bash
Error: build failed
```
**Solution**: Check the error logs in `errors.log` or Air output

#### 3. Port Already in Use
```bash
Error: bind: address already in use
```
**Solution**: Stop other instances:
```bash
pkill -f air
pkill -f "go run"
```

#### 4. Database Connection Issues
```bash
Error: database connection failed
```
**Solution**: Ensure development database is running:
```bash
make dev
```

### Debug Commands

```bash
# Check if Air is running
ps aux | grep air

# Check if backend is running
curl http://localhost:8080/health

# View Air logs
tail -f backend/errors.log

# Stop all Go processes
pkill -f "go run"
pkill -f air
```

## 🎯 Benefits

### With Air (Live Reloading)
- ✅ Automatic server restart on code changes
- ✅ Fast rebuild and restart cycle
- ✅ Colored output for better debugging
- ✅ Error logging and reporting
- ✅ Clean temporary files on exit

### Without Air (Simple Mode)
- ✅ Direct Go execution
- ✅ No additional dependencies
- ✅ Simpler debugging
- ❌ Manual restart required for changes

## 📚 Best Practices

### Development with Air

1. **Use Air for Active Development**: When making frequent code changes
2. **Monitor Build Logs**: Watch for compilation errors
3. **Keep Database Running**: Use `make dev` to keep PostgreSQL running
4. **Clean Restart**: Use `Ctrl+C` to stop Air cleanly

### File Organization

1. **Separate Test Files**: Keep tests in separate files with `_test.go` suffix
2. **Avoid Vendor Changes**: Don't modify files in vendor directory
3. **Use Git Ignore**: Ensure `.gitignore` includes Air temporary files

### Performance Tips

1. **Exclude Large Directories**: Add large directories to `exclude_dir`
2. **Use Specific Extensions**: Only watch necessary file extensions
3. **Monitor Resource Usage**: Air uses minimal resources but monitor if needed

## 🔗 Integration

### With Docker

Air works perfectly with the Docker development environment:
- Database runs in Docker container
- Backend runs locally with Air
- Migrations run automatically on startup

### With IDEs

Air works with any IDE or editor:
- VS Code
- GoLand
- Vim/Neovim
- Emacs

### With Git

Air respects `.gitignore` and doesn't interfere with Git operations.

## 📖 Additional Resources

- [Air Documentation](https://github.com/air-verse/air)
- [Air Configuration Reference](https://github.com/air-verse/air#configuration)
- [Go Development Best Practices](https://golang.org/doc/effective_go.html)

## 🆘 Getting Help

If you encounter issues:

1. Check Air logs: `tail -f backend/errors.log`
2. Verify Air installation: `air --version`
3. Check database connection: `make status`
4. Review this documentation
5. Check the [Docker Setup Guide](../DOCKER_SETUP.md)
