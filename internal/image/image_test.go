package image

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing-demo/internal/logging"
	"testing-demo/internal/middleware"
)

var ts *httptest.Server

func TestMain(m *testing.M) {
	logging.Configure("debug", "test", []string{"stdout"}, []string{"stdout"})

	r := gin.New()
	r.Use(middleware.ContextLogger(logging.GetLogger()))

	apiRouter := r.Group("/api")
	apiRouter.Use(middleware.ErrorHandler())
	RegisterHandlers(apiRouter)

	ts = httptest.NewServer(r)
	defer ts.Close()

	//Send a null string to configure the message queue for api testing, so it does not send mq messages
	exitVal := m.Run()

	//do any additional teardown here
	os.Exit(exitVal)
}

func TestImagePost(t *testing.T) {
	client := &http.Client{}
	uri := fmt.Sprintf("%s/%s", ts.URL, "api/images")

	imgPath := path.Join(os.Getenv("TESTDATA_PATH"), "puppy-asleep.jpg")

	file, err := os.Open(imgPath)
	require.NoError(t, err)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err1 := writer.CreateFormFile("image", filepath.Base(file.Name()))
	require.NoError(t, err1)
	_, err2 := io.Copy(part, file)
	require.NoError(t, err2)

	errWriter := writer.Close()
	require.NoError(t, errWriter)

	req, err := http.NewRequest("POST", uri, body)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer service-user-token")
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	img, imgErr := png.Decode(resp.Body)
	require.NoError(t, imgErr)

	assert.LessOrEqual(t, img.Bounds().Max.X, 256)
	assert.LessOrEqual(t, img.Bounds().Max.Y, 256)
}
