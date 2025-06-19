package main

import(
    "net/http"
    "log"
    "io/ioutil"
)

func main(){
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
    log.Fatal(http.ListenAndServeTLS(":443", "./cert/cert.pem", "./cert/key.pem", mux))
}

