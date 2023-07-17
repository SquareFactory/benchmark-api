package try

import (
	"log"
	"time"
)

func Do[T interface{}](
	fn func() (T, error),
	tries int,
	delay time.Duration,
) (result T, err error) {
	for try := 0; try < tries; try++ {
		result, err = fn()
		if err == nil {
			break
		}
		log.Printf("try failed: %s", err)
		time.Sleep(delay)
	}
	return result, err
}
