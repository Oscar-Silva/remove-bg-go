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
	modelPath     string
}

// NewApp creates a new App application struct
func NewApp() *App {
	modelPath := "models/RMBG-2.0/onnx/model.onnx"
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		absModelPath := filepath.Join(exeDir, "..", "..", "models", "RMBG-2.0", "onnx", "model.onnx")
		if _, err := os.Stat(absModelPath); err == nil {
			modelPath = absModelPath
		} else {
			// fallback check inside same dir
			absModelPath2 := filepath.Join(exeDir, "models", "RMBG-2.0", "onnx", "model.onnx")
			if _, err := os.Stat(absModelPath2); err == nil {
				modelPath = absModelPath2
			}
		}
	}

	return &App{
		preprocessor:  inference.NewPreprocessor(1024, []float32{0.5, 0.5, 0.5}, []float32{0.5, 0.5, 0.5}),
		postprocessor: inference.NewPostprocessor(),
		modelPath:     modelPath,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.EventsEmit(ctx, "status", "initializing")
}

// RemoveBackground removes the background from a base64 encoded image
func (a *App) RemoveBackground(imageBase64 string) (string, error) {
	// Emit loading status
	runtime.EventsEmit(a.ctx, "status", "loading_model")

	// Initialize session if not already done
	if a.session == nil {
		session, err := inference.NewSession(a.modelPath)
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
