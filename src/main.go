package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)


func main() {
	username := os.Getenv("FTP_USERNAME")
	host := os.Getenv("FTP_HOST")
	password := os.Getenv("FTP_PASSWORD")
	localDir := os.Getenv("LOCAL_DIR")
	remoteDir := os.Getenv("REMOTE_DIR")
	deleteExisting := os.Getenv("DELETE_EXISTING")

	var incomp []string
	if username == "" {
		incomp = append(incomp, "username")
	}
	if host == "" {
		incomp = append(incomp, "host")
	}
	if password == "" {
		incomp = append(incomp, "password")
	}
	if localDir == "" {
		incomp = append(incomp, "localDir")
	}
	if remoteDir == "" {
		incomp = append(incomp, "remoteDir")
	}

	 if len(incomp) > 0 {
			log.Fatalf("incompatible options: %s\n", strings.Join(incomp, ", "))
			return
		}

	githubActionPath := os.Getenv("GITHUB_ACTION_PATH")

	if githubActionPath != "" {
		localDir = filepath.Join(githubActionPath, localDir)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatalf("dialing host: %v", err)
		return
	}
	defer client.Close()
	fClient, err := sftp.NewClient(client)
	if err != nil {
		log.Fatalf("sftp client: %v", err)
		return
	}
	defer fClient.Close()

	if strings.ToLower(deleteExisting) == "true" {
		// TODO: Mark files for deletion that aren't tracked by git nor gitignored and delete them
	}

	files := make(map[string]string)
	err = filepath.Walk(localDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		originalPath := path
		path = path[len(localDir):]
		splat := strings.Split(path, string(os.PathSeparator))
		joined := fClient.Join(splat...)
		fmt.Println(originalPath, path, splat, joined, localDir)
		files[originalPath] = joined
		return nil
	})
	if err != nil {
		log.Fatalf("walking local_dir: %v", err)
		return
	}

	for localpath, path := range files {
		log.Printf("opening %s\n", localpath)
		localF, err := os.OpenFile(localpath, os.O_RDONLY, 0644)
		if err != nil {
			log.Fatalf("opening local file: %v\n", err)
			return
		}
		newPath := fClient.Join(remoteDir, path)
		dir, _ := filepath.Split(newPath)
		log.Println("making directories")
		err = fClient.MkdirAll(dir)
		if err != nil {
			log.Fatalf("making directories: %v\n", err)
			return
		}

		log.Printf("creating %s\n", newPath)
		f, err := fClient.Create(newPath)
		if err != nil {
			log.Fatalf("creating file: %v\n", err)
			return
		}
		log.Printf("copying to %s\n", newPath)
		_, err = io.Copy(f, localF)
		if err != nil {
			log.Fatalf("copying file: %v\n", err)
			return
		}
		err = f.Close()
		if err != nil {
			log.Fatalf("closing file: %v\n", err)
			return
		}
		log.Println("done")
	}
}
