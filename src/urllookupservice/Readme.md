## Architecture
![Architecture](./Architecture.jpg)

As the above diagram depicts its a very simple URL lookup service.

It does the CRUD operation and uses in memory cache for fast lookup before querying Couchbase.

### Read Request
1. Check if present in Cache, return if gets a hit.
2. Otherwise, query Couchbase and update Cache

### Write Request
1. Write to Couchbase, and update/add the cache.

### Delete Request
1. Deletes from Couchbase and invalidates cache

## Build & Unit Test

## Test Log
* Check if `www.blockedsite.com` is a malware url:
```
C:\Users\mohsheik>curl http://localhost:8080/v1/lookupservice/urls/www.blockedsite.com
false
C:\Users\mohsheik>
```
It returns `false` as expected.

* Let us add it to malware url list
    * Before that check if it exists in Couchbase
    ```
    C:\Users\mohsheik>curl http://localhost:8091/pools/default/buckets/badurls/docs/www.blockedsite.com
    
    C:\Users\mohsheik>
  ```
  * Add `www.blockedsite.com` to malware url
  ```
    C:\Users\mohsheik>curl -X POST http://localhost:8080/v2/lookupservice/urls/www.blockedsite.com -v
    *   Trying ::1...
    * TCP_NODELAY set
    * Connected to localhost (::1) port 8080 (#0)
    > POST /v2/lookupservice/urls/www.blockedsite.com HTTP/1.1
    > Host: localhost:8080
    > User-Agent: curl/7.55.1
    > Accept: */*
    >
    < HTTP/1.1 200 OK
    < Date: Sat, 10 Oct 2020 00:51:27 GMT
    < Content-Length: 0
    <
    * Connection #0 to host localhost left intact
    
    C:\Users\mohsheik>
  ```
  * Check Couchbase again
  ```
  C:\Users\mohsheik>curl http://admin:password@localhost:8091/pools/default/buckets/badurls/docs/www.blockedsite.com
  {"meta":{"id":"www.blockedsite.com","rev":"3-163c7b4125f500005f82958e02000000","expiration":1602393486,"flags":33554432},"json":"{\"IsMalware\":true}","xattrs":{}}
  C:\Users\mohsheik>
  ```
  * Lets check if it is a malware url
  ```
  C:\Users\mohsheik>curl http://localhost:8080/v2/lookupservice/urls/www.blockedsite.com
  true
  C:\Users\mohsheik>
  ```
  * Delete it from malware list
  ```
  C:\Users\mohsheik>curl -X "DELETE" http://localhost:8080/v2/lookupservice/urls/www.blockedsite.com
  
  C:\Users\mohsheik>
  ```
  * Let check if key the is removed from CB
  ```
  C:\Users\mohsheik>curl http://admin:password@localhost:8091/pools/default/buckets/badurls/docs/www.blockedsite.com  -v
  *   Trying ::1...
  * TCP_NODELAY set
  * Connected to localhost (::1) port 8091 (#0)
  * Server auth using Basic with user 'admin'
  > GET /pools/default/buckets/badurls/docs/www.blockedsite.com HTTP/1.1
  > Host: localhost:8091
  > Authorization: Basic YWRtaW46cGFzc3dvcmQ=
  > User-Agent: curl/7.55.1
  > Accept: */*
  >
  < HTTP/1.1 404 Object Not Found
  < X-XSS-Protection: 1; mode=block
  < X-Permitted-Cross-Domain-Policies: none
  < X-Frame-Options: DENY
  < X-Content-Type-Options: nosniff
  < Server: Couchbase Server
  < Pragma: no-cache
  < Expires: Thu, 01 Jan 1970 00:00:00 GMT
  < Date: Sat, 10 Oct 2020 00:56:45 GMT
  < Content-Length: 0
  < Cache-Control: no-cache,no-store,must-revalidate
  <
  * Connection #0 to host localhost left intact
  
  C:\Users\mohsheik>
  ```
  as expected `404 Object Not Found`
  * Check again if url is a bad url
  ```
  C:\Users\mohsheik>curl http://localhost:8080/v2/lookupservice/urls/www.blockedsite.com
  false
  C:\Users\mohsheik>
  ```
 
 ```
 NORTHAMERICA+mohsheik@MININT-GQ0AU3F MINGW64 ~/go/src/urllookupservice (master)
 $ go build
 
 NORTHAMERICA+mohsheik@MININT-GQ0AU3F MINGW64 ~/go/src/urllookupservice (master)
 $ go test ./... -v
 ?       urllookupservice        [no test files]
 === RUN   TestParser
 === RUN   TestParser/test0
 === RUN   TestParser/test1
 === RUN   TestParser/test2
 --- PASS: TestParser (0.00s)
     --- PASS: TestParser/test0 (0.00s)
     --- PASS: TestParser/test1 (0.00s)
     --- PASS: TestParser/test2 (0.00s)
 PASS
 ok      urllookupservice/common (cached)
 ?       urllookupservice/datastore      [no test files]
 ?       urllookupservice/routers        [no test files]
 === RUN   TestIsMalwareUrl
 === RUN   TestIsMalwareUrl/returns_Pepper's_score
 Checking if url:https://git-scm.com:9090/download/win belongs to malware URLs
 Checking if url:https://www.google.com/ belongs to malware URLs
 --- PASS: TestIsMalwareUrl (4.68s)
     --- PASS: TestIsMalwareUrl/returns_Pepper's_score (4.68s)
 === RUN   TestAddUrl
 === RUN   TestAddUrl/returns_Pepper's_score
 Received request to add URL:https://git-scm.com:9090/download/win to malware url list
 Received request to add URL:https://www.google.com/ to malware url list
 --- PASS: TestAddUrl (4.66s)
     --- PASS: TestAddUrl/returns_Pepper's_score (4.66s)
 === RUN   TestRemoveUrl
 === RUN   TestRemoveUrl/returns_Pepper's_score
 Received request to delete URL:https://git-scm.com:9090/download/win to malware url list
 --- PASS: TestRemoveUrl (4.67s)
     --- PASS: TestRemoveUrl/returns_Pepper's_score (4.67s)
 PASS
 ok      urllookupservice/server 1.269s
 
 NORTHAMERICA+mohsheik@MININT-GQ0AU3F MINGW64 ~/go/src/urllookupservice (master)
 ```
 
 ### Start Server
 
 ```
 ./<binary> -port ":8081"
 ```

 ## Some Proposed Improvements 
 
Security, Containerization, better caching.
 
1. Use Zookeper, memcache or similar for chaching and exposing it as Cache as a service.
2. Creating Datastore as a service for Couchbase

