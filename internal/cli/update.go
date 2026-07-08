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

func getDownloadUrl(tag string) string {
	osName, arch := getPlatformArch()
	if osName == "" || arch == "" {
		return ""
	}
	filename := fmt.Sprintf("voinznext-%s-%s", osName, arch)
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}
	return fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
		repoOwner, repoName, tag, filename)
}

func downloadBinary(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
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

		latestTag := release.TagName

		if !compareVersions(Version, latestTag) {
			fmt.Printf("  %s You're already on the latest version (%s)\n", style.SprintGreen("✔"), style.Value(latestTag))
			return nil
		}

		fmt.Printf("  %s Latest version: %s (current: %s)\n", style.SprintCyan("●"), style.Value(latestTag), style.Dimmed(Version))
		fmt.Printf("  %s Downloading update...\n", style.SprintCyan("●"))

		downloadUrl := getDownloadUrl(release.TagName)
		if downloadUrl == "" {
			fmt.Printf("  %s Unsupported platform: %s/%s\n", style.SprintRed("✘"), runtime.GOOS, runtime.GOARCH)
			return fmt.Errorf("unsupported platform")
		}

		exe, err := os.Executable()
		if err != nil {
			fmt.Printf("  %s Cannot determine executable path: %v\n", style.SprintRed("✘"), err)
			return err
		}

		tmpPath := exe + ".new"
		if err := downloadBinary(downloadUrl, tmpPath); err != nil {
			os.Remove(tmpPath)
			fmt.Printf("  %s Download failed: %v\n", style.SprintRed("✘"), err)
			return err
		}

		if runtime.GOOS != "windows" {
			os.Chmod(tmpPath, 0755)
		}

		fmt.Printf("  %s Installing update...\n", style.SprintCyan("●"))

		if runtime.GOOS == "windows" {
			return replaceViaScript(exe, tmpPath)
		}

		if err := os.Rename(tmpPath, exe); err != nil {
			os.Remove(tmpPath)
			fmt.Printf("  %s Failed to install update: %v\n", style.SprintRed("✘"), err)
			return err
		}

		fmt.Printf("  %s VoinzNext has been updated to %s!\n", style.SprintGreen("✔"), style.Value(latestTag))
		fmt.Printf("  %s Run %s to verify\n", style.SprintCyan("●"), style.Value("voinznext version"))
		return nil
	},
}

func replaceViaScript(exe, tmpPath string) error {
	// The .bat script:
	// 1. Retries while the current process still locks the exe
	// 2. Moves the new binary over the old one
	// 3. Gives up after a bounded wait and deletes the temp binary
	// 4. Deletes itself at the end (standard Windows self-delete pattern)
	scriptPath := exe + ".update.bat"
	script := fmt.Sprintf(`@echo off
set RETRY_COUNT=0
:retry
ping -n 2 127.0.0.1 > nul
move /Y "%s" "%s" > nul 2>&1
if not exist "%s" (
  echo VoinzNext updated successfully!
  goto cleanup
)
set /A RETRY_COUNT+=1
if %%RETRY_COUNT%% GEQ 10 (
  echo Update failed - could not replace binary.
  del "%s" > nul 2>&1
  goto cleanup
)
goto retry
:cleanup
(goto) 2>nul & del "%%~f0"
`, tmpPath, exe, tmpPath, tmpPath)

	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		os.Remove(tmpPath)
		fmt.Printf("  %s Cannot create update script: %v\n", style.SprintRed("✘"), err)
		return err
	}

	cmd := exec.Command("cmd", "/C", scriptPath)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		os.Remove(tmpPath)
		os.Remove(scriptPath)
		fmt.Printf("  %s Cannot start update script: %v\n", style.SprintRed("✘"), err)
		return err
	}

	fmt.Printf("  %s Update will complete in a moment...\n", style.SprintGreen("✔"))
	fmt.Printf("  %s Run %s after this window closes\n", style.SprintCyan("●"), style.Value("voinznext version"))
	return nil
}
