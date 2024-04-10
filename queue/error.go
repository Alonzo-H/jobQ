package queue

import "fmt"

var ErrIdCollision = fmt.Errorf("job with same id already exists")
var ErrNotFound = fmt.Errorf("not found")
var ErrFinalState = fmt.Errorf("final state was reached")
