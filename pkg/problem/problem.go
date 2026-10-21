package problem

import (
	"errors"
	"fmt"
	"path"
)

var _ error = (*ProblemError)(nil)

// ProblemError is implementation of RFC 9457 problem details.
type ProblemError struct {
	Type   string `json:"type"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`

	innerErr error
}

func New(t string, status int, title string) *ProblemError {
	return &ProblemError{
		Type:   t,
		Status: status,
		Title:  title,
		Detail: "",

		innerErr: nil,
	}
}
func NewFrom(err error, t string, status int, title string) *ProblemError {
	return &ProblemError{
		Type:   t,
		Status: status,
		Title:  title,
		Detail: "",

		innerErr: err,
	}
}
func newFromExisting(p *ProblemError) *ProblemError {
	return &ProblemError{
		Type:   p.Type,
		Status: p.Status,
		Title:  p.Title,
		Detail: p.Detail,

		innerErr: p,
	}
}

func (p *ProblemError) WithDetails(details string) *ProblemError {
	newErr := newFromExisting(p)

	if newErr.Detail != "" {
		newErr.Detail = fmt.Sprintf("%s: %s", details, newErr.Detail)
	} else {
		newErr.Detail = details
	}

	return newErr
}

func (p *ProblemError) WithError(err error) *ProblemError {
	newErr := newFromExisting(p)

	newErr.innerErr = errors.Join(newErr.innerErr, err)

	return newErr
}

func (p *ProblemError) Extend(t string, status int, title string) *ProblemError {
	newErr := newFromExisting(p)

	newErr.Type = path.Join(newErr.Type, t)
	newErr.Status = status
	newErr.Title = title

	return newErr
}

func (p *ProblemError) Error() string {
	err := fmt.Sprintf("%q: %s", p.Type, p.Title)
	if p.Detail != "" {
		err = fmt.Sprintf("%s: %s", err, p.Detail)
	}

	return err
}

func (p *ProblemError) Unwrap() error {
	return p.innerErr
}
