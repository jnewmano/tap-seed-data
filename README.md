[![Go](https://github.com/jnewmano/tap-seed-data/actions/workflows/go.yml/badge.svg)](https://github.com/jnewmano/tap-seed-data/actions/workflows/go.yml)

# tap-seed-data
Generate random lists of contact and appointment data.

## Singer.IO tap

This app will output a Singer.IO stream of random contact
and appointment data. 

Assuming Go (1.16 or newer) is installed in your local environment, the 
app can be run using

```
go run github.com/jnewmano/tap-seed-data
```

## Random contact and appointment generator as a package
Alternatively, the `seeddata` package can be used as a library
and embedded into your application. Both contacts and appointments
have some configurable options to shape the generated data into
how you may need it to look.

## Notes

Ages and names are randomly chosen using recent US Census data to
give a representative sample of ages and names.

