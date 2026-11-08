# Hand-In-4
 ## How to start our system

Our program works for three nodes, no more or less.
After starting up the first node, the two remaining nodes should be started within approximately 5 seconds.

The port should be specified in the command line when starting up a Node.

When starting the program open up 3 terminals.

Write in the first terminal:
```bash
go run ./main  -port=5050
```

Write in the second terminal:
```bash
go run ./main  -port=5051
```

Write in the third terminal:
```bash
go run ./main  -port=5052
```


