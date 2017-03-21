package main

import (
	"fmt"
	"os"
	"compress/gzip"
	"io/ioutil"
	"gopkg.in/mgo.v2"
	"log"
	"encoding/json"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				fmt.Printf("worker%d: Filename: %s\n", w.ID, work.FileName)
				processFilename(work.FileName)

			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func processFilename(fileName string){
	fmt.Printf("Filename: %s\n\n\n", fileName)

	session, err := mgo.Dial("mongodb")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	dbconn := session.DB("cloudtrail")
	dbconn.C("records").Insert()

	data := readZip(fileName)

	for record := range data["Records"] {
		//fmt.Printf("Request: %s\n", record["requestID"])
		storeResult(record)
	}

	//session, err := mgo.Dial("server1.example.com,server2.example.com")
	//if err != nil {
	//	panic(err)
	//}
	//defer session.Close()
	//
	//// Optional. Switch the session to a monotonic behavior.
	//session.SetMode(mgo.Monotonic, true)
	//
	//c := session.DB("test").C("people")
	//err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
	//	&Person{"Cla", "+55 53 8402 8510"})
	//if err != nil {
	//	log.Fatal(err)
	//}

	//result := Person{}
	//err = c.Find(bson.M{"name": "Ale"}).One(&result)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println("Phone:", result.Phone)
}

func storeResult(result map[string]string) {
	session, err := mgo.Dial("database")
	log.Print("Connected to db")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("workout_tracker").C("workouts")

	err = c.Insert(result)
	if err != nil {
		panic(err)
	}
}

func readZip(zipFile string) (map[string]string){
	handle, err := openFile(zipFile)
	if err != nil {
		fmt.Println("[ERROR] Opening file:", err)
	}

	zipReader, err := gzip.NewReader(handle)
	if err != nil {
		fmt.Println("[ERROR] New gzip reader:", err)
	}
	defer zipReader.Close()

	fileContents, err := ioutil.ReadAll(zipReader)
	if err != nil {
		fmt.Println("[ERROR] ReadAll:", err)
	}

	fmt.Printf("[INFO] Uncompressed contents: %s\n", fileContents)

	var i_json map[string]string
	err = json.Unmarshal(fileContents, &i_json)

	if err != nil {
		panic(err)
	}

	// ** Another way of reading the file **
	//
	// fileInfo, _ := handle.Stat()
	// fileContents := make([]byte, fileInfo.Size())
	// bytesRead, err := zipReader.Read(fileContents)
	// if err != nil {
	//     fmt.Println("[ERROR] Reading gzip file:", err)
	// }
	// fmt.Println("[INFO] Number of bytes read from the file:", bytesRead)

	closeFile(handle)
	return i_json
}

func openFile(fileToOpen string) (*os.File, error) {
	return os.OpenFile(fileToOpen, openFileOptions, openFilePermissions)
}

func closeFile(handle *os.File) {
	if handle == nil {
		return
	}

	err := handle.Close()
	if err != nil {
		fmt.Println("[ERROR] Closing file:", err)
	}
}

const openFileOptions int = os.O_CREATE | os.O_RDWR
const openFilePermissions os.FileMode = 0660