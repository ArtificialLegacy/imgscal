package states

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ArtificialLegacy/imgscal/pkg/libs/esrgan"
	"github.com/ArtificialLegacy/imgscal/pkg/state/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/utility/cli"
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

	if _, err = io.Copy(out, resp.Body); err != nil {
		transition(statemachine.ESRGAN_FAIL)
		return
	}
}

func unzip() error {
	pwd, _ := os.Getwd()

	reader, err := zip.OpenReader(fmt.Sprintf("%s\\%s.zip", pwd, os.Getenv("ESRGAN_FOLDER_NAME")))
	defer reader.Close()
	if err != nil {
		return err
	}

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

		zipped, err := file.Open()
		defer zipped.Close()

		if err != nil {
			return err
		}

		path := filepath.Join(fmt.Sprintf("%s\\esrgan-tool", pwd), file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		defer writer.Close()

		if err != nil {
			return err
		}

		_, err = io.Copy(writer, zipped)
		if err != nil {
			return err
		}
	}

	return nil
}

var esrganDownloadEnter statemachine.StateStepFunction = func(sm *statemachine.StateMachine) {
	println("-------- Downloading ESRGAN --------")

	pwd, _ := os.Getwd()

	download(sm.Transition)

	println(fmt.Sprintf("%sDownloaded ESRGAN. Unzipping...%s", cli.GREEN, cli.RESET))

	mkerr := os.Mkdir(fmt.Sprintf("%s\\esrgan-tool", pwd), 0777)
	if mkerr != nil {
		sm.Transition(statemachine.ESRGAN_FAIL)
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
		sm.Transition(statemachine.ESRGAN_FAIL)
		return
	}

	sm.Transition(statemachine.LANDING_MENU)
}

var ESRGANDownload = statemachine.NewState(
	statemachine.ESRGAN_DOWNLOAD,
	esrganDownloadEnter,
	[]statemachine.CliState{statemachine.LANDING_MENU, statemachine.ESRGAN_FAIL},
)
