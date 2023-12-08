package pick_loli

import (
	"context"
	"testing"

	"github.com/chiyoi/apricot/test"
)

func TestPick(t *testing.T) {
	image, err := Pick(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer image.Body.Close()

	test.AddAttachment(t, image.Filename, image.Body)
	test.AddAttachment(t, image.Filename+"-description.txt", image.Description)
}
