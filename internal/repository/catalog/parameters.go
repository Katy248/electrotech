package catalog

import "electrotech/internal/models"

func (r *Repo) GetParameters() ([]models.Parameter, error) {
	return r.parser.GetParameters()
}
