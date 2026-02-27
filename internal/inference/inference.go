package inference

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"

	ort "github.com/yalue/onnxruntime_go"
)

// Session represents ONNX Runtime session with loaded model
type Session struct {
	session      *ort.AdvancedSession
	inputTensor  *ort.Tensor[float32]
	outputTensor *ort.Tensor[float32]
	inputName    string
	outputName   string
	loaded       bool
	modelPath    string
}

// NewSession creates a new ONNX session and loads the model
func NewSession(modelPath string) (*Session, error) {
	// Initialize ONNX Runtime
	ort.SetSharedLibraryPath("/usr/local/lib/libonnxruntime.so.1")
	if err := ort.InitializeEnvironment(); err != nil {
		return nil, fmt.Errorf("failed to initialize ONNX environment: %w", err)
	}

	// Input shape: [1, 3, 1024, 1024]
	inputShape := ort.NewShape(1, 3, 1024, 1024)
	// Output shape: [1, 1, 1024, 1024]
	outputShape := ort.NewShape(1, 1, 1024, 1024)

	// Create input tensor (pre-allocated)
	inputData := make([]float32, 3*1024*1024)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	if err != nil {
		ort.DestroyEnvironment()
		return nil, fmt.Errorf("failed to create input tensor: %w", err)
	}

	// Create output tensor (pre-allocated)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		inputTensor.Destroy()
		ort.DestroyEnvironment()
		return nil, fmt.Errorf("failed to create output tensor: %w", err)
	}

	// The model input/output names need to be inferred from the model
	// For RMBG model, typically input is "input" and output is "sigmoid_0" or similar
	// We'll try common names and adjust if needed
	inputNames := []string{"input"}
	outputNames := []string{"sigmoid_0"}

	// Try to create session - if input/output names are wrong, we'll get an error
	// and can try alternative names
	session, err := ort.NewAdvancedSession(
		modelPath,
		inputNames,
		outputNames,
		[]ort.Value{inputTensor},
		[]ort.Value{outputTensor},
		nil,
	)

	if err != nil {
		// Try alternative names that might work
		inputNames = []string{"x"}
		outputNames = []string{"sigmoid"}
		session, err = ort.NewAdvancedSession(
			modelPath,
			inputNames,
			outputNames,
			[]ort.Value{inputTensor},
			[]ort.Value{outputTensor},
			nil,
		)
		if err != nil {
			inputNames = []string{"input.1"}
			outputNames = []string{"output.1"}
			session, err = ort.NewAdvancedSession(
				modelPath,
				inputNames,
				outputNames,
				[]ort.Value{inputTensor},
				[]ort.Value{outputTensor},
				nil,
			)
			if err != nil {
				inputTensor.Destroy()
				outputTensor.Destroy()
				ort.DestroyEnvironment()
				return nil, fmt.Errorf("failed to create ONNX session: %w", err)
			}
		}
	}

	return &Session{
		session:      session,
		inputTensor:  inputTensor,
		outputTensor: outputTensor,
		inputName:    inputNames[0],
		outputName:   outputNames[0],
		loaded:       true,
		modelPath:    modelPath,
	}, nil
}

// RunInference runs the ONNX model inference
func (s *Session) RunInference(inputData []float32) error {
	if !s.loaded {
		return errors.New("session not loaded")
	}

	// Copy input data to tensor
	copy(s.inputTensor.GetData(), inputData)

	// Run inference
	if err := s.session.Run(); err != nil {
		return fmt.Errorf("inference failed: %w", err)
	}

	return nil
}

// GetOutputData returns the output tensor data
func (s *Session) GetOutputData() []float32 {
	return s.outputTensor.GetData()
}

// Destroy cleans up the session
func (s *Session) Destroy() {
	if s.session != nil {
		s.session.Destroy()
		s.session = nil
	}
	if s.inputTensor != nil {
		s.inputTensor.Destroy()
		s.inputTensor = nil
	}
	if s.outputTensor != nil {
		s.outputTensor.Destroy()
		s.outputTensor = nil
	}
	ort.DestroyEnvironment()
	s.loaded = false
}

// GetInputShape returns expected input shape
func (s *Session) GetInputShape() []int64 {
	return []int64{1, 3, 1024, 1024}
}

// GetOutputShape returns expected output shape
func (s *Session) GetOutputShape() []int64 {
	return []int64{1, 1, 1024, 1024}
}

// Preprocessor handles image preprocessing
type Preprocessor struct {
	targetSize int
	mean       []float32
	std        []float32
}

// NewPreprocessor creates a new preprocessor
func NewPreprocessor(targetSize int, mean, std []float32) *Preprocessor {
	return &Preprocessor{targetSize: targetSize, mean: mean, std: std}
}

