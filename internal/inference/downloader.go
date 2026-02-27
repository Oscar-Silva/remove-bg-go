package inference

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DownloadProgress struct {
	TotalBytes      int64
	DownloadedBytes int64
}

// DownloadProgressTracker tracks the download and emits events
type DownloadProgressTracker struct {
	ctx          context.Context
	totalBytes   int64
	downloaded   int64
	lastEmitTime time.Time
}

func (pt *DownloadProgressTracker) Write(p []byte) (int, error) {
	n := len(p)
	pt.downloaded += int64(n)

	// Emit at most every 200ms
	if time.Since(pt.lastEmitTime) > 200*time.Millisecond {
		runtime.EventsEmit(pt.ctx, "download_progress", map[string]interface{}{
			"downloaded": pt.downloaded,
			"total":      pt.totalBytes,
		})
		pt.lastEmitTime = time.Now()
	}

	return n, nil
}

// DownloadModel downloads the requested model from Hugging Face
func DownloadModel(ctx context.Context, modelID, destPath string) error {
	url := fmt.Sprintf("https://huggingface.co/camenduru/RMBG-2.0/resolve/main/onnx/%s?download=true", modelID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// We no longer strictly need HF_TOKEN since we use a public mirror,
	// but we could still pass it if set.
	hfToken := os.Getenv("HF_TOKEN")
	if hfToken != "" {
		req.Header.Set("Authorization", "Bearer "+hfToken)
	}

	client := &http.Client{
		Timeout: 30 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download (HTTP %d). Make sure HF_TOKEN is set for gated models", resp.StatusCode)
	}

	totalBytes := resp.ContentLength

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create structural directories: %w", err)
	}

	// Create temp file for downloading
	tmpPath := destPath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpPath) // cleaned up if not renamed successfully

	tracker := &DownloadProgressTracker{
		ctx:        ctx,
		totalBytes: totalBytes,
	}

	// Emit initial progress event
	runtime.EventsEmit(ctx, "status", "downloading_model")
	runtime.EventsEmit(ctx, "download_progress", map[string]interface{}{
		"downloaded": 0,
		"total":      totalBytes,
	})

	_, err = io.Copy(io.MultiWriter(out, tracker), resp.Body)
	if err != nil {
		out.Close()
		return fmt.Errorf("interrupted during download: %w", err)
	}

	out.Close()

	// Ensure final 100% emission
	runtime.EventsEmit(ctx, "download_progress", map[string]interface{}{
		"downloaded": totalBytes,
		"total":      totalBytes,
	})

	// Download finished successfully, rename tmp to target
	if err := os.Rename(tmpPath, destPath); err != nil {
		return fmt.Errorf("failed to copy temporary file to final destination: %w", err)
	}

	return nil
}
