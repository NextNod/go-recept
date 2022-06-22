package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func ToInt(data string) (int, error) {
	i, err := strconv.Atoi(data)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func SaveFile(r *http.Request) (string, error) {
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()
	id := uuid.New().String()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", id+"."+name[1])
	f, err1 := os.OpenFile("./images/"+id+"."+name[1], os.O_WRONLY|os.O_CREATE, 0666)
	if err1 != nil {
		return "", err1
	}
	defer f.Close()
	io.Copy(f, file)
	return id + "." + name[1], nil
}
