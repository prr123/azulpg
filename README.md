# pgsql

A program that parses input. The input is submitted to postgresql for execution via pgx.

The program allows for multiple input lines. Each line is aggregated into one string until the input parser sees an emptly line.

All commands are treated as sql commands unless the first command is 'show'.


