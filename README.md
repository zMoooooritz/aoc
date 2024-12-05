## Running Locally
### Requirements
Go 1.16+ is required because [embed][embed] is used for input files.

Use `go run main.go -part <1 or 2>` will be usable to run the actual inputs for that day.

Use `go test -run RegExpToMatchFunctionNames .` to run examples and unit tests via the `main_test.go` files.

## Scripts
Makefile should be fairly self-documenting. Alternatively you can run the binaries yourself via `go run` or `go build`.

`make help` prints a help message.

### Create skeleton and input for a day

```bash
make setup DAY=10 YEAR=2020
```

### Make skeleton files
```bash
for ((i=1; i<26; i++)); do
    make skeleton DAY=$i YEAR=2020
done
```

Note that skeletons use [embed][embed] and __will not compile__ without an `input.txt` file located in the same folder. Input files can be made via `make input`.
```sh
make skeleton DAY=10 YEAR=2020
make input DAY=10 YEAR=2020 AOC_SESSION_COOKIE=your_cookie
```

### Fetch inputs and write to input.txt files
Requires passing your cookie from AOC from either `-cookie` flag, or `AOC_SESSION_COOKIE` env variable.
```bash
make input DAY=10 YEAR=2020
```

[embed]: https://golang.org/pkg/embed/
