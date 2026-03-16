Debug a Go issue. Argument: $ARGUMENTS

1. Reproduce the issue described
2. Use `go test -v -run <TestName>` to isolate failures
3. Read relevant source and test files
4. Check for common Go bugs: nil receivers, map initialization, slice aliasing, interface satisfaction
5. For Bubble Tea issues: verify Update returns correct Cmd, check View doesn't panic on zero-value model
6. Propose and apply fix, then verify with tests
