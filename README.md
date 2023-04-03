# go-adp
ADP Developer Library for Golang

## Installing
```bash
go get github.com/derickdiaz/go-adp@latest
```

## Worker API Service

### Authentication
Before initializing a service, the application must authentication using an ADPAuthenticationSystem. Current only OAuth is supported.
To obtain the required files, a CSR must be submitted to ADP: https://developers.adp.com/articles/general/generate-a-certificate-signing-request.
Once the CSR is approved, ADP will provide the following:

1. Signed PFX certificate
2. Key file

Use OpenSSL To convert the signed pfx certificate to a pem file.

```bash
openssl pkcs12 -in certificate.pfx -out certificate.pem -clcerts
```

Lastly, an ADP Representative will help create a service account. The base64 credentials follow the "username:password" format.

```golang
    auth, err := adp.NewOAuthAuthenticationSystem("certificate.pem", "certificate_key.key", "base64_credentials")
    if err != nil {
        log.Fatal(err)
    }
    // Generates OAuth Client Credentials
    err = auth.Authenticate()
```

## Using Worker API Service

Pulling all workers from the ADP Directory synchronously
```golang

import (
    "fmt"
    adp "github.com/derickdiaz/go-adp"
    "log"
)


func main() {
    auth, err := adp.NewOAuthAuthenticationSystem("certificate.pem", "certificate_key.key", "base64_credentials")
    if err != nil {
        log.Fatal(err)
    }
    err = auth.Authenticate()
    if err != nil {
        log.Fatal(err)
    }

    svc := adp.NewADPWorkerService(auth)
    workers, err := svc.ListWorkers()
    if err != nil {
        log.Fatal(err)
    }

    for _, worker := range workers.Workers {
        fmt.Printf(
            "AssociateOID: %s\n"+
                "First Name: %s\n"+
                "Middle Name: %s\n"+
                "Last Name: %s\n"+
                "Active: %v",
            worker.GetAssociateOID(),
            worker.GetFirstName(),
            worker.GetMiddleName(),
            worker.GetLastName(),
            worker.IsActive(),
        )

        assignment := worker.GetPrimaryWorkAssignment()
        if assignment != nil {
            fmt.Println("")
        }
        fmt.Printf(
            "Job Title: %s\n\n",
            assignment.GetJobTitle(),
        )
    } 
}
```

Pulling all workers from the ADP Directory using channels
``` golang
    ...
	svc := adp.NewADPWorkerService(auth)
	reciever := make(chan *adp.Worker)

	go svc.ListWorkersAsync(reciever, 50)

	for worker := range reciever {
		fmt.Println(worker.GetAssociateOID())
	}
    ...
```