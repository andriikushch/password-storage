package hanlder

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/andriikushch/password-storage/repository"
	"io/ioutil"
	"log"
	"net/http"
)

type key struct {
	Key string `json:"key"`
}

type List struct {
	Items []string `json:"items"`
}

func NewRouter(r repository.Repository) *http.ServeMux {
	accountPost := func(writer http.ResponseWriter, request *http.Request) {
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		tmp := key{}
		err = json.Unmarshal(b, &tmp)

		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmpKey := sha256.Sum256([]byte(tmp.Key))
		key := tmpKey[:]

		l, err := r.GetAccountsList(key)

		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		res := &List{}
		for _, v := range l {
			res.Items = append(res.Items, v)
		}

		data, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = writer.Write(data)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/accounts", accountPost)

	return mux
}
