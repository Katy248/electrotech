package catalog

import "electrotech/internal/models"

func (r *CatalogRepo) GetParameters() ([]models.Parameter, error) {
	return r.parser.GetParameters()
}
