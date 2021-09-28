package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
	"github.com/chai2010/webp"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/doc/", func(response http.ResponseWriter, request *http.Request) {
		http.ServeFile(response, request, request.URL.Path[1:])
	})

	mux.HandleFunc("/", processImage)
	log.Fatal(http.ListenAndServe(":80", mux))
}

func getParams(request *http.Request) (float32, float32, bool, bool, bool, bool, bool, int, int, int, int, int, int) {
	var qualityParam, blurParam float32 = 75, 0
	var convertToWebpParam, grayscaleParam, sepiaParam, fliphParam, flipvParam bool = false, false, false, false, false
	var widthParam, heightParam, cropWidthParam, cropHeightParam, cropWidthParamStart, cropHeightParamStart int = 0, 0, 0, 0, 0, 0

	auxQuality, _ := strconv.ParseFloat(request.URL.Query().Get("quality"), 32)
	if auxQuality > 0 {
		qualityParam = float32(auxQuality)
	}

	auxBlur, _ := strconv.ParseFloat(request.URL.Query().Get("blur"), 32)
	if auxBlur > 0 {
		blurParam = float32(auxBlur)
	}

	auxConvertToWebpParam, err := strconv.ParseBool(request.URL.Query().Get("webp"))
	if err == nil {
		convertToWebpParam = auxConvertToWebpParam
	}

	auxGrayscaleParam, err := strconv.ParseBool(request.URL.Query().Get("grayscale"))
	if err == nil {
		grayscaleParam = auxGrayscaleParam
	}

	auxSepiaParam, err := strconv.ParseBool(request.URL.Query().Get("sepia"))
	if err == nil {
		sepiaParam = auxSepiaParam
	}

	auxFliphParam, err := strconv.ParseBool(request.URL.Query().Get("fliph"))
	if err == nil {
		fliphParam = auxFliphParam
	}

	auxFlipvParam, err := strconv.ParseBool(request.URL.Query().Get("flipv"))
	if err == nil {
		flipvParam = auxFlipvParam
	}

	auxWidthParam, err := strconv.Atoi(request.URL.Query().Get("width"))
	if err == nil {
		widthParam = auxWidthParam
	}

	auxHeightParam, err := strconv.Atoi(request.URL.Query().Get("height"))
	if err == nil {
		heightParam = auxHeightParam
	}

	auxCropDimensions := strings.Split(request.URL.Query().Get("crop"), "x")
	if len(auxCropDimensions) > 1 {
		auxCropWidthParamStart, err := strconv.Atoi(auxCropDimensions[0])
		if err == nil {
			cropWidthParamStart = auxCropWidthParamStart
		}
		auxCropHeightParamStart, err := strconv.Atoi(auxCropDimensions[1])
		if err == nil {
			cropHeightParamStart = auxCropHeightParamStart
		}
		auxCropWidthParam, err := strconv.Atoi(auxCropDimensions[2])
		if err == nil {
			cropWidthParam = auxCropWidthParam
		}
		auxCropHeightParam, err := strconv.Atoi(auxCropDimensions[3])
		if err == nil {
			cropHeightParam = auxCropHeightParam
		}
	}

	return qualityParam, blurParam, convertToWebpParam, grayscaleParam, sepiaParam, fliphParam, flipvParam, widthParam, heightParam, cropHeightParam, cropWidthParam, cropHeightParamStart, cropWidthParamStart
}

