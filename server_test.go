package fmk

import (
	"testing"
)

func TestProcessRequestURIGarden(t *testing.T) {
	root := "/get/"
	uri := "/get/files/index.html"
	p, err := processURI(root, uri)
	if err != nil {
		t.Error(err)
		return
	}
	if p != "files/index.html" {
		t.Errorf("Bad result: %s", p)
		return
	}
}

func TestProcessRequestURIFailDoubleDot(t *testing.T) {
	root := "/get/"
	uri := "/get/files/../../passwd"
	_, err := processURI(root, uri)
	if err == nil {
		t.Error("Did not get error on illegal path")
		return
	}
	t.Log(err)
}

func TestProcessRequestURIFailDoubleSlash(t *testing.T) {
	root := "/get/"
	uri := "/get//etc/passwd"
	_, err := processURI(root, uri)
	if err == nil {
		t.Error("Did not get error on VERY illegal double slash path")
		return
	}
	t.Log(err)
}

func TestProcessRequestURIFailMismatch(t *testing.T) {
	root := "/get/"
	uri := "/files/index.html"
	_, err := processURI(root, uri)
	if err == nil {
		t.Error("Did not get error on mismatched URI")
		return
	}
	t.Log(err)
}
