package product

func toResponse(entity Entity) Response {
	return Response{
		ID:          entity.ID,
		TenantID:    entity.TenantID,
		Name:        entity.Name,
		Description: entity.Description,
		Price:       entity.Price,
		Stock:       entity.Stock,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func toResponseList(entities []Entity) []Response {
	result := make([]Response, 0, len(entities))
	for _, entity := range entities {
		result = append(result, toResponse(entity))
	}
	return result
}
