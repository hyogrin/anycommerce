package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"strconv"
	"strings"
)

var imageRootURL string
var missingImageFile = "product_image_coming_soon.png"

// initResponse
func initResponse(c echo.Context) {

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func fullyQualifyImageURLs(r *http.Request) bool {
	param := r.URL.Query().Get("fullyQualifyImageUrls")
	if len(param) == 0 {
		param = "1"
	}

	fullyQualify, _ := strconv.ParseBool(param)
	return fullyQualify
}

// fullyQualifyCategoryImageURL - fully qualifies image URL for a category
func fullyQualifyCategoryImageURL(r *http.Request, c *Category) {
	if fullyQualifyImageURLs(r) {
		if len(c.Image) > 0 && c.Image != missingImageFile {
			c.Image = imageRootURL + c.Name + "/" + c.Image
		} else {
			c.Image = imageRootURL + missingImageFile
		}
	} else if len(c.Image) == 0 || c.Image == missingImageFile {
		c.Image = missingImageFile
	}
}

// fullyQualifyCategoryImageURLs - fully qualifies image URL for categories
func fullyQualifyCategoryImageURLs(r *http.Request, categories *Categories) {
	for i := range *categories {
		category := &((*categories)[i])
		fullyQualifyCategoryImageURL(r, category)
	}
}

// fullyQualifyProductImageURL - fully qualifies image URL for a product
func fullyQualifyProductImageURL(r *http.Request, p *Product) {
	if fullyQualifyImageURLs(r) {
		if len(p.Image) > 0 && p.Image != missingImageFile {
			p.Image = imageRootURL + p.Category + "/" + p.Image
		} else {
			p.Image = imageRootURL + missingImageFile
		}
	} else if len(p.Image) == 0 || p.Image == missingImageFile {
		p.Image = missingImageFile
	}
}

// fullyQualifyProductImageURLs - fully qualifies image URLs for all products
func fullyQualifyProductImageURLs(r *http.Request, products *Products) {
	for i := range *products {
		product := &((*products)[i])
		fullyQualifyProductImageURL(r, product)
	}
}

// Index Handler
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the Products Web Service")
}

// ProductIndex Handler
func ProductIndex(c echo.Context) error {
	initResponse(c)

	ret := RepoFindALLProducts()

	fullyQualifyProductImageURLs(c.Request(), &ret)

	return c.JSON(http.StatusOK, ret)
}

// CategoryIndex Handler
func CategoryIndex(c echo.Context) error {
	initResponse(c)

	ret := RepoFindALLCategories()

	fullyQualifyCategoryImageURLs(c.Request(), &ret)

	return c.JSON(http.StatusOK, ret)
}

// ProductShow Handler
func ProductShow(c echo.Context) error {
	initResponse(c)

	productIds := strings.Split(c.Param("productIDs"), ",")

	if len(productIds) > MAX_BATCH_GET_ITEM {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Sprintf("Maximum number of product IDs per request is %d", MAX_BATCH_GET_ITEM))
	}

	if len(productIds) > 1 {
		ret := RepoFindMultipleProducts(productIds)

		fullyQualifyProductImageURLs(c.Request(), &ret)

		return c.JSON(http.StatusOK, ret)
	} else {
		ret := RepoFindProduct(productIds[0])

		if !ret.Initialized() {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}

		fullyQualifyProductImageURL(c.Request(), &ret)

		return c.JSON(http.StatusOK, ret)
	}
}

// CategoryShow Handler
func CategoryShow(c echo.Context) error {
	initResponse(c)

	ret := RepoFindCategory(c.Param("categoryID"))

	if !ret.Initialized() {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	fullyQualifyCategoryImageURL(c.Request(), &ret)

	return c.JSON(http.StatusOK, ret)
}

// ProductInCategory Handler
func ProductInCategory(c echo.Context) error {
	initResponse(c)

	categoryName := c.Param("categoryName")

	ret := RepoFindProductByCategory(categoryName)

	fullyQualifyProductImageURLs(c.Request(), &ret)

	return c.JSON(http.StatusOK, ret)
}

// ProductFeatured Handler
func ProductFeatured(c echo.Context) error {
	initResponse(c)

	ret := RepoFindFeatured()

	fullyQualifyProductImageURLs(c.Request(), &ret)

	return c.JSON(http.StatusOK, ret)
}

func validateProduct(product *Product) error {
	if len(product.Name) == 0 {
		return errors.New("Product name is required")
	}

	if product.Price < 0 {
		return errors.New("Product price cannot be a negative value")
	}

	if product.Stock < 0 {
		return errors.New("Product current stock cannot be a negative value")
	}

	if len(product.Category) > 0 {
		categories := RepoFindCategoriesByName(product.Category)
		if len(categories) == 0 {
			return errors.New("invalid product category; does not exist")
		}
	}

	return nil
}

// UpdateProduct - updates a product
func UpdateProduct(c echo.Context) error {
	initResponse(c)

	var product Product

	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := c.Request().Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &product); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid request payload")
	}

	if err := validateProduct(&product); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	existingProduct := RepoFindProduct(c.Param("productID"))
	if !existingProduct.Initialized() {
		// Existing product does not exist
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	if err := RepoUpdateProduct(&existingProduct, &product); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error updating product")
	}

	fullyQualifyProductImageURL(c.Request(), &product)

	return c.JSON(http.StatusOK, product)
}

// UpdateInventory - updates stock quantity for one item
func UpdateInventory(c echo.Context) error {
	initResponse(c)

	var inventory Inventory

	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := c.Request().Body.Close(); err != nil {
		panic(err)
	}
	log.Println("UpdateInventory Body ", body)

	if err := json.Unmarshal(body, &inventory); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid request payload")
	}

	log.Println("UpdateInventory --> ", inventory)

	// Get the current product
	product := RepoFindProduct(c.Param("productID"))
	if !product.Initialized() {
		// Existing product does not exist
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	if err := RepoUpdateInventoryDelta(&product, inventory.StockDelta); err != nil {
		panic(err)
	}

	fullyQualifyProductImageURL(c.Request(), &product)

	return c.JSON(http.StatusOK, product)
}

// NewProduct  - creates a new Product
func NewProduct(c echo.Context) error {
	initResponse(c)

	var product Product
	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := c.Request().Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &product); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid request payload")
	}

	log.Println("NewProduct  ", product)

	if err := validateProduct(&product); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid request payload")
	}

	if err := RepoNewProduct(&product); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error creating product")
	}

	fullyQualifyProductImageURL(c.Request(), &product)

	return c.JSON(http.StatusOK, product)
}

// DeleteProduct - deletes a single product
func DeleteProduct(c echo.Context) error {
	initResponse(c)

	// Get the current product
	product := RepoFindProduct(c.Param("productID"))
	if !product.Initialized() {
		// Existing product does not exist
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")

	}

	if err := RepoDeleteProduct(&product); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error deleting product")
	}

	return c.NoContent(http.StatusNoContent)
}
