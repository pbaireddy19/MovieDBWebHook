package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"MovieDBWebhook/shared"
	"fmt"
	"bytes"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/get-movie-details", getMovieDetails).Methods("POST")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func getMovieDetails(w http.ResponseWriter, r *http.Request) {

	defaultMovieName := "The Godfather"
	decoder := json.NewDecoder(r.Body)
	var request shared.MovieSearchReq

	err := decoder.Decode(&request)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()
	defaultMovieName = request.MovieName
	fmt.Println(request.MovieName)

	req, err := http.NewRequest("GET", "http://www.omdbapi.com/?t="+ defaultMovieName +"&apikey="+API_KEY, bytes.NewBuffer([]byte("")))
	fmt.Println(req)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	hookResp := shared.WebHookResp {
		"Something went wrong !",
		"Something went wrong !",
		"get-movie-details",
	}

	fmt.Println(resp)
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		mResp := shared.MovieDBResp{}
		json.Unmarshal(body, &mResp)
		fmt.Println(mResp.Title)
		speech := mResp.Title + "is a" + mResp.Actors + "starer" + mResp.Genre + "movie released in " + mResp.Year + "It was directed by " + mResp.Director
		displayText := mResp.Title + "is a" + mResp.Actors + "starer" + mResp.Genre + "movie released in " + mResp.Year + ".It was directed by " + mResp.Director
		hookResp = shared.WebHookResp {
			speech,
			displayText,
			"get-movie-details",
		}
	}

	json.NewEncoder(w).Encode(hookResp)
}