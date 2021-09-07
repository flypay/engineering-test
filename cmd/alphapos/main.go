package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

var (
	addrFlag = flag.String("addr", ":8081", "address to run alphapos on")
)

func main() {
	flag.Parse()

	http.HandleFunc("/menu/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		p := path.Base(r.URL.Path)
		switch p {
		case "categories":
			fmt.Fprintln(w, categories)
		case "products":
			fmt.Fprintln(w, products)
		case "ingredients":
			fmt.Fprintln(w, ingredients)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		req, _ := ioutil.ReadAll(r.Body)
		log.Printf("Order received:\n%s\n", req)

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "bad content type", http.StatusBadRequest)
			return
		}

		order := order{}
		if err := json.Unmarshal(req, &order); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		or, err := handleOrder(order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(or); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	})

	log.Printf("🙋‍♀️ AlphaPOS is running on http://%s\n", *addrFlag)
	_ = http.ListenAndServe(*addrFlag, nil)
}

func handleOrder(order order) (orderResp, error) {
	if order.OrderID == "" {
		return orderResp{}, ErrOrderMissingID
	}
	if len(order.Products) == 0 {
		return orderResp{}, ErrOrderMissingProducts
	}

	prods := struct {
		Products []struct {
			ID            string `json:"productId"`
			Name          string `json:"name"`
			DefaultSizeID string `json:"default_size_id"`
			Sizes         []struct {
				ID   string `json:"sizeId"`
				Name string `json:"name"`
			} `json:"sizes"`
		} `json:"products"`
	}{}
	_ = json.Unmarshal([]byte(products), &prods)

	ings := struct {
		Ingredients []struct {
			ID   string `json:"ingredientId"`
			Name string `json:"name"`
		} `json:"ingredients"`
	}{}
	_ = json.Unmarshal([]byte(ingredients), &ings)

	foundProducts := make([]orderProduct, 0, len(order.Products))
	log.Printf("handling order %q", order.OrderID)
	for i, op := range order.Products {
		foundProduct := false
		for _, p := range prods.Products {
			if op.ID != p.ID {
				continue
			}
			sizeID := op.SizeID
			if sizeID == "" {
				sizeID = p.DefaultSizeID
			}
			sizeName := ""
			for _, s := range p.Sizes {
				if s.ID == sizeID {
					sizeName = s.Name
				}
			}
			ingredientNames := []string{}
			for _, pi := range op.IngredientIDs {
				foundIngredient := false
				for _, ing := range ings.Ingredients {
					if ing.ID == pi {
						ingredientNames = append(ingredientNames, ing.Name)
						foundIngredient = true
						break
					}
				}
				if !foundIngredient {
					log.Printf("%q: product %d: ingredient %q not found", order.OrderID, i+1, pi)
					return orderResp{}, fmt.Errorf("ingredient %q for product %q not found", pi, op.ID)
				}
			}

			log.Printf(
				"%q: product %d: quantity=%d name=%q size=%q ingredients=%v",
				order.OrderID, i+1, op.Quantity, p.Name, sizeName, ingredientNames,
			)
			foundProducts = append(foundProducts, orderProduct{
				Name:        p.Name,
				Size:        sizeName,
				Quantity:    op.Quantity,
				Ingredients: ingredientNames,
			})
			foundProduct = true
			break
		}
		if !foundProduct {
			log.Printf("%q: product %d: %q not found", order.OrderID, i+1, op.ID)
			return orderResp{}, fmt.Errorf("product %q not found", op.ID)
		}
	}

	return orderResp{
		OrderID:  order.OrderID,
		Products: foundProducts,
	}, nil
}

type order struct {
	OrderID  string `json:"orderId"`
	Products []struct {
		ID            string   `json:"productId"`
		SizeID        string   `json:"sizeId"`
		IngredientIDs []string `json:"ingredientIds"`
		Quantity      int      `json:"quantity"`
	} `json:"products"`
}

type orderProduct struct {
	Name        string   `json:"name"`
	Size        string   `json:"size"`
	Ingredients []string `json:"ingredients"`
	Quantity    int      `json:"quantity"`
}
type orderResp struct {
	OrderID  string         `json:"orderId"`
	Products []orderProduct `json:"products"`
}

