# remotest

A CLI tool to generate authenticated GitHub repository URLs for cloning or pushing, using either your work or personal account credentials.

## Features
- Parse GitHub repository names from HTTP and SSH URLs.
- Support for configuring two GitHub accounts (work and personal).

## Prerequisites
- Go 1.18 or later installed on your system.

## Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   ```
2. Navigate to the project directory:
   ```bash
   cd GetSecureGitHubStr
   ```
3. Set up environment variables for your GitHub accounts:
   - `GITHUB_WORK_USER`: Your work GitHub username.
   - `GITHUB_WORK_TOKEN`: Your work GitHub personal access token.
   - `GITHUB_PERSONAL_USER`: Your personal GitHub username.
   - `GITHUB_PERSONAL_TOKEN`: Your personal GitHub personal access token.
   - 
4. Build the project:
   ```bash
   go build -ldflags="-s -w" -o remotest
   ```

## Usage
Run the program with the following flags:
- `-p`: Use personal account configuration.
- `-w`: Use work account configuration.

Example:
```bash
./remotest -p
```

## Project Structure
- `main.go`: The main application logic.
- `go.mod`: Go module file.
- `build.sh`: Build script for the project.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request.
