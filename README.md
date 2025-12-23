# REQUIREMENTS gator app

* go installed
* postgres installed

# INSTALL

For install `gator` you should use `go install` in the root of the project

## CONFIGURATION

You should have a `.gatoconfig` file with user information and db url as a json

i.e.
{"current_user_name":"alice","db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"}

## FIRST STEPS

- go run . register alice
- go run . login alice
- go run . addfeed "TechCrunch" https://techcrunch.com/feed/
