package main

import (
  "fmt"
  "os"
  "log"
  exec "os/exec"
  json "encoding/json"
  b64 "encoding/base64"
  io "io/ioutil"
  gjson "github.com/tidwall/gjson"

)

type SecretJSON struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Data       Data     `json:"data"`
  Metadata   MetaData `json:"metadata"`
}

type Data struct {
	Raw string `json:"raw"`
}

type MetaData struct {
	Name string `json:"name"`
}

func main(){

  if len(os.Args) == 1{
    fmt.Println("No secret provided")
  } else {

    // Get secret and base64 encode it
    secret := os.Args[1]
    b64encoded := b64.StdEncoding.EncodeToString([]byte(secret))


    // Generate Kubernetes Secret Manifest Struct
    data := Data{
  		Raw: b64encoded,
  	}
    metadata := MetaData{
      Name: "impenetrable",
    }
    secretStruct := SecretJSON{
      ApiVersion: "v1",
      Kind:       "Secret",
      Data:       data,
      Metadata:   metadata,
    }

    // Convert Manifest Struct to JSON
    secretJSON, err := json.Marshal(secretStruct)
    // fmt.Println(string(secretJSON))

    // Write Secret to tmpfile
    tmpSecret, err := io.TempFile("", "tmp-secret.json")
    if err != nil {
      log.Fatal(err)
    }
    defer os.Remove(tmpSecret.Name()) // clean up
    if _, err := tmpSecret.Write(secretJSON); err != nil {
      log.Fatal(err)
    }

    // Generate Sealed Destination
    tmpSealed, err := io.TempFile("", "tmp-sealed.json")
    if err != nil {
      log.Fatal(err)
    }
    defer os.Remove(tmpSealed.Name()) // clean up

    // Check environment for cert
    cert := os.Getenv("IMPENETRABLE_CERT")
    var txt string

    fmt.Println("")
    if cert != "" {
      fmt.Println("Using provided env cert..")
      txt = "kubeseal" + " --cert=" + cert + " --scope=cluster-wide " + "<" + tmpSecret.Name() + " >" + tmpSealed.Name()
    } else {
      fmt.Println("Will atempt to fetch cert from cluster..")
      txt = "kubeseal" + " --scope=cluster-wide " + "<" + tmpSecret.Name() + " >" + tmpSealed.Name()
    }

    // Execute kubeseal
    cmd := exec.Command(os.Getenv("SHELL"), "-c", txt)
    _ = cmd.Wait()
    cmd.CombinedOutput()

    // Read sealed secret
    sealedJSON, err := os.Open(tmpSealed.Name())
    if err != nil {
        log.Fatal(err)
    }
    defer sealedJSON.Close()
    body, err := io.ReadAll(sealedJSON)
     if err != nil {
      panic(err.Error())
    }
    bodyString := string(body)
    res := gjson.Get(bodyString, "spec.encryptedData.raw")
    fmt.Println("")
    if bodyString == "" {
      fmt.Println("Secret could not be sealed!")
    } else {
      fmt.Println(res)
    }
    fmt.Println("")
  }
}
