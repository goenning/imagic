package imagic_test

import (
	"image/color"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goenning/imagic"
)

var parseTestCases = []struct {
	fileName  string
	width     int
	height    int
	supported bool
}{
	{"./testdata/logo1.png", 300, 300, true},
	{"./testdata/logo2.jpg", 2624, 2184, true},
	{"./testdata/logo3.gif", 1165, 822, true},
	{"./testdata/logo4.png", 150, 150, true},
	{"./testdata/logo5.png", 200, 200, true},
	{"./testdata/logo6.jpg", 400, 400, true},
	{"./testdata/logo7.gif", 400, 400, true},
	{"./testdata/favicon.ico", 0, 0, false},
}

func TestImageParse(t *testing.T) {
	for _, testCase := range parseTestCases {
		bytes, err := ioutil.ReadFile(testCase.fileName)
		assert.Nil(t, err)

		file, err := imagic.Parse(bytes)
		if testCase.supported {
			assert.Nil(t, err)
			assert.Equal(t, file.Width, testCase.width)
			assert.Equal(t, file.Height, testCase.height)
			assert.Equal(t, file.Size, len(bytes))
		} else {
			assert.Equal(t, err, imagic.ErrNotSupported)
			assert.Nil(t, file)
		}
	}
}

var resizeTestCases = []struct {
	fileName        string
	resizedFileName string
	size            int
	padding         int
}{
	{"./testdata/logo1.png", "./testdata/logo1-200x200.png", 200, 0},
	{"./testdata/logo2.jpg", "./testdata/logo2-200w.jpg", 200, 0},
	{"./testdata/logo2.jpg", "./testdata/logo2-200w-pad20.jpg", 200, 20},
	{"./testdata/logo3.gif", "./testdata/logo3-200w.gif", 200, 0},
	{"./testdata/logo4.png", "./testdata/logo4-100x100.png", 100, 0},
	{"./testdata/logo5.png", "./testdata/logo5-200x200.png", 200, 0},
	{"./testdata/logo6.jpg", "./testdata/logo6-200x200.jpg", 200, 0},
	{"./testdata/logo7.gif", "./testdata/logo7-200x200.gif", 200, 0},
	{"./testdata/logo7.gif", "./testdata/logo7-1000-1000.gif", 1000, 0},
	{"./testdata/logo8.png", "./testdata/logo8-200h.gif", 200, 0},
}

func TestImageResize(t *testing.T) {
	for _, testCase := range resizeTestCases {
		t.Run(testCase.fileName, func(t *testing.T) {
			bytes, err := ioutil.ReadFile(testCase.fileName)
			assert.Nil(t, err)

			resized, err := imagic.Apply(
				bytes,
				imagic.Resize(testCase.size),
				imagic.Padding(testCase.padding),
			)
			assert.Nil(t, err)

			expected, err := ioutil.ReadFile(testCase.resizedFileName)
			assert.Nil(t, err)

			assert.Equal(t, resized, expected)
		})
	}
}

var bgColorTestCases = []struct {
	fileName           string
	whiteColorFileName string
	bgColor            color.Color
}{
	{"./testdata/logo1.png", "./testdata/logo1-white.png", color.White},
	{"./testdata/logo1.png", "./testdata/logo1-black.png", color.Black},
}

func TestImageChangeBackground(t *testing.T) {
	for _, testCase := range bgColorTestCases {
		bytes, err := ioutil.ReadFile(testCase.fileName)
		assert.Nil(t, err)

		withColor, err := imagic.Apply(bytes, imagic.ChangeBackground(testCase.bgColor))
		assert.Nil(t, err)

		expected, err := ioutil.ReadFile(testCase.whiteColorFileName)
		assert.Nil(t, err)

		assert.Equal(t, withColor, expected)
	}
}
