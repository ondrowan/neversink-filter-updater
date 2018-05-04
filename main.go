package main

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
)

var poePath = filepath.Join(os.Getenv("USERPROFILE"), "Documents/My Games/Path of Exile")
var dotFilePath = filepath.Join(poePath, ".neversink-updater")

func main() {
	var filterType string

	if len(os.Args) > 1 {
		filterType = os.Args[1]
	}

	release, err := getLatestRelease()

	if err != nil {
		exit(1, err.Error())
	}

	currentVersion := getCurrentVersion()

	if *release.TagName == currentVersion {
		exit(0, "There no need to update.")
	}

	zipFile, err := downloadZip(release.GetZipballURL())

	if err != nil {
		exit(1, err.Error())
	}

	tmpArchivePath := createTmpArchive(zipFile)
	unzipArchive(tmpArchivePath, filterType)
	writeToDotfile(*release.TagName)
	showReleaseNotes(release)

	// Clean up

	zipFile.Close()
	os.Remove(tmpArchivePath)

	exit(0, "")
}

func exit(code int, message string) {
	output := os.Stdout

	if code != 0 {
		output = os.Stderr
	}

	if message != "" {
		fmt.Fprintln(output, message)
	}

	fmt.Fprint(os.Stdout, "Press any key to continue...")
	fmt.Scanln()

	os.Exit(0)
}

func getLatestRelease() (*github.RepositoryRelease, error) {
	fmt.Fprint(os.Stdout, "Fetching the latest release... ")

	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(
		context.Background(), "NeverSinkDev", "NeverSink-Filter")

	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stdout, "It is: %s\n", *release.TagName)

	return release, nil
}

func downloadZip(url string) (io.ReadCloser, error) {
	fmt.Fprint(os.Stdout, "Downloading the archive... ")

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Zipfile couldn't be downloaded. This isn't probably your fault. Try again later.")
	}

	fmt.Fprintln(os.Stdout, "Done.")

	return resp.Body, nil
}

func createTmpArchive(content io.ReadCloser) string {
	tmpfile, err := ioutil.TempFile("", "neversink-updater.zip")

	if err != nil {
		log.Fatal(err)
	}

	defer tmpfile.Close()

	_, err = io.Copy(tmpfile, content)

	if err != nil {
		fmt.Println(err)
	}

	return tmpfile.Name()
}

func unzipArchive(path string, filterType string) {
	archiveReader, err := zip.OpenReader(path)

	if err != nil {
		log.Fatal(err)
	}

	defer archiveReader.Close()

	var fileFilter func(*zip.File) bool

	if filterType == "" {
		fileFilter = func(file *zip.File) bool {
			return strings.Count(file.Name, "/") == 1
		}
	} else {
		fileFilter = func(file *zip.File) bool {
			return strings.Split(file.Name, "/")[1] == filterTypeToFolder(filterType)
		}
	}

	for _, archiveFile := range archiveReader.File {
		if strings.HasSuffix(archiveFile.Name, ".filter") && fileFilter(archiveFile) {
			copyFileContent(archiveFile, poePath)
		}
	}
}

func filterTypeToFolder(filterType string) string {
	return fmt.Sprintf("(STYLE) %s", strings.ToUpper(filterType))
}

func copyFileContent(file *zip.File, path string) {
	rc, err := file.Open()

	if err != nil {
		log.Fatal(err)
	}

	filename := strings.Split(file.Name, "/")[1]

	f, err := os.OpenFile(filepath.Join(path, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())

	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(f, rc)

	if err != nil {
		log.Fatal(err)
	}

	rc.Close()
}

func getCurrentVersion() string {
	content, err := ioutil.ReadFile(dotFilePath)

	if err != nil {
		fmt.Fprintln(os.Stdout, "Couldn't determine the latest installed version.")
		return ""
	}

	fmt.Fprintf(os.Stdout, "Your current version is %s. ", content)

	return string(content)
}

func writeToDotfile(version string) {
	content := []byte(version)
	err := ioutil.WriteFile(dotFilePath, content, 0644)

	if err != nil {
		log.Fatal(err)
	}
}

func showReleaseNotes(release *github.RepositoryRelease) {
	fmt.Fprintf(os.Stdout, "\nRelease notes (%s):\n\n%s\n\n", *release.TagName, *release.Body)
}
