package httpui

import (
	"deploy-lander/internal/structs"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type tagOrBranch struct {
	Tag    string `json:"tag"`
	Branch string `json:"branch"`
}

func parseVersion(req *http.Request) (structs.BuildVersioner, error) {
	defer req.Body.Close()

	contentType := req.Header.Get("Content-Type")

	var value tagOrBranch
	var parseerr error
	if contentType == "application/json" {
		value, parseerr = parseJsonBody(req.Body)
	} else {
		value, parseerr = parseFormBody(req)
	}

	if parseerr != nil {
		return nil, parseerr
	}

	if len(value.Tag) > 0 {
		return structs.NewTag(value.Tag), nil
	}
	if len(value.Branch) > 0 {
		return structs.NewBranch(value.Branch), nil
	}

	return nil, errors.New("neither tag nor branch was specified")
}

func parseFormBody(req *http.Request) (tagOrBranch, error) {
	formerr := req.ParseForm()
	if formerr != nil {
		return tagOrBranch{}, formerr
	}

	return tagOrBranch{
		Tag:    req.PostForm.Get("tag"),
		Branch: req.PostForm.Get("branch"),
	}, nil
}

func parseJsonBody(r io.Reader) (tagOrBranch, error) {
	dec := json.NewDecoder(r)
	var value tagOrBranch
	jsonerr := dec.Decode(&value)
	if jsonerr != nil {
		return tagOrBranch{}, jsonerr
	}

	return value, nil
}
