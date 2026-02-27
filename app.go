package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	_ "golang.org/x/image/webp"

	"remove-bg-go/internal/inference"
)

// App struct
type App struct {
	ctx           context.Context
	session       *inference.Session
	preprocessor  *inference.Preprocessor
	postprocessor *inference.Postprocessor
	modelDir      string // base models dir
}

type ModelConfig struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SizeMB    int    `json:"sizeMB"`
	RAM       string `json:"ram"`
	Speed     string `json:"speed"`
	Quality   string `json:"quality"`
	IsDefault bool   `json:"isDefault"`
}

var AvailableModels = []ModelConfig{
	{ID: "model.onnx", Name: "Alta Precisão (FP32)", SizeMB: 1024, RAM: "~2.0 GB", Speed: "Lenta", Quality: "Excelente", IsDefault: false},
	{ID: "model_fp16.onnx", Name: "Equilibrado (FP16)", SizeMB: 513, RAM: "~1.0 GB", Speed: "Rápida", Quality: "Ótima", IsDefault: true},
	{ID: "model_quantized.onnx", Name: "Rápido (INT8)", SizeMB: 366, RAM: "~700 MB", Speed: "Muito Rápida", Quality: "Boa", IsDefault: false},
	{ID: "model_bnb4.onnx", Name: "Ultra Rápido (Q4)", SizeMB: 233, RAM: "~500 MB", Speed: "Mais Rápida", Quality: "Aceitável", IsDefault: false},
}

// NewApp creates a new App application struct
func NewApp() *App {
	modelDir := "models/RMBG-2.0/onnx"
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		absModelDir := filepath.Join(exeDir, "..", "..", "models", "RMBG-2.0", "onnx")
		if _, err := os.Stat(filepath.Join(absModelDir, ".")); err == nil {
			modelDir = absModelDir
		} else {
			// fallback check inside same dir
			absModelDir2 := filepath.Join(exeDir, "models", "RMBG-2.0", "onnx")
			modelDir = absModelDir2
		}
	}

	return &App{
		preprocessor:  inference.NewPreprocessor(1024, []float32{0.5, 0.5, 0.5}, []float32{0.5, 0.5, 0.5}),
		postprocessor: inference.NewPostprocessor(),
		modelDir:      modelDir,
	}
}

// GetModels returns the available models configuration
func (a *App) GetModels() []ModelConfig {
	return AvailableModels
}

// IsModelDownloaded checks if a specific model exists in local disk
func (a *App) IsModelDownloaded(modelID string) bool {
	modelPath := filepath.Join(a.modelDir, modelID)
	_, err := os.Stat(modelPath)
	return err == nil
}

// RemoveBackground removes the background from a base64 encoded image
func (a *App) RemoveBackground(imageBase64 string, modelID string) (string, error) {
	// Validate modelID
	valid := false
	for _, m := range AvailableModels {
		if m.ID == modelID {
			valid = true
			break
		}
	}
	if !valid {
		return "", fmt.Errorf("invalid model selected: %s", modelID)
	}

	modelPath := filepath.Join(a.modelDir, modelID)

	// Check if model exists, download if it doesn't
	if !a.IsModelDownloaded(modelID) {
		runtime.EventsEmit(a.ctx, "status", "downloading_model")
		if err := inference.DownloadModel(a.ctx, modelID, modelPath); err != nil {
			runtime.EventsEmit(a.ctx, "status", "error")
			return "", fmt.Errorf("failed to download requested model: %w", err)
		}
	}
	// Initialize session
	// If the user requested a different model than what is loaded, we must reload
	if a.session != nil && a.session.GetModelPath() != modelPath {
		a.session.Destroy()
		a.session = nil
	}

	if a.session == nil {
		runtime.EventsEmit(a.ctx, "status", "loading_model")
		session, err := inference.NewSession(modelPath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "status", "error")
			return "", fmt.Errorf("failed to initialize ONNX session: %w", err)
		}
		a.session = session
	}

	// Decode base64 image
	runtime.EventsEmit(a.ctx, "status", "decoding")
	imgData, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		runtime.EventsEmit(a.ctx, "status", "error")
		return "", fmt.Errorf("failed to decode base64 image: %w", err)
	}

	// Decode PNG
	img, format, err := image.Decode(strings.NewReader(string(imgData)))
	if err != nil {
		runtime.EventsEmit(a.ctx, "status", "error")
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Verify format
	if format != "png" && format != "jpeg" && format != "jpg" {
		runtime.EventsEmit(a.ctx, "status", "error")
		return "", fmt.Errorf("unsupported image format: %s", format)
	}

	// Get original dimensions
	origWidth, origHeight := inference.GetOriginalDimensions(img)

	// Preprocess
	runtime.EventsEmit(a.ctx, "status", "preprocessing")
	inputData, err := a.preprocessor.Preprocess(img)
	if err != nil {
		runtime.EventsEmit(a.ctx, "status", "error")
		return "", fmt.Errorf("failed to preprocess image: %w", err)
	}

	// Run inference
	runtime.EventsEmit(a.ctx, "status", "processing")
	err = a.session.RunInference(inputData)
	if err != nil {
		runtime.EventsEmit(a.ctx, "status", "error")
		return "", fmt.Errorf("inference failed: %w", err)
	}

	// Get output data
	outputData := a.session.GetOutputData()

	// Postprocess
	runtime.EventsEmit(a.ctx, "status", "finalizing")
	resultBytes, err := a.postprocessor.Postprocess(outputData, img, origWidth, origHeight)
	if err != nil {
		runtime.EventsEmit(a.ctx, "status", "error")
		return "", fmt.Errorf("failed to postprocess: %w", err)
	}

	// Encode result to base64
	resultBase64 := base64.StdEncoding.EncodeToString(resultBytes)

	// Emit completion status
	runtime.EventsEmit(a.ctx, "status", "done")

	return resultBase64, nil
}

// GetVersion returns the app version
func (a *App) GetVersion() string {
	return "1.0.0"
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.EventsEmit(ctx, "status", "initializing")
}
