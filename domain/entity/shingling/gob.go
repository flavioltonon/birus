package shingling

import "encoding/gob"

func init() {
	gob.Register(new(Shingling))
	gob.Register(new(Shingle))
	gob.Register(new(ShinglesCounter))
	gob.Register(new(ShinglesMapper))
}
