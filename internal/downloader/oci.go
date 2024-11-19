package downloader

import (
	"github.com/chris-cmsoft/concom/internal"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"os"
	"path"
)

type OciDownloader struct {
	source      string
	destination string

	// Reference is the processed OCI Name of the Source
	reference name.Reference
}

func NewOciDownloader(source, destination string) (Downloader, error) {
	reference, err := name.ParseReference(source)
	if err != nil {
		return nil, err
	}
	return &OciDownloader{
		source:      source,
		destination: destination,
		reference:   reference,
	}, nil
}

func (dl *OciDownloader) GetFinalDestination() (string, error) {
	return dl.getOutputDirectory()
}

// Download executes the download of the OCI artifact into memory, untars it and write it to a directory.
// This will need to be updated at some point when we are working with OCI artifacts rather than images,
// to take slightly different actions based on the artifact type we receive from the registry (image / binary / fs)
func (dl *OciDownloader) Download() error {
	img, err := remote.Image(dl.reference, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return err
	}

	outputDirectory, err := dl.getOutputDirectory()
	if err != nil {
		return err
	}

	err = os.MkdirAll(outputDirectory, 0755)
	if err != nil {
		return err
	}

	layers, err := img.Layers()
	for _, layer := range layers {
		layerReader, err := layer.Uncompressed()
		if err != nil {
			return err
		}
		err = internal.Untar(outputDirectory, layerReader)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dl *OciDownloader) getOutputDirectory() (string, error) {
	var outputDirectory string
	if path.IsAbs(dl.destination) {
		outputDirectory = path.Join(dl.destination, dl.reference.Context().RepositoryStr(), dl.reference.Identifier())
	} else {
		workDir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		outputDirectory = path.Join(workDir, dl.destination, dl.reference.Context().RepositoryStr(), dl.reference.Identifier())
	}
	return outputDirectory, nil
}
