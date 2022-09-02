# beatriz
Scrapper para uso de toner en impresores usando SNMP

## Construcción
Escoja el método: Desde un sistema equivalente al sistema destino con la misma versión de Go (1.18)
```bash
go build .
```

Usando `podman` o `docker`
```bash
podman run  -it  --rm -v "$PWD":/go/src/myapp -w /go/src/myapp golang:1.18-bullseye go build .
```

## Configuración
```bash
cat <<MAFI > /etc/beatriz.yaml
destino: http://stats.sanidad.gob.sv:8086
destino-token: 2-qbR-mdKDF6f9qO-QW-UftFSeuGnXUoc_R2W_UKEw6mC1mbndISAbnKyw40dCdgaQtfQH2dYFHlRtV0gWpHgA==
destino-organizacion: sanidad
destino-bucket: ambientales
MAFI
```