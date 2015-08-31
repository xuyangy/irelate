package irelate

import (
	"log"

	"github.com/brentp/irelate/interfaces"
	"github.com/brentp/vcfgo"
	"github.com/brentp/xopen"
)

type Variant struct {
	interfaces.IVariant
	source  uint32
	related []interfaces.Relatable
}

func NewVariant(v interfaces.IVariant, source uint32, related []interfaces.Relatable) *Variant {
	return &Variant{v, source, related}
}

func (v *Variant) AddRelated(r interfaces.Relatable) {
	if len(v.related) == 0 {
		v.related = make([]interfaces.Relatable, 0, 12)
	}
	v.related = append(v.related, r)
}

func (v *Variant) Related() []interfaces.Relatable {
	return v.related
}

func (v *Variant) SetSource(src uint32) { v.source = src }
func (v *Variant) Source() uint32       { return v.source }

func Vopen(f string) *vcfgo.Reader {
	rdr, err := xopen.Ropen(f)
	check(err)
	vcf, err := vcfgo.NewReader(rdr, true)
	if err != nil {
		log.Fatal(err)
	}
	return vcf
}

func StreamVCF(vcf *vcfgo.Reader) RelatableChannel {
	ch := make(RelatableChannel, 256)
	go func() {
		j := 0
		for {
			v := vcf.Read()
			if v == nil {
				break
			}
			ch <- &Variant{v, 0, nil}
			j++
			if j < 1000 {
				vcf.Clear()
				j = 0
			}
		}
		close(ch)
	}()
	return ch
}
