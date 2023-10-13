package main

import (
	"fmt"
	"path"
	"strings"
)

const (
	audioBucketName = "audios"
	fontBucketName  = "fonts"
)

func isAudio(tracker FileTracker) bool {
	return strings.HasSuffix(tracker.Entry.Name(), ".wav") ||
		strings.HasSuffix(tracker.Entry.Name(), ".mp3") ||
		strings.HasSuffix(tracker.Entry.Name(), ".ogg")
}

func isFont(tracker FileTracker) bool {
	return strings.HasSuffix(tracker.Entry.Name(), ".ttf")
}

func isBinaryResource(tracker FileTracker) bool {
	return isAudio(tracker) || isFont(tracker)
}

func detectResourceBucket(tracker FileTracker) (string, error) {
	if isAudio(tracker) {
		return audioBucketName, nil
	}

	if isFont(tracker) {
		return fontBucketName, nil
	}

	return "", fmt.Errorf("unknown resource type")
}

func resourceName(tracker FileTracker) string {
	return strings.TrimSuffix(path.Base(tracker.Entry.Name()),
		path.Ext(tracker.Entry.Name()))
}
