package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Account struct {
	Username string
	Token    string
}

type AccountConfig struct {
	WorkAccount     Account
	PersonalAccount Account
}

func loadEnv(filename string) (map[string]string, error) {
	env := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing file: %v\n", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		env[key] = value
	}

	return env, scanner.Err()
}

func loadAccountConfig(env map[string]string) (*AccountConfig, error) {
	config := &AccountConfig{}

	personalAccount, personalAccountUsernameExists := env["GITHUB_PERSONAL_USER"]
	personalAccountToken, personalAccountTokenExists := env["GITHUB_PERSONAL_TOKEN"]

	if !personalAccountUsernameExists || !personalAccountTokenExists {
		return nil, fmt.Errorf("missing personal account configuration: GITHUB_PERSONAL_USER and GITHUB_PERSONAL_TOKEN required")
	}

	config.PersonalAccount = Account{
		Username: personalAccount,
		Token:    personalAccountToken,
	}

	workAccountUsername, workAccountUsernameExists := env["GITHUB_WORK_USER"]
	workAccountToken, workAccountTokenExists := env["GITHUB_WORK_TOKEN"]

	if !workAccountUsernameExists || !workAccountTokenExists {
		return nil, fmt.Errorf("missing Work account configuration: GITHUB_WORK_USER and GITHUB_WORK_TOKEN required")
	}

	config.WorkAccount = Account{
		Username: workAccountUsername,
		Token:    workAccountToken,
	}

	return config, nil
}

func parseRepoName(input string) string {
	httpRegex := regexp.MustCompile(`https?://(?:.*@)?github\.com/(?:[^/]+)/([^/.]+)(?:\.git)?$`)
	sshRegex := regexp.MustCompile(`git@github\.com:(?:[^/]+)/([^/.]+)(?:\.git)?$`)

	if matches := httpRegex.FindStringSubmatch(input); len(matches) > 1 {
		return matches[1]
	}
	if matches := sshRegex.FindStringSubmatch(input); len(matches) > 1 {
		return matches[1]
	}

	return input
}

func main() {
	var (
		personalAccountFlag = flag.Bool("p", false, "Use personal account configuration")
		workAccountFlag     = flag.Bool("w", false, "Use work account configuration")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-p|-w] <repository-name>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  -p    Use personal account configuration\n")
		fmt.Fprintf(os.Stderr, "  -w    Use work account configuration\n")
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s -p Example\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -w MyRepo\n", os.Args[0])
	}

	flag.Parse()

	if *workAccountFlag && *personalAccountFlag {
		fmt.Fprintf(os.Stderr, "Error: -p and -w flags are mutually exclusive\n")
		flag.Usage()
		os.Exit(1)
	}

	if !*workAccountFlag && !*personalAccountFlag {
		fmt.Fprintf(os.Stderr, "Error: You must specify either -p or -w flag\n")
		flag.Usage()
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: You must provide exactly one repository name\n")
		flag.Usage()
		os.Exit(1)
	}

	repoName := parseRepoName(args[0])

	envFile := ".secure-git.env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		execPath, _ := os.Executable()
		envFile = filepath.Join(filepath.Dir(execPath), ".secure-git.env")
	}

	env, err := loadEnv(envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .secure-git.env file: %v\n", err)
		fmt.Fprintf(os.Stderr, "Make sure you have a .secure-git.env file with GITHUB_WORK_USER, "+
			"GITHUB_WORK_TOKEN, "+
			"GITHUB_PERSONAL_USER, "+
			"and GITHUB_PERSONAL_TOKEN\n")
		os.Exit(1)
	}

	config, err := loadAccountConfig(env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading account configuration: %v\n", err)
		os.Exit(1)
	}

	var selectedAccount Account
	if *workAccountFlag {
		selectedAccount = config.WorkAccount
	} else {
		selectedAccount = config.PersonalAccount
	}

	repoURL := fmt.Sprintf("https://%s@github.com/%s/%s.git",
		selectedAccount.Token,
		selectedAccount.Username,
		repoName)

	fmt.Println(repoURL)
}
