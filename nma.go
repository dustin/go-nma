// NotifyMyAndroid client for go.
//
// See https://www.notifymyandroid.com/api.jsp for API details.
package nma

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	API_SERVER  = "https://www.notifymyandroid.com"
	VERIFY_PATH = "/publicapi/verify"
	NOTIFY_PATH = "/publicapi/notify"

	VERIFY_URL = API_SERVER + VERIFY_PATH
	NOTIFY_URL = API_SERVER + NOTIFY_PATH
)

type Notification struct {
	Application string
	Description string
	Event       string
	Priority    int
}

type NMA struct {
	apiKey       []string
	developerKey string
	client       *http.Client
}

// Get a new NMA object with the given apiKey
func New(apiKey string) *NMA {
	return NewWithClient(apiKey, http.DefaultClient)
}

// Get a new NMA object with the given apiKey and http.Client
func NewWithClient(apiKey string, client *http.Client) *NMA {
	return &NMA{apiKey: []string{apiKey}, client: client}
}

// Add an API key to the list to try.
func (nma *NMA) AddKey(apiKey string) {
	nma.apiKey = append(nma.apiKey, apiKey)
}

type response struct {
	Err *struct {
		Code       int    `xml:"code,attr"`
		Resettimer int    `xml:"resettimer,attr"`
		Message    string `xml:",chardata"`
	} `xml:"error"`
	Succ *struct {
		Code       int `xml:"code,attr"`
		Remaining  int `xml:"remaining,attr"`
		Resettimer int `xml:"resettimer,attr"`
	} `xml:"success"`
}

func (e *response) Error() string {
	return e.Err.Message
}

func decodeResponse(def string, r io.Reader) (xres response, err error) {
	if xml.NewDecoder(r).Decode(&xres) != nil {
		err = errors.New(def)
	} else {
		if xres.Err != nil {
			err = &xres
		}
	}
	return
}

func (nma *NMA) handleResponse(def string, r io.Reader) error {
	_, err := decodeResponse(def, r)
	if err != nil {
		// Fill response stuff here.
	}
	return err
}

// Verify your credentials.
func (nma *NMA) Verify() (err error) {
	vals := url.Values{"apikey": {strings.Join(nma.apiKey, ",")}}
	var r *http.Response
	r, err = nma.client.Get(NOTIFY_URL + "?" + vals.Encode())
	if err == nil {
		defer r.Body.Close()
		err = nma.handleResponse(r.Status, r.Body)
	}
	return
}

// Send a notification.
func (nma *NMA) Notify(n *Notification) (err error) {

	vals := url.Values{
		"apikey":      {strings.Join(nma.apiKey, ",")},
		"application": {n.Application},
		"description": {n.Description},
		"event":       {n.Event},
	}

	r, err := nma.client.PostForm(NOTIFY_URL, vals)

	if err != nil {
		return
	} else {
		defer r.Body.Close()
		err = nma.handleResponse(r.Status, r.Body)
	}
	return
}
