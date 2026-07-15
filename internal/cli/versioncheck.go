package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/VoinzzZ/VoinzNext/internal/style"
)

const (
	repoOwner = "VoinzzZ"
	repoName  = "VoinzNext"
	cacheFile = "voinznext_update_check.json"
	cacheTTL  = 24 * time.Hour
)

type updateCache struct {
	CheckedAt time.Time `json:"checked_at"`
	LatestTag string    `json:"latest_tag"`
	LatestURL string    `json:"latest_url"`
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func getCachePath() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	return filepath.Join(cacheDir, "voinznext", cacheFile)
}

func loadCache() *updateCache {
	data, err := os.ReadFile(getCachePath())
	if err != nil {
		return nil
	}
	var c updateCache
	if json.Unmarshal(data, &c) != nil {
		return nil
	}
	return &c
}

func saveCache(tag, url string) {
	c := updateCache{
		CheckedAt: time.Now(),
		LatestTag: tag,
		LatestURL: url,
	}
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	path := getCachePath()
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, data, 0644)
}

func fetchLatestRelease() (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "voinznext-update-check")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

// parseVersion splits a version string like "0.10.3" into [0, 10, 3].
// Returns [0, 0, 0] if parsing fails.
func parseVersion(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	var result [3]int
	for i := 0; i < len(parts) && i < 3; i++ {
		n, err := strconv.Atoi(parts[i])
		if err != nil {
			return [3]int{}
		}
		result[i] = n
	}
	return result
}

// compareVersions returns true if current < latest (i.e. update available).
func compareVersions(current, latest string) bool {
	c := parseVersion(current)
	l := parseVersion(latest)
	for i := 0; i < 3; i++ {
		if c[i] < l[i] {
			return true
		}
		if c[i] > l[i] {
			return false
		}
	}
	return false
}

func CheckForUpdate() {
	cached := loadCache()

	var latestTag string
	var latestURL string

	if cached != nil && time.Since(cached.CheckedAt) < cacheTTL {
		latestTag = cached.LatestTag
		latestURL = cached.LatestURL
	} else {
		release, err := fetchLatestRelease()
		if err != nil {
			return
		}
		latestTag = release.TagName
		latestURL = release.HTMLURL
		saveCache(latestTag, latestURL)
	}

	if latestTag == "" {
		return
	}

	currentTag := "v" + strings.TrimPrefix(Version, "v")
	if compareVersions(currentTag, latestTag) {
		fmt.Println()
		fmt.Printf("  %s %s\n", style.SprintYellow("⚠"), style.Dimmed(fmt.Sprintf("A new version is available: %s", style.SprintWhite(latestTag))))
		fmt.Printf("  %s %s\n", style.SprintYellow("⚠"), style.Dimmed("Run \"voinznext update\" to update."))
		fmt.Println()
	} else if cached != nil && cached.LatestTag != currentTag {
		// Current version is newer or equal to cached latest — refresh cache
		// This handles the case where user updated via npm/go install but cache still shows old version
		saveCache(currentTag, "")
	}
}

func IsNpmInstall() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	return strings.Contains(exe, "node_modules") || runtime.GOOS == "windows" && strings.Contains(strings.ToLower(exe), "nvm")
}

func IsGoInstall() bool {
	_, err := os.Stat(filepath.Join(runtime.GOROOT(), "bin", "voinznext"))
	if err == nil {
		return true
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}
	_, err = os.Stat(filepath.Join(gopath, "bin", "voinznext"))
	return err == nil
}
