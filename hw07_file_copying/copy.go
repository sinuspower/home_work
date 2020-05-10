package main

import (
	"errors"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrGettingFileInfo       = errors.New("error getting file info")
	ErrCanNotCreateFile      = errors.New("can not create file")
	ErrCanNotWriteFile       = errors.New("can not write file")
	ErrNegativeParameter     = errors.New("negative numeric parameter passed")
	ErrSeekingFile           = errors.New("error seeking file")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	if offset < 0 || limit < 0 {
		return ErrNegativeParameter
	}

	in, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer in.Close()

	inFileInfo, err := in.Stat()
	if err != nil {
		return ErrGettingFileInfo
	}
	inSize := inFileInfo.Size()
	if offset > inSize {
		return ErrOffsetExceedsFileSize
	}

	toCopy := inSize - offset
	if limit != 0 && limit < toCopy {
		toCopy = limit
	}
	_, err = in.Seek(offset, io.SeekStart)
	if err != nil {
		return ErrSeekingFile
	}

	out, err := os.Create(toPath)
	if err != nil {
		return ErrCanNotCreateFile
	}
	defer out.Close()

	bar := pb.Full.Start64(toCopy)
	barReader := bar.NewProxyReader(in)
	_, err = io.CopyN(out, barReader, toCopy)
	bar.Finish()

	if err != nil {
		return ErrCanNotWriteFile
	}
	return nil
}
