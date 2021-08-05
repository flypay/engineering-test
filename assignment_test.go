package assignment

import (
	"testing"
)

func TestMenus(t *testing.T) {
	t.Run("I can fetch the menu for restaurant Alpha", func(t *testing.T) {
		t.Fatal("You can implement this test to verify the menu")
	})

	t.Run("I can fetch the menu for restaurant Beta", func(t *testing.T) {
		t.Fatal("You can implement this test to verify the menu")
	})
}

func TestOrdering(t *testing.T) {
	t.Run("I can order dinner from restarurant Alpha", func(t *testing.T) {
		// Place an order for
		// 	1x Individual Mushroom Burger
		//	 	Extra Mayo
		//		No Tomato
		// 	1x Individual Buttermilk Chicken Burger
		// 	2x Regular Fries
		//		Paprika Salt
		//  2x Large Fizzy Cola
		t.Fatal("You can implement this test to verify the order is placed")
	})

	t.Run("I can order lunch from restarurant Beta", func(t *testing.T) {
		// Place an order for
		// 	1x Falafel Wrap
		//	 	Extra Chilli Sauce
		// 	2x Jam Donut
		// 	1x Bottled Water
		t.Fatal("You can implement this test to verify the order is placed")
	})
}
