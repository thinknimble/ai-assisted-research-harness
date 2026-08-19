package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	updateRepoOwner = "thinknimble"
	updateRepoName  = "ai-assisted-research-harness"
)

// version is set at build time via -ldflags.
var version = "dev"

// httpClient abstracts HTTP requests for testability.
var httpClient interface {
	Do(req *http.Request) (*http.Response, error)
} = http.DefaultClient

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func runUpdate(checkOnly bool) error {
	if version == "dev" {
		return fmt.Errorf("cannot update a development build; install a released version first")
	}

	release, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	if latest == "" {
		return fmt.Errorf("no releases found on GitHub")
	}

	if !isNewer(latest, version) {
		fmt.Printf("Already on latest version (%s)\n", version)
		return nil
	}

	if checkOnly {
		fmt.Printf("Current version: %s\nLatest version:  %s\nRun 'research-assistant update' to install\n", version, latest)
		return nil
	}

	fmt.Printf("Updating %s → %s ...\n", version, latest)

	assetName := updateAssetName(runtime.GOOS, runtime.GOARCH)
	downloadURL := ""
	for _, a := range release.Assets {
		if a.Name == assetName {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("no binary found for %s/%s in release %s", runtime.GOOS, runtime.GOARCH, release.TagName)
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot determine executable path: %w", err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("cannot resolve executable path: %w", err)
	}

	if err := downloadAndReplace(downloadURL, exePath); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Printf("Updated successfully: %s → %s\n", version, latest)
	return nil
}

func fetchLatestRelease() (*ghRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", updateRepoOwner, updateRepoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var release ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
	}
	return &release, nil
}

func downloadAndReplace(url, destPath string) error {
	dir := filepath.Dir(destPath)
	tmp, err := os.CreateTemp(dir, ".research-assistant-update-*")
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("no write permission for %s — run: sudo research-assistant update, or reinstall to ~/.local/bin: curl -fsSL https://raw.githubusercontent.com/%s/%s/main/install.sh | sh", dir, updateRepoOwner, updateRepoName)
		}
		return fmt.Errorf("cannot create temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer func() {
		tmp.Close()
		os.Remove(tmpName)
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/octet-stream")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return fmt.Errorf("cannot write new binary: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("cannot close temp file: %w", err)
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(tmpName, 0755); err != nil {
			return fmt.Errorf("cannot chmod new binary: %w", err)
		}
	}

	if runtime.GOOS == "windows" {
		if err := os.Rename(destPath, destPath+".old"); err != nil {
			return fmt.Errorf("cannot move old binary: %w", err)
		}
	}

	if err := os.Rename(tmpName, destPath); err != nil {
		return fmt.Errorf("cannot replace binary: %w", err)
	}
	return nil
}

func updateAssetName(goos, goarch string) string {
	name := fmt.Sprintf("research-assistant-%s-%s", goos, goarch)
	if goos == "windows" {
		name += ".exe"
	}
	return name
}

func parseVersion(v string) []int {
	v = strings.TrimPrefix(v, "v")
	if idx := strings.IndexByte(v, '-'); idx != -1 {
		v = v[:idx]
	}
	parts := strings.Split(v, ".")
	nums := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			break
		}
		nums = append(nums, n)
	}
	if len(nums) == 0 {
		return []int{0}
	}
	return nums
}

func isNewer(a, b string) bool {
	va := parseVersion(a)
	vb := parseVersion(b)
	n := len(va)
	if len(vb) > n {
		n = len(vb)
	}
	for i := 0; i < n; i++ {
		ai, bi := 0, 0
		if i < len(va) {
			ai = va[i]
		}
		if i < len(vb) {
			bi = vb[i]
		}
		if ai != bi {
			return ai > bi
		}
	}
	return false
}
