// Goroutine - concepts
package main  

import (
	"fmt"
	"sync"
)



func main(){

	var wg sync.WaitGroup

	wg.Add(1)

	go func(){
		defer wg.Done()
		fmt.Println("hello")
	}()

	wg.Wait()
  
}
