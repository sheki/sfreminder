package handlers

import (
	"apikeys"

	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"

	"appengine"
	"appengine/mail"
	"appengine/urlfetch"
)

func init() {
	http.HandleFunc("/catan", catanHandler)
	http.HandleFunc("/handle", handler)
	http.HandleFunc("/prod/handle", prodHandler)
	http.HandleFunc("/prod/datastore", handleStorage)
}

// catanHandler is a test handler returns a test string.
func catanHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Catan Rules")
}

// handler is the main handler called by cron.
func handler(w http.ResponseWriter, r *http.Request) {
	handlerReal(false, w, r)
}

// prodhandler called every weekend sends mail to everyone.
func prodHandler(w http.ResponseWriter, r *http.Request) {
	handlerReal(true, w, r)
}

// handlerReal if prod is set to true send email to every one.
func handlerReal(prod bool, w http.ResponseWriter, r *http.Request) {
	cams, err := reqWebUnderground(appengine.NewContext(r))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	theOne := pickOne(filterSFCams(cams))
	//TODO fix this, use a correct field.
	imgurURL, err := copyToImgur(appengine.NewContext(r), theOne.CurentURL)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	theOne.CurentURL = imgurURL
	if prod {
		sendMail(appengine.NewContext(r), theOne, apikeys.Recipients)
	} else {
		sendMail(appengine.NewContext(r), theOne, apikeys.TestRecipients)
	}
	fmt.Fprintf(w, "Neihborhood : %s Current URL : %s\n", theOne.Neighborhood, theOne.CurentURL)

}

func pickOne(cams []webcams) webcams {
	index := rand.Intn(len(cams))
	return cams[index]
}

var blockedImages = []string{
	"I-80",
	"101",
	"W1",
	"N1",
	"W80",
	"Foo",
}

func filterSFCams(cams []webcams) []webcams {
	var r []webcams
	for _, c := range cams {
		if c.City == "San Francisco" {

			selected := true
			for _, x := range blockedImages {
				if strings.Contains(c.Neighborhood, x) {
					selected = false
				}
			}
			if selected {
				r = append(r, c)
			}
		}
	}
	return r
}

type webcams struct {
	City         string `json:"city"`
	UpdatedEpoch string `json:"updated_epoch"`
	CurentURL    string `json:"CURRENTIMAGEURL"`
	Neighborhood string `json:"neighborhood"`
}

func reqWebUnderground(ctx appengine.Context) ([]webcams, error) {
	curl := `http://api.wunderground.com/api/` + apikeys.WebUnderground + `/webcams/q/CA/San_Francisco.json`
	var b struct {
		W []webcams `json:"webcams"`
	}
	client := urlfetch.Client(ctx)
	resp, err := client.Get(curl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		return nil, err
	}
	return b.W, nil
}

func sendMail(ctx appengine.Context, c webcams, emails []string) {
	msg := &mail.Message{
		Sender:   "San Francisco <abhishek.kona@gmail.com>",
		To:       emails,
		Subject:  "SF weather be like",
		HTMLBody: emailBody(ctx, c),
	}
	if err := mail.Send(ctx, msg); err != nil {
		ctx.Errorf("Couldn't send email: %v", err)
	}
}

func emailBody(ctx appengine.Context, c webcams) string {
	t := template.New("everything")
	var err error
	t, err = t.Parse(body)
	if err != nil {
		ctx.Errorf("Couldn't render template: %v", err)

	}
	var buf bytes.Buffer
	t.Execute(&buf, c)
	return buf.String()

}

const body = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
<body> 
<p>
  Are you wondering how the weather in SF is right now ?
  Well it is great.
  This is how it looks from <b>{{.Neighborhood}}</b>
</p> 
<p>
  <img src="{{.CurentURL}}">
</p>

</body>
</html>
`
