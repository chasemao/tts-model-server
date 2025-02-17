package edge

import (
	"fmt"
	"testing"
)

func TestCreateSSML(t *testing.T) {
	i := &impl{}
	fmt.Println(i.createSSML("haha\"\\'<>&", "haha\"\\'<>&"))
}
