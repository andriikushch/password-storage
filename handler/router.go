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
			//todo handle
			log.Println(err)
			return
		}
		tmp := key{}
		err = json.Unmarshal(b, &tmp)

		log.Println(tmp.Key)
		if err != nil {
			//todo handle
			log.Println(err)
			return
		}

		tmpKey := sha256.Sum256([]byte(tmp.Key))
		key := tmpKey[:]

		l, err := r.GetAccountsList(key)

		if err != nil {
			//todo handle
			log.Println(err)
			return
		}
		res := &List{}
		for _, v := range l {
			res.Items = append(res.Items, v)
		}

		data, err := json.Marshal(res)
		if err != nil {
			//todo handle
			log.Println(err)
			return
		}

		_, err = writer.Write(data)
		if err != nil {
			//todo handle
			log.Println(err)
			return
		}

		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/accounts", accountPost)

	return mux
}