var categories = `{
	"categories": [
		{
			"categoryId": "1001",
			"name": "Burgers",
			"description": "",
			"subcategories": [
				{
					"subcategoryId": "2001",
					"name": "Veggie Burgers",
					"products": [
						"6001"
					]
				},
				{
					"subcategoryId": "2002",
					"name": "Chicken Burgers",
					"products": [
						"6002"
					]
				}
			]
		},
		{
			"categoryId": "1002",
			"name": "Sides",
			"description": "",
			"subcategories": [
				{
					"subcategoryId": "2003",
					"name": "Potato",
					"products": [
						"6003"
					]
				}
			]
		},
		{
			"categoryId": "1003",
			"name": "Drinks",
			"description": "",
			"subcategories": [
				{
					"subcategoryId": "2004",
					"name": "Soft",
					"products": [
						"6004"
					]
				}
			]
		}
	]
}`

var products = `{
    "products": [
        {
            "productId": "6001",
            "name": "Mushroom Burger",
            "description": "Whole fried mushroom with our unique spice blend",
            "image": "https://images.unsplash.com/photo-1516774266634-15661f692c19?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2532&q=80",
            "sizes": [
                {
                    "sizeId": "8001",
                    "name": "Individual",
                    "price": 14.95
                },
                {
                    "sizeId": "8002",
                    "name": "Double Up",
                    "price": 21.95
                }
            ],
            "defaultIngredients": [
                "9001",
                "9002",
                "9003",
                "9004"
            ],
            "extras": [
                {
                    "ingredientId": "9010",
                    "price": 0.99
                },
                {
                    "ingredientId": "9007",
                    "price": 0.99
                }
            ]
        },
        {
            "productId": "6002",
            "name": "Buttermilk Chicken Burger",
            "description": "Juicy tender chicken thigh coated in crispy buttermilk batter.",
            "image": "https://images.unsplash.com/photo-1597900121060-cf21f1cfa5e6?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1234&q=80",
            "sizes": [
                {
                    "sizeId": "8001",
                    "name": "Individual",
                    "price": 15.95
                }
            ],
            "defaultIngredients": [
                "9005",
                "9006",
                "9007",
                "9008"
            ],
            "extras": []
        },
        {
            "productId": "6003",
            "name": "Fries",
            "description": "Skinny cut locally sourced potato fries.",
            "image": "https://images.unsplash.com/photo-1541592106381-b31e9677c0e5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2550&q=80",
            "sizes": [
                {
                    "sizeId": "8005",
                    "name": "Regular",
                    "price": 3.95
                },
                {
                    "sizeId": "8006",
                    "name": "Large",
                    "price": 5.95
                }
            ],
            "defaultIngredients": [],
            "extras": [
                {
                    "ingredientId": "9009",
                    "price": 0.99
                }
            ]
        },
        {
            "productId": "6004",
            "name": "Fizzy Cola",
            "description": "",
            "image": "https://images.unsplash.com/photo-1592153995863-9fb8fe173740?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2734&q=80",
            "sizes": [
                {
                    "sizeId": "8007",
                    "name": "Medium",
                    "price": 2.95
                },
                {
                    "sizeId": "8006",
                    "name": "Large",
                    "price": 3.95
                }
            ],
            "defaultIngredients": [],
            "extras": []
        }
    ]
}`

var ingredients = `{
    "ingredients": [
        {
            "ingredientId": "9001",
            "name": "Seeded Bun",
            "groupDescription": "Bread"
        },
        {
            "ingredientId": "9002",
            "name": "Mushroom Patty",
            "groupDescription": "Patties"
        },
        {
            "ingredientId": "9003",
            "name": "Lettuce",
            "groupDescription": "Salad"
        },
        {
            "ingredientId": "9004",
            "name": "Tomato",
            "groupDescription": "Salad"
        },
        {
            "ingredientId": "9005",
            "name": "Brioche Bun",
            "groupDescription": "Bread"
        },
        {
            "ingredientId": "9006",
            "name": "Chicken Patty",
            "groupDescription": "Patties"
        },
        {
            "ingredientId": "9007",
            "name": "Pickles",
            "groupDescription": "Salad"
        },
        {
            "ingredientId": "9008",
            "name": "Slaw",
            "groupDescription": "Salad"
        },
        {
            "ingredientId": "9009",
            "name": "Paprika Salt",
            "groupDescription": "Condiments"
        },
        {
            "ingredientId": "9010",
            "name": "Extra Mayo",
            "groupDescription": "Condiments"
        }
    ]
}`

var (
	ErrOrderMissingID       = errors.New("order id missing from order")
	ErrOrderMissingProducts = errors.New("no products in order")
)
