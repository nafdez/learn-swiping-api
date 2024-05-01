package picture

import (
	"crypto/rand"
	"encoding/hex"
	"learn-swiping-api/erro"
	"os"
)

var (
	storageDir = "./data/pictures/"
)

func Store(extension string, image []byte) (string, error) {
	if extension != ".png" && extension != ".webp" && extension != ".jpeg" {
		return "", erro.ErrBadField
	}

	os.MkdirAll("./data/pictures/", 0777)

	// TODO: If png/jpeg Then convert to webp
	imageID, err := generateImageID()
	imageID = imageID + extension
	if err != nil {
		return "", err
	}

	f, err := os.Create(storageDir + imageID)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write(image)
	if err != nil {
		return "", err
	}

	return imageID, nil
}

func Picture(imageID string) ([]byte, error) {
	image, err := os.ReadFile(storageDir + imageID)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func Modify(imageID string, image []byte) error {
	return os.WriteFile(imageID, image, 0060)
}

func Remove(imageID string) error {
	return os.Remove(storageDir + imageID)
}

func generateImageID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
