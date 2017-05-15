*DECOMISSIONED*
See [Basic TME Transformer](https://github.com/Financial-Times/basic-tme-transformer) instead

# special-reports-transformer

[![Circle CI](https://circleci.com/gh/Financial-Times/special-reports-transformer/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/special-reports-transformer/tree/master)

Retrieves SpecialReports taxonomy from TME and transforms the special reports to the internal UP json model.
The service exposes endpoints for getting all the special reports and for getting special reports by uuid.

# usage
`go get github.com/Financial-Times/special-reports-transformer`

`$GOPATH/bin/special-reports-transformer --port=8080 --base-url="http://localhost:8080/transformers/special-reports/" --tme-base-url="https://tme.ft.com" --tme-username="user" --tme-password="pass" --token="token"`
```
export|set PORT=8080
export|set BASE_URL="http://localhost:8080/transformers/special-reports/"
export|set TME_BASE_URL="https://tme.ft.com"
export|set TME_USERNAME="user"
export|set TME_PASSWORD="pass"
export|set TOKEN="token"
$GOPATH/bin/special-reports-transformer
```

With Docker:

`docker build -t coco/special-reports-transformer .`

`docker run -ti --env BASE_URL=<base url> --env TME_BASE_URL=<structure service url> --env TME_USERNAME=<user> --env TME_PASSWORD=<pass> --env TOKEN=<token> coco/special-reports-transformer`
