package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Define command-line flags
    packageManager := flag.String("pm", "uv", "Package manager to use: uv or poetry")
    flag.Parse()

	// Validate package manager choice
    if *packageManager != "uv" && *packageManager != "poetry" {
        log.Fatalf("Invalid package manager: %s. Choose 'uv' or 'poetry'.", *packageManager)
    }

	// Prompt for root directory
	fmt.Println("Root directory where you want the repository to be cloned to (leave blank for current directory):")
	var rootDir string
	fmt.Scanln(&rootDir)
	if strings.TrimSpace(rootDir) == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("unable to determine current directory: %s", err)
		}
		rootDir = pwd
		fmt.Printf("Using current directory: %s\n", rootDir)
	}

	// Prompt user for repository
	fmt.Println("Repository you are cloning:")
	var repo string
	fmt.Scanln(&repo)

	// Change to root directory
	if err := os.Chdir(rootDir); err != nil {
		log.Fatalf("Error changing to directory: %s\n%s", rootDir, err)
	}

	// Extract repository name from URL (assuming format like https://github.com/user/repo.git)
    repoURLParts := strings.Split(repo, "/")
    repoName := strings.TrimSuffix(repoURLParts[len(repoURLParts)-1], ".git")
    clonedDir := filepath.Join(rootDir, repoName)

    // Check if repository directory already exists
    if _, err := os.Stat(clonedDir); err == nil {
        fmt.Printf("Repository directory %s already exists. Skipping clone and proceeding with environment setup.\n", clonedDir)
    } else {
        // Clone repository
        clone := exec.Command("git", "clone", repo)
        cloneOutput, err := clone.CombinedOutput()
        if err != nil {
            log.Fatalf("clone.CombinedOutput() failed with %s\n%s", err, string(cloneOutput))
        }
        fmt.Printf("Repository cloned at %s\n", clonedDir)
    }

	// Change to cloned repository directory
    if err := os.Chdir(clonedDir); err != nil {
        log.Fatalf("Error changing to cloned directory: %s\n%s", clonedDir, err)
    }

	// Setup environment based on package manager choice
    switch *packageManager {
    case "uv":
        setupUVEnvironment()
    case "poetry":
        setupPoetryEnvironment()
    default:
        log.Fatalf("Unknown package manager: %s", *packageManager)
    }

    fmt.Println("Environment setup complete.")
}

func setupUVEnvironment() {
    // Create virtual environment
    venvCmd := exec.Command("uv", "venv")
    venvOutput, err := venvCmd.CombinedOutput()
    if err != nil {
        log.Fatalf("uv venv failed with %s\n%s", err, string(venvOutput))
    }

    // Check for pyproject.toml first
    if _, err := os.Stat("pyproject.toml"); err == nil {
        syncCmd := exec.Command("uv", "sync")
        syncOutput, err := syncCmd.CombinedOutput()
        if err != nil {
            log.Fatalf("uv sync failed with %s\n%s", err, string(syncOutput))
        }
    } else if _, err := os.Stat("requirements.txt"); err == nil {
        installCmd := exec.Command("uv", "pip", "install", "-r", "requirements.txt")
        installOutput, err := installCmd.CombinedOutput()
        if err != nil {
            log.Fatalf("uv pip install failed with %s\n%s", err, string(installOutput))
        }
    } else {
        fmt.Println("No pyproject.toml or requirements.txt found; skipping dependency installation.")
    }
}

func setupPoetryEnvironment() {
    // Check for pyproject.toml
    if _, err := os.Stat("pyproject.toml"); err != nil {
        log.Fatalf("No pyproject.toml found. Poetry requires a pyproject.toml file.")
    }

    // Install dependencies using poetry
    installCmd := exec.Command("poetry", "install")
    installOutput, err := installCmd.CombinedOutput()
    if err != nil {
        log.Fatalf("poetry install failed with %s\n%s", err, string(installOutput))
    }

    fmt.Println("Poetry environment created and dependencies installed.")
}