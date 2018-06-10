=alamode=

Avoid using the browser UI for writing and running mode queries.

This app lets you edit SQL with vim and then run the query with jq.

Examples:

$ ./alamode spaces

1. fd55555c5555 Personal
2. fd55555c4444 Dash Reports
3. fd55555c3333 Company Wide

$ ./alamode spaces 1

(this selects the first space)

$ ./alamode reports

1. ed55555a5551 Report1
2. ed55555a4441 Report2
3. ed55555a3331 Great Report

$ ./alamode reports 3

(this selects the 3rd report)

$ ./alamode queries

1. dd55555b5552 Query by the sea
2. dd55555b4442 This one is good
3. dd55555b3332 Untitled

$ ./alamode queries 2

(this selected the 2nd query)

$ ./alamode sql

(this opens the selected query in vim)

$ ./alamode run | jq .

(this take changes you made in vim and runs the new sql, you can pipe output to jq)


