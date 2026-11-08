# Hand-In-4
How to start our system

Our program works for three nodes, no more or less.
After starting up the first node, the two remaining nodes should be started within approximately 5 seconds.

The port should be specified in the command line when starting up a Node:

-- Example on how to run 3 Nodes with the main-folder ---
go run ./main  -port=5050
go run ./main  -port=5051
go run ./main  -port=5052

