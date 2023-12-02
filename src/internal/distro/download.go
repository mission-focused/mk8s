package distro

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

func DownloadArtifacts(distro, arch, version string) error {
	// Purpose would be to download the artfiacts to the current directory
	// no return besides error required

	switch distro {
	case "rke2":
		return downloadRKE2(arch, version)
	// case "k3s":
	// 	return downloadK3s()
	default:
		return fmt.Errorf("unsupported distro %s", distro)
	}

}

func downloadRKE2(arch, version string) error {

	escapedVer := strings.Replace(version, "+", "%2B", -1)
	baseUrl := "https://github.com/rancher/rke2/releases/download/" + escapedVer + "/"
	var artifacts = []string{
		baseUrl + "rke2.linux-" + arch + ".tar.gz",
	}

	for _, artifact := range artifacts {
		err := DownloadFile(artifact, "rke2.linux-amd64.tar.gz")
		if err != nil {
			return err
		}
	}

	return nil
}

// func downloadK3s() error {

// 	return nil
// }

// ProgressBarWriter is a custom writer that wraps pb.ProgressBar.
type ProgressBarWriter struct {
	bar *pb.ProgressBar
}

func (p *ProgressBarWriter) Write(data []byte) (int, error) {
	p.bar.Add(len(data))
	return len(data), nil
}

// DownloadFile downloads a file from the given URL and displays progress.
func DownloadFile(url, destination string) error {
	// Create a progress bar
	resp, err := http.Head(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	size := resp.ContentLength
	bar := pb.Full.Start64(size)

	// Wrap the progress bar with a custom writer
	progressWriter := &ProgressBarWriter{bar: bar}

	// Create the file
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create a multi writer to write to file and progress bar simultaneously
	writer := io.MultiWriter(out, progressWriter)

	// Make the GET request
	resp, err = http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy the response body to the file and progress bar
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}

	// Finish the progress bar
	bar.Finish()

	return nil
}
