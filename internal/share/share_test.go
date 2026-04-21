package share

import "testing"

func TestDefaultOutputPath_withExt(t *testing.T) {
	t.Parallel()

	got := defaultOutputPath("session.wf")
	want := "session.shared.wf"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestDefaultOutputPath_nestedPath(t *testing.T) {
	t.Parallel()

	got := defaultOutputPath("tmp/run.wf")
	want := "tmp/run.shared.wf"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestDefaultOutputPath_withoutExt(t *testing.T) {
	t.Parallel()

	got := defaultOutputPath("session")
	want := "session.shared.wf"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
