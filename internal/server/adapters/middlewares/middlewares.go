package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/VoevodinAnton/metrics/internal/pkg/constants"
	"github.com/pkg/errors"
)

type MiddlewareManager interface {
	GzipCompressHandle(next http.Handler) http.Handler
	GzipDecompressHandle(next http.Handler) http.Handler
}

type middlewareManager struct {
}

func NewMiddlewareManager() *middlewareManager {
	return &middlewareManager{}
}

func (mw *middlewareManager) GzipCompressHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var supportsGzip bool
		for _, encodingHeader := range r.Header.Values(constants.AcceptEncodingHeader) {
			if strings.Contains(encodingHeader, constants.GzipEncoding) {
				supportsGzip = true
			}
		}
		if !supportsGzip {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		defer func() {
			_ = gz.Close()
		}()

		w.Header().Set(constants.ContentEncodingHeader, constants.GzipEncoding)
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func (mw *middlewareManager) GzipDecompressHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sendsGzip bool
		for _, encodingHeader := range r.Header.Values(constants.ContentEncodingHeader) {
			if strings.Contains(encodingHeader, constants.GzipEncoding) {
				sendsGzip = true
			}
		}

		if sendsGzip {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer func() {
				_ = reader.Close()
			}()

			r.Body = http.MaxBytesReader(w, reader, http.DefaultMaxHeaderBytes)
		}

		next.ServeHTTP(w, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	return n, errors.Wrap(err, "Writer.Write")
}
