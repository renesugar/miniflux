// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package rewrite

import "testing"

func TestRewriteWithNoMatchingRule(t *testing.T) {
	output := Rewriter("https://example.org/article", `Some text.`, ``)
	expected := `Some text.`

	if expected != output {
		t.Errorf(`Not expected output: got "%s" instead of "%s"`, output, expected)
	}
}

func TestRewriteWithYoutubeLink(t *testing.T) {
	output := Rewriter("https://www.youtube.com/watch?v=1234", `Video Description`, ``)
	expected := `<iframe width="650" height="350" frameborder="0" src="https://www.youtube-nocookie.com/embed/1234" allowfullscreen></iframe><p>Video Description</p>`

	if expected != output {
		t.Errorf(`Not expected output: got "%s" instead of "%s"`, output, expected)
	}
}

func TestRewriteWithInexistingCustomRule(t *testing.T) {
	output := Rewriter("https://www.youtube.com/watch?v=1234", `Video Description`, `some rule`)
	expected := `Video Description`
	if expected != output {
		t.Errorf(`Not expected output: got "%s" instead of "%s"`, output, expected)
	}
}

func TestRewriteWithXkcdLink(t *testing.T) {
	description := `<img src="https://imgs.xkcd.com/comics/thermostat.png" title="Your problem is so terrible, I worry that, if I help you, I risk drawing the attention of whatever god of technology inflicted it on you." alt="Your problem is so terrible, I worry that, if I help you, I risk drawing the attention of whatever god of technology inflicted it on you." />`
	output := Rewriter("https://xkcd.com/1912/", description, ``)
	expected := description + `<blockquote cite="https://xkcd.com/1912/">Your problem is so terrible, I worry that, if I help you, I risk drawing the attention of whatever god of technology inflicted it on you.</blockquote>`
	if expected != output {
		t.Errorf(`Not expected output: got "%s" instead of "%s"`, output, expected)
	}
}
func TestRewriteWithXkcdLinkAndNoImage(t *testing.T) {
	description := "test"
	output := Rewriter("https://xkcd.com/1912/", description, ``)
	expected := description
	if expected != output {
		t.Errorf(`Not expected output: got "%s" instead of "%s"`, output, expected)
	}
}

func TestRewriteWithXkcdAndNoImage(t *testing.T) {
	description := "test"
	output := Rewriter("https://xkcd.com/1912/", description, ``)
	expected := description

	if expected != output {
		t.Errorf(`Not expected output: got "%s" instead of "%s"`, output, expected)
	}
}

func TestRewriteWithPDFLink(t *testing.T) {
	description := "test"
	output := Rewriter("https://example.org/document.pdf", description, ``)
	expected := `<a href="https://example.org/document.pdf">PDF</a><br>test`

	if expected != output {
		t.Errorf(`Not expected output: got "%s" instead of "%s"`, output, expected)
	}
}
