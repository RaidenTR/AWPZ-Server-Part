///Artificial Intelligence powered by Machine Box

package recognizer

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strconv"

	"github.com/machinebox/sdk-go/facebox"
)

var faceRecognizor *facebox.Client

func InitializeRecognizor() {
	faceRecognizor = facebox.New("http://localhost:8080")
}
func Teach(photo string, id int64) {
	err := faceRecognizor.TeachBase64(photo, strconv.Itoa(int(id)), strconv.Itoa(int(id)))
	fmt.Println(err)
}

// Recognize function was mock without IoT
func Recognize(photo []byte) []int {
	result := make([]int, 0)
	//photoMat := gocv.IMRead("paul1.jpg", gocv.IMReadColor)
	//photoTest := photoMat.ToBytes()
	existingImageFile, _ := os.Open("paul1.jpg")
	defer existingImageFile.Close()
	imageData, _, _ := image.Decode(existingImageFile)
	var buff bytes.Buffer
	jpeg.Encode(&buff, imageData, nil)
	reader := bytes.NewReader(buff.Bytes())
	recognizedPeople, err := faceRecognizor.Check(reader)
	fmt.Println(err)
	for _, person := range recognizedPeople {
		id, _ := strconv.Atoi(person.ID)
		result = append(result, id)
	}
	return result
}

func RecognizeBase64(photo string) []int {
	result := make([]int, 0)
	recognizedPeople, err := faceRecognizor.CheckBase64(photo)
	fmt.Println(err)
	for _, person := range recognizedPeople {
		id, _ := strconv.Atoi(person.ID)
		result = append(result, id)
	}
	return result
}
