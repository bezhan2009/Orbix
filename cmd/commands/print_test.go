package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/common-nighthawk/go-figure"
)

func capturePrintOutput(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	os.Stdout = writer

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close stdout writer: %v", err)
	}

	os.Stdout = originalStdout

	var output bytes.Buffer
	if _, err := io.Copy(&output, reader); err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}

	return output.String()
}

func TestPrintSupportsSplitFontArgument(t *testing.T) {
	output := capturePrintOutput(t, func() {
		Print([]string{"ORBIX", "font", "3d"})
	})

	expected := figure.NewFigure("ORBIX", "larry3d", true).String()

	if !strings.Contains(output, expected) {
		t.Fatalf("expected 3d output, got %q", output)
	}

	if strings.Contains(output, "font 3d") {
		t.Fatalf("font argument leaked into output: %q", output)
	}
}

func TestPrintSupportsEqualsFontArgument(t *testing.T) {
	output := capturePrintOutput(t, func() {
		Print([]string{"ORBIX", "font=3d"})
	})

	expected := figure.NewFigure("ORBIX", "larry3d", true).String()

	if !strings.Contains(output, expected) {
		t.Fatalf("expected 3d output, got %q", output)
	}
}

func TestPrintSupports2DFont(t *testing.T) {
	output := capturePrintOutput(t, func() {
		Print([]string{"ORBIX", "font", "2d"})
	})

	expected := figure.NewFigure("ORBIX", "", true).String()

	if !strings.Contains(output, expected) {
		t.Fatalf("expected 2d output, got %q", output)
	}
}

func TestPrintSupportsColorWithSplitFontArgument(t *testing.T) {
	output := capturePrintOutput(t, func() {
		Print([]string{"cyan:ORBIX", "font", "3d"})
	})

	expected := figure.NewFigure("ORBIX", "larry3d", true).String()

	if !strings.Contains(output, expected) {
		t.Fatalf("expected colored 3d output, got %q", output)
	}

	if strings.Contains(output, "font 3d") {
		t.Fatalf("font argument leaked into output: %q", output)
	}
}

func TestPrintSupportsAnimateArgument(t *testing.T) {
	output := capturePrintOutput(t, func() {
		Print([]string{"animate", "Hi"})
	})

	if !strings.Contains(output, "Hi") {
		t.Fatalf("expected animated text output, got %q", output)
	}

	if strings.Contains(output, "animate") {
		t.Fatalf("animate argument leaked into output: %q", output)
	}
}
