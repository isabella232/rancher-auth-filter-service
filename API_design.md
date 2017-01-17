#Domian Validation

##Test Environment Setup
- RANCHER SERVER : http://54.255.182.226:8080/
- DNS SERVER: 54.255.182.226
- HOST SERVER: 54.169.69.238

###DNS server setup






##1. API Document
    
    
    
- Create: POST /v1-domains/domains 
    - input `{domainName: "foo.com",projectid: "1a1"}` I will get the accountid from token
    - output `{domainName: "foo.com",domainName: "foo.com"}`
    

- List: GET /v1-domains/domains
- Get by ID: GET /v1-domains/domains/:id
- Have the user retry validation: POST /v1-domains/domain/:id?action=validate
- Delete: DELETE /v1-domains/domains/id
- Errors have specific fields, {type: "error", status: 422, code: "DomainAlreadyInUse", message: "domain.com has already been validated by another account"}


- Error Message `{Type: "error",status: "401", message: "error message"}`    



