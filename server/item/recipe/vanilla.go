package recipe

import (
	// Insure all blocks and items are registered before trying to load vanilla recipes.
	_ "github.com/df-mc/dragonfly/server/block"
	_ "github.com/df-mc/dragonfly/server/item"

	_ "embed"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

//go:embed crafting_data.nbt
var vanillaCraftingData []byte

func init() {
	var vanillaRecipes struct {
		Shaped []struct {
			Input    inputItems `nbt:"input"`
			Output   outputItem `nbt:"output"`
			Priority int32      `nbt:"priority"`
			Width    int32      `nbt:"width"`
			Height   int32      `nbt:"height"`
		} `nbt:"shaped"`
		Shapeless []struct {
			Input    inputItems `nbt:"input"`
			Output   outputItem `nbt:"output"`
			Priority int32      `nbt:"priority"`
		} `nbt:"shapeless"`
	}

	if err := nbt.Unmarshal(vanillaCraftingData, &vanillaRecipes); err != nil {
		panic(err)
	}

	for _, s := range vanillaRecipes.Shapeless {
		input, ok := s.Input.toInputItems()
		if !ok {
			continue
		}
		output, ok := s.Output.ToStack()
		if !ok {
			continue
		}
		Register(ShapelessRecipe{
			Inputs: input,
			Output: output,
		})
	}

	for _, s := range vanillaRecipes.Shaped {
		input, ok := s.Input.toInputItems()
		if !ok {
			continue
		}
		output, ok := s.Output.ToStack()
		if !ok {
			continue
		}
		Register(ShapedRecipe{
			Inputs:     input,
			Output:     output,
			Dimensions: Dimensions{int(s.Width), int(s.Height)},
		})
	}
}
