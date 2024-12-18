package utils

type Temperatura struct {
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
}

func ConverterTemperaturas(celsius float64) Temperatura {
	return Temperatura{
		Celsius:    celsius,
		Fahrenheit: celsius*1.8 + 32,
		Kelvin:     celsius + 273,
	}
}
