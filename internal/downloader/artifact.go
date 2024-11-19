package downloader

type ArtifactDownloader struct {
	source      string
	destination string
}

func NewArtifactDownloader(source, destination string) Downloader {
	return &ArtifactDownloader{
		source:      source,
		destination: destination,
	}
}

func (dl *ArtifactDownloader) GetFinalDestination() (string, error) {
	return "", nil
}

func (dl *ArtifactDownloader) Download() error {
	return nil
}
