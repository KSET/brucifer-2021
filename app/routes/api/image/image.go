package image

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	mp "brucosijada.kset.org/app/providers/minio"
	"brucosijada.kset.org/app/repo"
	"brucosijada.kset.org/app/response"
)

func ShowImage(ctx *fiber.Ctx) (err error) {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid ID",
				err,
			),
		)
	}

	model, err := repo.Image().GetVariationByUuid(id)

	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	browserEtag := ctx.Get("If-None-Match")

	if strings.ToLower(browserEtag) == strings.ToLower(model.Etag) {
		return ctx.SendStatus(fiber.StatusNotModified)
	}

	mc := mp.MinioProvider()

	obj, err := mc.Client().GetObject(
		ctx.Context(),
		mc.BucketName(),
		model.MinioKey,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	ctx.Set("Content-Type", model.MimeType)
	ctx.Set("Etag", model.Etag)
	ctx.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int((365*24*time.Hour).Seconds())))

	return ctx.SendStream(obj)
}
