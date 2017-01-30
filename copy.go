package chunked

import "io"

func Copy(dst io.Writer, src io.Reader) (n int64, err error) {
	if w, ok := src.(io.WriterTo); ok {
		return w.WriteTo(dst)
	}
	if r, ok := dst.(io.ReaderFrom); ok {
		return r.ReadFrom(src)
	}
	c := GetLargestChunk()
	c = c[:cap(c)]

	for a, b := 0, 0; err == nil; {
		a, err = src.Read(c)
		if a > 0 {
			if err == nil {
				b, err = dst.Write(c[:a])
			} else {
				b, _ = dst.Write(c[:a])
			}
			n += int64(b)

			if a != b {
				err = io.ErrShortWrite
			}
		}
	}
	if err == io.EOF {
		err = nil
	}
	PutChunk(c)
	return
}

func CopyN(dst io.Writer, src io.Reader, num int64) (n int64, err error) {
	n, err = Copy(dst, io.LimitReader(src, num))
	if n < num && err == nil {
		err = io.EOF
	}
	return
}
