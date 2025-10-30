package database

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidOrder = errors.New("invalid order")

var (
	allowedOrderFields = []string{
		"name",
		"created_at", // used as default order in service layer
		// notification field
		"integration_type",
	}
	allowedOrderDirections = []string{
		"", // default sort
		"asc",
		"desc",
		"asc nulls first",
		"desc nulls first",
		"asc nulls last",
		"desc nulls last",
	}

	sortsMapping map[string]struct{}
)

func init() {
	sortsMapping = make(map[string]struct{}, len(allowedOrderFields)*len(allowedOrderDirections))

	for _, dir := range allowedOrderDirections {
		for _, sort := range allowedOrderFields {
			order := sort
			if dir != "" {
				order = fmt.Sprintf("%s %s", sort, dir)
			}

			sortsMapping[order] = struct{}{}
		}
	}
}

func mapOrderToSQL(order []string) ([]string, error) {
	res := make([]string, 0, len(order))
	for _, o := range order {
		value := strings.ToLower(o)

		_, ok := sortsMapping[value]
		if !ok {
			return nil, fmt.Errorf("unsupported order value: %s", value)
		}

		res = append(res, value)

	}
	return res, nil
}

func orderToSlice(order any) ([]string, error) {
	ss, ok := order.(string)
	if !ok {
		return nil, fmt.Errorf("unsupported order type: %T", order)
	}

	s := strings.Split(ss, ",")
	res := make([]string, 0, len(s))

	for _, o := range s {
		res = append(res, strings.TrimSpace(o))
	}

	return res, nil
}

func sanitizeOrder(order any) (any, error) {
	if order == nil || order == "" {
		return nil, nil
	}

	orderArr, err := orderToSlice(order)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidOrder, err)
	}

	mapped, err := mapOrderToSQL(orderArr)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidOrder, err)
	}

	return strings.Join(mapped, ","), nil
}
