package imgur

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// UploadImage uploads the image to imgur
// image                Can be a binary file, base64 data, or a URL for an image. (up to 10MB)
// album       optional The id of the album you want to add the image to.
//                      For anonymous albums, album should be the deletehash that is returned at creation.
// dtype                The type of the file that's being sent; file, base64 or URL
// title       optional The title of the image.
// description optional The description of the image.
// returns image info, status code of the upload, error
func (client *Client) UploadImage(image []byte, album string, dtype string, title string, description string) (*ImageInfo, int, error) {
	if dtype != "binary" && dtype != "base64" && dtype != "URL" {
		return nil, -1, errors.New("Passed invalid dtype: " + dtype + ". Please use binary/base64/URL.")
	}

	form := url.Values{}

	if dtype == "binary" {
		// TODO test binary
		form.Add("image", base64.StdEncoding.EncodeToString(image))
		form.Add("type", "base64")
	}
	if dtype == "base64" {
		form.Add("image", string(image[:]))
		form.Add("type", "base64")
	}
	if dtype == "URL" {
		form.Add("image", string(image[:]))
		form.Add("type", "URL")
	}

	if album != "" {
		form.Add("album", album)
	}
	if title != "" {
		form.Add("title", title)
	}
	if description != "" {
		form.Add("description", description)
	}
	client.Log.Infof("Form %v\n", form)

	URL := apiEndpoint + "image"
	req, err := http.NewRequest("POST", URL, bytes.NewBufferString(form.Encode()))
	client.Log.Infof("Posting to URL %v\n", URL)
	if err != nil {
		return nil, -1, errors.New("Could create request for " + URL + " - " + err.Error())
	}

	req.Header.Add("Authorization", "Client-ID "+client.ImgurClientID)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, -1, errors.New("Could not post " + URL + " - " + err.Error())
	}
	defer res.Body.Close()

	// Read the whole body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, -1, errors.New("Problem reading the body for " + URL + " - " + err.Error())
	}

	client.Log.Debugf("%v\n", string(body[:]))

	dec := json.NewDecoder(bytes.NewReader(body))
	var img imageInfoDataWrapper
	if err = dec.Decode(&img); err != nil {
		return nil, -1, errors.New("Problem decoding json result from image upload - " + err.Error())
	}

	if !img.Success {
		return nil, img.Status, errors.New("Upload to imgur failed with status: " + strconv.Itoa(img.Status))
	}

	rl, err := extractRateLimits(res.Header)
	if err != nil {
		client.Log.Infof("Problem with extracting reate limits: %v", err)
	} else {
		img.Ii.Limit = rl
	}

	return img.Ii, img.Status, nil
}
