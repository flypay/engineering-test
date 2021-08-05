# BetaPOS

## Endpoints

* `GET /menu` returns the structure of the categories and items within the POS.
    * Categories can be found in the top-level response
	* CategoriesOrder defines that order the catagories should be displayed
	* Items can be found in a Category
	* AddOns can be found in an Item
	* NOTE: Each Category and Item contains an ID, and the ID is the JSON key.
	* Menu Payload
	```
	{ 
		"CategoriesOrder": [
			"category-id",
		],
		"Categories": {
		{
			"category-id": {
				"Name": "Category",
				"Items": {
					"item-id": {
						"Name": "Item"
					}
				}
			}
		}
	}
	```

* `POST /orders/create` submits an order to the POS to be made by the kitchen.
    * Request Payload
        ```
        {
			"Id": "order-id-1234",
			"Items": [
				{"CategoryId": "category-id", "ItemId": "item-id", "Quantity": 1},
				{
					"CategoryId": "category-id", 
					"ItemId": "item-id", 
					"Quantity": 2,
					"AddOns": ["addon-id"]
				}
			]
		}
        ```

	* Response Payload
		```
		{
			"OrderId": "order-id-1234",
			"Items": [
				{
					"Name": "Falafel Wrap",
					"AddOns": [
						"Extra Chilli Sauce"
					],
					"Quantity": 1
				},
				{
					"Name": "Bottled Water",
					"AddOns": [],
					"Quantity": 1
				}
			]
		}
		```
