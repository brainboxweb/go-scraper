package main

import (
	"fmt"
	"os"
	"github.com/headzoo/surf"
	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
	"encoding/json"
	"regexp"
	"strconv"
)


type Data struct {
	Products []Product  `json:"results"`
	Total       float64 `json:"total"`
}

type Product struct {
	Title       string  `json:"title"`
	Size        string  `json:"size"`
	UnitPrice   float64 `json:"unit_price"`
	Description string  `json:"description"`
}


func main() {

	fmt.Println("\n\n------STARTING------\n\n")

	//Struct to collect the product data
	var data Data

	targetUrl := os.Args[1]
	//todo - check input params

	bow := surf.NewBrowser()
	err := bow.Open(targetUrl)
	if err != nil {
		panic(err)
	}

	//Find links on the target page
	fmt.Println("\n\nGetting links for " + targetUrl +  "\n\n")
	bow.Find("#productLister .productInner h3 a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		product := getProduct(url)
		data.Products = append(data.Products, product)

		//increment the total
		data.Total += product.UnitPrice
	})

	//Create json
	b, _ := json.Marshal(data)

	//output json
	fmt.Println("\n\nJSON output:\n\n")
	fmt.Println(string(b))

	fmt.Println("\n\n------FINISHED------\n\n")
}

func getProduct(url string) Product {

	var product Product

	fmt.Println("Scraping " +  url)

	bow := surf.NewBrowser()
	err := bow.Open(url)
	if err != nil {
		panic(err)
	}

	//Page size
	body := bow.Body()  //Too simplistic?
	size := len(body) //bytes v. runes issue?
	human := humanize.SI(float64(size), "b")
	product.Size = human

	//Title
	bow.Find(".productTitleDescriptionContainer h1").Each(func(_ int, s *goquery.Selection) {
		title := s.Text()
		product.Title = title
	})

	//PricePerUnit
	bow.Find(".pricePerUnit").Each(func(_ int, s *goquery.Selection) {
		price := s.Text()
		product.UnitPrice = stringToFloat(price)
	})

	//Description
	bow.Find(".productText").Each(func(_ int, s *goquery.Selection) {
		description := s.Text()
		product.Description = description
	})

	return product
}

func stringToFloat(price string) float64 {
	r, _ := regexp.Compile("[^0-9.]") //Strip unwanted characters
	price = r.ReplaceAllString(price, "")
	unitPriceFloat, _ := strconv.ParseFloat(price, 64)

	return unitPriceFloat
}

//go run main.go 'http://www.sainsburys.co.uk/webapp/wcs/stores/servlet/CategoryDisplay?listView=true&orderBy=FAVOURITES_FIRST&parent_category_rn=12518&top_category=12518&langId=44&beginIndex=0&pageSize=20&catalogId=10137&searchTerm=&categoryId=185749&listId=&storeId=10151&promotionId=#langId=44&storeId=10151&catalogId=10137&categoryId=185749&parent_category_rn=12518&top_category=12518&pageSize=20&orderBy=FAVOURITES_FIRST&searchTerm=&beginIndex=0&hideFilters=true'
