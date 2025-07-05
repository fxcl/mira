package upload

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"mira/config"
	"net/textproto"
	"os"
	"strings"
	"time"
)

// Upload file
type Upload struct {
	Config *Config
	File   *File
}

var (
	UploadLocalDriver = "local"
	UploadOssDriver   = "oss"
)

type UploadOption func(*Config)

// Upload configuration
type Config struct {
	Driver     string   // Upload driver
	SavePath   string   // Save path
	UrlPath    string   // URL path
	LimitSize  int      // Limit file size
	LimitType  []string // Limit file type
	RandomName bool     // Use random file name
}

// File information
type File struct {
	FileName    string               // File name
	FileSize    int                  // File size
	FileType    string               // File type
	FileHeader  textproto.MIMEHeader // File header
	FileContent []byte               // File content
}

// Result
type Result struct {
	OriginalName string `json:"originalName"`
	FileName     string `json:"fileName"`
	FileSize     int    `json:"fileSize"`
	FileType     string `json:"fileType"`
	SavePath     string `json:"savePath"`
	UrlPath      string `json:"urlPath"`
	Url          string `json:"url"`
}

// Initialize upload object
func New(options ...UploadOption) *Upload {
	todayPath := time.Now().Format("20060102") + "/"

	// Configure default driver
	config := &Config{
		Driver:     UploadLocalDriver,
		UrlPath:    config.Data.Ruoyi.UploadPath + todayPath,
		SavePath:   config.Data.Ruoyi.UploadPath + todayPath,
		RandomName: false,
	}

	for _, option := range options {
		option(config)
	}

	return &Upload{
		Config: config,
	}
}

// Set upload driver
func SetDriver(driver string) UploadOption {
	return func(config *Config) {
		config.Driver = driver
	}
}

// Set save path
func SetSavePath(savePath string) UploadOption {
	return func(config *Config) {
		config.SavePath = savePath
	}
}

// Set URL path
func SetUrlPath(urlPath string) UploadOption {
	return func(config *Config) {
		config.UrlPath = urlPath
	}
}

// Set limit file size
func SetLimitSize(limitSize int) UploadOption {
	return func(config *Config) {
		config.LimitSize = limitSize
	}
}

// Set limit file type
func SetLimitType(limitType []string) UploadOption {
	return func(config *Config) {
		config.LimitType = limitType
	}
}

// Use random file name
func SetRandomName(isRandomName bool) UploadOption {
	return func(config *Config) {
		config.RandomName = isRandomName
	}
}

// Set upload file
func (u *Upload) SetFile(file *File) *Upload {
	u.File = file

	return u
}

// Save file
func (u *Upload) Save() (*Result, error) {
	var err error
	var domain string

	if config.Data.Ruoyi.Domain == "" {
		return nil, errors.New("domain not found, cannot generate access address")
	}

	if config.Data.Ruoyi.SSL {
		domain = "https://" + config.Data.Ruoyi.Domain
	} else {
		domain = "http://" + config.Data.Ruoyi.Domain
	}

	if u.File == nil || len(u.File.FileContent) <= 0 {
		return nil, errors.New("upload file data is incomplete and cannot be saved")
	}

	// Get the file suffix and generate a hash file name
	fileName := strings.Split(u.File.FileName, ".")
	if len(fileName) != 2 {
		return nil, errors.New("file missing suffix")
	}

	// Splice random file name
	randomName := u.File.FileName
	if u.Config.RandomName {
		randomName = u.generateRandomName() + "." + fileName[1]
	}

	if err = u.checkLimitSize(); err != nil {
		return nil, err
	}

	if err = u.checkLimitType(); err != nil {
		return nil, err
	}

	switch u.Config.Driver {
	case UploadLocalDriver:
		err = u.saveToLocal(randomName)
	case UploadOssDriver:
		err = u.saveToOss()
	default:
		err = u.saveToLocal(randomName)
	}

	if err != nil {
		return nil, err
	}

	return &Result{
		OriginalName: u.File.FileName,
		FileName:     randomName,
		FileSize:     u.File.FileSize,
		FileType:     u.File.FileType,
		SavePath:     u.Config.SavePath,
		UrlPath:      u.Config.UrlPath,
		Url:          domain + "/" + u.Config.UrlPath + randomName,
	}, err
}

// Check file size
func (u *Upload) checkLimitSize() error {
	if u.Config.LimitSize > 0 && u.File.FileSize > 0 && u.Config.LimitSize < u.File.FileSize {
		return errors.New("file size exceeds the limit")
	}

	return nil
}

// Check file type
func (u *Upload) checkLimitType() error {
	if len(u.Config.LimitType) <= 0 || u.File.FileType == "" {
		return nil
	}

	for _, limitType := range u.Config.LimitType {
		if limitType == u.File.FileType {
			return nil
		}
	}

	return errors.New("invalid file format")
}

// Generate random string
func (u *Upload) generateRandomName() string {
	// Create a new random number generator instance
	r := rand.New(rand.NewSource(int64(len(base64.StdEncoding.EncodeToString([]byte(u.File.FileName))))))

	// Define the possible character set, including letters and numbers
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Generate random string
	var randomName string
	for i := 0; i < 64; i++ {
		// Randomly select a character from the character set
		randomChar := chars[r.Intn(len(chars))]
		randomName = randomName + string(randomChar)
	}

	return randomName
}

// Save to local
func (u *Upload) saveToLocal(randomName string) error {
	if _, err := os.Stat(u.Config.SavePath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(u.Config.SavePath, 0o644); err != nil {
				return err
			}
		}
	}

	return os.WriteFile(u.Config.SavePath+randomName, u.File.FileContent, 0o644)
}

// Save to Oss
func (u *Upload) saveToOss() error {
	// TODO

	return nil
}
