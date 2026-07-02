package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/VoinzzZ/VoinzNext/internal/style"
	"github.com/spf13/cobra"
)

const binaryName = "voinznext"

func getTargetName() string {
	if runtime.GOOS == "windows" {
		return binaryName + ".exe"
	}
	return binaryName
}

func getPlatformArch() (string, string) {
	archMap := map[string]string{
		"amd64": "amd64",
		"arm64": "arm64",
	}
	osMap := map[string]string{
		"windows": "windows",
		"darwin":  "darwin",
		"linux":   "linux",
	}
	return osMap[runtime.GOOS], archMap[runtime.GOARCH]
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VoinzNext to the latest version",
	Long: `Check for updates and install the latest version of VoinzNext.

Auto-detects installation method:
  - npm: runs "npm install -g voinznext@latest"
  - binary: downloads latest release from GitHub
  - go: runs "go install"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("  %s Checking for updates...\n", style.SprintCyan("●"))

		release, err := fetchLatestRelease()
		if err != nil {
			fmt.Printf("  %s Failed to check for updates: %v\n", style.SprintRed("✘"), err)
			return err
		}

		currentTag := "v" + Version
		latestTag := release.TagName

		if !compareVersions(currentTag, latestTag) {
			fmt.Printf("  %s You're already on the latest version (%s)\n", style.SprintGreen("✔"), style.Value(latestTag))
			return nil
		}

		fmt.Printf("  %s Latest version: %s (current: %s)\n", style.SprintCyan("●"), style.Value(latestTag), style.Dimmed(currentTag))
		fmt.Printf("  %s Downloading update...\n", style.SprintCyan("●"))

		if IsNpmInstall() {
			return updateViaNpm()
		}

		if IsGoInstall() {
			return updateViaGo()
		}

		return updateViaBinary(release.TagName)
	},
}

func updateViaNpm() error {
	cmd := exec.Command("npm", "install", "-g", "voinznext@latest")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("  %s npm update failed: %v\n", style.SprintRed("✘"), err)
		fmt.Printf("  %s\n", style.Dimmed(string(output)))
		return err
	}
	fmt.Printf("  %s VoinzNext has been updated!\n", style.SprintGreen("✔"))
	fmt.Printf("  %s Run %s to verify\n", style.SprintCyan("●"), style.Value("voinznext version"))
	return nil
}

func updateViaGo() error {
	cmd := exec.Command("go", "install", "github.com/VoinzzZ/VoinzNext/cmd/voinznext@latest")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("  %s Update failed: %v\n", style.SprintRed("✘"), err)
		fmt.Printf("  %s\n", style.Dimmed(string(output)))
		return err
	}
	fmt.Printf("  %s VoinzNext has been updated!\n", style.SprintGreen("✔"))
	fmt.Printf("  %s Run %s to verify\n", style.SprintCyan("●"), style.Value("voinznext version"))
	return nil
}

func updateViaBinary(tag string) error {
	osName, arch := getPlatformArch()
	if osName == "" || arch == "" {
		fmt.Printf("  %s Unsupported platform: %s/%s\n", style.SprintRed("✘"), runtime.GOOS, runtime.GOARCH)
		return fmt.Errorf("unsupported platform")
	}

	filename := fmt.Sprintf("voinznext-%s-%s", osName, arch)
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}

	downloadUrl := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
		repoOwner, repoName, tag, filename)

	exe, err := os.Executable()
	if err != nil {
		fmt.Printf("  %s Cannot determine executable path: %v\n", style.SprintRed("✘"), err)
		return err
	}

	tmpPath := exe + ".new"
	resp, err := http.Get(downloadUrl)
	if err != nil {
		fmt.Printf("  %s Download failed: %v\n", style.SprintRed("✘"), err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("  %s Download failed with status %d\n", style.SprintRed("✘"), resp.StatusCode)
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(tmpPath)
	if err != nil {
		fmt.Printf("  %s Cannot create temp file: %v\n", style.SprintRed("✘"), err)
		return err
	}

	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		os.Remove(tmpPath)
		fmt.Printf("  %s Download failed: %v\n", style.SprintRed("✘"), err)
		return err
	}

	if runtime.GOOS != "windows" {
		os.Chmod(tmpPath, 0755)
	}

	if err := os.Rename(tmpPath, exe); err != nil {
		if runtime.GOOS == "windows" {
			_ = os.Remove(exe)
			if err := os.Rename(tmpPath, exe); err != nil {
				os.Remove(tmpPath)
				fmt.Printf("  %s Failed to replace binary: %v\n", style.SprintRed("✘"), err)
				fmt.Printf("  %s New binary saved at: %s\n", style.SprintYellow("⚠"), tmpPath)
				return err
			}
		} else {
			os.Remove(tmpPath)
			fmt.Printf("  %s Failed to replace binary: %v\n", style.SprintRed("✘"), err)
			return err
		}
	}

	fmt.Printf("  %s VoinzNext has been updated to %s!\n", style.SprintGreen("✔"), style.Value(tag))
	fmt.Printf("  %s Run %s to verify\n", style.SprintCyan("●"), style.Value("voinznext version"))
	return nil
}


