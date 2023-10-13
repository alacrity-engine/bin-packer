package main

import (
	"flag"
	"os"
	"path"

	"github.com/golang-collections/collections/queue"
	bolt "go.etcd.io/bbolt"
)

var (
	projectPath      string
	resourceFilePath string
)

// TODO: create a font atlas packer.
// An atlas consists of a font texture
// and a set of frames (one for each character).

// TODO: create a text packer.

func parseFlags() {
	flag.StringVar(&projectPath, "project", ".",
		"Path to the project to pack raw resources for.")
	flag.StringVar(&resourceFilePath, "out", "./stage.res",
		"Resource file to store raw binary resources.")

	flag.Parse()
}

func main() {
	parseFlags()

	// Open the resource file.
	resourceFile, err := bolt.Open(resourceFilePath, 0666, nil)
	handleError(err)
	defer resourceFile.Close()

	entries, err := os.ReadDir(projectPath)
	handleError(err)

	traverseQueue := queue.New()

	if len(entries) <= 0 {
		return
	}

	for _, entry := range entries {
		traverseQueue.Enqueue(FileTracker{
			EntryPath: projectPath,
			Entry:     entry,
		})
	}

	for traverseQueue.Len() > 0 {
		fsEntry := traverseQueue.Dequeue().(FileTracker)

		if fsEntry.Entry.IsDir() {
			entries, err = os.ReadDir(path.Join(fsEntry.EntryPath, fsEntry.Entry.Name()))
			handleError(err)

			for _, entry := range entries {
				traverseQueue.Enqueue(FileTracker{
					EntryPath: path.Join(fsEntry.EntryPath, fsEntry.Entry.Name()),
					Entry:     entry,
				})
			}

			continue
		}

		if !isBinaryResource(fsEntry) {
			continue
		}

		binData, err := os.ReadFile(path.Join(
			fsEntry.EntryPath, fsEntry.Entry.Name()))
		handleError(err)
		resBucketName, err := detectResourceBucket(fsEntry)
		handleError(err)
		resName := resourceName(fsEntry)

		err = resourceFile.Update(func(tx *bolt.Tx) error {
			buck, err := tx.CreateBucketIfNotExists([]byte(resBucketName))

			if err != nil {
				return err
			}

			err = buck.Put([]byte(resName), binData)

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
