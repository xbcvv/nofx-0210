
package main
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)
func main() {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://fapi.binance.com/fapi/v1/ticker/24hr")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Status:", resp.StatusCode)
}

