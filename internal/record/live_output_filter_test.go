package record

import (
	"bytes"
	"testing"
)

func TestLiveOutputFilter_passthroughNoMarker(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	f := &liveOutputFilter{}

	f.write(&out, []byte("$ ls\nfile1\n"), false)
	f.write(&out, nil, true)

	got := out.String()
	want := "$ ls\nfile1\n"
	if got != want {
		t.Fatalf("output mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestLiveOutputFilter_stripFullMarkerLine(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	f := &liveOutputFilter{}

	f.write(&out, []byte("$ __TT_RC__:0\n$ ls\n"), false)
	f.write(&out, nil, true)

	got := out.String()
	want := "$ $ ls\n"
	if got != want {
		t.Fatalf("output mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestLiveOutputFilter_markerSplitAcrossChunks(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	f := &liveOutputFilter{}

	f.write(&out, []byte("$ __TT_"), false)
	f.write(&out, []byte("RC__:0\n$ pwd\n"), false)
	f.write(&out, nil, true)

	got := out.String()
	want := "$ $ pwd\n"
	if got != want {
		t.Fatalf("output mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestLiveOutputFilter_incompleteMarkerFlushed(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	f := &liveOutputFilter{}

	f.write(&out, []byte("hello __TT_RC__"), false)
	f.write(&out, nil, true)

	got := out.String()
	want := "hello __TT_RC__"
	if got != want {
		t.Fatalf("output mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestLiveOutputFilter_multipleMarkers(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	f := &liveOutputFilter{}

	f.write(&out, []byte("$ __TT_RC__:0\n"), false)
	f.write(&out, []byte("$ __TT_RC__:1\n"), false)
	f.write(&out, []byte("$ echo hi\nhi\n"), false)
	f.write(&out, nil, true)

	got := out.String()
	want := "$ $ $ echo hi\nhi\n"
	if got != want {
		t.Fatalf("output mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestLongestMarkerPrefixSuffix(t *testing.T) {
	t.Parallel()

	if got := longestMarkerPrefixSuffix("abc"); got != 0 {
		t.Fatalf("got %d want 0", got)
	}
	if got := longestMarkerPrefixSuffix("__TT_"); got != len("__TT_") {
		t.Fatalf("got %d want %d", got, len("__TT_"))
	}
	if got := longestMarkerPrefixSuffix("x__TT_"); got != len("__TT_") {
		t.Fatalf("got %d want %d", got, len("__TT_"))
	}
}
