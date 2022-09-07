package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"sanidad.gob.sv/alortiz/beatriz/peticiones"
)

func principal(contexto *cli.Context) error {
	origen := contexto.String("origen")

	envio := contexto.Bool("envio")
	destino := contexto.String("destino")
	destinoToken := contexto.String("destino-token")
	destinoOrganizacion := contexto.String("destino-organizacion")
	destinoBucket := contexto.String("destino-bucket")

	toners, err := peticiones.ObtenerUsoToners(origen)
	if err != nil {
		return err
	}

	if envio {
		backend := peticiones.NewBackend(destino, destinoToken, destinoOrganizacion, destinoBucket)
		return backend.Enviar(toners)
	} else {
		fmt.Printf("\nImpresora en %s\n", origen)
		for _, t := range toners {
			fmt.Println("\t", t.String())
		}
	}

	return nil
}

func main() {
	origen := &cli.StringFlag{Name: "origen", Usage: "IP del impresor a scrappear", Required: true}
	envio := &cli.BoolFlag{Name: "envio", Usage: "Indica si los datos deben enviarse o no"}
	destino := altsrc.NewStringFlag(&cli.StringFlag{Name: "destino", Usage: "Backend InfluxDB2 para almacenar datos"})
	destinoToken := altsrc.NewStringFlag(&cli.StringFlag{Name: "destino-token", Usage: "Token para conectarse a Backend"})
	destinoOrganizacion := altsrc.NewStringFlag(&cli.StringFlag{Name: "destino-organizacion", Usage: "Organización en Backend"})
	destinoBucket := altsrc.NewStringFlag(&cli.StringFlag{Name: "destino-bucket", Usage: "Bucket dentro de la organización"})

	banderas := []cli.Flag{origen, envio, destino, destinoToken, destinoOrganizacion, destinoBucket}
	app := &cli.App{
		Name:  "beatriz",
		Usage: "Scrapper para impresoras con SNMP",
		Flags: banderas,
		Before: altsrc.InitInputSourceWithContext(banderas,
			func(context *cli.Context) (altsrc.InputSourceContext, error) {
				return altsrc.NewYamlSourceFromFile("/etc/beatriz.yaml")
			}),
		Action: principal,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