// Preprocess resizes and normalizes the image for ONNX input
func (p *Preprocessor) Preprocess(img image.Image) ([]float32, error) {
	// Resize to target size
	resized := resizeImage(img, p.targetSize, p.targetSize)

	// Convert to float32 array in NCHW format
	inputData := make([]float32, 3*p.targetSize*p.targetSize)

	// Get RGBA data
	rgba, ok := resized.(*image.RGBA)
	if !ok {
		// Convert to RGBA if needed
		rgba = convertToRGBA(resized)
	}

	// Fill input tensor in NCHW format (RGB channels)
	for y := 0; y < p.targetSize; y++ {
		for x := 0; x < p.targetSize; x++ {
			idx := y*p.targetSize + x
			offset := rgba.PixOffset(x, y)
			r := float32(rgba.Pix[offset]) / 255.0
			g := float32(rgba.Pix[offset+1]) / 255.0
			b := float32(rgba.Pix[offset+2]) / 255.0

			// Normalize with mean and std
			inputData[0*p.targetSize*p.targetSize+idx] = (r - p.mean[0]) / p.std[0] // R
			inputData[1*p.targetSize*p.targetSize+idx] = (g - p.mean[1]) / p.std[1] // G
			inputData[2*p.targetSize*p.targetSize+idx] = (b - p.mean[2]) / p.std[2] // B
		}
	}

	return inputData, nil
}

// resizeImage resizes image to target size using nearest neighbor interpolation
func resizeImage(img image.Image, width, height int) image.Image {
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	srcBounds := img.Bounds()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * srcBounds.Dx() / width
			srcY := y * srcBounds.Dy() / height
			newImg.Set(x, y, img.At(srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY))
		}
	}

	return newImg
}

// convertToRGBA converts any image to RGBA
func convertToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	return rgba
}

// GetOriginalDimensions returns original image dimensions
func GetOriginalDimensions(img image.Image) (int, int) {
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

// CalculateScale calculates resize dimensions maintaining aspect ratio
func CalculateScale(origWidth, origHeight, targetSize int) (int, int) {
	if origWidth == 0 || origHeight == 0 {
		return targetSize, targetSize
	}
	scale := float64(targetSize) / float64(max(origWidth, origHeight))
	return int(float64(origWidth) * scale), int(float64(origHeight) * scale)
}

// Postprocessor handles output postprocessing
type Postprocessor struct{}

// NewPostprocessor creates a new postprocessor
func NewPostprocessor() *Postprocessor {
	return &Postprocessor{}
}

// Postprocess converts the ONNX output mask to a transparent PNG
func (p *Postprocessor) Postprocess(outputData []float32, originalImg image.Image, origWidth, origHeight int) ([]byte, error) {
	// Create mask image (1024x1024)
	mask := image.NewGray(image.Rect(0, 0, 1024, 1024))

	// Find min/max for normalization
	minVal := float32(math.MaxFloat32)
	maxVal := float32(-math.MaxFloat32)

	for _, v := range outputData {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	// Normalize and set mask pixels
	rangeVal := maxVal - minVal
	if rangeVal < 1e-6 {
		rangeVal = 1
	}

	for y := 0; y < 1024; y++ {
		for x := 0; x < 1024; x++ {
			idx := y*1024 + x
			normalized := (outputData[idx] - minVal) / rangeVal
			maskValue := uint8(normalized * 255)
			mask.Set(x, y, color.Gray{maskValue})
		}
	}

	// Resize mask to original dimensions
	resizedMask := resizeMask(mask, origWidth, origHeight)

	// Create RGBA image with alpha from mask
	rgba := image.NewNRGBA(image.Rect(0, 0, origWidth, origHeight))

	originalRGBA := convertToRGBA(originalImg)

	for y := 0; y < origHeight; y++ {
		for x := 0; x < origWidth; x++ {
			origOffset := originalRGBA.PixOffset(x, y)
			rgbaOffset := y*rgba.Stride + x*4

			// Set RGB from original
			rgba.Pix[rgbaOffset] = originalRGBA.Pix[origOffset]
			rgba.Pix[rgbaOffset+1] = originalRGBA.Pix[origOffset+1]
			rgba.Pix[rgbaOffset+2] = originalRGBA.Pix[origOffset+2]

			// Set alpha from mask
			grayVal := resizedMask.GrayAt(x, y).Y
			rgba.Pix[rgbaOffset+3] = grayVal
		}
	}

	// Encode to PNG
	var buf bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(&buf, rgba); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
}

// resizeMask resizes the mask to target dimensions
func resizeMask(mask *image.Gray, width, height int) *image.Gray {
	newMask := image.NewGray(image.Rect(0, 0, width, height))

	srcBounds := mask.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * srcWidth / width
			srcY := y * srcHeight / height
			newMask.Set(x, y, mask.GrayAt(srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY))
		}
	}

	return newMask
}
