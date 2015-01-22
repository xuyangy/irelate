// Implements Relatable for Bam

package irelate

import (
	"io"

	boom "code.google.com/p/biogo.boom"
)

type Bam struct {
	*boom.Record
	source  uint32
	related []Relatable
	chrom   *string
}

func (a *Bam) Chrom() string {
	return *(a.chrom)
}

// cast to 32 bits.
func (a *Bam) Start() uint32 {
	return uint32(a.Record.Start())
}

func (a *Bam) End() uint32 {
	return uint32(a.Record.End())
}

func (a *Bam) Source() uint32 {
	return a.source
}

func (a *Bam) SetSource(src uint32) {
	a.source = src
}

func (a *Bam) AddRelated(b Relatable) {
	a.related = append(a.related, b)
}
func (a *Bam) Related() []Relatable {
	return a.related
}

func (a *Bam) Less(other Relatable) bool {
	if a.Chrom() != other.Chrom() {
		return a.Chrom() < other.Chrom()
	}
	return a.Start() < other.Start()
}

func BamToRelatable(file string) RelatableChannel {

	ch := make(chan Relatable, 16)

	go func() {
		b, err := boom.OpenBAM(file)
		if err != nil {
			panic(err)
		}
		defer b.Close()
		names := b.RefNames()
		for {
			rec, _, err := b.Read()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					panic(err)
				}
			}
			if rec.RefID() == -1 { // unmapped
				continue
			}
			bam := Bam{Record: rec, chrom: &(names[rec.RefID()]),
				related: make([]Relatable, 0, 2)}
			ch <- &bam
		}
		close(ch)
	}()
	return ch
}
