package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/boltdb/bolt"
)

var (
	binPath          string
	resourceFilePath string
	collectionName   string
)

func parseFlags() {
	flag.StringVar(&binPath, "dir", "./res",
		"Path to the directory where raw files are stored.")
	flag.StringVar(&resourceFilePath, "out", "./stage.res",
		"Resource file to store raw binary resources.")
	flag.StringVar(&collectionName, "collection", "bin",
		"Collection to store the raw resources from the specified directory.")

	flag.Parse()
}

func main() {
	parseFlags()

	// Get binaries from the directory.
	binaries, err := ioutil.ReadDir(binPath)
	handleError(err)
	// Open the resource file.
	resourceFile, err := bolt.Open(resourceFilePath, 0666, nil)
	handleError(err)
	defer resourceFile.Close()

	// Create collections for raw binary resources.
	err = resourceFile.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(collectionName))

		if err != nil {
			return err
		}

		return nil
	})
	handleError(err)

	for _, binaryInfo := range binaries {
		if binaryInfo.IsDir() {
			fmt.Println("Error: directory found in the binaries folder.")
			os.Exit(1)
		}

		binaryFile, err := os.Open(path.Join(binPath, binaryInfo.Name()))
		handleError(err)
		defer binaryFile.Close()
		binary, err := ioutil.ReadAll(binaryFile)
		handleError(err)

		// Save the binary resource to the database.
		binaryID := strings.TrimSuffix(path.Base(binaryInfo.Name()),
			path.Ext(binaryInfo.Name()))
		err = resourceFile.Update(func(tx *bolt.Tx) error {
			buck := tx.Bucket([]byte(collectionName))

			if buck == nil {
				return fmt.Errorf("no %s bucket present", collectionName)
			}

			err = buck.Put([]byte(binaryID), binary)

			if err != nil {
				return err
			}

			return nil
		})
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
