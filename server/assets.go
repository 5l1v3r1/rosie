package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/packr"
)

// SetupAssets - Extract or create local assets
// TODO: Add some type of version awareness and a -force option
func SetupAssets() {
	rosieDir := GetRosieDir()
	assetsBox := packr.NewBox("./assets")

	SetupCerts(rosieDir)
	setupGo(rosieDir, assetsBox)
	setupCodenames(rosieDir, assetsBox)
}

// SetupCerts - Creates directories for certs
func SetupCerts(rosieDir string) {
	os.MkdirAll(path.Join(rosieDir, "certs"), os.ModePerm)
	GenerateCertificateAuthority("pivots", true)
	GenerateCertificateAuthority("clients", true)
}

// SetupGo - Unzip Go compiler assets
func setupGo(rosieDir string, assetsBox packr.Box) error {

	log.Printf("Unpacking assets ...")

	// Go compiler and stdlib
	goZip, err := assetsBox.MustBytes(path.Join(runtime.GOOS, "go.zip"))
	if err != nil {
		log.Printf("static asset not found: go.zip")
		return err
	}

	goZipPath := path.Join(rosieDir, "go.zip")
	defer os.Remove(goZipPath)
	ioutil.WriteFile(goZipPath, goZip, 0644)
	_, err = unzip(goZipPath, rosieDir)
	if err != nil {
		log.Printf("Failed to unzip file %s -> %s", goZipPath, rosieDir)
		return err
	}

	goSrcZip, err := assetsBox.MustBytes("src.zip")
	if err != nil {
		log.Printf("static asset not found: src.zip")
		return err
	}
	goSrcZipPath := path.Join(rosieDir, "src.zip")
	defer os.Remove(goSrcZipPath)
	ioutil.WriteFile(goSrcZipPath, goSrcZip, 0644)
	_, err = unzip(goSrcZipPath, path.Join(rosieDir, goDirName))
	if err != nil {
		log.Printf("Failed to unzip file %s -> %s/go", goSrcZipPath, rosieDir)
		return err
	}

	// GOPATH setup
	goPathSrc := path.Join(GetGoPathDir(), "src")
	if _, err := os.Stat(goPathSrc); os.IsNotExist(err) {
		log.Printf("Creating GOPATH directory: %s", goPathSrc)
		os.MkdirAll(goPathSrc, os.ModePerm)
	}

	// Protobuf dependencies
	protobufBox := packr.NewBox("../protobuf")
	pbGoSrc, err := protobufBox.MustBytes("rosie.pb.go")
	if err != nil {
		log.Printf("static asset not found: src.zip")
		return err
	}
	protobufDir := path.Join(goPathSrc, "rosie", "protobuf")
	os.MkdirAll(protobufDir, os.ModePerm)
	ioutil.WriteFile(path.Join(protobufDir, "rosie.pb.go"), pbGoSrc, 0644)

	// GOPATH 3rd party dependencies
	err = unzipGoDependency("github.com.zip", goPathSrc, assetsBox)
	if err != nil {
		log.Fatalf("Failed to unzip go dependency: %v", err)
	}

	return nil
}

func unzipGoDependency(fileName string, goPathSrc string, assetsBox packr.Box) error {
	log.Printf("Unpacking go dependency %s -> %s", fileName, goPathSrc)
	rosieDir := GetRosieDir()
	godep, err := assetsBox.MustBytes(fileName)
	if err != nil {
		log.Printf("static asset not found: %s", fileName)
		return err
	}

	godepZipPath := path.Join(rosieDir, fileName)
	defer os.Remove(godepZipPath)
	ioutil.WriteFile(godepZipPath, godep, 0644)
	_, err = unzip(godepZipPath, goPathSrc)
	if err != nil {
		log.Printf("Failed to unzip file %s -> %s", godepZipPath, rosieDir)
		return err
	}

	return nil
}

func setupCodenames(rosieDir string, assetsBox packr.Box) error {
	nouns, err := assetsBox.MustBytes("nouns.txt")
	adjectives, err := assetsBox.MustBytes("adjectives.txt")

	err = ioutil.WriteFile(path.Join(rosieDir, "nouns.txt"), nouns, 0600)
	if err != nil {
		log.Printf("Failed to write noun data to: %s", rosieDir)
		return err
	}

	err = ioutil.WriteFile(path.Join(rosieDir, "adjectives.txt"), adjectives, 0600)
	if err != nil {
		log.Printf("Failed to write adjective data to: %s", rosieDir)
		return err
	}
	return nil
}

func unzip(src string, dest string) ([]string, error) {

	var filenames []string

	reader, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer reader.Close()

	for _, file := range reader.File {

		rc, err := file.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, file.Name)
		filenames = append(filenames, fpath)

		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return filenames, err
			}
			_, err = io.Copy(outFile, rc)

			outFile.Close()

			if err != nil {
				return filenames, err
			}

		}
	}
	return filenames, nil
}

// copyFileContents - Copy/overwrite src to dst
func copyFileContents(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	out.Sync()
	return nil
}
