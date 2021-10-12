package repo

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	svgChecker "github.com/h2non/go-is-svg"
	"github.com/minio/minio-go/v7"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
	minio2 "brucosijada.kset.org/app/providers/minio"
	"brucosijada.kset.org/app/util/async"
)

type imageRepo struct {
}

type ImageCreateOptions struct {
	File     *multipart.FileHeader
	Uploader models.User
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

func (r imageRepo) Create(data ImageCreateOptions) (im *models.Image, err error) {
	if data.File.Size > 5*1024*1024 {
		return nil, errors.New("Logo must be smaller than 5MB")
	}

	logoFile, err := data.File.Open()
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

	isImage := kind.MIME.Type == "image"

	isSvg := false
	if !isImage {
		buf, _ := ioutil.ReadAll(logoFile)
		logoFile.Seek(0, io.SeekStart)
		isSvg = svgChecker.IsSVG(buf)
	}

	if !isImage && !isSvg {
		return nil, errors.New("Logo must be image")
	}

	uploadFolder := fmt.Sprintf(
		"uploads/%d/%d",
		data.Uploader.ID,
		time.Now().UTC().UnixNano(),
	)

	var variations *[]models.ImageVariation

	if isSvg {
		variations, err = r.uploadSvg(&logoFile, data.File, uploadFolder)
	} else {
		variations, err = r.uploadImage(&logoFile, &kind, uploadFolder)
	}

	if err != nil {
		return
	}

	im = &models.Image{
		UploaderID: data.Uploader.ID,
		Key:        uploadFolder,
		Variations: *variations,
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

func (r imageRepo) uploadSvg(
	file *multipart.File,
	fileHeader *multipart.FileHeader,
	uploadFolder string,
) (variations *[]models.ImageVariation, err error) {
	mc := minio2.MinioProvider()
	ctx := context.Background()

	dimensions, err := getSvgDimensions(file)
	if err != nil {
		return
	}

	var info minio.UploadInfo
	info, err = mc.Client().PutObject(
		ctx,
		mc.BucketName(),
		fmt.Sprintf("%s/default.svg", uploadFolder),
		*file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: "image/svg+xml",
		},
	)
	if err != nil {
		return
	}

	variation := models.ImageVariation{
		MinioKey: info.Key,
		Etag:     info.ETag,
		Width:    dimensions.Width,
		Height:   dimensions.Height,
		MimeType: "image/svg+xml",
	}

	vars := make([]models.ImageVariation, 1)
	vars[0] = variation
	variations = &vars

	return
}

func (r imageRepo) uploadImage(
	file *multipart.File,
	fileType *types.Type,
	uploadFolder string,
) (variations *[]models.ImageVariation, err error) {
	img, err := imaging.Decode(*file)

	if err != nil {
		return
	}

	mc := minio2.MinioProvider()
	ctx := context.Background()

	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("brucifer-%d-*", time.Now().UTC().UnixNano()))
	defer os.RemoveAll(tmpDir)

	imageDimensions := img.Bounds().Max
	originalRendered := false

	variationDimensions := r.VariationDimensions()
	fns := make([]func() (interface{}, error), 0)
	for _, dim := range variationDimensions {
		dim := dim

		if imageDimensions.X <= dim || imageDimensions.Y <= dim {
			if originalRendered {
				break
			} else {
				originalRendered = true
			}
		}

		fn := func() (variation interface{}, err error) {
			resizedImage := imaging.Fit(img, dim, dim, imaging.Lanczos)

			imgPath := path.Join(tmpDir, fmt.Sprintf("%d.%s", dim, fileType.Extension))
			if err := imaging.Save(resizedImage, imgPath); err != nil {
				return nil, err
			}

			info, err := mc.Client().FPutObject(
				ctx,
				mc.BucketName(),
				fmt.Sprintf("%s/%d.%s", uploadFolder, dim, fileType.Extension),
				imgPath,
				minio.PutObjectOptions{
					ContentType: fileType.MIME.Value,
				},
			)

			if err == nil {
				variation = models.ImageVariation{
					MinioKey: info.Key,
					Etag:     info.ETag,
					Width:    uint(resizedImage.Bounds().Max.X),
					Height:   uint(resizedImage.Bounds().Max.Y),
					MimeType: fileType.MIME.Value,
				}
			}

			return variation, err
		}

		fns = append(fns, fn)
	}

	variationsI, err := async.Async().RunInParallel(fns...).All()

	if err != nil {
		return
	}

	vars := make([]models.ImageVariation, len(variationsI))
	for i, v := range variationsI {
		vars[i] = v.(models.ImageVariation)
	}
	variations = &vars

	return
}

type svgDimensions struct {
	Width  uint
	Height uint
}

func getSvgDimensions(file *multipart.File) (dimensions *svgDimensions, err error) {
	var byteValue []byte

	if byteValue, err = ioutil.ReadAll(*file); err != nil {
		return
	}
	if _, err = (*file).Seek(0, io.SeekStart); err != nil {
		return
	}

	var svgInfo struct {
		XMLName xml.Name `xml:"svg"`
		Height  string   `xml:"height,attr"`
		Width   string   `xml:"width,attr"`
		ViewBox string   `xml:"viewBox,attr"`
	}
	err = xml.Unmarshal(byteValue, &svgInfo)
	if err != nil {
		return
	}

	width, _ := strconv.ParseUint(svgInfo.Width, 10, 64)
	height, _ := strconv.ParseUint(svgInfo.Height, 10, 64)
	dimensions = &svgDimensions{
		Width:  uint(width),
		Height: uint(height),
	}

	if width != 0 && height != 0 {
		return
	}

	viewBox := strings.Fields(svgInfo.ViewBox)

	if len(viewBox) < 4 {
		return nil, errors.New("Invalid viewbox dimensions")
	}

	viewBoxWidth, _ := strconv.ParseFloat(viewBox[2], 64)
	viewBoxHeight, _ := strconv.ParseFloat(viewBox[3], 64)

	if viewBoxWidth <= 0 || viewBoxHeight <= 0 {
		return nil, errors.New("Can't calculate SVG size")
	}

	if width > 0 {
		dimensions.Height = uint(float64(width) * (viewBoxHeight / viewBoxWidth))
		return
	}

	if height > 0 {
		dimensions.Width = uint(float64(height) * (viewBoxWidth / viewBoxHeight))
		return
	}

	viewBoxHeight, viewBoxWidth = normalizeSvgViewBoxDimensions(viewBoxHeight, viewBoxWidth)
	dimensions.Width = uint(viewBoxWidth)
	dimensions.Height = uint(viewBoxHeight)

	return
}

func normalizeSvgViewBoxDimensions(a, b float64) (A, B float64) {
	if a < b {
		return a * 65535 / b, 65535
	}

	return 65535, b * 65535 / a
}
