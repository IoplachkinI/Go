package generator

type IdGenerator interface {
	GenerateId(name string) int
}

type PolynomialHasher struct {
	Mod int
	Mult int
}

func (hasher PolynomialHasher) GenerateId(str string) int {
	var hash, deg = 0, 1
	for _, char := range str {
		charCode := (int(char) - '0')
		hash = (hash + (charCode * deg) % hasher.Mod) % hasher.Mod
		deg = (deg * hasher.Mult) % hasher.Mod
	}
	return hash
}

type SimpleHasher struct {}

func (hasher SimpleHasher) GenerateId(str string) int {
	sum := 0
	for _, char := range str {
		charCode := (int(char) - '0')
		sum += charCode * charCode
	}
	return sum
}