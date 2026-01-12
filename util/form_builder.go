package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path"
	"strings"
)

type FormBuilder interface {
	CreateFormFile(fieldname string, file *os.File) error
	CreateFormFileHeader(fieldname string, fileHeader *multipart.FileHeader) error
	CreateFormFileReader(fieldname string, r io.Reader, filename string) error
	WriteField(fieldname, value string) error
	Close() error
	FormDataContentType() string
}

type DefaultFormBuilder struct {
	writer *multipart.Writer
}

func NewFormBuilder(body io.Writer) *DefaultFormBuilder {
	return &DefaultFormBuilder{
		writer: multipart.NewWriter(body),
	}
}

func (fb *DefaultFormBuilder) CreateFormFile(fieldname string, file *os.File) error {
	return fb.createFormFile(fieldname, file, file.Name())
}

func (fb *DefaultFormBuilder) CreateFormFileHeader(fieldname string, fileHeader *multipart.FileHeader) error {
	return fb.createFormFileHeader(fieldname, fileHeader)
}

func (fb *DefaultFormBuilder) CreateFormFileReader(fieldname string, r io.Reader, filename string) error {
	return fb.createFormFile(fieldname, r, path.Base(filename))
}

func (fb *DefaultFormBuilder) createFormFile(fieldname string, r io.Reader, filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	fieldWriter, err := fb.writer.CreateFormFile(fieldname, filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(fieldWriter, r)
	if err != nil {
		return err
	}

	return nil
}
func (fb *DefaultFormBuilder) createFormFileHeader(fieldname string, fileHeader *multipart.FileHeader) error {

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(fieldname), escapeQuotes(fileHeader.Filename)))
	h.Set("Content-Type", fileHeader.Header.Get("Content-Type"))

	fieldWriter, err := fb.writer.CreatePart(h)
	if err != nil {
		return err
	}

	_, err = io.Copy(fieldWriter, file)
	if err != nil {
		return err
	}

	return nil
}

func (fb *DefaultFormBuilder) WriteField(fieldname, value string) error {
	return fb.writer.WriteField(fieldname, value)
}

func (fb *DefaultFormBuilder) Close() error {
	return fb.writer.Close()
}

func (fb *DefaultFormBuilder) FormDataContentType() string {
	return fb.writer.FormDataContentType()
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
