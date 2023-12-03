package distro

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/brandtkeller/mk8s/src/types"
	"github.com/cheggaaa/pb/v3"
)

func DownloadArtifacts(config types.MultiConfig) error {
	// Purpose would be to download the artifacts to the current directory
	// no return besides error required

	switch config.Distro {
	case "rke2":
		return downloadRKE2(config.Arch, config.Version, config.Artifacts)
	// case "k3s":
	// 	return downloadK3s()
	default:
		return fmt.Errorf("unsupported distro %s", config.Distro)
	}

}

func downloadRKE2(arch, version string, artifacts map[string]types.Artifact) error {

	allArtifacts, err := addDefaultArtifacts("rke2", arch, version, artifacts)
	if err != nil {
		return err
	}

	if exist, _ := dirOrFileExists("artifacts/"); !exist {
		// create artifacts directory if it doesn't exist
		err := os.Mkdir("artifacts/", 0755)
		if err != nil {
			return err
		}
	}

	for id, artifact := range allArtifacts {

		if exist, _ := dirOrFileExists("artifacts/" + artifact.Name); !exist {
			fmt.Println("Downloading", id)
			err := DownloadFile(artifact.URL, "artifacts/"+artifact.Name)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Artifact", id, "already exists")
		}

	}

	return nil
}

// func downloadK3s() error {

// 	return nil
// }

func addDefaultArtifacts(distro, arch, version string, artifacts map[string]types.Artifact) (map[string]types.Artifact, error) {

	switch distro {
	case "rke2":
		escapedVer := strings.Replace(version, "+", "%2B", -1)
		baseUrl := "https://github.com/rancher/rke2/releases/download/" + escapedVer + "/"

		if len(artifacts) == 0 {
			artifacts = make(map[string]types.Artifact)
		}

		if _, ok := artifacts["checksums"]; !ok {
			// Images key does not exist
			artifacts["checksums"] = types.Artifact{
				Name: "sha256sum-" + arch + ".txt",
				URL:  baseUrl + "sha256sum-" + arch + ".txt",
			}
		}

		if _, ok := artifacts["binary"]; !ok {
			// Images key does not exist
			artifacts["binary"] = types.Artifact{
				Name: "rke2.linux-" + arch + ".tar.gz",
				URL:  baseUrl + "rke2.linux-" + arch + ".tar.gz",
			}
		}

		if _, ok := artifacts["images"]; !ok {
			// Images key does not exist
			artifacts["images"] = types.Artifact{
				Name: "rke2-images.linux-" + arch + ".tar.zst",
				URL:  baseUrl + "rke2-images.linux-" + arch + ".tar.zst",
			}
		}

		if _, ok := artifacts["installScript"]; !ok {
			// Images key does not exist
			artifacts["installScript"] = types.Artifact{
				Name: "install.sh",
				URL:  "https://get.rke2.io",
			}
		}

		return artifacts, nil
	default:
		return artifacts, fmt.Errorf("unsupported distro: %s", distro)
	}

}

// Check if a file or directory exists
func dirOrFileExists(path string) (bool, error) {
	// Stat returns file information. If there is an error, it will be of type *PathError.
	// Check if the error is nil to determine if the file exists.
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		// File does not exist
		return false, nil
	} else {
		// Some other error occurred (permission issues, etc.)
		return false, err
	}
}

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
