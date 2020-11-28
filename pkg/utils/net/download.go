package net

import (
	"fmt"
	"github.com/AntonBorzenko/RestrictedPassportService/config"
	"io"
	"net/http"
	"os"
	"strings"
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func humanizeBytes(bytes uint64) string {
	return fmt.Sprintf("%.1f MB", float64(bytes) / 1000_000)
}

func (wc WriteCounter) PrintProgress() {
	if !config.Cfg.DBUpdateVerbose {
		return
	}
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanizeBytes(wc.Total))
}

func DownloadFile(filename string, url string) (err error) {
	out, err := os.Create(filename)
	if err != nil {
		return
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, &counter)); err != nil {
		return
	}

	// The progress use the same line so print a new line once it's finished downloading
	if config.Cfg.DBUpdateVerbose {
		fmt.Println()
	}

	return nil
}