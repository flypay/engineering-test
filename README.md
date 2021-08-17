# Backend Engineering Technical Assignment

Hello 👋 - thanks for your time on the phone just now - I really enjoyed our chat!

Next, we would like to get a sense of your technical skills. Yes - we hate these take home tests too - but it's the best way we know to learn about you and how you'll code on the job. We've kept it relevant to the type of work that you'll be doing at Flyt; so you get a taste of what our teams do every day.

## Getting Setup

* Fork this repository, everything you need is bundled here. 
* We've included applications for two APIs that mimick two different Point of Sale (POS) systems we work with. You will find these under `cmd/{alphapos,betapos}/main.go`.
* Each POS can be run with `go run cmd/{alphapos,betapos}/main.go`
* You can find documentation here for [alpha POS](cmd/alphapos/README.md) and [beta POS](cmd/betapos/README.md)

## The Task

* Your task is to create a new API that can do two things:
  1. receive request to fetch menu data from either POS and return this in a common format of your choice.
  2. receive orders in the example format below, and send the order to the relevant point of sale system.
* We have included the test which you must satisfy. You can be run with the standard Go test tool. We encourage additional edge case scenarios to be considered.

```bash
$ go test -v .
```

The existing tests do the following:

* call your menu API for `alpha` which results in a representation of the menu from `alphapos` with all categories and items in your standard schema defined
* call your menu API for `beta` which results in a representation of the menu from `betapos` with all categories and items in your standard schema defined
* call your orders API with [this basket](#order-payload) for restaurant `alpha` results in `alphapos` receiving the correct order
* call your orders API with [this basket](#order-payload) for restaurant `beta` results in `betapos` receiving the correct order

## Once You Are Done
* Please record a video talking us through your code for no more than 5 minutes
* Create an archive (zip) of your solution, then send it over to the Flyt recruiter

## The Process
* You can submit your assignment anytime in the next 7 days as we want to give you plenty of time to work your way through it
* One of our engineers will review your solution (usually 2-3 days) and send you feedback regardless of how you progress on the job application
* We might also reach out to discuss your solution in more detail before progressing you to the next stage

# Examples

## Order Payload

The following is the order payload that your API is expected to recieve, we have created a [Go type](internal/types/order.go) for it as well:

```

{
	"id": "order-id-1234",
	"pos": "alpha",
	"items": [
		{
			"id": "item-id", // required
			"size_id": "123", // required
			"ingredient_ids: [ "ingredient-id-1", "ingredient-id-2" ], // optional
			"extras": [ "extra-id-1", "extra-id-2" ] // optional
		}
	]
}
```

- The POS will be passed through on the order, this should be used to dynamically route to the correct POS API
- All items are expected to have a size
- If the POS doesn't support ingredients then this can be left empty


## Menu API

This can be any schema you want, and it is left to you to decide on the best schema to use. You
should be able to use the menu to help you when creating orders.

- If a POS doesn't support sizes, then a "Regular" size should be created for the item
- If a POS doesn't support ingredients then you can ignore this when it comes through on an order
- Extras are additional add-ons for an item, usually coming with an additional price

## Response Codes
### Response Codes

```
200: Success
400: Bad request
404: Cannot be found
405: Method not allowed
```

### Error Codes Details

```
1000: Bad Request
1001: Parameter Error
1002: Item Missing
```

### Example Error Message
```
http code 404
{
    "errorCode": 1002,
    "errorMessage": "CategoryId missing"
}
```
