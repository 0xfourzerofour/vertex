package persistance

import (
	"bytes"
	"context"
	"govertex/domain/schemas"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	gql "github.com/mattdamon108/gqlmerge/lib"
	"golang.org/x/sync/errgroup"
)

type schemaImp struct {
	s3Conn *s3.S3
}

func SchemaPersistance(s3Conn *s3.S3) schemas.SchemaRepository {
	return &schemaImp{
		s3Conn,
	}
}

func (s *schemaImp) Merge(ctx context.Context, schemaList [][]byte) error {

	eg, _ := errgroup.WithContext(ctx)

	for i, file := range schemaList {
		func(ind int, fileData []byte) {
			eg.Go(func() error {
				err := os.WriteFile("/tmp/schema"+strconv.Itoa(ind)+".graphql", fileData, 0644)

				if err != nil {
					return err
				}

				return nil

			})
		}(i, file)

	}

	if err := eg.Wait(); err != nil {
		return err
	}

	merged := gql.Merge("\t", "/tmp")

	log.Print(*merged)

	s3In := s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("SCHEMA_BUCKET")),
		Key:    aws.String("master.graphql"),
		Body:   bytes.NewReader([]byte(*merged)),
	}

	_, err := s.s3Conn.PutObjectWithContext(ctx, &s3In)

	if err != nil {
		return err
	}

	return nil

}

func (s *schemaImp) GetMaster(ctx context.Context) error {

	return nil

}

func (s *schemaImp) ListSubSchemas(ctx context.Context) ([][]byte, error) {

	eg, _ := errgroup.WithContext(ctx)

	s3In := s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("SUB_SCHEMA_BUCKET")),
	}

	s3List, err := s.s3Conn.ListObjectsV2WithContext(ctx, &s3In)

	if err != nil {
		return nil, err
	}

	schemaArr := make([][]byte, len(s3List.Contents))

	if len(s3List.Contents) > 0 {
		for i, object := range s3List.Contents {
			func(ind int, s3Obj *s3.Object) {

				eg.Go(func() error {

					s3In := s3.GetObjectInput{
						Bucket: aws.String(os.Getenv("SUB_SCHEMA_BUCKET")),
						Key:    s3Obj.Key,
					}

					out, err := s.s3Conn.GetObjectWithContext(ctx, &s3In)

					if err != nil {
						return err
					}

					outBytes, err := ioutil.ReadAll(out.Body)

					if err != nil {
						return err
					}

					schemaArr[ind] = outBytes

					return nil

				})
			}(i, object)

		}
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return schemaArr, nil

}
