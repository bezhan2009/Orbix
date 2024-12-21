package _chan

var pkgChan = make(map[string]interface{})

func GetChan() map[string]interface{} {
	return pkgChan
}
