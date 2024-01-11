package states

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ArtificialLegacy/imgscal/modules/cli"
	"github.com/ArtificialLegacy/imgscal/modules/esrgan"
	statemachine "github.com/ArtificialLegacy/imgscal/modules/state_machine"
)

func download(transition func(to statemachine.CliState) error) {
	pwd, _ := os.Getwd()

	out, err := os.Create(fmt.Sprintf("%s\\%s.zip", pwd, os.Getenv("ESRGAN_FOLDER_NAME")))
	defer out.Close()
	if err != nil {
		transition(statemachine.ESRGAN_FAIL)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s%s.zip", os.Getenv("ESRGAN_DOWNLOAD_URL"), os.Getenv("ESRGAN_FOLDER_NAME")))
	if err != nil {
		transition(statemachine.ESRGAN_FAIL)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		transition(statemachine.ESRGAN_FAIL)
		return
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		transition(statemachine.ESRGAN_FAIL)
		return
	}
}

func unzip() error {
	pwd, _ := os.Getwd()

	reader, unzipResult := zip.OpenReader(fmt.Sprintf("%s\\%s.zip", pwd, os.Getenv("ESRGAN_FOLDER_NAME")))
	defer reader.Close()

	for _, file := range reader.File {
		switch file.Name {
		case "input.jpg":
			continue
		case "input2.jpg":
			continue
		case "onepiece_demo.mp4":
			continue
		case "README_windows.md":
			continue
		}

		zipped, _ := file.Open()
		defer zipped.Close()

		path := filepath.Join(fmt.Sprintf("%s\\esrgan-tool", pwd), file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		writer, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		defer writer.Close()

		io.Copy(writer, zipped)
	}

	return unzipResult
}

var esrganDownloadEnter statemachine.StateEnterFunction = func(from statemachine.CliState, transition func(to statemachine.CliState) error) {
	println("-------- Downloading ESRGAN --------")

	pwd, _ := os.Getwd()

	download(transition)

	println(fmt.Sprintf("%sDownloaded ESRGAN. Unzipping...%s", cli.GREEN, cli.RESET))

	mkerr := os.Mkdir(fmt.Sprintf("%s\\esrgan-tool", pwd), 0777)
	if mkerr != nil {
		transition(statemachine.ESRGAN_FAIL)
		return
	}

	unzipResult := unzip()

	rmerr := os.Remove(fmt.Sprintf("%s\\%s.zip", pwd, os.Getenv("ESRGAN_FOLDER_NAME")))
	if rmerr != nil {
		println(rmerr.Error())
	}

	if unzipResult != nil {
		println(unzipResult.Error())
		esrgan.Remove()
		transition(statemachine.ESRGAN_FAIL)
		return
	}

	transition(statemachine.LANDING_MENU)
}

var ESRGANDownload = statemachine.NewState(statemachine.ESRGAN_DOWNLOAD, esrganDownloadEnter, nil, []statemachine.CliState{statemachine.LANDING_MENU, statemachine.ESRGAN_FAIL})
