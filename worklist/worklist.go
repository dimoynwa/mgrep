package worklist

type Entry struct {
	// Path to the file
	Path string
}

type Worklist struct {
	jobs chan Entry
}

func (w *Worklist) Add(e Entry) {
	w.jobs <- e
}

func (w *Worklist) Next() Entry {
	j := <-w.jobs
	return j
}

func New(buffSize int) Worklist {
	return Worklist{make(chan Entry, buffSize)}
}

func NewJob(path string) Entry {
	return Entry{path}
}

func (w *Worklist) Finalize(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		w.Add(Entry{""})
	}
}
