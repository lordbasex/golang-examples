package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func uploadToS3(localPath, bucketName string) error {
	// Configuración de la sesión de AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Reemplaza con tu región
	})
	if err != nil {
		return err
	}

	// Crear un cliente de S3
	svc := s3.New(sess)

	// Crear un uploader de S3
	uploader := s3manager.NewUploaderWithClient(svc)

	// Obtener la lista de archivos en el bucket
	bucketFiles, err := getListOfBucketFiles(svc, bucketName)
	if err != nil {
		return err
	}

	// Obtener la lista de archivos en el directorio local
	localFiles, err := getListOfLocalFiles(localPath)
	if err != nil {
		return err
	}

	// Identificar archivos que deben subirse
	filesToUpload := findFilesToUpload(localFiles, bucketFiles)

	// Subir archivos identificados
	for _, file := range filesToUpload {
		filePath := localPath + file
		key := file // Puedes ajustar el nombre del objeto en S3 según tus necesidades

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   file,
		})

		if err != nil {
			return err
		}

		fmt.Printf("Archivo %s subido con éxito a S3\n", filePath)
	}

	fmt.Println("Backup sincronizado con éxito en el bucket:", bucketName)
	return nil
}

// Obtener la lista de archivos en el bucket
func getListOfBucketFiles(svc *s3.S3, bucketName string) ([]string, error) {
	var files []string

	err := svc.ListObjectsPages(&s3.ListObjectsInput{Bucket: &bucketName}, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		for _, obj := range page.Contents {
			files = append(files, *obj.Key)
		}
		return !lastPage
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// Obtener la lista de archivos en el directorio local
func getListOfLocalFiles(localPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(localPath, path)
			if err != nil {
				return err
			}
			files = append(files, strings.ReplaceAll(relPath, `\`, `/`))

		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// Identificar archivos que deben subirse
func findFilesToUpload(localFiles, bucketFiles []string) []string {
	var filesToUpload []string

	// Crear un mapa para buscar archivos en el bucket
	bucketFileMap := make(map[string]bool)
	for _, file := range bucketFiles {
		bucketFileMap[file] = true
	}

	// Identificar archivos que no están en el bucket
	for _, file := range localFiles {
		if _, exists := bucketFileMap[file]; !exists {
			filesToUpload = append(filesToUpload, file)
		}
	}

	return filesToUpload
}

func main() {
	localPath := "/tmp/bkp/"           // Reemplaza con tu ruta local
	bucketName := "tu-bucket-s3" // Reemplaza con el nombre de tu bucket en S3

	err := uploadToS3(localPath, bucketName)
	if err != nil {
		fmt.Println("Error al subir el backup a S3:", err)
		os.Exit(1)
	}
}
