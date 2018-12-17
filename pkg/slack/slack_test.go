package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestEncoding(t *testing.T) {
	body := errorResp(errors.New("it failed"))
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}
