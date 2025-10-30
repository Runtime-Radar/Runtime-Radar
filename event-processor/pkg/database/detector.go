package database

import (
	"context"

	"github.com/runtime-radar/runtime-radar/event-processor/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// defaultDetectorOrder is used in methods which require deterministic order,
	// so that returned records can be used for things that are order-sensitive (e.g. calculating hashes).
	defaultDetectorOrder = "id"

	detectorWasmHashField = "wasm_hash"
)

type DetectorRepository interface {
	Add(ctx context.Context, ss ...*model.Detector) error
	Delete(ctx context.Context, id string, version uint) error
	GetPage(ctx context.Context, filter, order any, pageSize int, pageNum int, preloadData bool) ([]*model.Detector, error)
	GetCount(ctx context.Context, filter any) (int, error)
	GetAll(ctx context.Context, filter, order any, preloadData bool) ([]*model.Detector, error)
	GetAllBins(ctx context.Context, filter any) ([][]byte, error)
	GetAllHashes(ctx context.Context, filter any) ([]string, error)
}

type DetectorDatabase struct {
	*gorm.DB
}

// Add adds new entry to the database, it can add multiple instances at once.
func (dd *DetectorDatabase) Add(ctx context.Context, ds ...*model.Detector) error {
	if len(ds) == 0 {
		return nil
	}

	return dd.WithContext(ctx).Create(ds).Error
}

// Delete deletes scope if it doesn't have assigned rule. If it has one, ErrScopeHasRule is returned.
func (dd *DetectorDatabase) Delete(ctx context.Context, id string, version uint) error {
	return dd.WithContext(ctx).
		Where(map[string]any{
			"id":      id,
			"version": version,
		}).
		Delete(&model.Detector{}).
		Error
}

func (dd *DetectorDatabase) GetPage(ctx context.Context, filter, order any, pageSize int, pageNum int, preloadData bool) ([]*model.Detector, error) {
	ds := []*model.Detector{}

	if filter == nil {
		filter = ""
	}

	sanitizedOrder, err := sanitizeOrder(order)
	if err != nil {
		return nil, err
	}

	err = dd.preloadData(ctx, preloadData).
		Where(filter).
		Order(sanitizedOrder).
		Limit(pageSize).
		Offset(pageSize * (pageNum - 1)).
		Find(&ds).
		Error

	return ds, err
}

func (dd *DetectorDatabase) GetCount(ctx context.Context, filter any) (int, error) {
	if filter == nil {
		filter = ""
	}

	var count int64

	err := dd.WithContext(ctx).
		Model(&model.Detector{}).
		Where(filter).
		Count(&count).
		Error

	return int(count), err
}

func (dd *DetectorDatabase) GetAll(ctx context.Context, filter, order any, preloadData bool) ([]*model.Detector, error) {
	ds := []*model.Detector{}

	if filter == nil {
		filter = ""
	}

	sanitizedOrder, err := sanitizeOrder(order)
	if err != nil {
		return nil, err
	}

	err = dd.preloadData(ctx, preloadData).
		Where(filter).
		Order(sanitizedOrder).
		Find(&ds).
		Error

	return ds, err
}

func (dd *DetectorDatabase) GetAllBins(ctx context.Context, filter any) ([][]byte, error) {
	ds, err := dd.GetAll(ctx, filter, defaultDetectorOrder, false)
	if err != nil {
		return nil, err
	}

	bins := [][]byte{}
	for _, d := range ds {
		bins = append(bins, d.WasmBinary)
	}

	return bins, nil
}

func (dd *DetectorDatabase) GetAllHashes(ctx context.Context, filter any) ([]string, error) {
	res := []string{}

	if filter == nil {
		filter = ""
	}

	err := dd.WithContext(ctx).
		Model(&model.Detector{}).
		Select(detectorWasmHashField).
		Where(filter).
		Order(defaultDetectorOrder).
		Find(&res).
		Error

	return res, err
}

func (dd *DetectorDatabase) preloadData(ctx context.Context, preloadData bool) *gorm.DB {
	if preloadData {
		return dd.WithContext(ctx).
			// This should load all associations without nested: https://gorm.io/docs/preload.html#Preload-All
			Preload(clause.Associations)
	}
	return dd.WithContext(ctx)
}
