**OUTDATED**
___

| Error Code | Meaning                                 | Comments                       |
|------------|-----------------------------------------|--------------------------------|
| 11         | MkdirAll error                          | main.go                        |
| 12         | Unable to touch logfile                 | main.go                        |
| 13         | Unable to return to pwd after execution | main.go                        |
| 14         | Unable to --MARK logfile                | main.go                        |
| 22         | Missing archive name                    |                                |
| 101        | Connection error                        | Mainly called from `Connect()` |
| 201        | Missing archive name                    | backup.go                      |
| 202        | Failed to create archive                | backup.go                      | 
| 203        | Missing database name                   | backup.go                      | 
| 204        | 
| 301        | Error while fetching primary keys       | constraints.go                 |
| 302        | Primary key scan failed                 | constraints.go                 |
| 303        | Foreign key query failed                | constraints.go                 |
| 304        | Foreign key scan failed                 | constraints.go                 |
| 305        | Table constraint query failed           | constraints.go                 |
| 306        | Table constraint scan failed            | constraints.go                 |
| 307        | Unique index query failed               | constraints.go                 |
| 308        | Unique index scan failed                | constraints.go                 |
| 401        | Ownership lookup failed                 | dbmetadata.go                  |
| 402        | SHOW failed                             | dbmetadata.go                  |
| 403        | Failed to query sequences               | dbmetadata.go                  |
| 404        | Failed to scan sequence row             | dbmetadata.go                  |
