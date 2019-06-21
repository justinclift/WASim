PostgreSQL database schema

This contains the PostgreSQL database for storing the operation logging data.

It was created using:

  $ pg_dump -c --if-exists wasim -C -s -O -x > schema.sql

To create the database, make sure you have PostgreSQL running locally
and your user has permission to connect, then run:

  $ psql template1 < schema.sql

This will:

  1. Connect to the PostgreSQL "template1" database (it doesn't get modified)
  2. Create the "wasim" database using that connection
  3. Re-connect to PostgreSQL, this time to the new wasim database
  4. Create the tables and structures needed for the operation log
