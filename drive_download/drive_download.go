package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

var (
	findingName = flag.String("i", "default", "download file id")
	baseDir     = flag.String("b", "./", "download base dir")
)

func main() {
	flag.Parse()

	if *findingName == "default" {
		fmt.Println("please arg -i")
		return
	}

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	r, err := srv.Files.Get(*findingName).Fields("id, name, webContentLink").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if r == nil {
		fmt.Println("No files found.")
	} else {
		downloadToPath := path.Join(*baseDir, r.Name)
		if r.WebContentLink == "" {
			os.Mkdir(downloadToPath, 0744)
			fileList, err := srv.Files.List().Q(fmt.Sprintf("'%s' in  parents", r.Id)).Do()
			if err != nil {
				log.Fatalf("no files: %v", err)
			}
			for _, i := range fileList.Files {
				downloadChildpath := path.Join(path.Join(*baseDir, r.Name), i.Name)
				err := downloadFromURL(downloadChildpath, i.Id, srv)
				if err != nil {
					log.Fatalf("Unable to retrieve files: %v", err)
				}
				fmt.Printf("Download to %s from %s\n", downloadChildpath, i.WebContentLink)
			}
		} else {
			err = downloadFromURL(downloadToPath, r.Id, srv)
			if err != nil {
				log.Fatalf("Unable to download files: %v", err)
			}
			fmt.Printf("Download to %s from %s\n", downloadToPath, r.WebContentLink)
		}

	}
}

func downloadFromURL(filepath string, fileId string, srv *drive.Service) error {

	resp, err := srv.Files.Get(fileId).Download()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
