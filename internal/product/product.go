package product

import (
	"errors"
	"github.com/globalsign/mgo/bson"
)

// WIP
// This should contain specific functionality to deal with products
// (the only domain for our sample)

var (
	// Errors
	ErrNotImplemented = errors.New("not implemented")

	// TODO: Add database handle, which should be initialized from main
)

type Product struct {
	Name string `bson:"name" json:"name" xml:"Name"`
	Desc string `bson:"desc" json:"desc" xml:"Desc"`
}

func FindProduct(id bson.ObjectId) (Product, error) {
	var p Product
	return p, ErrNotImplemented
}

func AddProduct(p *Product) (bson.ObjectId, error) {
	return "", ErrNotImplemented
}

func UpdateProduct(id bson.ObjectId, p *Product) (bool, error) {
	return false, ErrNotImplemented
}

func DeleteProduct(id bson.ObjectId) (bool, error) {
	return false, ErrNotImplemented
}
