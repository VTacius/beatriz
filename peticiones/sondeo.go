package peticiones

import (
	"fmt"
	"strconv"
	"strings"

	g "github.com/gosnmp/gosnmp"
)

func obtenerDatos(origen string, oidRoot string) ([]g.SnmpPDU, error) {
	// Realiza la conexión inicial. TODO: Cuidar la configuración a futuro
	g.Default.Target = origen

	err := g.Default.Connect()
	if err != nil {
		return []g.SnmpPDU{}, fmt.Errorf("conexión:> %v", err)
	}

	defer g.Default.Conn.Close()

	// Obtenemos los datos que queremos usar
	resultado, err := g.Default.BulkWalkAll(oidRoot)
	if err != nil {
		return []g.SnmpPDU{}, fmt.Errorf("get:> %v", err)
	}

	return resultado, nil
}

func estructurarDatos(resultado []g.SnmpPDU) map[string]map[string]any {
	// Le damos estructura a los datos obtenidos
	mapa := map[string]map[string]any{}
	for _, item := range resultado {
		aspecto, indice := mapearOid(item.Name)
		valor := valorDefecto(item)
		if _, existe := mapa[aspecto]; existe {
			mapa[aspecto][indice] = valor
		} else {
			elemento := map[string]any{indice: valor}
			mapa[aspecto] = elemento
		}
	}

	return mapa
}

func ObtenerUsoToners(origen string) ([]Toner, error) {

	pduUso, err := obtenerDatos(origen, ".1.3.6.1.2.1.43.11")
	if err != nil {
		return []Toner{}, err
	}

	pduColores, err := obtenerDatos(origen, ".1.3.6.1.2.1.43.12")
	if err != nil {
		return []Toner{}, err
	}

	datosUso := estructurarDatos(pduUso)
	datosColores := estructurarDatos(pduColores)

	var toners []Toner
	indices := obtenerIndices(datosUso["2"])
	for _, j := range indices {
		var color string
		if c, existe := datosColores["4"][fmt.Sprintf("%v", j)]; existe {
			color = fmt.Sprintf("%v", c)
		} else {
			color = "Desconocido"
		}

		//.1.3.6.1.2.1.43.11.1.1.8
		capacidad := conversionAnyToInt(datosUso["8"][fmt.Sprintf("%v", j)])
		//.1.3.6.1.2.1.43.11.1.1.9
		uso := conversionAnyToInt(datosUso["9"][fmt.Sprintf("%v", j)])
		// .1.3.6.1.2.1.43.11.1.1.6
		modelo := datosUso["6"][fmt.Sprintf("%v", j)]

		toner := NewToner(origen, capacidad, uso, fmt.Sprintf("%v", modelo), color)
		if toner.Capacidad > 0 {
			toners = append(toners, *toner)
		}
	}

	return toners, nil
}

func conversionAnyToInt(valor any) int64 {
	if n, e := strconv.Atoi(fmt.Sprintf("%v", valor)); e == nil {
		return int64(n)
	} else {
		return 0
	}
}

func mapearOid(oid string) (aspecto string, indice string) {
	elementos := strings.Split(oid, ".")

	indice = elementos[len(elementos)-1]
	aspecto = elementos[len(elementos)-3]

	return
}

func valorDefecto(entrada g.SnmpPDU) any {
	switch entrada.Type {
	case g.OctetString:
		bytes := entrada.Value.([]byte)
		return string(bytes)
	default:
		return g.ToBigInt(entrada.Value).Int64()
	}
}

func obtenerIndices(lista map[string]any) (resultado []any) {
	for elemento := range lista {
		resultado = append(resultado, elemento)
	}

	return
}
