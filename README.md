# Checkers-machine-learning

Current state:

- implemented checkers game
- added some ml parts to code


## ML model

Current model for loss function of player:

## Help

To run sample game with random model parameters run:

```bash
go run main.go
```


To see more options and commands:

```bash
go run main.go -help
```

## Test

Can be run with either

```bash
go test ./...
```

or when goconvey is installed with (start web dashboard with test results):

```bash
goconvey
```

## For web preview

```bash
go run main.go -drawer web
open html-view/index.html
```
