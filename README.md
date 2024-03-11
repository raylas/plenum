# plenum

[![Main](https://github.com/raylas/plenum/actions/workflows/main.yaml/badge.svg)](https://github.com/raylas/plenum/actions/workflows/main.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/raylas/plenum)](https://goreportcard.com/report/github.com/raylas/plenum)

A package for membership and consensus that uses [Serf](https://github.com/hashicorp/serf) and [Raft](https://github.com/hashicorp/raft).

## Configuration

- `PLENUM_DEBUG`: Exposes Serf and Raft logging to `stdout`

## Usage

### Package

```go
bindAddress := "10.0.0.10"
bindPort := 6789
members := []string{"10.0.0.20:6789", "10.0.0.30:6789"}

p, err := plenum.Convene(nil, bindAddress, bindPort, members)
if err != nil {
  log.Fatal(err)
}
go p.Procede()

// Wait for, and act on leader designation
go func(*plenum.Plenum) {
  for <-p.IsLeader {
    fmt.Println("Members: ", p.Members())
    fmt.Println("Leader: ", p.Consensus())
    fmt.Println("Elected leader!")
  }
}(p)
```

### Example CLI

1. Start three members, in three terminals, with different ports:
```bash
# Terminal A
go run cmd/plenum/main.go -port 1111

# Terminal B
go run cmd/plenum/main.go -port 2222 -members=127.0.0.1:1111

# Terminal C
go run cmd/plenum/main.go -port 3333 -members=127.0.0.1:2222
```

2. Kill the member (`CTRL+C`) in **terminal A**, and watch consensus adjust.

3. Rejoin **terminal A** with:
```bash
go run cmd/plenum/main.go -port 1111 -members=127.0.0.1:2222
```

Each member also starts an HTTP listener with the followings routes for discovery:

- `/members`: Returns list of active membership
- `/leader`: Returns active leader