func processImage(response http.ResponseWriter, request *http.Request) {

	if request.URL.Path == "/" {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	qualityParam, blurParam, convertToWebpParam, grayscaleParam, sepiaParam, fliphParam, flipvParam, widthParam, heightParam, cropHeightParam, cropWidthParam, cropHeightParamStart, cropWidthParamStart := getParams(request)

	// Decode image
	decodedImage, contentType := decodeImageContent(request)
	if convertToWebpParam {
		contentType = "image/webp"
	}

	// Resize image
	newImage := decodedImage
	if widthParam > 0 || heightParam > 0 {
		newImage = resizeImage(decodedImage, widthParam, heightParam)
	}

	// Crop Image
	if cropHeightParam > 0 && cropWidthParam > 0 && cropHeightParamStart > 0 && cropWidthParamStart > 0 {
		newImage = cropImage(newImage, cropWidthParam, cropHeightParam, cropWidthParamStart, cropHeightParamStart)
	}

	// Apply image effects
	if grayscaleParam || blurParam > 0 || sepiaParam || fliphParam || flipvParam {
		newImage = applyEffects(newImage, grayscaleParam, sepiaParam, fliphParam, flipvParam, float64(blurParam))
	}

	// New image to response
	imageBuffer := encodeImage(contentType, newImage, qualityParam)

	// Response with new image
	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Lenght", strconv.Itoa(len(imageBuffer.Bytes())))
	response.Header().Set("Content-Type", contentType)
	response.Header().Set("Pragma", "no-cache")
	response.Write(imageBuffer.Bytes())
}

func encodeImage(newContentType string, decodedImage image.Image, qualityParam float32) bytes.Buffer {
	var imageBuffer bytes.Buffer

	if strings.Contains(newContentType, "jpeg") {
		if err := jpeg.Encode(&imageBuffer, decodedImage, &jpeg.Options{Quality: int(qualityParam)}); err != nil {
			log.Printf("Error encoding JPG Image: %s", err)
			return bytes.Buffer{}
		}
	} else if strings.Contains(newContentType, "png") {
		if err := png.Encode(&imageBuffer, decodedImage); err != nil {
			log.Printf("Error encoding PNG Image: %s", err)
			return bytes.Buffer{}
		}
	} else if strings.Contains(newContentType, "gif") {
		if err := gif.Encode(&imageBuffer, decodedImage, nil); err != nil {
			log.Printf("Error encoding Gif Image: %s", err)
			return bytes.Buffer{}
		}
	} else if strings.Contains(newContentType, "webp") {
		if err := webp.Encode(&imageBuffer, decodedImage, &webp.Options{Lossless: false, Quality: qualityParam, Exact: true}); err != nil {
			log.Printf("Error encodgin wepb image: %s", err)
			return bytes.Buffer{}
		}
	}

	return imageBuffer
}

func decodeImageContent(request *http.Request) (image.Image, string) {
	var originalImage image.Image

	var imageResponse, err = http.Get(getOriginalImage(request.URL.Path))
	if err != nil {
		log.Printf("Error getting original image: %s", err)
	}

	var contentType = imageResponse.Header.Get("Content-type")
	data, err := ioutil.ReadAll(imageResponse.Body)
	if err != nil {
		log.Printf("Error reading original image: %s", err)
		return originalImage, ""
	}

	if strings.Contains(contentType, "jpeg") {
		originalImage, _ = jpeg.Decode(bytes.NewReader(data))
	} else if strings.Contains(contentType, "png") {
		originalImage, _ = png.Decode(bytes.NewReader(data))
	} else if strings.Contains(contentType, "gif") {
		originalImage, _ = gif.Decode(bytes.NewReader(data))
	}

	if originalImage == nil {
		return originalImage, ""
	}

	return originalImage, contentType
}

func calcFactors(width, height uint, oldWidth, oldHeight float64) (scaleX, scaleY float64) {
	if width == 0 {
		if height == 0 {
			scaleX = 1.0
			scaleY = 1.0
		} else {
			scaleY = oldHeight / float64(height)
			scaleX = scaleY
		}
	} else {
		scaleX = oldWidth / float64(width)
		if height == 0 {
			scaleY = scaleX
		} else {
			scaleY = oldHeight / float64(height)
		}
	}
	return
}

func resizeImage(newImage image.Image, width int, height int) image.Image {

	if width == 0 || height == 0 {
		uWidth := uint(width)
		uHeight := uint(height)

		scaleX, scaleY := calcFactors(uWidth, uHeight, float64(newImage.Bounds().Dx()), float64(newImage.Bounds().Dy()))
		if uWidth == 0 {
			uWidth = uint(0.7 + float64(newImage.Bounds().Dx())/scaleX)
		}
		if height == 0 {
			uHeight = uint(0.7 + float64(newImage.Bounds().Dy())/scaleY)
		}

		width = int(uWidth)
		height = int(uHeight)
	}

	return transform.Resize(newImage, width, height, transform.NearestNeighbor)
}

func cropImage(newImage image.Image, cropWidth int, cropHeight int, cropWidthParamStart int, cropHeightParamStart int) image.Image {
	if cropWidth > newImage.Bounds().Dx() {
		cropWidth = newImage.Bounds().Dx()
	}

	if cropHeight > newImage.Bounds().Dy() {
		cropHeight = newImage.Bounds().Dy()
	}

	if cropWidthParamStart > newImage.Bounds().Dx() {
		cropWidthParamStart = 0
	}

	if cropHeightParamStart > newImage.Bounds().Dy() {
		cropHeightParamStart = 0
	}

	return transform.Crop(newImage, image.Rect(cropWidthParamStart, cropHeightParamStart, cropWidth, cropHeight))
}

func applyEffects(newImage image.Image, grayscaleParam, sepiaParam, fliphParam, flipvParam bool, blurParam float64) image.Image {

	if grayscaleParam {
		newImage = effect.Grayscale(newImage)
	}

	if blurParam > 0 {
		newImage = blur.Gaussian(newImage, blurParam)
	}

	if sepiaParam {
		newImage = effect.Sepia(newImage)
	}

	if fliphParam {
		newImage = transform.FlipH(newImage)
	}

	if flipvParam {
		newImage = transform.FlipV(newImage)
	}

	return newImage
}

func getOriginalImage(urlPath string) string {
	return getEnv("imageBaseUrl") + urlPath
}

func getEnv(key string) string {
	var env = os.Getenv(key)
	if env == "" {
		log.Fatalf("Error loading .env file")
	}

	return env
}
