package main

import (
	"io"
	"net/http"
	"os"

	"fmt"

	"strconv"

	"github.com/labstack/echo"
)

const (
	//ImageDir 画像の保存フォルダ
	ImageDir = "image"
	//ImageExt 画像の拡張子
	ImageExt = ".jpg"
)

//RegistRequest 登録リクエストデータ
type RegistRequest struct {
	URL    string `json:"url" xml:"url" form:"url" query:"url"`
	Number int    `json:"number" xml:"number" form:"number" query:"number"`
	Title  string `json:"title" xml:"title" form:"title" query:"title"`
}

//RegistResponce 登録レスポンスデータ
type RegistResponce struct {
	Status int
}

//SelectRequest データ選択リクエストデータ
type SelectRequest struct {
	Number int `json:"number" xml:"number" form:"number" query:"number"`
}

//SelectResponce データ選択レスポンスデータ
type SelectResponce struct {
	Status int
	URL    string
	Number int
	Title  string
}

//RegistHandler 登録ハンドラ
func RegistHandler(c echo.Context) error {
	req := new(RegistRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Printf("request=%v", *req)

	if err := saveURLImage(req.URL, req.Number); err != nil {
		return err
	}

	if err := InsertStory(req.Number, req.Title); err != nil {
		return err
	}

	res := new(RegistResponce)
	res.Status = 0
	return c.JSON(http.StatusOK, res)
}

//saveURLImage 指定した画像URLをファイルとして保存する
func saveURLImage(imageURL string, number int) error {
	if imageURL == "" || number == 0 {
		return fmt.Errorf("registURLImage パラメーターエラー")
	}

	if _, err := os.Stat(ImageDir); err != nil {
		os.Mkdir(ImageDir, 0777)
	}
	imageFilePath := getImageFilePath(ImageDir, number)
	if _, err := os.Stat(imageFilePath); err == nil {
		return nil //すでに保存済み
	}

	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(imageFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	io.Copy(file, response.Body)
	return nil
}

//getImageFilePath 画像の相対パスを取得する
func getImageFilePath(dir string, number int) string {
	return dir + "/" + strconv.Itoa(number) + ImageExt
}

//SelectHandler データ取得ハンドラ
func SelectHandler(c echo.Context) error {
	req := new(SelectRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Printf("request=%v", *req)

	record := SelectStory(req.Number)
	if record.ID == 0 {
		return fmt.Errorf("SelectHandler データなし")
	}

	res := new(SelectResponce)
	res.Status = 0
	res.Number = record.Number
	res.Title = record.Title
	res.URL = c.Scheme() + "://" + c.Request().Host + "/api/image/" + strconv.Itoa(record.Number)
	return c.JSON(http.StatusOK, res)
}

//
func ImageHandler(c echo.Context) error {
	numberText := c.Param("number")
	number, err := strconv.Atoi(numberText)
	if err != nil {
		return c.JSON(http.StatusOK, "NG")
	}
	return c.File(getImageFilePath(ImageDir, number))
}
