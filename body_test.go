package webhook

import (
	"bytes"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCaptureBodyRejectsLimitThatCannotReserveOverflowByte(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("body"))
	if _, err := CaptureBody(request, math.MaxInt64); !errors.Is(err, ErrInvalidConfiguration) {
		t.Fatalf("CaptureBody(MaxInt64) error = %v", err)
	}
}

func TestCaptureBodyPreservesExactBytesAndRestoresRequest(t *testing.T) {
	t.Parallel()

	want := []byte{0x00, 0xff, '\r', '\n', 'x'}
	request := &http.Request{
		Body:          io.NopCloser(bytes.NewReader(want)),
		ContentLength: int64(len(want)),
	}

	got, err := CaptureBody(request, int64(len(want)))
	if err != nil {
		t.Fatalf("CaptureBody() error = %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("CaptureBody() = %v, want %v", got, want)
	}
	restored, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatalf("reading restored body: %v", err)
	}
	if !bytes.Equal(restored, want) {
		t.Fatalf("restored body = %v, want %v", restored, want)
	}
}

func TestCaptureBodyRejectsDeclaredOversizeBeforeReading(t *testing.T) {
	t.Parallel()

	body := &observedBody{reader: bytes.NewReader([]byte("payload"))}
	request := &http.Request{Body: body, ContentLength: 8}

	_, err := CaptureBody(request, 7)
	if !errors.Is(err, ErrBodyTooLarge) {
		t.Fatalf("CaptureBody() error = %v, want ErrBodyTooLarge", err)
	}
	if body.reads != 0 {
		t.Fatalf("CaptureBody() performed %d reads before rejecting content length", body.reads)
	}
	if !body.closed {
		t.Fatal("CaptureBody() did not close rejected body")
	}
}

func TestCaptureBodyBoundsUnknownLengthBeforeAllocation(t *testing.T) {
	t.Parallel()

	body := &observedBody{reader: bytes.NewReader([]byte("123456789"))}
	request := &http.Request{Body: body, ContentLength: -1}

	_, err := CaptureBody(request, 7)
	if !errors.Is(err, ErrBodyTooLarge) {
		t.Fatalf("CaptureBody() error = %v, want ErrBodyTooLarge", err)
	}
	if body.bytesRead != 8 {
		t.Fatalf("CaptureBody() read %d bytes, want max+1", body.bytesRead)
	}
	if !body.closed {
		t.Fatal("CaptureBody() did not close oversized body")
	}
}

func TestCaptureBodyRejectsUnsafeArguments(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		request *http.Request
		limit   int64
	}{
		"nil request": {request: nil, limit: 1},
		"nil body":    {request: &http.Request{}, limit: 1},
		"zero limit":  {request: &http.Request{Body: http.NoBody}, limit: 0},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if _, err := CaptureBody(test.request, test.limit); !errors.Is(err, ErrInvalidConfiguration) {
				t.Fatalf("CaptureBody() error = %v, want ErrInvalidConfiguration", err)
			}
		})
	}
}

type observedBody struct {
	reader    *bytes.Reader
	reads     int
	bytesRead int
	closed    bool
}

func (b *observedBody) Read(buffer []byte) (int, error) {
	b.reads++
	count, err := b.reader.Read(buffer)
	b.bytesRead += count

	return count, err
}

func (b *observedBody) Close() error {
	b.closed = true

	return nil
}
