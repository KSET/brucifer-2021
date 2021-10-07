package repo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/minio/minio-go/v7"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
	minio2 "brucosijada.kset.org/app/providers/minio"
	"brucosijada.kset.org/app/util/async"
)

type imageRepo struct {
}

func Image() imageRepo {
	return imageRepo{}
}

func (r imageRepo) VariationDimensions() []int {
	return []int{
		32,
		256,
		800,
		1200,
		1920,
	}
}

func (r imageRepo) Create(logo *multipart.FileHeader, uploaderId uint) (im *models.Image, err error) {
	if logo.Size > 5*1024*1024 {
		return nil, errors.New("Logo must be smaller than 5MB")
	}

	logoFile, err := logo.Open()
	if err != nil {
		return nil, err
	}
	defer logoFile.Close()

	head := make([]byte, 261)
	if _, err := logoFile.Read(head); err != nil {
		return nil, err
	}
	if _, err := logoFile.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	kind, err := filetype.Match(head)
	if err != nil {
		return nil, err
	}

	if kind.MIME.Type != "image" {
		return nil, errors.New("Logo must be image")
	}
	img, err := imaging.Decode(logoFile)

	if err != nil {
		return nil, err
	}

	mc := minio2.MinioProvider()
	ctx := context.Background()

	uploadTime := time.Now().UTC().UnixNano()
	uploadFolder := fmt.Sprintf("uploads/%d/%d", uploaderId, uploadTime)
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("brucifer-%d-*", uploadTime))
	defer os.RemoveAll(tmpDir)

	variationDimensions := r.VariationDimensions()
	fns := make([]func() (interface{}, error), len(variationDimensions))
	for i, dim := range variationDimensions {
		i, dim := i, dim

		fns[i] = func() (variation interface{}, err error) {
			resizedImage := imaging.Fit(img, dim, dim, imaging.Lanczos)

			imgPath := path.Join(tmpDir, fmt.Sprintf("%d.jpg", dim))
			if err := imaging.Save(resizedImage, imgPath); err != nil {
				return nil, err
			}

			info, err := mc.Client().FPutObject(
				ctx,
				mc.BucketName(),
				fmt.Sprintf("%s/%d.%s", uploadFolder, dim, kind.Extension),
				imgPath,
				minio.PutObjectOptions{
					ContentType: kind.MIME.Value,
				},
			)

			if err == nil {
				variation = models.ImageVariation{
					MinioKey: info.Key,
					Etag:     info.ETag,
					Width:    uint(resizedImage.Bounds().Max.X),
					Height:   uint(resizedImage.Bounds().Max.Y),
					MimeType: kind.MIME.Value,
				}
			}

			return variation, err
		}
	}

	variationsI, err := async.Async().RunInParallel(fns...).All()

	if err != nil {
		return nil, err
	}

	variations := make([]models.ImageVariation, len(variationsI))
	for i, v := range variationsI {
		variations[i] = v.(models.ImageVariation)
	}

	im = &models.Image{
		UploaderID: uploaderId,
		Key:        uploadFolder,
		Variations: variations,
	}

	return im, err
}

func (r imageRepo) GetVariationByUuid(uuid uuid.UUID) (*models.ImageVariation, error) {
	im := models.ImageVariation{}

	err := database.DatabaseProvider().Client().Where(
		"uuid = ?",
		uuid.String(),
	).First(&im).Error

	return &im, err
}
