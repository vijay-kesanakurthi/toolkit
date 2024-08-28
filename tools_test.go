package toolkit

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	randString := testTools.RandomString(10)
	if len(randString) != 10 {
		t.Error("wrong length of random string")
	}
}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{
		"allowed no rename", []string{"image/png"}, false, false,
	},
	{
		"allowed rename", []string{"image/png"}, true, false,
	},
}

func TestTools_UploadFiles(t *testing.T) {
	for _, test := range uploadTests {
		pr, pw := io.Pipe()

		writer := multipart.NewWriter(pw)

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			part, err := writer.CreateFormFile("image", "./testdata/img.png")
			if err != nil {
				t.Error(err)
			}

			img, err := os.Open("./testdata/img.png")
			if err != nil {
				t.Error(err)
			}
			defer img.Close()

			decode, _, err := image.Decode(img)
			if err != nil {
				t.Error(err)
			}

			err = png.Encode(part, decode)
			if err != nil {
				t.Error(err)
			}

		}()

		req, err := http.NewRequest("POST", "/", pr)

		if err != nil {
			t.Error(err)
		}
		req.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools

		testTools.AllowedFileTypes = test.allowedTypes

		files, err := testTools.UploadFiles(req, "./testdata/uploads/", test.renameFile)
		if err != nil && !test.errorExpected {
			t.Error(err)
		}
		if !test.errorExpected {
			_, err2 := os.Stat(fmt.Sprintf("./testdata/uploads/%s", files[0].NewFileName))
			if err2 != nil {
				t.Error(err2)
			}
			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", files[0].NewFileName))

		}
		if !test.errorExpected && err != nil {
			t.Error(err)
		}

	}
}

func TestTools_UploadFile(t *testing.T) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("image", "./testdata/img.png")

		if err != nil {
			t.Errorf("%v cannot create form file", err)
		}
		img, err := os.Open("./testdata/img.png")
		if err != nil {
			t.Errorf("%v cannot open image file", err)
		}
		defer img.Close()

		decode, _, err := image.Decode(img)
		if err != nil {
			t.Error(err)
		}

		err = png.Encode(part, decode)
		if err != nil {
			t.Error(err)
		}

	}()

	req, err := http.NewRequest("POST", "/", pr)

	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	var testTools Tools

	files, err := testTools.UploadFile(req, "./testdata/uploads/", true)
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(fmt.Sprintf("./testdata/uploads/%s", files.NewFileName))
	if err != nil {
		t.Error(err)
	}
	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", files.NewFileName))

	if err != nil {
		t.Error(err)
	}

}
