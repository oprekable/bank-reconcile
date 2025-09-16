[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=oprekable_bank-reconcile&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=oprekable_bank-reconcile)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=oprekable_bank-reconcile&metric=bugs)](https://sonarcloud.io/summary/new_code?id=oprekable_bank-reconcile)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=oprekable_bank-reconcile&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=oprekable_bank-reconcile)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=oprekable_bank-reconcile&metric=coverage)](https://sonarcloud.io/summary/new_code?id=oprekable_bank-reconcile)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=oprekable_bank-reconcile&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=oprekable_bank-reconcile)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/oprekable/bank-reconcile)

# What is this?

This code repository demonstrates how to create an application using the Go programming language while adhering to common software development best practices.

The implementation covers the following aspects:

- Logic abstractions using SQL (SQLite) to handle the business logic of transaction reconciliation. For enhanced scalability, the database can be switched to another DB.
- Implementation of unit tests to maximize code coverage across various aspects, including SQL mocks, common interface mocks, file system mocks, time mocks, and IO mocks.
- Optimization of the codebase to meet [SonarQube](https://sonarcloud.io) standards, focusing on code quality, bug detection, code smells, duplicate lines, and coverage.
- Utilization of development tools such as:
  - Code linter ([golangci-lint](https://github.com/golangci/golangci-lint/cmd/golangci-lint))
  - Dependency injection ([wire](https://github.com/google/wire/cmd/wire))
  - Code mock ([mockery](https://github.com/vektra/mockery/v2))
  - Dead code checker ([deadcode](https://golang.org/x/tools/cmd/deadcode))
  - Go import checker ([goimports](https://golang.org/x/tools/cmd/goimports))
  - Go struct sorting ([fieldalignment](https://golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment))
  - Code and dependency vulnerability check ([govulncheck](https://golang.org/x/vuln/cmd/govulncheck))
  - Static code check ([staticcheck](https://honnef.co/go/tools/cmd/staticcheck)).
- SonarQube analysis via github action.
- Release application binaries using [goreleaser](https://github.com/goreleaser/goreleaser) via github action.
- Improved and simplified code structures.
- Generation of Go profiler (pprof files) to analyze the performance of the Go application.

# What service or application is it?

The application is designed and implemented as a transaction reconciliation system that identifies unmatched and discrepant transactions between internal data (system transactions) and external data (multiple bank statements), ensuring scalability.

The application will generate a command-line interface (CLI) application with the following features:

- Generate sample internal transactions and multiple bank transaction statements in CSV files. The files will be structured as follows:

```shell
.
└── sample
    ├── bank
    │   ├── bca
    │   │   └── bca_1744030099.csv
    │   ├── bni
    │   │   └── bni_1744030099.csv
    │   ├── bri
    │   │   └── bri_1744030099.csv
    │   ├── danamon
    │   │   └── danamon_1744030099.csv
    │   └── mandiri
    │       └── mandiri_1744030099.csv
    └── system
        └── 1744030099.csv
```

- Assuming that different banks have varying CSV formats, the current code accommodates this by providing samples for Bank BCA and Bank BNI, each with their own CSV header formats. Other banks will use the default format.
- Sample internal transaction CSV file:

```shell
TrxID,TransactionTime,Type,Amount
0004340e6a526fa4eaa35926241795f6,2025-04-07 12:24:03,DEBIT,3700
0039773a35b1ec6ebee0066fdef9b684,2025-04-07 10:05:53,CREDIT,78700
```

- Sample default format bank statement csv file:

```shell
UniqueIdentifier,Date,Amount
danamon-efbe04d119377bc5103af2234abd0188,2025-04-07,84000
danamon-86664183e0a3f2b96339fa71dabe9479,2025-04-07,-61000
```

- Sample format BCA bank statement csv file:

```shell
BCAUniqueIdentifier,BCADate,BCAAmount
bca-8a2684155034b82f6e042572aa788709,2025-04-07,78700
bca-fe9d6560f91287df981fac2a4fc1c773,2025-04-07,-47400
```

- Sample format BNI bank statement csv file:

```shell
BNIUniqueIdentifier,BNIDate,BNIAmount
bni-e84915bf6bf6f7d0325e628ee252b2a5,2025-04-07,82800
bni-80f73e01ab6ef44742bd56051a16f9f1,2025-04-07,14200
```

- The generated CSV files are structured based on specific configurations. For more details, please refer to the manual.
- The application should perform the reconciliation process using CSV files generated from sample commands or real transaction files. The reconciliation rules are as follows:
  - `TransactionTime` of internal transaction (in `datetime` format) == `Date` or `BCADate` or `BNIDate` of bank statement (in `date` format)
  - And `Amount` + `Type` in internal transaction == `Amount` or `BCAAmount` or `BNIAmount` of bank statement (`Type` DEBIT in internal transaction == negative value of `Amount` in bank statement)
- The results of the reconciliation process will display the following information:
  - Total number of transactions processed
  - Total number of matched transactions
  - Total number of unmatched transactions
    - Details of unmatched transactions:
      - System transaction details if missing in bank statement(s)
      - Bank statement details if missing in system transactions (grouped by bank)
    - Total discrepancies (sum of absolute differences in amount between matched transactions)
  - Generate CSV files for matched and unmatched transactions, structured as follows:

```shell
.
└── report
├── bank
│   └── not_matched
│       ├── bca_1744030812.csv
│       ├── bni_1744030812.csv
│       ├── bri_1744030812.csv
│       ├── danamon_1744030812.csv
│       └── mandiri_1744030812.csv
└── system
   ├── matched
   │       └── matched_1744030812.csv
   └── not_matched
       └── not_matched_1744030812.csv
```

# How to use the application?

## How to install?

### Download binary release

Please download the latest binary release from https://github.com/oprekable/bank-reconcile/releases and select the release file appropriate for your operating system.

### Via go install

Install application via command `go install github.com/oprekable/bank-reconcile@latest`.

## Application manual

- Run application with command `bank-reconcile` and we can see instruction such:

```shell
Simple Bank Reconciliation command line tool

Usage:
  bank-reconcile [flags]
  bank-reconcile [command]

Examples:
Generate sample 
	bank-reconcile sample --systemtrxpath=/tmp/data/sample/system --banktrxpath=/tmp/data/sample/bank --listbank=bca,bni,mandiri,bri,danamon --percentagematch=100 --amountdata=10000 --from=2025-04-14 --to=2025-04-14
Process data 
	bank-reconcile process --systemtrxpath=/tmp/data/sample/system --banktrxpath=/tmp/data/sample/bank --reportpath==/tmp/data/report --listbank=bca,bni,mandiri,bri,danamon --from=2025-04-14 --to=2025-04-14


Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  process     Process reconciliation data
  sample      Generate sample reconciliation data
  version     Get application version

Flags:
  -h, --help   help for bank-reconcile

Use "bank-reconcile [command] --help" for more information about a command.
```

- To get help of call sample feature, call `bank-reconcile sample --help`:

```shell
Generate sample reconciliation data of System Transactions and Bank Transactions

Usage:
  bank-reconcile sample [flags]

Aliases:
  sample, sa, s

Examples:
Generate sample 
	bank-reconcile sample --systemtrxpath=/tmp/data/sample/system --banktrxpath=/tmp/data/sample/bank --listbank=bca,bni,mandiri,bri,danamon --percentagematch=100 --amountdata=10000 --from=2025-04-14 --to=2025-04-14


Flags:
  -a, --amountdata int         amount system trx data sample to generate, bank trx will be 2 times of this amount (default 1000)
  -b, --banktrxpath string     Path location of Bank Transaction directory (default "/tmp/data/sample/bank")
  -g, --debug                  debug mode
  -d, --deleteoldfile          delete old sample files (default true)
  -f, --from string            from date (YYYY-MM-DD) (default "2025-04-14")
  -h, --help                   help for sample
  -l, --listbank strings       List bank accepted (default [bca,bni,mandiri,bri,danamon])
  -p, --percentagematch int    percentage of matched trx for data sample to generate (default 100)
  -i, --profiler               pprof active mode
  -o, --showlog                show logs
  -s, --systemtrxpath string   Path location of System Transaction directory (default "/tmp/data/sample/system")
  -z, --time_zone string       time zone settings (default "Asia/Jakarta")
  -t, --to string              to date (YYYY-MM-DD) (default "2025-04-14")
```

- To get help of call reconciliation process feature, call `bank-reconcile process --help`:

```shell
Process reconciliation data of System Transactions and Bank Transactions

Usage:
  bank-reconcile process [flags]

Aliases:
  process, pr, p

Examples:
Process data 
	bank-reconcile process --systemtrxpath=/tmp/data/sample/system --banktrxpath=/tmp/data/sample/bank --reportpath==/tmp/data/report --listbank=bca,bni,mandiri,bri,danamon --from=2025-04-14 --to=2025-04-14


Flags:
  -b, --banktrxpath string     Path location of Bank Transaction directory (default "/tmp/data/sample/bank")
  -g, --debug                  debug mode
  -d, --deleteoldfile          delete old report files (default true)
  -f, --from string            from date (YYYY-MM-DD) (default "2025-04-14")
  -h, --help                   help for process
  -l, --listbank strings       List bank accepted (default [bca,bni,mandiri,bri,danamon])
  -i, --profiler               pprof active mode
  -r, --reportpath string      Path location of Archive directory (default "/tmp/data/report")
  -o, --showlog                show logs
  -s, --systemtrxpath string   Path location of System Transaction directory (default "/tmp/data/sample/system")
  -z, --time_zone string       time zone settings (default "Asia/Jakarta")
  -t, --to string              to date (YYYY-MM-DD) (default "2025-04-14")
```

Here details of subcommand and arguments:

| Sub Command | Available Flags       | Default Value                                                         | Description                                                                                                                                                       |
|-------------|-----------------------|-----------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|             |                       |                                                                       | will display help information, instructions how to run the application                                                                                            |
|             | --help                |                                                                       | will display help information, instructions how to run the application                                                                                            |
| sample      | --help                |                                                                       | will display help information with list available arguments of `sample` subcommand                                                                                |
| sample      | -f, --from            | current date with format YYYY-MM-DD (2025-04-08)                      | start date to generate sample data                                                                                                                                |
| sample      | -t, --to              | current date with format YYYY-MM-DD (2025-04-08)                      | end date to generate sample data (if equals with start date means data for one day)                                                                               |
| sample      | -l, --listbank        | bca,bni,mandiri,bri,danamon                                           | list of bank in sample data)                                                                                                                                      |
| sample      | -p, --percentagematch | 100                                                                   | at least of percentage matched transaction (internal transaction vs bank statement), if sets 10 matched transactions will be 10% or more                          |
| sample      | -a, --amountdata      | 1000                                                                  | total internal transaction generated                                                                                                                              |
| sample      | -b, --banktrxpath     | Current working directory + `sample/bank` (/tmp/data/sample/bank)     | root directory path of generated sample of bank statements data files located                                                                                     |
| sample      | -s, --systemtrxpath   | Current working directory + `sample/system` (/tmp/data/sample/system) | root directory path of generated sample of internal transaction data files located                                                                                |
| sample      | -d, --deleteoldfile   | true                                                                  | when value == true, delete previous any directory or files in `--banktrxpath` or `--systemtrxpath`                                                                |
| sample      | -i, --profiler        | false                                                                 | when value == true, turn on profiler, will generate files `mem.pprof, mutex.pprof, cpu.pprof  trace.pprof, block.pprof, goroutine.pprof` in current working directory |
| sample      | -o, --showlog         | false                                                                 | when value == true, turn on verbose logs                                                                                                                          |
| sample      | -g, --debug           | false                                                                 | when value == true, generate SQLite file `sample.db`                                                                                                              |
| process     | --help                |                                                                       | will display help information with list available arguments of `process` subcommand                                                                               |
| process     | -f, --from            | current date with format YYYY-MM-DD (2025-04-08)                      | start date to process reconciliation data                                                                                                                         |
| process     | -t, --to              | current date with format YYYY-MM-DD (2025-04-08)                      | end date to process reconciliation data (if equals with start date means data for one day)                                                                        |
| process     | -l, --listbank        | bca,bni,mandiri,bri,danamon                                           | list of bank in process reconciliation)                                                                                                                           |
| process     | -s, --systemtrxpath   | Current working directory + `sample/system` (/tmp/data/sample/system) | root directory path of internal transaction source data files located                                                                                             |
| process     | -b, --banktrxpath     | Current working directory + `sample/bank` (/tmp/data/sample/bank)     | root directory path of bank statements source data files located                                                                                                  |
| process     | -r, --reportpath      | Current working directory + `report` (/tmp/data/report)               | root directory path of reconciliation result data files located                                                                                                   |
| process     | -d, --deleteoldfile   | true                                                                  | when value == true, delete previous any directory or files in `--reportpath`                                                                                      |
| process     | -i, --profiler        | false                                                                 | when value == true, turn on profiler, will generate files `mem.pprof, mutex.pprof, cpu.pprof  trace.pprof, block.pprof, goroutine.pprof` in current working directory |
| process     | -o, --showlog         | false                                                                 | when value == true, turn on verbose logs                                                                                                                          |
| process     | -g, --debug           | false                                                                 | when value == true, generate SQLite file `reconciliation.db`                                                                                                      |
| version     |                       |                                                                       | will display application version                                                                                                                                  |

### Example syntax of `sample` sub command :
```shell
mkdir -p /tmp/data
cd /tmp/data
bank-reconcile sample --from=2025-03-29 --to=2025-04-08 --deleteoldfile=true --showlog=true --listbank=bca,bni,mandiri,bri,danamon --profiler=true --debug=true --percentagematch=10 --amountdata=10 -s=/tmp/data/sample/system -b=/tmp/data/sample/bank
```
Application will generate sample of internal transaction and bank statement data from `2025-03-29` to `2025-04-08` for banks `bca,bni,mandiri,bri,danamon` with at least 10% reconciliation matched, CSV files of internal transactions will located at path `/tmp/data/sample/system` and bank statements will located at path `/tmp/data/sample/bank` any old directories and files under both paths will wiped, go profiler files `mem.pprof, mutex.pprof, cpu.pprof  trace.pprof, block.pprof, goroutine.pprof` will generated at `/tmp/data`, SQLite file `sample.db` will generated at `/tmp/data/sample.db` and verbose logs will displayed in terminal.

Expected application output will be:
```shell
2025-04-08T02:46:40.721952+07:00 | INFO  | *** sqlite connection loaded **** | component:NEWDBSQLITE | uptime:"85.125ΜS"

2025-04-08T02:46:40.722008+07:00 | INFO  | *** [start] application **** | uptime:"140.875ΜS"
         CONFIG        |            VALUE             
-----------------------+------------------------------
  -f --from            | 2025-03-29                   
  -t --to              | 2025-04-08                   
  -s --systemtrxpath   | /tmp/data/sample/system      
  -b --banktrxpath     | /tmp/data/sample/bank        
  -l --listbank        | bca,bni,mandiri,bri,danamon  
  -o --showlog         | true                         
  -g --debug           | true                         
  -i --profiler        | true                         
  -a --amountdata      | 10                           
  -p --percentagematch | 10                           
  -d --deleteoldfile   | true                         
2025-04-08T02:46:40.723269+07:00 | INFO  | *** [sample.NewSvc] DeleteDirectory SystemTRXPath **** | component:"SAMPLE SERVICE" | uptime:1.400792MS
2025-04-08T02:46:40.723326+07:00 | INFO  | *** [sample.NewSvc] DeleteDirectory BankTRXPath **** | component:"SAMPLE SERVICE" | uptime:1.457875MS
■ [1/5] Pre Process Generate Sample...  [0s] 2025-04-08T02:46:40.727494+07:00 | INFO  | *** [sample.NewDB] Exec Pre method in db **** | component:"SAMPLE SERVICE" | uptime:5.626542MS
2025-04-08T02:46:40.727516+07:00 | INFO  | *** [sample.NewSvc] RepoSample.Pre executed **** | component:"SAMPLE SERVICE" | uptime:5.647667MS
■ [2/5] Populate Trx Data...  [0s] 2025-04-08T02:46:40.728078+07:00 | INFO  | *** [sample.NewDB] Exec GetData method in db **** | component:"SAMPLE SERVICE" | uptime:6.2095MS
2025-04-08T02:46:40.728088+07:00 | INFO  | *** [sample.NewSvc] RepoSample.GetTrx executed **** | component:"SAMPLE SERVICE" | uptime:6.219958MS
■ [3/5] Post Process Generate Sample...  [0s] 2025-04-08T02:46:40.72811+07:00 | INFO  | *** [sample.NewSvc] RepoSample.Post executed **** | component:"SAMPLE SERVICE" | uptime:6.24175MS
■ [4/5] Parse Sample Data...  [0s] 2025-04-08T02:46:40.728133+07:00 | INFO  | *** [sample.NewSvc] populate systemTrxData & bankTrxData executed **** | component:"SAMPLE SERVICE" | uptime:6.265MS
■ [5/5] Export Sample Data to CSV files...  [0s] 2025-04-08T02:46:40.728648+07:00 | INFO  | *** [sample.NewSvc] save csv file /tmp/data/sample/system/1744055200.csv executed **** | component:"SAMPLE SERVICE" | uptime:6.780625MS
2025-04-08T02:46:40.728708+07:00 | INFO  | *** [sample.NewSvc] save csv file /tmp/data/sample/bank/bca/bca_1744055200.csv executed **** | component:"SAMPLE SERVICE" | uptime:6.84MS
2025-04-08T02:46:40.728963+07:00 | INFO  | *** [sample.NewSvc] save csv file /tmp/data/sample/bank/bni/bni_1744055200.csv executed **** | component:"SAMPLE SERVICE" | uptime:7.094583MS
2025-04-08T02:46:40.728976+07:00 | INFO  | *** [sample.NewSvc] save csv file /tmp/data/sample/bank/danamon/danamon_1744055200.csv executed **** | component:"SAMPLE SERVICE" | uptime:7.108875MS
2025-04-08T02:46:40.728976+07:00 | INFO  | *** [sample.NewSvc] save csv file /tmp/data/sample/bank/bri/bri_1744055200.csv executed **** | component:"SAMPLE SERVICE" | uptime:7.107917MS
2025-04-08T02:46:40.728995+07:00 | INFO  | *** [sample.NewSvc] save csv file /tmp/data/sample/bank/mandiri/mandiri_1744055200.csv executed **** | component:"SAMPLE SERVICE" | uptime:7.127292MS
                                                 
+------------+---------+-----------+------------------------------------------------------+
|  TYPE TRX  |  BANK   |   TITLE   |                                                      |
+------------+---------+-----------+------------------------------------------------------+
| System Trx | -       | Total Trx | 10                                                   |
+            +         +-----------+------------------------------------------------------+
|            |         | File Path | /tmp/data/sample/system/1744055200.csv               |
+------------+---------+-----------+------------------------------------------------------+
| Bank Trx   | danamon | Total Trx | 2                                                    |
+            +         +-----------+------------------------------------------------------+
|            |         | File Path | /tmp/data/sample/bank/danamon/danamon_1744055200.csv |
+            +---------+-----------+------------------------------------------------------+
|            | bca     | Total Trx | 3                                                    |
+            +         +-----------+------------------------------------------------------+
|            |         | File Path | /tmp/data/sample/bank/bca/bca_1744055200.csv         |
+            +---------+-----------+------------------------------------------------------+
|            | bri     | Total Trx | 1                                                    |
+            +         +-----------+------------------------------------------------------+
|            |         | File Path | /tmp/data/sample/bank/bri/bri_1744055200.csv         |
+            +---------+-----------+------------------------------------------------------+
|            | bni     | Total Trx | 2                                                    |
+            +         +-----------+------------------------------------------------------+
|            |         | File Path | /tmp/data/sample/bank/bni/bni_1744055200.csv         |
+            +---------+-----------+------------------------------------------------------+
|            | mandiri | Total Trx | 2                                                    |
+            +         +-----------+------------------------------------------------------+
|            |         | File Path | /tmp/data/sample/bank/mandiri/mandiri_1744055200.csv |
+------------+---------+-----------+------------------------------------------------------+

■ Done  [0s]                                     

-------- Memory Dump --------

     DESCRIPTION     | VALUE   
---------------------+---------
  Allocated          | 2.8 MB  
  Total Allocated    | 5.1 MB  
  Memory Allocations | 54 kB   
  Memory Frees       | 30 kB   
  Heap Allocated     | 2.8 MB  
  Heap System        | 7.6 MB  
  Heap In Use        | 5.5 MB  
  Heap Idle          | 2.1 MB  
  Heap OS Related    | 1.3 MB  
  Heap Objects       | 24 kB   
  Stack In Use       | 754 kB  
  Stack System       | 754 kB  
  Stack Span In Use  | 126 kB  
  Stack Cache In Use | 9.7 kB  
  Next GC cycle      | 5ms     
  Last GC cycle      | now     

2025-04-08T02:46:40.730998+07:00 | INFO  | *** [shutdown] application **** | uptime:9.129625MS
2025-04-08T02:46:40.731013+07:00 | INFO  | *** [shutdown] cli **** | uptime:9.144292MS
□ Done  [0s]
```


### Example syntax of `process` sub command :
```shell
cd /tmp/data
bank-reconcile process --from=2025-03-29 --to=2025-04-08 --deleteoldfile=true  --showlog=true --listbank=bca,bni,mandiri,bri,danamon --profiler=true --debug=true -s=/tmp/data/sample/system -b=/tmp/data/sample/bank -r=/tmp/data/report
```
Application will process transaction reconciliation matching of internal transaction and bank statement data from `2025-03-29` to `2025-04-08` for banks `bca,bni,mandiri,bri,danamon`, source of CSV files of internal transactions will loaded from path `/tmp/data/sample/system` and bank statements will loaded from path `/tmp/data/sample/bank`, go profiler files `mem.pprof, mutex.pprof, cpu.pprof  trace.pprof, block.pprof, goroutine.pprof` will generated report of matched and not matched transaction at `/tmp/data/report`, any old directories and files under `/tmp/data/report` will wiped, SQLite file `reconciliation.db` will generated at `/tmp/data/reconciliation.db` and verbose logs will displayed in terminal.

Expected application output will be:
```shell
2025-04-08T02:48:45.865412+07:00 | INFO  | *** sqlite connection loaded **** | component:NEWDBSQLITE | uptime:"130ΜS"
2025-04-08T02:48:45.865512+07:00 | INFO  | *** [start] application **** | uptime:"230.292ΜS"
        CONFIG       |            VALUE             
---------------------+------------------------------
  -f --from          | 2025-03-29                   
  -t --to            | 2025-04-08                   
  -s --systemtrxpath | /tmp/data/sample/system      
  -b --banktrxpath   | /tmp/data/sample/bank        
  -l --listbank      | bca,bni,mandiri,bri,danamon  
  -o --showlog       | true                         
  -g --debug         | true                         
  -i --profiler      | true                         
  -r --reportpath    | /tmp/data/report             

■ [1/7] Pre Process Generate Reconciliation...  [0s] 2025-04-08T02:48:45.869125+07:00 | INFO  | *** [process.NewDB] Exec Pre method in db **** | component:"PROCESS SERVICE" | uptime:3.843167MS
2025-04-08T02:48:45.869154+07:00 | INFO  | *** [process.NewSvc] GenerateReconciliation RepoProcess.Pre executed **** | component:"PROCESS SERVICE" | uptime:3.870959MS
■ [2/7] Parse System/Bank Trx Files...  [0s] 2025-04-08T02:48:45.869664+07:00 | INFO  | *** [process.NewSvc] parseSystemTrxFile parse - '/tmp/data/sample/system/1744055200.csv' executed **** | component:"PROCESS SERVICE" | uptime:4.381334MS
2025-04-08T02:48:45.869944+07:00 | INFO  | *** [process.NewSvc] parseSystemTrxFiles executed **** | component:"PROCESS SERVICE" | uptime:4.661542MS
2025-04-08T02:48:45.869972+07:00 | INFO  | *** [process.NewSvc] GenerateReconciliation parseSystemTrxFiles executed **** | component:"PROCESS SERVICE" | uptime:4.689042MS
2025-04-08T02:48:45.870017+07:00 | INFO  | *** [process.NewSvc] parseBankTrxFile parse (MANDIRI) - 'mandiri' executed **** | component:"PROCESS SERVICE" | uptime:4.734584MS
2025-04-08T02:48:45.870048+07:00 | INFO  | *** [process.NewSvc] parseBankTrxFile parse (BNI) - 'bni' executed **** | component:"PROCESS SERVICE" | uptime:4.764959MS
2025-04-08T02:48:45.870095+07:00 | INFO  | *** [process.NewSvc] parseBankTrxFile parse.ToBankTrxData (BNI) executed **** | component:"PROCESS SERVICE" | uptime:4.811875MS
2025-04-08T02:48:45.87022+07:00 | INFO  | *** [process.NewSvc] parseBankTrxFile parse (DANAMON) - 'danamon' executed **** | component:"PROCESS SERVICE" | uptime:4.937209MS
2025-04-08T02:48:45.870097+07:00 | INFO  | *** [process.NewSvc] parseBankTrxFile parse.ToBankTrxData (MANDIRI) executed **** | component:"PROCESS SERVICE" | uptime:4.813917MS
2025-04-08T02:48:45.870282+07:00 | INFO  | *** [process.NewSvc] parseBankTrxFile parse.ToBankTrxData (DANAMON) executed **** | component:"PROCESS SERVICE" | uptime:4.999292MS
2025-04-08T02:48:45.870181+07:00 | INFO  | *** [process.NewSvc] parseBankTrxFile parse (BRI) - 'bri' executed **** | component:"PROCESS SERVICE" | uptime:4.898417MS
2025-04-08T02:48:45.870329+07:00 | INFO  | *** [process.NewSvc] parseBankTrxFile parse.ToBankTrxData (BRI) executed **** | component:"PROCESS SERVICE" | uptime:5.045625MS
2025-04-08T02:48:45.870237+07:00 | INFO  | *** [process.NewSvc] parseBankTrxFile parse (BCA) - 'bca' executed **** | component:"PROCESS SERVICE" | uptime:4.956375MS
2025-04-08T02:48:45.870579+07:00 | INFO  | *** [process.NewSvc] parseBankTrxFile parse.ToBankTrxData (BCA) executed **** | component:"PROCESS SERVICE" | uptime:5.2965MS
2025-04-08T02:48:45.87062+07:00 | INFO  | *** [process.NewSvc] parseBankTrxFiles executed **** | component:"PROCESS SERVICE" | uptime:5.337334MS
2025-04-08T02:48:45.870643+07:00 | INFO  | *** [process.NewSvc] GenerateReconciliation parseBankTrxFiles executed **** | component:"PROCESS SERVICE" | uptime:5.36MS
■ [3/7] Import System Trx to DB...  [0s] 2025-04-08T02:48:45.871198+07:00 | INFO  | *** [process.NewDB] Exec ImportSystemTrx : range data (0 - 1) method in db **** | component:"PROCESS SERVICE" | uptime:5.915584MS
2025-04-08T02:48:45.871539+07:00 | INFO  | *** [process.NewDB] Exec ImportSystemTrx : range data (1 - 2) method in db **** | component:"PROCESS SERVICE" | uptime:6.256709MS
2025-04-08T02:48:45.871815+07:00 | INFO  | *** [process.NewDB] Exec ImportSystemTrx : range data (2 - 3) method in db **** | component:"PROCESS SERVICE" | uptime:6.531792MS
2025-04-08T02:48:45.872165+07:00 | INFO  | *** [process.NewDB] Exec ImportSystemTrx : range data (3 - 4) method in db **** | component:"PROCESS SERVICE" | uptime:6.882625MS
2025-04-08T02:48:45.872447+07:00 | INFO  | *** [process.NewDB] Exec ImportSystemTrx : range data (4 - 5) method in db **** | component:"PROCESS SERVICE" | uptime:7.164125MS
2025-04-08T02:48:45.872706+07:00 | INFO  | *** [process.NewDB] Exec ImportSystemTrx : range data (5 - 6) method in db **** | component:"PROCESS SERVICE" | uptime:7.423042MS
2025-04-08T02:48:45.872966+07:00 | INFO  | *** [process.NewDB] Exec ImportSystemTrx : range data (6 - 7) method in db **** | component:"PROCESS SERVICE" | uptime:7.682917MS
2025-04-08T02:48:45.873225+07:00 | INFO  | *** [process.NewDB] Exec ImportSystemTrx : range data (7 - 8) method in db **** | component:"PROCESS SERVICE" | uptime:7.941792MS
2025-04-08T02:48:45.873462+07:00 | INFO  | *** [process.NewDB] Exec ImportSystemTrx : range data (8 - 9) method in db **** | component:"PROCESS SERVICE" | uptime:8.179292MS
2025-04-08T02:48:45.873695+07:00 | INFO  | *** [process.NewDB] Exec ImportSystemTrx : range data (9 - 10) method in db **** | component:"PROCESS SERVICE" | uptime:8.412542MS
2025-04-08T02:48:45.873709+07:00 | INFO  | *** [process.NewSvc] GenerateReconciliation importReconcileSystemDataToDB executed **** | component:"PROCESS SERVICE" | uptime:8.425834MS
■ [4/7] Import Bank Trx to DB...  [0s] 2025-04-08T02:48:45.874045+07:00 | INFO  | *** [process.NewDB] Exec ImportBankTrx : range data (0 - 1) method in db **** | component:"PROCESS SERVICE" | uptime:8.762042MS
2025-04-08T02:48:45.87428+07:00 | INFO  | *** [process.NewDB] Exec ImportBankTrx : range data (1 - 2) method in db **** | component:"PROCESS SERVICE" | uptime:8.997542MS
2025-04-08T02:48:45.874524+07:00 | INFO  | *** [process.NewDB] Exec ImportBankTrx : range data (2 - 3) method in db **** | component:"PROCESS SERVICE" | uptime:9.240875MS
2025-04-08T02:48:45.874756+07:00 | INFO  | *** [process.NewDB] Exec ImportBankTrx : range data (3 - 4) method in db **** | component:"PROCESS SERVICE" | uptime:9.472875MS
2025-04-08T02:48:45.874999+07:00 | INFO  | *** [process.NewDB] Exec ImportBankTrx : range data (4 - 5) method in db **** | component:"PROCESS SERVICE" | uptime:9.716MS
2025-04-08T02:48:45.87523+07:00 | INFO  | *** [process.NewDB] Exec ImportBankTrx : range data (5 - 6) method in db **** | component:"PROCESS SERVICE" | uptime:9.947167MS
2025-04-08T02:48:45.875467+07:00 | INFO  | *** [process.NewDB] Exec ImportBankTrx : range data (6 - 7) method in db **** | component:"PROCESS SERVICE" | uptime:10.1845MS
2025-04-08T02:48:45.875699+07:00 | INFO  | *** [process.NewDB] Exec ImportBankTrx : range data (7 - 8) method in db **** | component:"PROCESS SERVICE" | uptime:10.416209MS
2025-04-08T02:48:45.875928+07:00 | INFO  | *** [process.NewDB] Exec ImportBankTrx : range data (8 - 9) method in db **** | component:"PROCESS SERVICE" | uptime:10.644959MS
2025-04-08T02:48:45.876159+07:00 | INFO  | *** [process.NewDB] Exec ImportBankTrx : range data (9 - 10) method in db **** | component:"PROCESS SERVICE" | uptime:10.876084MS
2025-04-08T02:48:45.876171+07:00 | INFO  | *** [process.NewSvc] GenerateReconciliation importReconcileBankDataToDB executed **** | component:"PROCESS SERVICE" | uptime:10.888375MS
■ [5/7] Mapping Reconciliation Data...  [0s] 2025-04-08T02:48:45.876458+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (0 - 4831) method in db **** | component:"PROCESS SERVICE" | uptime:11.174625MS
2025-04-08T02:48:45.876652+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (4831 - 9662) method in db **** | component:"PROCESS SERVICE" | uptime:11.369MS
2025-04-08T02:48:45.876842+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (9662 - 14493) method in db **** | component:"PROCESS SERVICE" | uptime:11.558625MS
2025-04-08T02:48:45.877072+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (14493 - 19324) method in db **** | component:"PROCESS SERVICE" | uptime:11.788834MS
2025-04-08T02:48:45.877246+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (19324 - 24155) method in db **** | component:"PROCESS SERVICE" | uptime:11.962584MS
2025-04-08T02:48:45.877421+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (24155 - 28986) method in db **** | component:"PROCESS SERVICE" | uptime:12.138042MS
2025-04-08T02:48:45.877589+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (28986 - 33817) method in db **** | component:"PROCESS SERVICE" | uptime:12.306209MS
2025-04-08T02:48:45.877748+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (33817 - 38648) method in db **** | component:"PROCESS SERVICE" | uptime:12.465292MS
2025-04-08T02:48:45.877904+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (38648 - 43479) method in db **** | component:"PROCESS SERVICE" | uptime:12.621042MS
2025-04-08T02:48:45.878069+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (43479 - 48311) method in db **** | component:"PROCESS SERVICE" | uptime:12.785625MS
2025-04-08T02:48:45.878224+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (48311 - 53142) method in db **** | component:"PROCESS SERVICE" | uptime:12.940917MS
2025-04-08T02:48:45.878388+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (53142 - 57973) method in db **** | component:"PROCESS SERVICE" | uptime:13.104875MS
2025-04-08T02:48:45.878558+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (57973 - 62804) method in db **** | component:"PROCESS SERVICE" | uptime:13.275167MS
2025-04-08T02:48:45.878722+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (62804 - 67635) method in db **** | component:"PROCESS SERVICE" | uptime:13.439167MS
2025-04-08T02:48:45.878878+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (67635 - 72466) method in db **** | component:"PROCESS SERVICE" | uptime:13.595209MS
2025-04-08T02:48:45.879029+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (72466 - 77297) method in db **** | component:"PROCESS SERVICE" | uptime:13.746375MS
2025-04-08T02:48:45.879184+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (77297 - 82128) method in db **** | component:"PROCESS SERVICE" | uptime:13.901334MS
2025-04-08T02:48:45.879336+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (82128 - 86959) method in db **** | component:"PROCESS SERVICE" | uptime:14.053084MS
2025-04-08T02:48:45.879481+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (86959 - 91790) method in db **** | component:"PROCESS SERVICE" | uptime:14.1985MS
2025-04-08T02:48:45.879743+07:00 | INFO  | *** [process.NewDB] Exec GenerateReconciliationMap : range amount (91790 - 96621) method in db **** | component:"PROCESS SERVICE" | uptime:14.459959MS
2025-04-08T02:48:45.879755+07:00 | INFO  | *** [process.NewSvc] GenerateReconciliation importReconcileMapToDB executed **** | component:"PROCESS SERVICE" | uptime:14.471459MS
■ [6/7] Generate Reconciliation Report Files...  [0s] 2025-04-08T02:48:45.879952+07:00 | INFO  | *** [process.NewDB] Exec GetReconciliationSummary method from db **** | component:"PROCESS SERVICE" | uptime:14.668584MS
2025-04-08T02:48:45.880099+07:00 | INFO  | *** [process.NewDB] Exec GetMatchedTrx method from db **** | component:"PROCESS SERVICE" | uptime:14.81575MS
2025-04-08T02:48:45.880467+07:00 | INFO  | *** [process.NewSvc] save csv file /tmp/data/report/system/matched/matched_1744055325.csv executed **** | component:"PROCESS SERVICE" | uptime:15.183625MS
2025-04-08T02:48:45.880588+07:00 | INFO  | *** [process.NewDB] Exec GetNotMatchedSystemTrx method from db **** | component:"PROCESS SERVICE" | uptime:15.304709MS
2025-04-08T02:48:45.88082+07:00 | INFO  | *** [process.NewSvc] save csv file /tmp/data/report/system/not_matched/not_matched_1744055325.csv executed **** | component:"PROCESS SERVICE" | uptime:15.536959MS
2025-04-08T02:48:45.880938+07:00 | INFO  | *** [process.NewDB] Exec GetNotMatchedBankTrx method from db **** | component:"PROCESS SERVICE" | uptime:15.654792MS
2025-04-08T02:48:45.881221+07:00 | INFO  | *** [process.NewSvc] save csv file /tmp/data/report/bank/not_matched/bca_1744055325.csv executed **** | component:"PROCESS SERVICE" | uptime:15.937834MS
2025-04-08T02:48:45.88131+07:00 | INFO  | *** [process.NewSvc] save csv file /tmp/data/report/bank/not_matched/bri_1744055325.csv executed **** | component:"PROCESS SERVICE" | uptime:16.027959MS
2025-04-08T02:48:45.881337+07:00 | INFO  | *** [process.NewSvc] save csv file /tmp/data/report/bank/not_matched/danamon_1744055325.csv executed **** | component:"PROCESS SERVICE" | uptime:16.053917MS
2025-04-08T02:48:45.881371+07:00 | INFO  | *** [process.NewSvc] save csv file /tmp/data/report/bank/not_matched/bni_1744055325.csv executed **** | component:"PROCESS SERVICE" | uptime:16.088167MS
2025-04-08T02:48:45.881486+07:00 | INFO  | *** [process.NewSvc] save csv file /tmp/data/report/bank/not_matched/mandiri_1744055325.csv executed **** | component:"PROCESS SERVICE" | uptime:16.203584MS
2025-04-08T02:48:45.88153+07:00 | INFO  | *** [process.NewSvc] GenerateReconciliation RepoProcess.GetReconciliationSummary executed **** | component:"PROCESS SERVICE" | uptime:16.247375MS
2025-04-08T02:48:45.881543+07:00 | INFO  | *** [process.NewSvc] GenerateReconciliation generateReconciliationSummaryAndFiles executed **** | component:"PROCESS SERVICE" | uptime:16.25975MS
                                                      
                DESCRIPTION                |   VALUE     
-------------------------------------------+-------------
  Total number of transactions processed   | 10          
  Total number of matched transactions     | 1           
  Total number of not matched transactions | 9           
  Sum amount all transactions              | 530.400,00  
  Sum amount matched transactions          | 96.500,00   
  Total discrepancies                      | 433.900,00  


               DESCRIPTION              |                           FILE PATH                             
----------------------------------------+-----------------------------------------------------------------
  Matched system transaction data       | /tmp/data/report/system/matched/matched_1744055325.csv          
  Missing system transaction data       | /tmp/data/report/system/not_matched/not_matched_1744055325.csv  
  Missing bank statement data - bca     | /tmp/data/report/bank/not_matched/bca_1744055325.csv            
  Missing bank statement data - bni     | /tmp/data/report/bank/not_matched/bni_1744055325.csv            
  Missing bank statement data - danamon | /tmp/data/report/bank/not_matched/danamon_1744055325.csv        
  Missing bank statement data - bri     | /tmp/data/report/bank/not_matched/bri_1744055325.csv            
  Missing bank statement data - mandiri | /tmp/data/report/bank/not_matched/mandiri_1744055325.csv        

■ Done  [0s]                                          

-------- Memory Dump --------

     DESCRIPTION     | VALUE   
---------------------+---------
  Allocated          | 3.8 MB  
  Total Allocated    | 6.0 MB  
  Memory Allocations | 66 kB   
  Memory Frees       | 30 kB   
  Heap Allocated     | 3.8 MB  
  Heap System        | 12 MB   
  Heap In Use        | 5.8 MB  
  Heap Idle          | 5.9 MB  
  Heap OS Related    | 5.7 MB  
  Heap Objects       | 36 kB   
  Stack In Use       | 918 kB  
  Stack System       | 918 kB  
  Stack Span In Use  | 136 kB  
  Stack Cache In Use | 9.7 kB  
  Next GC cycle      | 5ms     
  Last GC cycle      | now     

2025-04-08T02:48:45.884036+07:00 | INFO  | *** [shutdown] application **** | uptime:18.754292MS
2025-04-08T02:48:45.884083+07:00 | INFO  | *** [shutdown] cli **** | uptime:18.800375MS
▪ Done  [0s]
```

# What are the make commands that this code uses?
- Run `make` to display all available commands
```shell
help                           Show this help
download                       Download go.mod dependencies
install-tools                  Install required command line tools
generate                       Run go generate google wire dependency injection and mock files
go-lint                        Run golangci-lint (dry run)
staticcheck                    Run staticcheck
govulncheck                    Run govulncheck to check code vulnerability
godeadcode                     Run deadcode to check dead codes
development-checks             Download dependencies, install tools, generate codes, linter, code check (use it in code development cycle)
test                           Run unit tests and open coverage page
run                            Build and run application
echo-sample-args               Generate command syntax to run application to generate "sample"
run-sample                     Build and run application to generate "sample"
echo-process-args              Generate command syntax to run application to "process" data
run-process                    Build and run application to "process" data
run-version                    Build and run application to show application version
go-version                     To check current golang version in machine
go-env                         To check current golang environment variables in machine
check-profiler-block           To open pprof data of block profile
check-profiler-cpu             To open pprof data of cpu profile
check-profiler-memory          To open pprof data of memory profile
check-profiler-mutex           To open pprof data of mutex profile
check-profiler-trace           To open pprof data of trace profile
```
- Run `make development-checks` will help us in developments
- Any make command to chek go profiler are alias of:
  - `make check-profiler-block` = `go tool pprof -http=:8080 block.pprof`
  - `make check-profiler-cpu` = `go tool pprof -http=:8080 cpu.pprof`
  - `make check-profiler-memory` = `go tool pprof -http=:8080 mem.pprof`
  - `make check-profiler-mutex` = `go tool pprof -http=:8080 mutex.pprof`
  - `make check-profiler-trace` = `go tool trace -http=:8080 trace.pprof`