package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	// Credenciales de AWS
	region := "us-east-1" // Cambia esto a tu región deseada
	accessKeyID := "TU_ACCESS_KEY_ID"
	secretAccessKey := "TU_SECRET_ACCESS_KEY"
	bucketName := "TU_BUCKET_NAME"

	// Inicializar una nueva sesión de AWS
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		fmt.Println("Error al crear la sesión de AWS:", err)
		return
	}

	// Crear un cliente de S3
	s3Client := s3.New(sess)

	// Nombre del archivo local y ruta de la imagen a subir
	localFilePath := "./test.png" // Cambia esto a la ruta de tu imagen local

	// Nombre del archivo en el bucket de S3
	s3FileName := "test-go.png" // Cambia esto al nombre que desees para el archivo en S3

	// Leer el archivo local
	file, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Configuración de la subida del archivo
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName), //bucket de S3
		Key:    aws.String(s3FileName),
		Body:   file,
	}

	// Subir el archivo a S3
	_, err = s3Client.PutObject(params)
	if err != nil {
		fmt.Println("Error al subir el archivo a S3:", err)
		return
	}

	fmt.Println("La imagen se ha subido correctamente a S3.")
	fmt.Printf("URL: https://%s.s3.amazonaws.com/%s\n", bucketName, s3FileName)

}
