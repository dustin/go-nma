// Package nma is a NotifyMyAndroid client for go.
//
// See https://www.notifymyandroid.com/api.jsp for API details.
package nma

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	apiServer  = "https://www.notifymyandroid.com"
	verifyPath = "/publicapi/verify"
	notifyPath = "/publicapi/notify"

	verifyURL = apiServer + verifyPath
	notifyURL = apiServer + notifyPath
)

// PriorityLevel defines the priority of a notification.
type PriorityLevel int

// Priority levels
const (
	PRIORITY_VERYLOW   PriorityLevel = -2
	PRIORITY_MODERATE                = -1
	PRIORITY_NORMAL                  = 0
	PRIORITY_HIGH                    = 1
	PRIORITY_EMERGENCY               = 2
)

// ContentType specifies the content type of a message.
type ContentType string

// Available content types
const (
	CONTENT_TYPE_HTML ContentType = "text/html"
	CONTENT_TYPE_TEXT             = "text/plain"
)

// A Notification contains all the information to deliver.
type Notification struct {
	Application string
	Description string
	Event       string
	Priority    PriorityLevel
	URL         string
	ContentType ContentType
}

// NMA is the entry point for all API calls.
type NMA struct {
	apiKey       []string
	developerKey string
	client       *http.Client
}

// New gets a new NMA object with the given apiKey
func New(apiKey string) *NMA {
	return NewWithClient(apiKey, http.DefaultClient)
}

// NewWithClient gets a new NMA object with the given apiKey and
// http.Client
func NewWithClient(apiKey string, client *http.Client) *NMA {
	return &NMA{apiKey: []string{apiKey}, client: client}
}

// AddKey adds an API key to the list to try.
func (nma *NMA) AddKey(apiKey string) {
	nma.apiKey = append(nma.apiKey, apiKey)
}

// SetDeveloperKey sets the Developer key for the NMA object
func (nma *NMA) SetDeveloperKey(devKey string) {
	nma.developerKey = devKey
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

func decodeResponse(r io.Reader) (xres response, err error) {
	if err = xml.NewDecoder(r).Decode(&xres); err != nil {
		return response{}, err
	}
	if xres.Err != nil {
		err = &xres
	}
	return
}

func (nma *NMA) handleResponse(def string, r io.Reader) error {
	_, err := decodeResponse(r)
	if err != nil {
		// Fill response stuff here.
	}
	return err
}

// Verify your credentials.
func (nma *NMA) Verify(apikey string) (err error) {
	vals := url.Values{"apikey": {apikey}}

	if nma.developerKey != "" {
		vals["developerkey"] = []string{nma.developerKey}
	}

	var r *http.Response
	r, err = nma.client.Get(verifyURL + "?" + vals.Encode())
	if err == nil {
		defer r.Body.Close()
		err = nma.handleResponse(r.Status, r.Body)
	}
	return
}

// Notify sends a notification.
func (nma *NMA) Notify(n *Notification) (err error) {

	vals := url.Values{
		"apikey":      {strings.Join(nma.apiKey, ",")},
		"application": {n.Application},
		"description": {n.Description},
		"event":       {n.Event},
	}

	if n.Priority != 0 {
		vals["priority"] = []string{strconv.Itoa(int(n.Priority))}
	}

	if n.URL != "" {
		vals["url"] = []string{n.URL}
	}

	if n.ContentType != "" {
		vals["content-type"] = []string{string(n.ContentType)}
	}

	if nma.developerKey != "" {
		vals["developerkey"] = []string{nma.developerKey}
	}

	r, err := nma.client.PostForm(notifyURL, vals)

	if err != nil {
		return
	}

	defer r.Body.Close()
	return nma.handleResponse(r.Status, r.Body)
}
