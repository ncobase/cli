package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CheckAndInstallAtlas checks if atlas is installed, and if not, asks to install it
func CheckAndInstallAtlas() error {
	_, err := exec.LookPath("atlas")
	if err == nil {
		return nil // Atlas is already installed
	}

	fmt.Println("Atlas CLI is required but not found in your PATH.")
	fmt.Print("Would you like to install it via 'go install ariga.io/atlas/cmd/atlas@latest'? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		return fmt.Errorf("atlas is required to run this command. Please install it manually")
	}

	fmt.Println("Installing Atlas...")
	cmd := exec.Command("go", "install", "ariga.io/atlas/cmd/atlas@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install atlas: %v", err)
	}

	// Verify installation
	if _, err := exec.LookPath("atlas"); err != nil {
		// If still not found in PATH, try to locate it in GOBIN
		goBin, err := getGoBin()
		if err == nil {
			atlasPath := filepath.Join(goBin, "atlas")
			if _, err := os.Stat(atlasPath); err == nil {
				fmt.Printf("Atlas installed successfully to %s.\n", atlasPath)
				fmt.Println("Warning: This directory is not in your PATH. You may need to add it.")
				return nil
			}
			// check windows extension
			if _, err := os.Stat(atlasPath + ".exe"); err == nil {
				fmt.Printf("Atlas installed successfully to %s.exe.\n", atlasPath)
				fmt.Println("Warning: This directory is not in your PATH. You may need to add it.")
				return nil
			}
		}
		return fmt.Errorf("atlas installed but not found in PATH. Please add $GOPATH/bin or $GOBIN to your PATH")
	}

	fmt.Println("Atlas installed successfully!")
	return nil
}

func getGoBin() (string, error) {
	cmd := exec.Command("go", "env", "GOBIN")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	gobin := strings.TrimSpace(string(out))
	if gobin != "" {
		return gobin, nil
	}

	cmd = exec.Command("go", "env", "GOPATH")
	out, err = cmd.Output()
	if err != nil {
		return "", err
	}
	gopath := strings.TrimSpace(string(out))
	if gopath == "" {
		return "", fmt.Errorf("GOPATH is not set")
	}
	return filepath.Join(gopath, "bin"), nil
}
