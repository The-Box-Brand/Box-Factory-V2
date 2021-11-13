package boxes

import wr "github.com/mroth/weightedrand"

//
// Options
//

var ExtrasConfig = map[string]wr.Choice{
	"strap": {
		Item:   "strap",
		Weight: 25,
	},
	"adhesive": {
		Item:   "adhesive",
		Weight: 25,
	},
	"cutout": {
		Item:   "cutout",
		Weight: 4,
	},
	"label": {
		Item:   "label",
		Weight: 4,
	},
}

// 50/50 for now
var NumberOfTraitsConfig = []wr.Choice{
	{
		Item:   1,
		Weight: 30,
	},
	{
		Item:   2,
		Weight: 20,
	},
	{
		Item:   3,
		Weight: 10,
	},
	{
		Item:   4,
		Weight: 4,
	},
	{
		Item:   5,
		Weight: 2,
	},
}
