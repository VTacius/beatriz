package peticiones

import (
	"fmt"

	g "github.com/gosnmp/gosnmp"
)

type ImpresoraOids struct {
	Usado string
	Total string
}

// TODO: Tener en cuenta la configuración a futuro
func ObtenerUsoToner(origen string, oidCfg ImpresoraOids) (float64, error) {

	oids := []string{oidCfg.Total, oidCfg.Usado}

	g.Default.Target = origen
	err := g.Default.Connect()
	if err != nil {
		return 0.0, fmt.Errorf("conexión:> %v", err)
	}
	defer g.Default.Conn.Close()

	result, errGet := g.Default.Get(oids)
	if errGet != nil {
		return 0.0, fmt.Errorf("get:> %v", errGet)
	}

	// TODO: Pues hay que cuidar el índice, y que el valor no sea negativo
	total := float64(g.ToBigInt(result.Variables[0].Value).Int64())
	usado := float64(g.ToBigInt(result.Variables[1].Value).Int64())

	resultado := (usado / total) * 100
	return resultado, nil
}
