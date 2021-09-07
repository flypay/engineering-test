package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	addrFlag = flag.String("addr", ":8082", "address to run betapos on")
)

func main() {
	flag.Parse()

	http.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintln(w, menu)
	})

	http.HandleFunc("/orders/create", func(w http.ResponseWriter, r *http.Request) {
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

	fmt.Printf("🙋‍♂️ BetaPOS is running on http://%s\n", *addrFlag)
	_ = http.ListenAndServe(*addrFlag, nil)
}

var (
	ErrOrderMissingID    = errors.New("id missing from order")
	ErrOrderMissingItems = errors.New("no items in order")
)

func handleOrder(order order) (orderResponse, error) {
	if order.ID == "" {
		return orderResponse{}, ErrOrderMissingID
	}
	if len(order.Items) == 0 {
		return orderResponse{}, ErrOrderMissingItems
	}

	m := struct {
		Categories map[string]struct {
			Name  string `json:"Name"`
			Items map[string]struct {
				Name     string  `json:"Name"`
				Price    float32 `json:"Price"`
				Quantity int     `json:"Quantity"`
				AddOns   []struct {
					ID    string  `json:"Id"`
					Name  string  `json:"Name"`
					Price float32 `json:"Price"`
				} `json:"AddOns"`
			} `json:"Items"`
		} `json:"Categories"`
	}{}

	_ = json.Unmarshal([]byte(menu), &m)

	categories := m.Categories

	log.Printf("handling order %q", order.ID)

	foundItems := make([]orderItem, 0, len(order.Items))
	for i, item := range order.Items {
		cat, ok := categories[item.CategoryID]
		if !ok {
			log.Printf("no category for %q found", item.CategoryID)
			continue
		}
		catItem, ok := cat.Items[item.ItemID]
		if !ok {
			log.Printf("no item for %q found in category %q", item.ItemID, item.CategoryID)
			continue
		}
		addOns := []string{}
		for _, ao := range item.AddOns {
			for _, catAO := range catItem.AddOns {
				if catAO.ID == ao {
					addOns = append(addOns, catAO.Name)
					break
				}
			}
		}

		log.Printf(
			"%q: item %d: quantity=%d name=%q addons=%v",
			order.ID, i+1, item.Quantity, catItem.Name, addOns,
		)
		foundItems = append(foundItems, orderItem{
			Name:     catItem.Name,
			Quantity: item.Quantity,
			AddOns:   addOns,
		})
	}
	return orderResponse{
		OrderID: order.ID,
		Items:   foundItems,
	}, nil
}

type order struct {
	ID    string `json:"Id"`
	Items []struct {
		CategoryID string   `json:"CategoryId"`
		ItemID     string   `json:"ItemId"`
		Quantity   int      `json:"Quantity"`
		AddOns     []string `json:"AddOns"`
	} `json:"Items"`
}

type orderItem struct {
	Name     string   `json:"Name"`
	AddOns   []string `json:"AddOns"`
	Quantity int      `json:"Quantity"`
}

type orderResponse struct {
	OrderID string      `json:"OrderId"`
	Items   []orderItem `json:"Items"`
}

func init() {
	// Test that the menuJSON is valid
	var m map[string]interface{}

	if err := json.Unmarshal([]byte(menu), &m); err != nil {
		log.Fatal(err)
	}
}

const menu = `{
	"CategoriesOrder": [
		"qqdluj",
		"dmcshb",
		"gokpww",
		"lzoaud"
	],
	"Categories": {
		"qqdluj":{
			"Name": "Lunch Menu",
			"Items": {
				"hjrlho": {
					"Name": "Falafel Wrap",
					"Description": "Made with real chickpeas",
					"Price": 10.1,
					"AddOns": [
						{
							"Id": "oecwhs",
							"Name": "Extra Falafel Ball",
							"Price": 1
						},
						{
							"Id": "warnlj",
							"Name": "Extra Hummus",
							"Price": 0.5
						},
						{
							"Id": "yzhllj",
							"Name": "Extra Tabouli",
							"Price": 0.25
						},
						{
							"Id": "cbnufj",
							"Name": "Extra Chilli Sauce",
							"Price": 0.99
						}
					]
				},
				"qmdehd": {
					"Name": "Beef Shawarma",
					"Price": 12,
					"AddOns": [
						{
							"Id": "duaweu",
							"Name": "Extra Rice",
							"Price": 2.5
						}
					]
				}
			}
		},
		"dmcshb": {
			"Name": "Dessert",
			"Items": {
				"ndytqi": {"Name": "Baklawa", "Price": 13}
			}
		},
		"gokpww": {
			"Name": "Sides",
			"Items": {
				"gtbiop": {"Name": "Hummus", "Price": 10.20},
				"dwaecr": {"Name": "Soup", "Price": 4.99},
				"ukhqnd": {"Name": "Jam Donut", "Price": 1.35},
				"fllrsi": {"Name": "Samosa", "Price": 3}
			}
		},
		"lzoaud": {
			"Name": "Drinks",
			"Items": {
				"ugqcbb": {"Name": "Bottled Water", "Price": 4},
				"pgzigb": {"Name": "Juice", "Price": 5},
				"sjalnl": {"Name": "Bottled Beer", "Price": 6}
			}
		}
	}
}`
