# Get Started

## Specification
- language : go 1.18
- database : postgres v14
- utils :
  - makefile
 
## Service Structure
I use minimal heararcy of directory because of golang way to read the package. So on the root of the project, you can see many list of directories. Here the details
- `config`: to put all config code, such as env variable and db setup
- `contant`: to put all contants file
- `controller`: to handle all gin controllers and the request body validation, which will called by gin route
- `middleware`: to store all middlewares, currently only auth middlewares
- `migrations`: all db migration scripts
- `models`: all struct that represent database table as a models and its utility/helper
- `request`: all request struct for the apis
- `response`: all response struct for the apis
- `service`: all logic and business flow of the api, it will called by controllers
- `test`: directory to put all testcase
- `utils`: all utility code are here, such as helper to filter sql error, or error message to be retuned if some error are thrown


## Instalation
1. Run go mod download

After cloning the repository, go to the directory of the project and then run this command on terminal to install all used library to your local machine
  ```
  go mod download
  ```

2. Create postgres database

if you want to use default config, you have to name the database `dealls-dating-app`

3. Dulicate `.env.example` and change the name to `.env`
4. Change `DB_CONN` to your database connection
 ```
 DB_CONN="postgres://{{your_postgres_user}}:{{your_postgres_password}}@{{postgres_host}}:{{postgres_port}}/{{db_name}}?sslmode=disable"
 ```
5. if you already installed `make` on your local
   - run `make server` on terminal to run server for the first time and apply database migration
   - run `make t` to test all available testcase
   - run `make tv` to test with verbose log
6. if not, do manual run using go cli, go to project root directory 
   - run `go run main.go` to run server
   - change directory to test dir by running `cd test` command. And `go test . -count=1` to run all test
   - append `-v` to the end of test command. so it would look like `go test . -count=1 -v` to run test with verbose log
  
7. the test result with verbose
  ```
  go test ./test/ -count=1 -v
  === RUN   TestExplore
  --- PASS: TestExplore (0.02s)
  === RUN   TestExplore_NoProfile
  --- PASS: TestExplore_NoProfile (0.00s)
  === RUN   TestSwipe
  --- PASS: TestSwipe (0.00s)
  === RUN   TestSwipe_SwipeMySelf
  --- PASS: TestSwipe_SwipeMySelf (0.00s)
  === RUN   TestSwipe_InactiveUser
  --- PASS: TestSwipe_InactiveUser (0.00s)
  === RUN   TestSwipe_OnceADay
  --- PASS: TestSwipe_OnceADay (0.00s)
  === RUN   TestSwipe_SwipeLimit
  --- PASS: TestSwipe_SwipeLimit (0.00s)
  === RUN   TestSwipe_SwipeLimit_Premium
  --- PASS: TestSwipe_SwipeLimit_Premium (0.00s)
  === RUN   TestGetPackageByID
  --- PASS: TestGetPackageByID (0.00s)
  === RUN   TestGetPackagesList
  --- PASS: TestGetPackagesList (0.00s)
  === RUN   TestGetPackageByID_NotFound
  --- PASS: TestGetPackageByID_NotFound (0.00s)
  === RUN   TestBuyPackage
  --- PASS: TestBuyPackage (0.07s)
  === RUN   TestLoginAndRegister
  --- PASS: TestLoginAndRegister (0.13s)
  === RUN   TestRegister_DuplicateEmail
  --- PASS: TestRegister_DuplicateEmail (0.00s)
  === RUN   TestLogin_WrongPassword
  --- PASS: TestLogin_WrongPassword (0.13s)
  === RUN   TestLogin_WrongEmail
  --- PASS: TestLogin_WrongEmail (0.07s)
  PASS
  ok      github.com/praadit/dating-apps/test     0.657s
  ```
