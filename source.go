package erikson
import "time"
import "sync"
import "errors"


type Node struct {
	Host string 
	Metrics map[string]float64 

}




type Source interface {
	C() chan []Node 
}



type Scraper interface {
	Scrape() ([]Node, error) 
}


func ScrapeAsync(s Scraper)  chan []Node {

	out := make(chan []Node) 
	
	f := func() {

	
		p, e := s.Scrape()

		if e != nil {
			out <- []Node{}
			return 
		}

		out <- p 
	}


	go f() 

	return out 
}



type joinScraper struct {
	a, b Scraper
} 

// fetches all scrapers within asynchronously and returns output
func (s *joinScraper) Scrape() ([]Node, error) {
	l :=  ScrapeAsync(s.a) 
	r := ScrapeAsync(s.b) 
	

	a := <- l 
	b := <- r 

	out := append(a, b...) 


	if len(out) == 0 {
		return out, errors.New("Empty State")
	}
	return out, nil 
}


//the behaviour of the outputted scraper
// is to run both scrapers aynchronously await both and join the result 
// remember you can also join the outputted scraper to eachother which scales very nicely
// so you can run 15 instances and only be bounded by the slowest call
// it ignores failures because there are multiple 
func JoinScrapers(a Scraper, b Scraper) Scraper {
	return &joinScraper{a, b}
}



type ScrapedSource struct {

	c chan []Node 
	t *time.Ticker

	scraper Scraper
}


func (s *ScrapedSource) Push() {
	p, e := s.scraper.Scrape()
	
	if e != nil {
		return 
	}

	s.c <- p 

}



func (s *ScrapedSource) loop() {

        for _ = range s.t.C {
                s.Push()
        }
}


func (s *ScrapedSource) Stop() {
	s.t.Stop()
	close(s.c) 
}

func (s *ScrapedSource) C() chan []Node {
	return s.C()
} 

func NewScrapedSource(i time.Duration, s Scraper)  ScrapedSource {
	

	c := make(chan []Node)
	t := time.NewTicker(i)

	return ScrapedSource {c, t, s} 	
}










