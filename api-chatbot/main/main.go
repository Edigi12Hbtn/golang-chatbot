package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-chatbot/api-chatbot/mlp"
	"github.com/kljensen/snowball"
)

// Data - phrases and targets for training.
type Data struct {
	phrase []string
	target [][]string
}

// loadTrainingData - loads data to train NN.
func loadTrainingData(path string, data *Data) {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// formatting training data
	text := strings.ReplaceAll(string(content), "?", " ?")
	text = strings.ReplaceAll(text, "\x00", "") //removes null character
	text = strings.ReplaceAll(text, "  ", " ")
	text = strings.ReplaceAll(text, " ", "_")
	trainingData := strings.Split(text, "#")
	lentrainingData := len(trainingData)

	// obtaining training phrases and its respective targets
	for i := 1; i < lentrainingData; i++ {
		expression := strings.Split(trainingData[i], "(")
		data.phrase = append(data.phrase, expression[0])
		vals := strings.Split(strings.Split(expression[1], ")")[0], ",")
		data.target = append(data.target, vals)
	}
}

// steamData - Steams training or user chats.
func steamData(data *Data) {

	phrase := data.phrase
	lenPhrase := len(phrase)
	for i := 0; i < lenPhrase; i++ {
		words := strings.Split(phrase[i], "_")
		lenWords := len(words)
		for j := 0; j < lenWords; j++ {
			stemmed, err := snowball.Stem(words[j], "spanish", true)
			if err == nil {
				words[j] = stemmed
			} else {
				fmt.Println(words[j], "couldn't be stemmed.")
			}
		}
		data.phrase[i] = strings.Join(words, "_")
	}
}

// requestBody - to decode users messages.
type requestBody struct {
	Um string `json:"user-message"`
}

// launchAPI - lauch api on http://localhost:3000/sky-restaurant
func launchAPI(data *Data, net *mlp.Network, baseVects *baseVectors) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// setting cors to allow web-chatbot requests.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "user-message"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// handling post requests to /sky-restaurant route
	r.Post("/sky-restaurant", func(w http.ResponseWriter, r *http.Request) {

		var body requestBody

		// obtaining user message
		json.NewDecoder(r.Body).Decode(&body)
		um := body.Um
		fmt.Println("Phrase: ", um)

		data.phrase[0] = strings.ReplaceAll(um, "?", " ?")
		data.phrase[0] = strings.ReplaceAll(data.phrase[0], " ", "_")
		data.phrase[0] = strings.ReplaceAll(data.phrase[0], "'", "")

		// random delay to make chatbot more natural
		n := 4 + rand.Intn(3)
		time.Sleep(time.Duration(n) * time.Second)

		steamData(data)
		answer := answerUserMsg(data, net, baseVects)
		w.Write([]byte(answer))
		w.WriteHeader(http.StatusOK)
	})
	http.ListenAndServe(":3000", r)
}

// answerUserMsg - uses NN to predict user request and returns an answer.
func answerUserMsg(data *Data, net *mlp.Network, baseVects *baseVectors) string {

	vectorizedVals := vectorizer(data.phrase[0], nil, *baseVects)
	output := net.Predict(vectorizedVals.phrase)
	targets := make([]string, 0)

	Rows, _ := output.Dims()
	for k := 0; k < Rows; k++ {
		if output.At(k, 0) > 0.75 {
			targets = append(targets, baseVects.BaseTargets[k])
		}
	}
	sort.Strings(targets)

	fmt.Println("Stemmed phrase: ", data.phrase[0])
	fmt.Println("Predicted targets: ", targets)

	answer := getAnswer(targets)

	return answer
}

// getAnswer - interprets the user intention and returns response.
func getAnswer(targets []string) string {

	if len(targets) == 1 {
		if targets[0] == "greeting" {
			return "Hola, ¡bienvenid@ a Sky Restaurant! Cuentanos en qué te podemos ayudar"
		} else if targets[0] == "liked" {
			return "Nos alegra que hayas disfrutado de nuestros servicios, esperamos verte pronto de nuevo"
		} else if targets[0] == "disliked" {
			return "Lamentamos que tu experiencia no haya sido la mejor, tomaremos nota de ello para asegurarnos que tu próxima experiencia sea mejor"
		}
	} else if len(targets) == 3 {
		opts := make([][]string, 4)
		opts[0] = []string{"food", "order", "pizza"}
		opts[1] = []string{"food", "hamburger", "order"}
		opts[2] = []string{"food", "order", "salad"}
		opts[3] = []string{"food", "order", "soda"}

		for i := 0; i < 4; i++ {
			areEqual := true
			for j := 0; j < 3; j++ {
				if targets[j] != opts[i][j] {
					areEqual = false
					break
				}
			}
			if areEqual {
				var ans string

				if targets[2] == "pizza" {
					ans = "¡Claro que si! Tendremos tu pizza lista en unos minutos. Esperamos la disfrutes"
				} else if targets[2] == "salad" {
					ans = "Por supuesto, tu ensalada estará lista en unos instantes. ¡Que la disfrutes!"
				} else if targets[2] == "soda" {
					ans = "En unos instantes te entregaremos tu gaseosa, ¡disfrútala!"
				} else {
					ans = "Tu hamburguesa estará prontamente terminada, esperamos que te guste"
				}
				return ans
			}
		}
	}

	return "No entendimos tu comentario. Cuéntanos nuevamente si deseas ordenar algo o opinar sobre tu experiencia"
}

func main() {

	data := Data{phrase: make([]string, 0), target: make([][]string, 0)}
	baseVects := baseVectors{}

	mode := flag.String("mode", "p", "Insert mode: 't' for training 'p' for predict")
	flag.Parse()

	net := mlp.CreateNetwork(33, 27, 9, 0.1)

	if *mode == "t" {
		loadTrainingData("data/chats", &data)
		steamData(&data)
		baseVects = genBaseVectors(&data)

		lentrainingData := len(data.phrase)
		for j := 0; j < 200; j++ {
			for i := 0; i < lentrainingData; i++ {
				vectorizedVals := vectorizer(data.phrase[i], data.target[i], baseVects)
				net.Train(vectorizedVals.phrase, vectorizedVals.targets)
			}
		}
		SaveTrainedNN(net, baseVects)
		fmt.Println("The neural network was trained.")
	} else {
		rand.Seed(time.Now().UnixNano())
		data.phrase = append(data.phrase, "")
		LoadTrainedNN(&net, &baseVects)
		launchAPI(&data, &net, &baseVects)
	}
}
