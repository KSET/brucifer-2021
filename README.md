# Brucifer 2021
Aplikacija za stranicu za Brucošijadu FER-a 2021.

Trenutno se sastoji od samo placeholder stranica, ali s vremenom će se pretvoriti u punu stranicu.



## Kako pokrenuti
Aplikaciju je prvo potrebno buildati.  

U izvršnoj datoteci je sadržano sve što je potrebno da aplikacija normalno radi (svi asseti su embedani unutar izvršne datoteke).

### Kako buildati
 1. Osigurati da je instaliran [Go](https://golang.org/)
 2. Pokrenuti `make build` ili `make compact`

`make compact` također zahtijeva da je instaliran [`upx`](https://upx.github.io/) paket.

### CLI parametri
  - `-p $PORT`, `--port $PORT` - na kojim vratima će aplikacija slušati (default `3000`)
  - `-h $HOST`, `--host $HOST` - na koji host će se aplikacija vezati (defaut `0.0.0.0`)

### Primjer
```bash
./bin/brucifer --host '0.0.0.0' --port 3000
```
