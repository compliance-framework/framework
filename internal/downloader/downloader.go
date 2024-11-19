package downloader

import (
	"regexp"
)

type Downloader interface {
	// Download is expected to download any artifacts, and store them on the local fs.
	// Download may create subDirectory structures, depending on what is being downloaded.
	// ex. The OCI downloader may create subdirectories such as [image]/[version]/[oci-output].
	Download() error

	// GetFinalDestination is expected to return the final directory where the downloaded artifacts are stored for
	// usage.
	// GetFinalDestination() commonly outputs the subDirectory structures which may be created during Download()
	//
	// Common usage of this method will be similar to:
	// exec.Command(dl.GetFinalDestination() + "plugin")
	GetFinalDestination() (string, error)
}

// Download guesses the type of artifact being requested, and calls the appropriate downloader to handle it.
// This is simply a utility method, and is not strictly required.
func Download(source string, destination string) (Downloader, error) {
	var downloader Downloader
	var err error
	if isOci(source) {
		downloader, err = NewOciDownloader(source, destination)
		if err != nil {
			return nil, err
		}
	} else {
		downloader = NewArtifactDownloader(source, destination)
	}
	return downloader, downloader.Download()
}

func isOci(source string) bool {
	// Check whether this looks like an OCI endpoint
	// You can see the verification for the regex at https://regex101.com/r/Z8172m
	r := regexp.MustCompile(`(?i)^((http|https|oci)?:*/*)?([a-zA-Z.]*)+\.([a-zA-Z]*)/([\-_/a-zA-Z]*)(:.*)?$`)
	return r.MatchString(source)
}
