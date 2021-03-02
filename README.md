# Homework
The repository covers a lab, based on  "Git Audit Tool" homework.

## Solution details

### Core package
The business logic, which covers communication with the github.com API, is located in the package `github.com/polisko/gitcommits`.
I choose `GitHub GraphQL API` (https://docs.github.com/en/graphql) just because everybody talks about GraphQL as a future, so I wanted to touch it too :).

A small complication for using this API is the necessity of using the authorization token, even for "anonymous, read-only public access". For both clients (REST and CLI), it must be exported as environment variable AUTH_TOKEN. I've created one and shared via mail.

Unit tests are not fully implemented, for the lab purposes, there is only one indicated.

### REST server
For HTTP service (package `github.com/polisko/gitcommits/rest`), I choose simple REST API implementation. Using my favorite Gorilla/mux, I expose API:

```
/{owner}/{repo}/{branch}/{commit}
```
which return JSON with the results. The results can be consumed by Unmarshaling the type Result, existed in the core package, or any other JSON-consuming method.

### CLI
CLI is implemented in the `github.com/polisko/gitcommits/cli` package. Usage::
```
Usage of gc-cli:
  -b string
    	Branch
  -c string
    	commit hash (OID)
  -l	Whether to try locate local git repository in the current tree (default true)
  -logLevel int
    	Loglevel, default 4=Info (default 4)
  -o string
    	Repository owner
  -r string
    	Repository name
  -s	Whether to print shorter or longer output
  
  ```
It possible to specify all the parameters using options or it is also functionality implemented by using Go native git client (https://github.com/go-git/go-git), which tries to find local repository in the current path and fullfil parameters not provided with the arguments.

### Dockerfiles

There are Docker build files for both use cases, which uses builder image based on Go official docker to compile and runtime image with only produced binary, based on tiny alpine image to provide as small resulting image.

Exmaple using:

```bash
# build cli
docker build --tag quay.io/polisko/gitcommits-cli -f cli/Dockerfile .
#build REST server
docker build --tag quay.io/polisko/gitcommits -f rest/Dockerfile .

#run cli
docker run --name cli -v $(pwd):/app --rm -e AUTH_TOKEN=<YOUR AUTH_TOKEN here> quay.io/polisko/gitcommits-cli -s=false

Owner: polisko
Repository: gitcommits
Branch: main
Commit: 2e1aa67153f91d6d5fd042f23be4d685dedf86e0
Commits count: 1
********************
2e1aa67 2021-03-01 17:05:46 Pavel
Homework
********************

#run REST server
docker run -d --name rest -p 8080:8080 --rm -e AUTH_TOKEN=<YOUR AUTH_TOKEN here> quay.io/polisko/gitcommits

curl http://localhost:8080/wandera/puppet-prometheus/master/1093b8
```
```json
{
   "repository":{
      "ref":{
         "target":{
            "commit":{
               "nodes":[
                  {
                     "oid":"1093b8722651fb6c8cc1f1e943b77a1cad066707",
                     "committed_date":"2021-02-19T10:28:48Z",
                     "message":"Fix mongo exporter for newer versions (#4)\n\n* Create the extract folder for newer versions of mongo_exporter\r\n\r\n* Create the extract folder for newer versions of mongo_exporter\r\n\r\n* Create the extract folder for newer versions of mongo_exporter\r\n\r\n* Add create_extract_folder param to mongodb_exporter\r\n\r\n* fixup! Add create_extract_folder param to mongodb_exporter\r\n\r\n* Update mongo_exporter default version\r\n\r\n* fixup! Update mongo_exporter default version",
                     "author":{
                        "name":"Paraic Chung",
                        "email":"49403519+wandera-pchung@users.noreply.github.com"
                     }
                  }
               ],
               "total_count":1,
               "page_info":{
                  "end_cursor":"1093b8722651fb6c8cc1f1e943b77a1cad066707 0"
               }
            }
         }
      }
   }
}

```

### Kubernetes (Openshift)

Since my current work is on OpenShift (I somehow skipped pure K8s), I created openshift manifest yaml.
https://github.com/polisko/gitcommits/blob/main/rest/gc-openshift.yaml
I believe, that the only significant change to run on k8s cluster is to change openshift's native deploymentconfig to k8s deployment.

Service exposed temporarily:

http://gc-gc.th4ga-ocp-no-production-4f4fa0110bb0538f4a66f0b939386d36-0000.eu-de.containers.appdomain.cloud/wandera/puppet-prometheus/master/1093b8

## Conclusion

Due to restricted time to make this homework, code and documents provided is raw and simple. There's lot of space for refactoring (missing interface usage, better CLI implementation using Cobra/Viper standard, secure routes by https etc etc). Maybe choosing GrapQL was con from time scheduling point of view, because I spent hours to study new stuff, instead better coding.

PT