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
	fmt.Println("Root directory where you want the repository to be cloned to:")
	var root_dir string
	fmt.Scanln(&root_dir)

	// Prompt user for repository
	fmt.Println("Repository you are cloning:")
	var repo string
	fmt.Scanln(&repo)

	// Change to root directory
	if err := os.Chdir(root_dir); err != nil {
		log.Fatalf("Error changing to directory: %s\n%s", root_dir, err)
	}

	// Extract repository name from URL (assuming format like https://github.com/user/repo.git)
    repoURLParts := strings.Split(repo, "/")
    repoName := strings.TrimSuffix(repoURLParts[len(repoURLParts)-1], ".git")
    clonedDir := filepath.Join(root_dir, repoName)

    // Check if repository directory already exists
    if _, err := os.Stat(clonedDir); err == nil {
        fmt.Printf("Repository directory %s already exists. Skipping clone and proceeding with environment setup.\n", clonedDir)
    } else {
        // Clone repository
        clone := exec.Command("git", "clone", repo)
        clone_output, err := clone.CombinedOutput()
        if err != nil {
            log.Fatalf("clone.CombinedOutput() failed with %s\n%s", err, string(clone_output))
        }
        fmt.Printf("Repository cloned at %s\n", clonedDir)
    }

	// Change to cloned repository directory
    if err := os.Chdir(clonedDir); err != nil {
        log.Fatalf("Error changing to cloned directory: %s\n%s", clonedDir, err)
    }

	// Setup environment based on package manager choice
    if *packageManager == "uv" {
        setupUVEnvironment()
    } else if *packageManager == "poetry" {
        setupPoetryEnvironment()
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