# AlphaPOS

## Endpoints

* `GET /menu/categories` returns the structure of the products within the POS.
    * Products can be found in the `products` response.
* `GET /menu/products` returns the products available on the POS which are referenced by the categories.
    * Product `extras` can be found in the `ingredients` response.
* `GET /menu/ingredients` returns the extras and ingredients found within products, these things are not orderable items within the menu but can be used as add-ons or indicated that they are included with the item.

* `POST /orders` submits an order to the POS to be made by the kitchen. 
    * Request Payload
        ```
        {
			"orderId": "some-id",
			"products": [
				{"productId": "product-id-1", "sizeId": "size-id", "ingredientIds": ["ingredient-id-1"], "quantity": 1},
				{"productId": "product-id-2", "ingredientIds": ["ingredient-id-2", "ingredient-id-3"], "quantity": 5}
			]
		}
        ```
	* Response Payload
		```
		{
			"orderId": "some-id",
			"products": [
				{
					"name": "Mushroom Burger",
					"size": "Individual",
					"ingredients": [
						"Seeded Bun",
						"Mushroom Patty",
						"Lettuce",
						"Extra Mayo"
					],
					"quantity": 1
				},
				{
					"name": "Fizzy Cola",
					"size": "Large",
					"ingredients": [],
					"quantity": 2
				}
			]
		}
		```
