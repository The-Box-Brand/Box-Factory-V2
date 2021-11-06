package boxes

import wr "github.com/mroth/weightedrand"

//
// Options
//

var ExtrasConfig = map[string]wr.Choice{
	"binding": {
		Item:   "binding",
		Weight: 25,
	},
	"cutout": {
		Item:   "cutout",
		Weight: 25,
	},
}

// 50/50 for now
var NumberOfTraitsConfig = []wr.Choice{
	{
		Item:   1,
		Weight: 33,
	},
	{
		Item:   2,
		Weight: 33,
	},
	{
		Item:   3,
		Weight: 27,
	},
}
