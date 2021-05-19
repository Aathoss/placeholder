package main

import (
	"net/http"
	"strconv"
	"strings"

	fonction "web/placeholder_web/function"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "02-01-2006 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	http.HandleFunc("/", placeholder)
	log.Info("Service démarré avec succès | Port : 9950")
	err := http.ListenAndServe(":9950", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
		return
	}
}

func placeholder(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	list := strings.Split(strings.Replace(url, "/", "", 1), "/")
	buffer, err := fonction.Do(list)
	if err != nil {
		log.Error(err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Error("unable to write image.")
		return
	}

	log.WithFields(logrus.Fields{
		"Requete": r.Host + r.URL.Path,
	}).Info("Génération réussi.")

}
