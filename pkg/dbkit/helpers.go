package dbkit

import (
	"fmt"
	"strings"
)

const (
	UpdateClause  = "Update"
	DirectionDESC = "DESC"
	DirectionASC  = "ASC"

	_defaultDirection = "DESC"
	_defaultLimit     = 10
	_defaultOffset    = 0
	_defaultPage      = 1
	_defaultSortBy    = "created_at"
)

func PrimaryOrDefaultOrderBy(primaryDirection, primarySortBy string) string {
	direction := PrimaryOrDefaultDirection(primaryDirection)
	sortBy := PrimaryOrDefaultSortBy(primarySortBy)

	return fmt.Sprintf("%s %s", sortBy, direction)
}

func PrimaryOrDefaultOrderByPtr(primaryDirection, primarySortBy *string) string {
	direction := PrimaryOrDefaultDirectionPtr(primaryDirection)
	sortBy := PrimaryOrDefaultSortByPtr(primarySortBy)

	return fmt.Sprintf("%s %s", sortBy, direction)
}

func PrimaryOrDefaultDirection(primary string) string {
	if primary != "" && (primary == DirectionASC || primary == DirectionDESC) {
		return primary
	}

	return _defaultDirection
}

func PrimaryOrDefaultDirectionPtr(primary *string) string {
	if primary != nil && *primary != "" && (strings.EqualFold(*primary, DirectionASC) || strings.EqualFold(*primary, DirectionDESC)) {
		return strings.ToUpper(*primary)
	}

	return _defaultDirection
}

func PrimaryOrDefaultSortBy(primary string) string {
	if primary != "" {
		return primary
	}

	return _defaultSortBy
}

func PrimaryOrDefaultSortByPtr(primary *string) string {
	if primary != nil && *primary != "" {
		return *primary
	}

	return _defaultSortBy
}

func PrimaryOrDefaultLimit(primary uint32) uint32 {
	if primary > 0 {
		return primary
	}

	return _defaultLimit
}

func PrimaryOrDefaultLimitPtr(primary *uint32) uint32 {
	if primary != nil && *primary > 0 {
		return *primary
	}

	return _defaultLimit
}

func PrimaryOrDefaultOffset(primary uint64) uint64 {
	if primary > 0 {
		return primary
	}

	return _defaultOffset
}

func PrimaryOrDefaultOffsetPtr(primary *uint64) uint64 {
	if primary != nil && *primary > 0 {
		return *primary
	}

	return _defaultOffset
}

func PrimaryOrCalculateOffset(primaryOffset uint64, limit, page uint32) uint64 {
	if primaryOffset > 0 {
		return primaryOffset
	}

	return CalculateOffset(limit, page)
}

func PrimaryOrCalculateOffsetPtr(primaryOffset *uint64, limit uint32, page *uint32) uint64 {
	if primaryOffset != nil && *primaryOffset > 0 {
		return *primaryOffset
	}

	var pageValue uint32
	if page != nil {
		pageValue = *page
	}

	return CalculateOffset(limit, pageValue)
}

func CalculateOffset(limit, page uint32) uint64 {
	if limit == 0 {
		limit = _defaultLimit
	}

	if page == 0 {
		page = _defaultPage
	}

	return uint64(limit * (page - 1))
}
