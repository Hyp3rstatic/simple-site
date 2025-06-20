package main

import(
    "net/http"
    "log"
    "io/ioutil"
    "os/exec"
    "fmt"
    "regexp"
    "os"
)

func termCmd (cmd string) {
  args := regexp.MustCompile(` `).Split(cmd, -1)
  cmdOut := exec.Command(args[0], args[1:]...)
  stdout, err := cmdOut.Output()
  if err != nil {
      fmt.Println(err.Error())
      return
  }
  fmt.Println(string(stdout))
}

func main(){

	_, err := os.Stat("./cert")
  if ! os.IsNotExist(err) {
    termCmd("rm -rf ./cert")
  }

  termCmd("mkdir ./cert")
  termCmd("openssl genpkey -algorithm RSA -out ./cert/key.pem")
  termCmd("openssl req -new -batch -x509 -key ./cert/key.pem -out ./cert/cert.pem -days 7")
  termCmd("openssl x509 -text -noout -in ./cert/cert.pem")

  go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
      http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
  }))
  var mux *http.ServeMux = http.NewServeMux()
  mux.Handle("/", http.FileServer((http.Dir("./frontend"))))
  mux.HandleFunc("/verify-login", func(w http.ResponseWriter, r *http.Request){
      body, err := ioutil.ReadAll(r.Body)
      if err != nil{
          panic(err)
      }
      log.Println(string(body))
      http.Redirect(w, r, "https://"+r.Host+"/", http.StatusMovedPermanently)
  })

  go func () {
    log.Fatal(http.ListenAndServeTLS(":443", "./cert/cert.pem", "./cert/key.pem", mux))
  }()

  var input string
  for {
    fmt.Printf(":")
    fmt.Scanln(&input)
    if input == "exit" {
      break;
    }
  }
  
  termCmd("rm -rf ./cert")

}

