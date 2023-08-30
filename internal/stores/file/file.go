package file

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io/fs"
	"os"

	storeschema "github.com/compliance-framework/configuration-service/internal/stores/schema"
)

type FileDriver struct {
	Path string
}

func (f *FileDriver) Update(_ context.Context, collection, id string, object interface{}) error {
	// TODO - Implement proper upsert. A method 'MergeFrom' on the BaseModel is needed
	dirPath := f.Path + "/" + collection
	filePath := dirPath + "/" + id + ".gob"
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}
	dataFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	dataEncoder := gob.NewEncoder(dataFile)
	return dataEncoder.Encode(object)
}

func (f *FileDriver) Create(_ context.Context, collection, id string, object interface{}) error {
	dirPath := f.Path + "/" + collection
	filePath := dirPath + "/" + id + ".gob"
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}
	dataFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	dataEncoder := gob.NewEncoder(dataFile)
	return dataEncoder.Encode(object)
}

func (f *FileDriver) Delete(_ context.Context, collection, id string) error {
	dirPath := f.Path + "/" + collection
	filePath := dirPath + "/" + id + ".gob"
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}
	return os.Remove(filePath)
}

func (f *FileDriver) Get(_ context.Context, collection, id string, object interface{}) error {
	dirPath := f.Path + "/" + collection
	filePath := dirPath + "/" + id + ".gob"
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}
	dataFile, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return storeschema.NotFoundErr{}
		}
		return fmt.Errorf("failed to open file: %w", err)
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(object)
	return err
}

func init() {
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
	storeschema.MustRegister("file", &FileDriver{})
}
