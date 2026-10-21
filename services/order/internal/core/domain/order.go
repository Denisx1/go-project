package domain

import (
	"fmt"
	error_type "go_project/infrastructure/error"
)

type Order struct {
	ID      string
	Version int

	Title    string
	Text     string
	IsActive bool
}

func NewOrder(
	id string,
	version int,
	title string,
	text string,
	isActive bool,
) Order {
	return Order{
		ID:       id,
		Version:  version,
		Title:    title,
		Text:     text,
		IsActive: isActive,
	}
}
func NewOrderUninitialized(
	title string,
	text string,
) Order {
	return NewOrder(
		UninitializedID,
		UninitializedVersion,
		title,
		text,
		UninitializedActive,
	)
}
func (o *Order) Validate() error {
	titleLen := len([]rune(o.Title))
	if titleLen < 10 || titleLen > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			titleLen,
			error_type.ErrInvalidArgument)
	}
	textLen := len([]rune(o.Text))
	if textLen < 10 || textLen > 15 {
		return fmt.Errorf(
			"Invalid `PhoneNumber` len: %d: %w",
			textLen,
			error_type.ErrInvalidArgument,
		)
	}
	return nil
}
