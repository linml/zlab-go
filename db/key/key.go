package key

import "fmt"

func PlayerID(id string) string {
	return fmt.Sprintf("plyid:%s", id)
}
