package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/koffeinsource/go-imgur"
	"github.com/koffeinsource/go-klogger"
)

func main() {
	imgurClientID := flag.String("id", "", "Your imgur client id. REQUIRED!")
	url := flag.String("url", "", "Gets information based on the URL passed.")
	upload := flag.String("upload", "", "Filepath to an image that will be uploaded to imgur.")
	image := flag.String("image", "", "The image ID to be queried.")
	album := flag.String("album", "", "The album ID to be queried.")
	gimage := flag.String("gimage", "", "The gallery image ID to be queried.")
	galbum := flag.String("galbum", "", "The gallery album ID to be queried.")
	rate := flag.Bool("rate", false, "Get the current rate limit.")
	flag.Parse()

	// Check if there is anything todo
	if *imgurClientID == "" || (*image == "" && *album == "" && *gimage == "" && *galbum == "" && *upload == "" && *url == "" && *rate == false) {
		flag.PrintDefaults()
		return
	}

	client := new(imgur.Client)
	client.HTTPClient = new(http.Client)
	client.Log = new(klogger.CLILogger)
	client.ImgurClientID = *imgurClientID

	if *upload != "" {
		client.Log.Infof("*** IMAGE UPLOAD ***\n")
		f, err := os.Open(*upload)
		if err != nil {
			client.Log.Errorf("Could not open file %v - Error: %v", *upload, err)
			return
		}
		defer f.Close()
		fileinfo, err := f.Stat()
		if err != nil {
			client.Log.Errorf("Could not stat file %v - Error: %v", *upload, err)
			return
		}
		size := fileinfo.Size()
		b := make([]byte, size)
		n, err := f.Read(b)
		if err != nil || int64(n) != size {
			client.Log.Errorf("Could not read file %v - Error: %v", *upload, err)
			return
		}

		img, _, err := client.UploadImage(b, "", "binary", "test upload", "test desc")
		if err != nil {
			client.Log.Errorf("Error in UploadImage: %v\n", err)
			return
		}
		client.Log.Infof("%v\n", img)
	}

	if *rate {
		client.Log.Infof("*** RATE LIMIT ***\n")
		rl, err := client.GetRateLimit()
		if err != nil {
			client.Log.Errorf("Error in GetRateLimit: %v\n", err)
			return
		}
		client.Log.Infof("%v\n", *rl)
	}

	if *image != "" {
		client.Log.Infof("*** IMAGE ***\n")
		img, _, err := client.GetImageInfo(*image)
		if err != nil {
			client.Log.Errorf("Error in GetImageInfo: %v\n", err)
			return
		}
		client.Log.Infof("%v\n", img)
	}

	if *album != "" {
		client.Log.Infof("*** ALBUM ***\n")
		img, _, err := client.GetAlbumInfo(*album)
		if err != nil {
			client.Log.Errorf("Error in GetAlbumInfo: %v\n", err)
			return
		}
		client.Log.Infof("%v\n", img)
	}

	if *gimage != "" {
		client.Log.Infof("*** GALLERY IMAGE ***\n")
		img, _, err := client.GetGalleryImageInfo(*gimage)
		if err != nil {
			client.Log.Errorf("Error in GetGalleryImageInfo: %v\n", err)
			return
		}
		client.Log.Infof("%v\n", img)
	}

	if *galbum != "" {
		client.Log.Infof("*** GALLERY ALBUM ***\n")
		img, _, err := client.GetGalleryAlbumInfo(*galbum)
		if err != nil {
			client.Log.Errorf("Error in GetGalleryAlbumInfo: %v\n", err)
			return
		}
		client.Log.Infof("%v\n", img)
	}

	if *url != "" {
		client.Log.Infof("*** URL ***\n")
		img, _, err := client.GetInfoFromURL(*url)
		if err != nil {
			client.Log.Errorf("Error in GetInfoFromURL: %v\n", err)
			return
		}
		client.Log.Infof("%v\n", img)
	}
}
