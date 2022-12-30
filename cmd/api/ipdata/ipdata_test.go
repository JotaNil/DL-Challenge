package ipdata

import (
	"fmt"
	"strconv"
	"testing"
)

func TestIPTransformer(t *testing.T) {
	ip1 := "5.181.131.180"

	fmt.Printf(strconv.FormatInt(stringIPToDecimal(ip1), 10))

}

func TestTestIPTransformer2(t *testing.T) {
	ip := 95781817

	fmt.Printf(decimalIPToString(int64(ip)))
}

//5.181.131.178
