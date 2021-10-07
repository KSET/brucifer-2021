package models

import (
	"context"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"

	m "brucosijada.kset.org/app/providers/minio"
	"brucosijada.kset.org/app/util/async"
)

type Image struct {
	Base
	UploaderID uint             `json:"-"`
	Uploader   User             `json:"uploader"`
	Key        string           `json:"-"`
	Variations []ImageVariation `json:"variations"`
	OwnerID    uint             `json:"-"`
	OwnerType  string           `json:"-"`
}

// BeforeDelete hook defined for cascade delete
func (r *Image) BeforeDelete(tx *gorm.DB) (err error) {
	err = tx.Model(&ImageVariation{}).Where(
		"image_id = ?",
		r.ID,
	).Delete(&ImageVariation{}).Error

	if err != nil {
		return
	}

	mp := m.MinioProvider()

	ctx := context.Background()

	objects := mp.Client().ListObjects(
		ctx,
		mp.BucketName(),
		minio.ListObjectsOptions{
			Prefix:    r.Key,
			Recursive: true,
		},
	)

	fns := make([]func() (interface{}, error), 0)
	for obj := range objects {
		key := obj.Key
		fns = append(
			fns, func() (interface{}, error) {
				return nil, mp.Client().RemoveObject(
					ctx,
					m.MinioProvider().BucketName(),
					key,
					minio.RemoveObjectOptions{
						ForceDelete: true,
					},
				)
			},
		)
	}

	_, err = async.Async().RunInParallel(fns...).All()

	return
}
