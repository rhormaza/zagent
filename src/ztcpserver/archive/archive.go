package archive

import (
  "archive/zip"
  "io"
  "log"
  "os"
  "bytes"
  "io/ioutil"
  "path/filepath"
)


func UnzipMem(arc []byte, size int64) error{
  var (
    err error
    file *os.File
    rc io.ReadCloser
    r *zip.Reader
  )
  ra := bytes.NewReader(arc)

  if r, err = zip.NewReader(ra, size); err != nil {
    log.Fatal(err)
    return err
  }

  for _, f := range r.File {
    if f.Mode().IsDir() { // if file is a directory
      if err = os.MkdirAll(f.Name, f.Mode().Perm()); err != nil {
        return err
      }
      continue
    }

    if rc, err = f.Open(); err != nil {
      return err
    }

    if file, err = os.Create(f.Name); err != nil {
      return err
    }

    if _, err = io.Copy(file, rc); err != nil {
      return err
    }

    if err = file.Chmod( f.Mode().Perm() ); err != nil {
      return err
    }

    file.Close()
    rc.Close()
  }
  return nil
}


func UnzipFile(name string) error{
  file, err := ioutil.ReadFile(name)
  if err != nil {
    return err
  }
  err = UnzipMem(file, int64(len(file)));
  if err != nil {
    return err
  }
  return nil
}



func ZipDir(dir string, zip_name string) error{

  zip_file, err := os.Create(zip_name) 
  if err != nil {
    return err
  }
  w := zip.NewWriter(zip_file)

  addFile := func (path string, info os.FileInfo, err error) error {

    if err != nil {
      return err
    }

    fh, err := zip.FileInfoHeader(info)
    if info.IsDir() {
      fh.Name = path + "/"
    } else {
      fh.Name = path
    }

    if err != nil {
      return err
    }

    out, err := w.CreateHeader(fh);
    if err != nil {
      return err
    }

    if !info.IsDir() { //only copy content when file is regular file
      in, err := os.Open(path)
      if err != nil {
        return err
      }

      if _, err = io.Copy(out, in); err != nil {
        return err
      }
      in.Close()
    }

    return nil
  }

  if err = filepath.Walk(dir, addFile); err != nil {
    return err
  }
  err = w.Close()
  return err
}


