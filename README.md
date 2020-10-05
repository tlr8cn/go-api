## Overview

This is my response to PEX's Web Backend Technical Challenge. This API steps through the fibonacci sequence on every 20 seconds. This value can be configured through code by changing the `const TIMEOUT` to the value in seconds of your choosing. 

## Testing

In order to run the code, you'll need to ensure that you have a local MySQL server installed and running locally. Then, you'll can create a database locally called `fibonacci`.

After that, simply clone this repo and run

   `cd go-api && ./go-api`

Then, you'll be able to make GET requests to the `/previous`, `/current`, and `/next` endpoints at `localhost:8080` to get the associated values at that time.

## Design

I built this API using Go with a MySQL database back-end. The application consists of 2 threads, one for handling requests, and one for stepping through the fibonacci sequence. I'm also using the GORM library to enforce the relational database schema through code using its automigration feature. 

GORM only allows for creation of tables through code. So I have 3 tables, one for each endpoint's resource, that each have an integer field holding the value of the resource at any given time. This model coupled with Gorilla's `RecoveryHandler` allow for the service to safely restart processing where it left off in case of failure.

If you keep the service running for long enough, the fibonacci sequence will eventually restart after an integer overflow is detected. 


