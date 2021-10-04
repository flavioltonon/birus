package classifier

import "encoding/gob"

func init() {
	gob.Register(new(Classifier))
	gob.Register(new(classifierOptions))
}
