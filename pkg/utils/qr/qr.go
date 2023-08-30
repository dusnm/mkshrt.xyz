package qr

import (
	"bytes"
	"encoding/base64"
	"io"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type (
	nopCloser struct {
		io.Writer
	}
)

func (nopCloser) Close() error {
	return nil
}

func Generate(data string) ([]byte, error) {
	q, err := qrcode.NewWith(
		data,
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest),
	)

	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(nil)
	c := nopCloser{Writer: b}
	w := standard.NewWithWriter(
		c,
		standard.WithQRWidth(40),
		standard.WithBuiltinImageEncoder(standard.JPEG_FORMAT),
	)

	if err = q.Save(w); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func GenerateWithBase64Encoding(data string) (string, error) {
	b, err := Generate(data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
